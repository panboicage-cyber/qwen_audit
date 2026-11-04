package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}
func TestCuratedWorkspaceBaseline(t *testing.T) {
	root := repoRoot(t)
	_, _, queue, _, stages, err := check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(stages) != 2 || !queue.Items[0].Enabled || queue.Items[1].Enabled || stages[0].Stage.TargetHead != "0af4707d66bb30aef09fe80316c408568e256e53" {
		t.Fatalf("unexpected staged baseline: %+v", stages)
	}
	var out, errs bytes.Buffer
	if code := dispatch([]string{"-root", root, "-profile", "gen-type-go-pr6", "context"}, &out, &errs); code != 0 {
		t.Fatalf("context exit=%d %s", code, errs.String())
	}
	var ctx Context
	if err := json.Unmarshal(out.Bytes(), &ctx); err != nil {
		t.Fatal(err)
	}
	if ctx.ExecutionState != "NOT_EXECUTED" || ctx.SourceVerified || ctx.OwnerApproval || ctx.MemoryAccess != "READ_ONLY" {
		t.Fatal("unsafe context claim")
	}
	out.Reset()
	errs.Reset()
	if code := dispatch([]string{"-root", root, "-profile", "kol-merkle-pr4", "context"}, &out, &errs); code != 2 {
		t.Fatal("inactive reference stage accepted")
	}
}
func cloneMinimal(t *testing.T) string {
	t.Helper()
	root := repoRoot(t)
	tmp := t.TempDir()
	files := []string{"staging_audit/workspace.json", "staging_audit/queue.json", "staging_audit/branches/gen-type-go-pr6.json", "staging_audit/branches/kol-merkle-pr4.json", "profiles/gen-type-go-pr6/target.lock.json", "scripts/replay-gen-type-go-pr6.sh", "scripts/replay-kol-merkle-pr4.sh", "prompts/gen-type-go-pr6.md", "prompts/staged-gen-type-go-pr6.md", "prompts/kol-merkle-pr4.md", "memory/catalog.json", "memory/SHA256SUMS", "scripts/setup-master-audit-workspace.sh", "scripts/setup-auditgo.sh", "scripts/auditgo-smoke.sh", "scripts/repo-memory-smoke.sh", "README.md", "AGENTS.md", "AUDIT_POLICY.md"}
	for _, f := range files {
		b, err := os.ReadFile(filepath.Join(root, f))
		if err != nil {
			t.Fatal(err)
		}
		p := filepath.Join(tmp, f)
		if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	return tmp
}
func TestRejectTamperedStageFile(t *testing.T) {
	root := cloneMinimal(t)
	p := filepath.Join(root, "profiles/gen-type-go-pr6/target.lock.json")
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(" ")
	f.Close()
	if _, _, _, _, _, err := check(root); err == nil || !strings.Contains(err.Error(), "SHA256 mismatch") {
		t.Fatalf("expected fail closed, got %v", err)
	}
}
func TestRejectDuplicateEscapedKey(t *testing.T) {
	root := cloneMinimal(t)
	p := filepath.Join(root, "staging_audit/queue.json")
	b, _ := os.ReadFile(p)
	b = bytes.Replace(b, []byte(`"strategy":`), []byte(`"\u0073trategy": "ORDERED_EXPLICIT_SELECTION", "strategy":`), 1)
	os.WriteFile(p, b, 0600)
	if _, _, _, _, _, err := check(root); err == nil || !strings.Contains(err.Error(), "duplicate decoded key") {
		t.Fatalf("expected duplicate key rejection, got %v", err)
	}
}
func TestRejectForgedAuthority(t *testing.T) {
	root := cloneMinimal(t)
	p := filepath.Join(root, "staging_audit/workspace.json")
	b, _ := os.ReadFile(p)
	b = bytes.Replace(b, []byte(`"allow_merge": false`), []byte(`"allow_merge": true`), 1)
	os.WriteFile(p, b, 0600)
	if _, _, _, _, _, err := check(root); err == nil || !strings.Contains(err.Error(), "workspace contract") {
		t.Fatalf("expected authority mismatch, got %v", err)
	}
}
func TestRejectSymlinkComponents(t *testing.T) {
	root := cloneMinimal(t)
	os.Remove(filepath.Join(root, "staging_audit/queue.json"))
	if err := os.Symlink(filepath.Join(repoRoot(t), "staging_audit/queue.json"), filepath.Join(root, "staging_audit/queue.json")); err != nil {
		t.Skip(err)
	}
	if _, _, _, _, _, err := check(root); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected symlink rejection, got %v", err)
	}
}
func TestUnsupportedCommandsAreRefused(t *testing.T) {
	for _, cmd := range []string{"save", "learn", "capture", "execute", "approve", "merge"} {
		if code := dispatch([]string{cmd}, &bytes.Buffer{}, &bytes.Buffer{}); code != 2 {
			t.Fatalf("unsupported %s accepted", cmd)
		}
	}
}
