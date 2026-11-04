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

// Rewrite one boolean in an isolated copy of the baseline. For a target lock,
// rebind its SHA in the stage to test the actual schema check rather than an
// unrelated digest mismatch. No canonical staging or source files are edited.
func rewriteBooleanFixture(t *testing.T, root, path, field string, itemIndex int, replacement string) {
	t.Helper()
	file := filepath.Join(root, path)
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		t.Fatal(err)
	}
	if itemIndex >= 0 {
		var items []map[string]json.RawMessage
		if err := json.Unmarshal(obj["items"], &items); err != nil {
			t.Fatal(err)
		}
		if itemIndex >= len(items) {
			t.Fatalf("queue index %d out of bounds", itemIndex)
		}
		if _, ok := items[itemIndex][field]; !ok {
			t.Fatalf("fixture key missing: %s", field)
		}
		if replacement == "" {
			delete(items[itemIndex], field)
		} else {
			items[itemIndex][field] = json.RawMessage(replacement)
		}
		obj["items"], err = json.Marshal(items)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		if _, ok := obj[field]; !ok {
			t.Fatalf("fixture key missing: %s", field)
		}
		if replacement == "" {
			delete(obj, field)
		} else {
			obj[field] = json.RawMessage(replacement)
		}
	}
	changed, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	changed = append(changed, '\n')
	if err := os.WriteFile(file, changed, 0600); err != nil {
		t.Fatal(err)
	}
	if path == "profiles/gen-type-go-pr6/target.lock.json" {
		stageFile := filepath.Join(root, "staging_audit/branches/gen-type-go-pr6.json")
		stageData, err := os.ReadFile(stageFile)
		if err != nil {
			t.Fatal(err)
		}
		var stage map[string]json.RawMessage
		if err := json.Unmarshal(stageData, &stage); err != nil {
			t.Fatal(err)
		}
		stage["lock_sha256"], err = json.Marshal(sha(changed))
		if err != nil {
			t.Fatal(err)
		}
		stageData, err = json.MarshalIndent(stage, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(stageFile, append(stageData, '\n'), 0600); err != nil {
			t.Fatal(err)
		}
	}
}

func TestAllBooleanContractsRequireExplicitJSONLiterals(t *testing.T) {
	fields := []struct {
		scope string
		path  string
		key   string
		item  int
	}{
		{"workspace", "staging_audit/workspace.json", "allow_auto_execute", -1},
		{"workspace", "staging_audit/workspace.json", "allow_runtime_memory_writes", -1},
		{"workspace", "staging_audit/workspace.json", "allow_merge", -1},
		{"workspace", "staging_audit/workspace.json", "allow_publish", -1},
		{"queue-active", "staging_audit/queue.json", "enabled", 0},
		{"queue-disabled", "staging_audit/queue.json", "enabled", 1},
		{"stage-active", "staging_audit/branches/gen-type-go-pr6.json", "requires_explicit_execution", -1},
		{"stage-active", "staging_audit/branches/gen-type-go-pr6.json", "allow_upstream_mutation", -1},
		{"stage-active", "staging_audit/branches/gen-type-go-pr6.json", "allow_merge", -1},
		{"stage-active", "staging_audit/branches/gen-type-go-pr6.json", "allow_release", -1},
		{"stage-disabled", "staging_audit/branches/kol-merkle-pr4.json", "requires_explicit_execution", -1},
		{"stage-disabled", "staging_audit/branches/kol-merkle-pr4.json", "allow_upstream_mutation", -1},
		{"stage-disabled", "staging_audit/branches/kol-merkle-pr4.json", "allow_merge", -1},
		{"stage-disabled", "staging_audit/branches/kol-merkle-pr4.json", "allow_release", -1},
		{"lock", "profiles/gen-type-go-pr6/target.lock.json", "upstream_mutations_authorized", -1},
		{"lock", "profiles/gen-type-go-pr6/target.lock.json", "merge_authorized", -1},
		{"lock", "profiles/gen-type-go-pr6/target.lock.json", "release_authorized", -1},
		{"lock", "profiles/gen-type-go-pr6/target.lock.json", "new_authority_authorized", -1},
	}
	mutations := []struct {
		name        string
		replacement string
		want        string
	}{
		{"missing", "", "missing required boolean"},
		{"null", "null", "must be explicit JSON boolean"},
		{"quoted-false", `"false"`, "must be explicit JSON boolean"},
		{"numeric-zero", "0", "must be explicit JSON boolean"},
	}
	for _, f := range fields {
		for _, m := range mutations {
			t.Run(f.scope+"/"+f.key+"/"+m.name, func(t *testing.T) {
				root := cloneMinimal(t)
				rewriteBooleanFixture(t, root, f.path, f.key, f.item, m.replacement)
				_, _, _, _, _, err := check(root)
				if err == nil || !strings.Contains(err.Error(), f.key) {
					t.Fatalf("expected rejection for %s, got %v", f.key, err)
				}
				// Go already rejects strings and numbers for bool fields during
				// typed decode; the repair closes missing/null bypasses. Preserve
				// either valid rejection diagnostic rather than requiring a new one.
				diagnostic := err.Error()
				switch m.name {
				case "missing":
					if !strings.Contains(diagnostic, "missing required boolean") && !strings.Contains(diagnostic, "missing lock field") {
						t.Fatalf("expected missing-field rejection, got %v", err)
					}
				case "null":
					if !strings.Contains(diagnostic, m.want) {
						t.Fatalf("expected null rejection, got %v", err)
					}
				default:
					if !strings.Contains(diagnostic, m.want) && !strings.Contains(diagnostic, "cannot unmarshal") {
						t.Fatalf("expected wrong-type rejection, got %v", err)
					}
				}
			})
		}
	}
}

func TestLockBooleanNullRejectedWithoutHashGuard(t *testing.T) {
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, "profiles/gen-type-go-pr6/target.lock.json"))
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	raw["merge_authorized"] = json.RawMessage("null")
	data, err = json.Marshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	var lock Lock
	if err := decodeLock(data, &lock); err == nil || !strings.Contains(err.Error(), "merge_authorized") {
		t.Fatalf("decodeLock accepted null: %v", err)
	}
}
