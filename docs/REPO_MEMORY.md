> **Master Audit Workspace integration:** For coder.qwen.ai use the single setup command `sh scripts/setup-master-audit-workspace.sh`, not the older component-specific commands below. Those remain documented as standalone diagnostics. See [`README.md`](../README.md) and `MASTER_AUDIT_WORKSPACE.md` for the authoritative staging sequence.

# RepoMemory v0.1.0 — deterministic context without agent capture

**Purpose:** Give Qwen Coder brief predefined repository facts at startup while preventing the bundled memory interface from accepting Qwen-generated observations. This is a reader, not Engram's persistent memory service.

**Location:** `memory/catalog.json`, pinned by `memory/SHA256SUMS`. `tools/repomemory/main.go` is a pure-Go source-only command available with Go 1.23.2.

## Interaction

Run `sh scripts/repo-memory-smoke.sh` after the Qwen shell preflight. To load the safe short-form context, run `go run ./tools/repomemory/main.go context`. For a precise detail, run `go run ./tools/repomemory/main.go get audit.authorization.boundary`.

The returned source paths and commits are historical evidence locators; they do not prove current GitHub head state or justify actions. Memory cannot approve, merge, release, or edit audited projects. No read command automatically updates the catalog.

## Governance and attack model

A compromised model that can write arbitrary files can rewrite both `catalog.json` and `SHA256SUMS`. Hash checks detect accidental/unpaired modification, not adversarial edits by a repository writer. For a hard no-write guarantee, run Qwen with read-only repository permissions and review all changes through a separately controlled, protected branch. Merely saying "do not write memory" in AGENTS.md is not an access-control boundary.

The code rejects unknown commands (including `save`, `append`, `capture`, `learn`, `update`, and `delete`), encoded duplicate JSON object keys, unknown schema properties, malformed provenance paths, and tampered catalog bytes. Missing or stale catalog data blocks context rather than inventing replacement facts.

No network listener, database, external storage, agent identity tracking, session memory, model-generated summaries, or self-learning behavior is present.

## Proof and deployment

- `go test -count=1 tools/repomemory/main.go tools/repomemory/main_test.go`
- `go test -race -count=1 tools/repomemory/main.go tools/repomemory/main_test.go`
- `go vet tools/repomemory/main.go tools/repomemory/main_test.go`
- `sh scripts/repo-memory-smoke.sh`

This is a candidate overlay for `panboicage-cyber/qwen_audit`; it should be added to a **normal branch descended from the latest `main`**, with the AuditGo integration copied from its orphan bootstrap branch, before PR review. Do not merge or replace repository history to deploy it.
