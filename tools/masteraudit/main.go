// MasterAudit is a read-only, deterministic selector for a curated Qwen audit workspace.
// It cannot write memories, clone repositories, execute profiles or promote authority.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	workspaceSchema = "qwen-master-audit.workspace/v1"
	queueSchema     = "qwen-master-audit.queue/v1"
	stageSchema     = "qwen-master-audit.stage/v1"
	lockSchema      = "qwen-audit-target-lock-v1"
	maxJSON         = 64 * 1024
)

var (
	idPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{2,63}$`)
	hex40     = regexp.MustCompile(`^[0-9a-f]{40}$`)
	hex64     = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

type Workspace struct {
	Schema                   string `json:"schema"`
	Name                     string `json:"name"`
	Setup                    string `json:"setup"`
	AuditGo                  string `json:"auditgo"`
	Memory                   string `json:"memory"`
	Queue                    string `json:"queue"`
	AllowAutoExecute         bool   `json:"allow_auto_execute"`
	AllowRuntimeMemoryWrites bool   `json:"allow_runtime_memory_writes"`
	AllowMerge               bool   `json:"allow_merge"`
	AllowPublish             bool   `json:"allow_publish"`
}
type Queue struct {
	Schema   string      `json:"schema"`
	Strategy string      `json:"strategy"`
	Items    []QueueItem `json:"items"`
}
type QueueItem struct {
	ID      string `json:"id"`
	Stage   string `json:"stage"`
	Enabled bool   `json:"enabled"`
}
type Stage struct {
	Schema                    string `json:"schema"`
	ID                        string `json:"id"`
	Status                    string `json:"status"`
	SourceBranch              string `json:"source_branch"`
	SourceBranchCommit        string `json:"source_branch_commit"`
	TargetRepository          string `json:"target_repository"`
	TargetHead                string `json:"target_head"`
	LockPath                  string `json:"lock_path"`
	LockSHA256                string `json:"lock_sha256"`
	RunnerPath                string `json:"runner_path"`
	RunnerSHA256              string `json:"runner_sha256"`
	PromptPath                string `json:"prompt_path"`
	PromptSHA256              string `json:"prompt_sha256"`
	RequiresExplicitExecution bool   `json:"requires_explicit_execution"`
	AllowUpstreamMutation     bool   `json:"allow_upstream_mutation"`
	AllowMerge                bool   `json:"allow_merge"`
	AllowRelease              bool   `json:"allow_release"`
}
type Lock struct {
	Schema                      string `json:"schema"`
	TargetRepo                  string `json:"target_repo"`
	SourceHead                  string `json:"source_head"`
	UpstreamMutationsAuthorized bool   `json:"upstream_mutations_authorized"`
	MergeAuthorized             bool   `json:"merge_authorized"`
	ReleaseAuthorized           bool   `json:"release_authorized"`
	NewAuthorityAuthorized      bool   `json:"new_authority_authorized"`
}
type CheckedStage struct {
	Stage       Stage  `json:"stage"`
	StageSHA256 string `json:"stage_sha256"`
}
type Context struct {
	Schema          string       `json:"schema"`
	WorkspaceName   string       `json:"workspace_name"`
	WorkspaceSHA256 string       `json:"workspace_sha256"`
	QueueSHA256     string       `json:"queue_sha256"`
	Profile         CheckedStage `json:"profile"`
	MemoryAccess    string       `json:"memory_access"`
	ExecutionState  string       `json:"execution_state"`
	SourceVerified  bool         `json:"source_verified"`
	OwnerApproval   bool         `json:"owner_approval"`
	NextAction      string       `json:"next_action"`
}

// Decoder.Token() unescapes JSON object keys. Reject duplicate decoded keys before
// decoding into structs, including equivalent spellings such as "a" and "\\u0061".
func rejectDuplicateKeys(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	var walk func() error
	walk = func() error {
		token, err := d.Token()
		if err != nil {
			return err
		}
		delim, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				token, err := d.Token()
				if err != nil {
					return err
				}
				key, ok := token.(string)
				if !ok {
					return errors.New("object key is not a string")
				}
				if seen[key] {
					return fmt.Errorf("duplicate decoded key %q", key)
				}
				seen[key] = true
				if err := walk(); err != nil {
					return err
				}
			}
			_, err := d.Token()
			return err
		case '[':
			for d.More() {
				if err := walk(); err != nil {
					return err
				}
			}
			_, err := d.Token()
			return err
		default:
			return fmt.Errorf("unexpected delimiter %q", delim)
		}
	}
	if err := walk(); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		if err == nil {
			return errors.New("trailing JSON values")
		}
		return err
	}
	return nil
}
func sha(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }
func safePath(p string) bool {
	if p == "" || filepath.IsAbs(p) || filepath.Clean(p) != p {
		return false
	}
	for _, part := range strings.Split(filepath.ToSlash(p), "/") {
		if part == ".." || part == "." || part == "" {
			return false
		}
	}
	return true
}
func readRegular(root, p string) ([]byte, error) {
	if !safePath(p) {
		return nil, fmt.Errorf("unsafe path %q", p)
	}
	// Symlink paths are rejected at every component, not just the final file.
	cur := root
	for _, part := range strings.Split(filepath.ToSlash(p), "/") {
		cur = filepath.Join(cur, part)
		st, err := os.Lstat(cur)
		if err != nil {
			return nil, err
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("symlink component blocked: %s", p)
		}
	}
	st, err := os.Lstat(cur)
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsRegular() || st.Size() == 0 || st.Size() > maxJSON {
		return nil, fmt.Errorf("invalid regular file size/type: %s", p)
	}
	return os.ReadFile(cur)
}
func decodeStrict(data []byte, out any) error {
	if err := rejectDuplicateKeys(data); err != nil {
		return err
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(out); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		if err == nil {
			return errors.New("trailing value")
		}
		return err
	}
	return nil
}

// requireExplicitBooleans is called after duplicate-key and unknown-field
// validation. Go's encoding/json otherwise treats both omitted fields and JSON
// null as the zero value of bool, blurring a deliberate false with no decision.
func requireExplicitBooleans(data []byte, keys ...string) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if raw == nil {
		return errors.New("expected JSON object for boolean contract")
	}
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			return fmt.Errorf("missing required boolean field %q", key)
		}
		value = bytes.TrimSpace(value)
		if !bytes.Equal(value, []byte("true")) && !bytes.Equal(value, []byte("false")) {
			return fmt.Errorf("field %q must be explicit JSON boolean (true or false)", key)
		}
	}
	return nil
}

// The queue has nested booleans. Require an explicit enabled decision for
// every entry, including disabled reference-only entries.
func requireQueueBooleans(data []byte) error {
	var envelope struct {
		Items []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return err
	}
	for i, item := range envelope.Items {
		if err := requireExplicitBooleans(item, "enabled"); err != nil {
			return fmt.Errorf("queue items[%d]: %w", i, err)
		}
	}
	return nil
}

func fileWithHash(root, p, expected string) (string, error) {
	if !hex64.MatchString(expected) {
		return "", fmt.Errorf("invalid expected SHA256 for %s", p)
	}
	data, err := readRegular(root, p)
	if err != nil {
		return "", err
	}
	actual := sha(data)
	if actual != expected {
		return actual, fmt.Errorf("SHA256 mismatch %s", p)
	}
	return actual, nil
}
func check(root string) (Workspace, string, Queue, string, []CheckedStage, error) {
	var w Workspace
	var q Queue
	data, err := readRegular(root, "staging_audit/workspace.json")
	if err != nil {
		return w, "", q, "", nil, err
	}
	if err := decodeStrict(data, &w); err != nil {
		return w, "", q, "", nil, err
	}
	if err := requireExplicitBooleans(data, "allow_auto_execute", "allow_runtime_memory_writes", "allow_merge", "allow_publish"); err != nil {
		return w, "", q, "", nil, fmt.Errorf("workspace: %w", err)
	}
	wsha := sha(data)
	if w.Schema != workspaceSchema || w.Name != "master_audit_workspace" || w.Setup != "scripts/setup-master-audit-workspace.sh" || w.AuditGo != "scripts/setup-auditgo.sh" || w.Memory != "memory/catalog.json" || w.Queue != "staging_audit/queue.json" || w.AllowAutoExecute || w.AllowRuntimeMemoryWrites || w.AllowMerge || w.AllowPublish {
		return w, "", q, "", nil, errors.New("workspace contract mismatch")
	}
	data, err = readRegular(root, w.Queue)
	if err != nil {
		return w, "", q, "", nil, err
	}
	if err := decodeStrict(data, &q); err != nil {
		return w, "", q, "", nil, err
	}
	if err := requireQueueBooleans(data); err != nil {
		return w, "", q, "", nil, err
	}
	qsha := sha(data)
	if q.Schema != queueSchema || q.Strategy != "ORDERED_EXPLICIT_SELECTION" || len(q.Items) == 0 || len(q.Items) > 16 {
		return w, "", q, "", nil, errors.New("invalid queue contract")
	}
	seen := map[string]bool{}
	active := 0
	checked := make([]CheckedStage, 0, len(q.Items))
	for _, entry := range q.Items {
		if !idPattern.MatchString(entry.ID) || seen[entry.ID] {
			return w, "", q, "", nil, fmt.Errorf("invalid/duplicate queue ID %q", entry.ID)
		}
		seen[entry.ID] = true
		if entry.Stage != "staging_audit/branches/"+entry.ID+".json" {
			return w, "", q, "", nil, errors.New("noncanonical stage path")
		}
		b, err := readRegular(root, entry.Stage)
		if err != nil {
			return w, "", q, "", nil, err
		}
		var s Stage
		if err := decodeStrict(b, &s); err != nil {
			return w, "", q, "", nil, err
		}
		if err := requireExplicitBooleans(b, "requires_explicit_execution", "allow_upstream_mutation", "allow_merge", "allow_release"); err != nil {
			return w, "", q, "", nil, fmt.Errorf("stage %s: %w", entry.ID, err)
		}
		if s.Schema != stageSchema || s.ID != entry.ID || !hex40.MatchString(s.SourceBranchCommit) || !hex40.MatchString(s.TargetHead) || s.TargetRepository == "" || s.SourceBranch == "" || !s.RequiresExplicitExecution || s.AllowUpstreamMutation || s.AllowMerge || s.AllowRelease {
			return w, "", q, "", nil, fmt.Errorf("invalid stage contract %s", entry.ID)
		}
		if entry.Enabled {
			active++
			if s.Status != "STAGED" {
				return w, "", q, "", nil, fmt.Errorf("active stage not STAGED: %s", entry.ID)
			}
		} else if s.Status != "REFERENCE_ONLY" {
			return w, "", q, "", nil, fmt.Errorf("inactive stage must be REFERENCE_ONLY: %s", entry.ID)
		}
		for _, p := range []struct{ path, hash string }{{s.RunnerPath, s.RunnerSHA256}, {s.PromptPath, s.PromptSHA256}} {
			if _, err := fileWithHash(root, p.path, p.hash); err != nil {
				return w, "", q, "", nil, err
			}
		}
		if entry.Enabled {
			if _, err := fileWithHash(root, s.LockPath, s.LockSHA256); err != nil {
				return w, "", q, "", nil, err
			}
			lockbytes, _ := readRegular(root, s.LockPath)
			var lock Lock
			if err := decodeLock(lockbytes, &lock); err != nil {
				return w, "", q, "", nil, err
			}
			if lock.Schema != lockSchema || lock.TargetRepo != s.TargetRepository || lock.SourceHead != s.TargetHead || lock.UpstreamMutationsAuthorized || lock.MergeAuthorized || lock.ReleaseAuthorized || lock.NewAuthorityAuthorized {
				return w, "", q, "", nil, fmt.Errorf("target lock does not match stage %s", entry.ID)
			}
		} else if s.LockPath != "" || s.LockSHA256 != "" {
			return w, "", q, "", nil, fmt.Errorf("reference stage must not claim a lock: %s", entry.ID)
		}
		checked = append(checked, CheckedStage{Stage: s, StageSHA256: sha(b)})
	}
	if active != 1 {
		return w, "", q, "", nil, fmt.Errorf("expected exactly one active staged entry, got %d", active)
	}
	// Verify memory bytes and hash. The dedicated RepoMemory utility verifies semantics.
	manifest, err := readRegular(root, "memory/SHA256SUMS")
	if err != nil {
		return w, "", q, "", nil, err
	}
	parts := strings.Split(strings.TrimSuffix(string(manifest), "\n"), "  ")
	if len(parts) != 2 || parts[1] != "catalog.json" || !strings.HasSuffix(string(manifest), "\n") {
		return w, "", q, "", nil, errors.New("invalid memory manifest")
	}
	if _, err := fileWithHash(root, w.Memory, parts[0]); err != nil {
		return w, "", q, "", nil, err
	}
	// Static workspace sources: protect known scripts from missing/symlink states.
	for _, p := range []string{w.Setup, w.AuditGo, "scripts/auditgo-smoke.sh", "scripts/repo-memory-smoke.sh", "AUDIT_POLICY.md", "AGENTS.md", "README.md"} {
		if _, err := readRegular(root, p); err != nil {
			return w, "", q, "", nil, err
		}
	}
	return w, wsha, q, qsha, checked, nil
}

// The preexisting target lock has additional fields that are deliberately not
// interpreted as mutable runtime state; all recognized authorization fields are checked.
func decodeLock(data []byte, l *Lock) error {
	if err := rejectDuplicateKeys(data); err != nil {
		return err
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	allowed := map[string]bool{"schema": true, "profile_id": true, "target_repo": true, "pull_request": true, "source_head": true, "source_tree": true, "base_main": true, "original_v04_head": true, "parents": true, "reconciliation_paths": true, "tracked_files_expected": true, "manifest_entries_expected": true, "auditgo": true, "audit_types": true, "modules": true, "claims": true, "authority_state": true, "upstream_mutations_authorized": true, "merge_authorized": true, "release_authorized": true, "new_authority_authorized": true}
	for k := range raw {
		if !allowed[k] {
			return fmt.Errorf("unknown lock field %q", k)
		}
	}
	// Required fields must exist to avoid omitting a flag and interpreting false as approval.
	for _, k := range []string{"schema", "target_repo", "source_head", "upstream_mutations_authorized", "merge_authorized", "release_authorized", "new_authority_authorized"} {
		if _, ok := raw[k]; !ok {
			return fmt.Errorf("missing lock field %q", k)
		}
	}
	if err := requireExplicitBooleans(data, "upstream_mutations_authorized", "merge_authorized", "release_authorized", "new_authority_authorized"); err != nil {
		return fmt.Errorf("target lock: %w", err)
	}
	return json.Unmarshal(data, l)
}
func dispatch(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("masteraudit", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("root", ".", "workspace root")
	profile := fs.String("profile", "", "exact staged ID (only for context)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	rest := fs.Args()
	if len(rest) != 1 {
		fmt.Fprintln(stderr, "usage: masteraudit [-root ROOT] [-profile ID] verify|list|context")
		return 2
	}
	cmd := rest[0]
	if cmd != "verify" && cmd != "list" && cmd != "context" {
		fmt.Fprintln(stderr, "unsupported command: read-only interface")
		return 2
	}
	rootabs, err := filepath.Abs(*root)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	w, wsha, q, qsha, checked, err := check(rootabs)
	if err != nil {
		fmt.Fprintln(stderr, "MASTER_AUDIT_BLOCKED:", err)
		return 1
	}
	switch cmd {
	case "verify":
		fmt.Fprintf(stdout, "MASTER_AUDIT_VERIFY_PASS workspace_sha256=%s queue_sha256=%s staged=%d total=%d\n", wsha, qsha, 1, len(checked))
	case "list":
		for i, s := range checked {
			fmt.Fprintf(stdout, "%s\t%s\tenabled=%t\ttarget=%s\n", s.Stage.ID, s.Stage.Status, q.Items[i].Enabled, s.Stage.TargetHead)
		}
	case "context":
		if *profile == "" {
			fmt.Fprintln(stderr, "context requires -profile ID")
			return 2
		}
		for i, s := range checked {
			if s.Stage.ID == *profile && q.Items[i].Enabled {
				ctx := Context{Schema: "qwen-master-audit.context/v1", WorkspaceName: w.Name, WorkspaceSHA256: wsha, QueueSHA256: qsha, Profile: s, MemoryAccess: "READ_ONLY", ExecutionState: "NOT_EXECUTED", SourceVerified: false, OwnerApproval: false, NextAction: "Inspect the frozen target lock and mission; obtain explicit approval before executing the separate sandbox replay."}
				enc := json.NewEncoder(stdout)
				enc.SetIndent("", "  ")
				if err := enc.Encode(ctx); err != nil {
					fmt.Fprintln(stderr, err)
					return 1
				}
				return 0
			}
		}
		fmt.Fprintln(stderr, "profile not staged/enabled; fail closed")
		return 2
	}
	return 0
}
func main() { os.Exit(dispatch(os.Args[1:], os.Stdout, os.Stderr)) }
