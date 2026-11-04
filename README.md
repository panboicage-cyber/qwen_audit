# Qwen Audit — Master Audit Workspace

**One Qwen repository. One setup script. One deterministic staging queue.**

This is a review harness for [coder.qwen.ai](https://coder.qwen.ai), not a mirror of the repositories being audited. It integrates the pinned **AuditGo v0.1.3** CLI, Engram-inspired **read-only RepoMemory**, a Go-validated **audit stage queue**, and an optional TypeScript presentation interface. Original historical reviews remain in `prompts/`, `scripts/` and `reports/`.

## Qwen cloud environment settings

| Field | Value |
| --- | --- |
| Repository | `panboicage-cyber/qwen_audit` |
| Branch | `master_audit_workspace` **after it has been published** |
| Code execution | Enabled |
| Environment variables | None required |
| Setup script | `sh scripts/setup-master-audit-workspace.sh` |

The setup script **does not clone, inspect, test, mutate, approve, publish, or merge any audit target**. It installs the pinned AuditGo archive, validates RepoMemory, verifies the frozen queue, and performs a planning-only CLI smoke test. A working shell, Go, and Linux AMD64 are required. TypeScript/Node is optional and **never part of the mandatory Qwen setup**.

## Standard session sequence

```sh
echo QWEN_SHELL_OK
sh scripts/setup-master-audit-workspace.sh
go run ./tools/masteraudit/main.go list
sh scripts/stage-audit.sh gen-type-go-pr6
cat audit-artifacts/staging/gen-type-go-pr6/context.json
```

Stage preparation is **not** audit execution. It writes a derived handoff artifact under the ignored `audit-artifacts/` directory; it does not record new memories or modify the staged lock. Only after separate explicit authorization, and only in an isolated environment without privileged credentials, the pinned target replay may be run:

```sh
# Optional, separately authorized target replay; not a startup command.
bash scripts/replay-gen-type-go-pr6.sh
```

Read `prompts/MASTER_AUDIT_INSTRUCTIONS.md`, the selected profile and target lock before reviewing source. Audit outputs must distinguish `EXECUTED`, `OBSERVED`, `REPORTED`, and `UNVERIFIED`; the final result must use `AUDIT_PASS`, `AUDIT_FAIL`, or `AUDIT_INCOMPLETE`, with separate owner approval and publication gates.

## Stage queue

| ID | State | Frozen target | Auto-run |
| --- | --- | --- | --- |
| `gen-type-go-pr6` | `STAGED` | `es-3581100/Gen_Type-go` PR #6 @ `0af4707d66bb30aef09fe80316c408568e256e53` | **No** |
| `kol-merkle-pr4` | `REFERENCE_ONLY` | `donCannoli-burns/kol_merkle` PR #4 @ `c64502cb5db3efd9316991ce1c0756d97729a6c3` | **No** |

- `staging_audit/workspace.json` is the closed authority boundary for the environment.
- `staging_audit/queue.json` is the deliberately ordered list of profiles.
- `staging_audit/branches/*.json` contain pinned provenance, content hashes, and execution denial flags.
- `profiles/*/target.lock.json` contain frozen target revisions, not "latest" branch aliases.
- `tools/masteraudit/main.go` validates JSON key uniqueness, known fields, SHA-256 identities, sources, and locked authority flags. It only supports `verify`, `list`, `context`. No `save`, `learn`, or `execute` command exists.
- `tools/masteraudit/queue.ts` consumes an already validated Go context for a typed UI/agent handoff; it is not a second authorization implementation.

## Baseline and provenance

This unified workspace is assembled on **`main` commit `af6609fae097efe80742025b0b5bf27b9acfccea`**, importing:

- RepoMemory source branch `feat/auditgo-curated-memory-v0.1.0` @ `57a08f95db8f297136db6e2a5eed8730d0df6fc1`.
- Gen Type-Go staging source branch `feat/gen-type-go-pr6-deterministic-audit` @ `f046b362a482d1fa4d557df284db3cf39bf0a55d`.

Both are ordinary descendants of `main`; the older orphan bootstrap branch is **not merged**. `OVERLAY_SHA256SUMS.txt` pins the overlay files. `memory/SHA256SUMS` independently pins the curated catalog. Run `sha256sum -c OVERLAY_SHA256SUMS.txt` after applying the overlay.

For architecture, maintenance, and test gates see [docs/MASTER_AUDIT_WORKSPACE.md](docs/MASTER_AUDIT_WORKSPACE.md). For historical Qwen shell recovery see [QWEN_RECOVERY.md](QWEN_RECOVERY.md).

## Boundaries

The bundled memory interfaces cannot save or capture observations, but an agent with general Git/file-write access can still edit project files. Branch protection and CODEOWNERS reviews are required for a hard no-write memory governance boundary. A successful bootstrap or green AuditGo seed is **not** completed source verification or approval. No automatic PR comments, `APPROVE`, merges, releases or authority promotion are performed.

## Historical first assignment

The original seven-file `kol_merkle` PR #4 assignment is preserved in `prompts/kol-merkle-pr4.md`, `scripts/setup-qwen-target.sh`, `scripts/replay-kol-merkle-pr4.sh`, and `reports/REVIEW_TEMPLATE.md`. It is currently disabled in the unified queue. Run it only through an explicitly reviewed staging change or directly under separate authorization; never substitute its SHA for the active Gen Type-Go target.
