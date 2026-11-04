# Active staged audit — Gen Type-Go PR #6

This is the **master workspace wrapper** around the frozen PR #6 task prompt, not a rewrite of its historical audit claims.

1. Read `README.md`, `AGENTS.md`, `AUDIT_POLICY.md`, `staging_audit/workspace.json`, `staging_audit/branches/gen-type-go-pr6.json`, and `profiles/gen-type-go-pr6/target.lock.json`.
2. Run shell preflight and `sh scripts/setup-master-audit-workspace.sh`. Capture actual stdout/exit code. Stop on missing files or hash mismatch; do not guess another revision.
3. Prepare the immutable handoff with `sh scripts/stage-audit.sh gen-type-go-pr6`; read `audit-artifacts/staging/gen-type-go-pr6/context.json` and its checksum.
4. Read `prompts/gen-type-go-pr6.md` and `docs/PR6_STAGED_HANDOFF.md` for *scope*. Treat them as historical instructions, not permission to bypass the master setup or to perform network actions.
5. **Do not execute the target replay without separate explicit owner authorization**. When authorized, use a disposable, credentials-free sandbox and `bash scripts/replay-gen-type-go-pr6.sh` with no upstream mutation. Record real output, tests, findings and provenance.
6. Generate an honest audit receipt with source identity, technical findings, coverage gaps, actual evidence artifact paths and SHA-256, plus a concise summary. An AuditGo checklist, staged context or passing CI is not approval.
7. No automatic PR comment, GitHub `APPROVE`, merge, release, or authority promotion. Default disposition remains `AUDIT_INCOMPLETE` until enough evidence is collected.

Frozen target SHA: `0af4707d66bb30aef09fe80316c408568e256e53`. No branch-tip fallback. This wrapper does not modify the underlying lock or replay script.
