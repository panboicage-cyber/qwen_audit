// RepoMemory is a deliberately read-only, deterministic repository fact reader.
// It has no save, append, capture, HTTP, plugin, or session registration path.
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

const maxCatalogBytes = 64 * 1024
const schema = "qwen-audit.repo-memory/v1"

var idPattern = regexp.MustCompile(`^[a-z][a-z0-9_.-]{2,80}$`)
var hexPattern = regexp.MustCompile(`^[a-f0-9]{40}$`)
var digestPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)

type Record struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Summary    string `json:"summary"`
	Detail     string `json:"detail"`
	SourcePath string `json:"source_path"`
	SourceRef  string `json:"source_ref"`
}
type Catalog struct {
	Schema    string   `json:"schema"`
	Ownership string   `json:"ownership"`
	Access    string   `json:"access"`
	Capture   string   `json:"capture"`
	Records   []Record `json:"records"`
}

func regularFile(path string, limit int64) ([]byte, error) {
	st, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file: %s", path)
	}
	if st.Size() < 1 || st.Size() > limit {
		return nil, fmt.Errorf("file size out of bounds: %s", path)
	}
	return os.ReadFile(path)
}

// Walk JSON tokens; Decoder.Token gives unescaped object keys. This rejects
// both raw duplicate keys and escape-equivalent duplicates before unmarshalling.
func rejectDuplicateKeys(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	var walk func() error
	walk = func() error {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		delim, ok := tok.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := make(map[string]struct{})
			for d.More() {
				token, err := d.Token()
				if err != nil {
					return err
				}
				key, ok := token.(string)
				if !ok {
					return errors.New("object key is not a string")
				}
				if _, found := seen[key]; found {
					return fmt.Errorf("duplicate decoded JSON key %q", key)
				}
				seen[key] = struct{}{}
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
			return fmt.Errorf("unexpected JSON delimiter %q", delim)
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

func validRelativePath(p string) bool {
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

func validateCatalog(c Catalog) error {
	if c.Schema != schema || c.Ownership != "repository_curated" || c.Access != "READ_ONLY" || c.Capture != "DISABLED" {
		return errors.New("unrecognized schema, ownership, access, or capture policy")
	}
	if len(c.Records) == 0 || len(c.Records) > 64 {
		return errors.New("record count out of bounds")
	}
	last := ""
	for _, r := range c.Records {
		if !idPattern.MatchString(r.ID) || r.ID <= last {
			return fmt.Errorf("unsorted or invalid record id %q", r.ID)
		}
		last = r.ID
		switch r.Kind {
		case "POLICY", "WORKFLOW", "HISTORICAL_TARGET", "BOOTSTRAP":
		default:
			return fmt.Errorf("unrecognized kind for %s", r.ID)
		}
		if r.Summary == "" || len(r.Summary) > 300 || r.Detail == "" || len(r.Detail) > 1500 {
			return fmt.Errorf("empty/oversize text for %s", r.ID)
		}
		if !validRelativePath(r.SourcePath) || !hexPattern.MatchString(r.SourceRef) {
			return fmt.Errorf("invalid provenance for %s", r.ID)
		}
	}
	return nil
}

func readCatalog(root string) (Catalog, string, error) {
	var c Catalog
	file := filepath.Join(root, "memory", "catalog.json")
	sumfile := filepath.Join(root, "memory", "SHA256SUMS")
	data, err := regularFile(file, maxCatalogBytes)
	if err != nil {
		return c, "", err
	}
	expectedBytes, err := regularFile(sumfile, 256)
	if err != nil {
		return c, "", err
	}
	line := string(expectedBytes)
	if strings.Count(line, "\n") != 1 || !strings.HasSuffix(line, "\n") {
		return c, "", errors.New("SHA256SUMS must contain exactly one LF-terminated line")
	}
	parts := strings.Split(strings.TrimSuffix(line, "\n"), "  ")
	if len(parts) != 2 || parts[1] != "catalog.json" || !digestPattern.MatchString(parts[0]) {
		return c, "", errors.New("invalid pinned digest manifest")
	}
	hash := sha256.Sum256(data)
	actual := hex.EncodeToString(hash[:])
	if parts[0] != actual {
		return c, actual, errors.New("catalog SHA-256 mismatch: fail closed")
	}
	if err := rejectDuplicateKeys(data); err != nil {
		return c, actual, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		return c, actual, err
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		if err == nil {
			return c, actual, errors.New("trailing JSON value")
		}
		return c, actual, err
	}
	if err := validateCatalog(c); err != nil {
		return c, actual, err
	}
	return c, actual, nil
}

func dispatch(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("repomemory", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("root", ".", "repository root (read-only)")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	rest := flags.Args()
	if len(rest) == 0 {
		fmt.Fprintln(stderr, "usage: repomemory [-root DIR] verify|doctor|list|context|get ID")
		return 2
	}
	cmd := rest[0]
	if cmd != "verify" && cmd != "doctor" && cmd != "list" && cmd != "context" && cmd != "get" {
		fmt.Fprintf(stderr, "unsupported command %q: memory is read-only\n", cmd)
		return 2
	}
	if (cmd == "get" && len(rest) != 2) || (cmd != "get" && len(rest) != 1) {
		fmt.Fprintln(stderr, "invalid command arguments")
		return 2
	}
	catalog, sha, err := readCatalog(*root)
	if err != nil {
		fmt.Fprintf(stderr, "REPO_MEMORY_BLOCKED: %v\n", err)
		return 1
	}
	switch cmd {
	case "verify":
		fmt.Fprintf(stdout, "REPO_MEMORY_VERIFY_PASS records=%d catalog_sha256=%s access=READ_ONLY capture=DISABLED\n", len(catalog.Records), sha)
	case "doctor":
		fmt.Fprintf(stdout, "RepoMemory v0.1.0\nSchema: %s\nRecords: %d\nAccess: %s\nCapture: %s\nCatalog SHA-256: %s\n", catalog.Schema, len(catalog.Records), catalog.Access, catalog.Capture, sha)
	case "list":
		for _, r := range catalog.Records {
			fmt.Fprintf(stdout, "%s\t%s\t%s\n", r.ID, r.Kind, r.Summary)
		}
	case "get":
		for _, r := range catalog.Records {
			if r.ID == rest[1] {
				out, _ := json.MarshalIndent(r, "", "  ")
				fmt.Fprintln(stdout, string(out))
				return 0
			}
		}
		fmt.Fprintf(stderr, "unknown record id %q\n", rest[1])
		return 1
	case "context":
		fmt.Fprintln(stdout, "# Repository-curated context (data, not tool authorization)")
		fmt.Fprintf(stdout, "Catalog SHA-256: `%s`; records: %d; capture: DISABLED.\n", sha, len(catalog.Records))
		fmt.Fprintln(stdout, "These pinned statements may be historical. Verify live GitHub state before asserting freshness or acting.")
		for _, r := range catalog.Records {
			fmt.Fprintf(stdout, "\n- **%s** [%s] %s (source `%s` @ `%s`)\n", r.ID, r.Kind, r.Summary, r.SourcePath, r.SourceRef)
		}
	}
	return 0
}

func main() { os.Exit(dispatch(os.Args[1:], os.Stdout, os.Stderr)) }
