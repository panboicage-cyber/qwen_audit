# Qwen audit bootstrap + curated read-only memory

This is a **candidate repo overlay**, to apply to a new feature branch descended from `main`. It includes the pinned AuditGo v0.1.3 Linux bootstrap and a separate static Go RepoMemory reader. It does not overwrite the existing `README.md`, `AUDIT_POLICY.md`, `QWEN_RECOVERY.md`, target scripts, or report templates.

## coder.qwen.ai environment configuration

- Code execution: enabled, only in a disposable no-credentials workspace.
- Environment variables: none required.
- Setup script:

```sh
sh scripts/setup-auditgo.sh && sh scripts/repo-memory-smoke.sh
```

After setup, from the repo root:

```sh
$HOME/.local/bin/auditgo doctor
go run ./tools/repomemory/main.go context
go run ./tools/repomemory/main.go get audit.evidence.classes
```

The reader is intentionally read-only and cannot record Qwen's messages or observations. It uses no database, network, MCP memory-write endpoint, or automatic capture. The `context` command returns only the nine owner-curated facts that are pinned in the repository. Historical facts require live corroboration before current decisions.

## Git history warning

An earlier AuditGo bootstrap branch `feat/auditgo-qwen-bootstrap-v0.1.3` was created as an orphan root, not a child of `main`. **Do not merge that branch as-is** or use its base as the new feature branch. Create a new feature branch from the latest `main`, then copy this overlay into it. The original `main` files must remain unchanged and present.

## Verification and approvals

```sh
sha256sum -c OVERLAY_SHA256SUMS.txt
sh scripts/repo-memory-smoke.sh
go test -count=1 tools/repomemory/main.go tools/repomemory/main_test.go
go test -race -count=1 tools/repomemory/main.go tools/repomemory/main_test.go
```

A passing check validates the packaging and read path. It is not a source-code audit or permission to merge/release. Any edit to `memory/catalog.json` or `memory/SHA256SUMS` requires separate owner review. A model with repository write permission can change both files, so truly preventing unauthorized writes requires GitHub branch protection or a read-only repository mount; the CLI alone cannot enforce that boundary.
