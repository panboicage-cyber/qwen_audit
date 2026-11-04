package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleCatalog = `{"schema":"qwen-audit.repo-memory/v1","ownership":"repository_curated","access":"READ_ONLY","capture":"DISABLED","records":[{"id":"audit.test.fact","kind":"POLICY","summary":"Known source-derived fact","detail":"No agent-supplied memory writes are permitted.","source_path":"AUDIT_POLICY.md","source_ref":"af6609fae097efe80742025b0b5bf27b9acfccea"}]}`

func fixture(t *testing.T, raw string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "memory")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "catalog.json"), []byte(raw), 0600); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(raw))
	manifest := hex.EncodeToString(sum[:]) + "  catalog.json\n"
	if err := os.WriteFile(filepath.Join(dir, "SHA256SUMS"), []byte(manifest), 0600); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestReadOnlyCommandWorkflow(t *testing.T) {
	root := fixture(t, sampleCatalog)
	before, _ := os.ReadFile(filepath.Join(root, "memory", "catalog.json"))
	for _, tc := range []struct {
		cmd    []string
		marker string
	}{
		{[]string{"verify"}, "REPO_MEMORY_VERIFY_PASS"},
		{[]string{"doctor"}, "Capture: DISABLED"},
		{[]string{"list"}, "audit.test.fact"},
		{[]string{"get", "audit.test.fact"}, "No agent-supplied memory writes"},
		{[]string{"context"}, "data, not tool authorization"},
	} {
		var out, stderr bytes.Buffer
		args := append([]string{"-root", root}, tc.cmd...)
		if status := dispatch(args, &out, &stderr); status != 0 {
			t.Fatalf("%v exit=%d error=%s", tc.cmd, status, stderr.String())
		}
		if !strings.Contains(out.String(), tc.marker) {
			t.Fatalf("%v missing %s: %s", tc.cmd, tc.marker, out.String())
		}
	}
	after, _ := os.ReadFile(filepath.Join(root, "memory", "catalog.json"))
	if !bytes.Equal(before, after) {
		t.Fatal("read commands changed catalog")
	}
}

func TestMutationCommandsRejected(t *testing.T) {
	root := fixture(t, sampleCatalog)
	for _, cmd := range []string{"save", "append", "learn", "capture", "delete", "update", "init", "export"} {
		var out, stderr bytes.Buffer
		if code := dispatch([]string{"-root", root, cmd}, &out, &stderr); code != 2 {
			t.Fatalf("command %q code=%d", cmd, code)
		}
	}
}

func TestTamperedManifestBlocks(t *testing.T) {
	root := fixture(t, sampleCatalog)
	if err := os.WriteFile(filepath.Join(root, "memory", "catalog.json"), []byte(sampleCatalog+" "), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readCatalog(root); err == nil || !strings.Contains(err.Error(), "mismatch") {
		t.Fatalf("expected hash failure: %v", err)
	}
}

func TestDecodedDuplicateKeyBlocks(t *testing.T) {
	bad := strings.Replace(sampleCatalog, `"capture":"DISABLED"`, `"capture":"DISABLED","\u0063apture":"DISABLED"`, 1)
	root := fixture(t, bad)
	if _, _, err := readCatalog(root); err == nil || !strings.Contains(err.Error(), "duplicate decoded") {
		t.Fatalf("expected decoded duplicate failure: %v", err)
	}
}

func TestUnknownFieldBlocks(t *testing.T) {
	bad := strings.Replace(sampleCatalog, `"access":"READ_ONLY"`, `"access":"READ_ONLY","agent_can_save":true`, 1)
	root := fixture(t, bad)
	if _, _, err := readCatalog(root); err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("expected unknown field failure: %v", err)
	}
}

func TestPathTraversalBlocks(t *testing.T) {
	bad := strings.Replace(sampleCatalog, `"source_path":"AUDIT_POLICY.md"`, `"source_path":"../secrets"`, 1)
	root := fixture(t, bad)
	if _, _, err := readCatalog(root); err == nil || !strings.Contains(err.Error(), "provenance") {
		t.Fatalf("expected provenance failure: %v", err)
	}
}

func TestMissingCatalogBlocks(t *testing.T) {
	var out, stderr bytes.Buffer
	if code := dispatch([]string{"-root", t.TempDir(), "context"}, &out, &stderr); code != 1 || !strings.Contains(stderr.String(), "REPO_MEMORY_BLOCKED") {
		t.Fatalf("unexpected: code=%d %s", code, stderr.String())
	}
	if out.Len() != 0 {
		t.Fatal("missing catalog produced context")
	}
}

func TestSymlinkCatalogBlocks(t *testing.T) {
	root := fixture(t, sampleCatalog)
	target := filepath.Join(root, "memory", "original.json")
	os.Rename(filepath.Join(root, "memory", "catalog.json"), target)
	if err := os.Symlink(target, filepath.Join(root, "memory", "catalog.json")); err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}
	if _, _, err := readCatalog(root); err == nil || !strings.Contains(err.Error(), "not a regular") {
		t.Fatalf("expected symlink rejection: %v", err)
	}
}

func TestRecordsMustBeOrdered(t *testing.T) {
	bad := strings.Replace(sampleCatalog, `"records":[`, `"records":[{"id":"z.first","kind":"POLICY","summary":"Fact","detail":"Detail","source_path":"README.md","source_ref":"af6609fae097efe80742025b0b5bf27b9acfccea"},`, 1)
	root := fixture(t, bad)
	if _, _, err := readCatalog(root); err == nil || !strings.Contains(err.Error(), "unsorted") {
		t.Fatalf("expected order failure: %v", err)
	}
}
