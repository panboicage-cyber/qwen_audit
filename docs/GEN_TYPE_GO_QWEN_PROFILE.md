# Gen Type-Go PR #6 deterministic Qwen audit profile

This overlay is designed for `panboicage-cyber/qwen_audit`. It is a profile of the **existing** AuditGo v0.1.3 bootstrap, not a substitute implementation of AuditGo. It does not mirror target source into the audit-harness repository: target Git is cloned into a separate disposable directory and pinned by SHA.

## Repository topology

- Base audit-harness `main`: `af6609fae097efe80742025b0b5bf27b9acfccea` at initial inspection.
- Existing AuditGo bootstrap: `feat/auditgo-qwen-bootstrap-v0.1.3`, commit `c1ed4143d9ceceb091a97baccb5e3afd12ea74b2` (independent root history; copy its *files*, do not force-merge unrelated history).
- Target under review: `es-3581100/Gen_Type-go`, PR #6, exact HEAD `0af4707d66bb30aef09fe80316c408568e256e53`, base `7b071c7c0ed5037b107f4ae9a263f3c06853c46f`.
- New proposed audit-harness branch: `feat/gen-type-go-pr6-deterministic-audit`. Never merge the target into the audit harness.

## Install this overlay on a new branch

From a working, authenticated developer checkout of `panboicage-cyber/qwen_audit`:

```bash
# Verify fresh working tree and upstream refs first.
git fetch origin main feat/auditgo-qwen-bootstrap-v0.1.3
git switch -c feat/gen-type-go-pr6-deterministic-audit origin/main

# Import only the bootstrap files from the separate-root branch; no merge.
git restore --source=origin/feat/auditgo-qwen-bootstrap-v0.1.3 --staged --worktree -- \
  .github/workflows/auditgo-bootstrap.yml AGENTS.md AUDITGO_SETUP.md \
  OVERLAY_SHA256SUMS.txt docs/AUDITGO-QWEN.md \
  scripts/auditgo-smoke.sh scripts/setup-auditgo.sh \
  vendor/AuditGo_v0.1.3_Crush_Hermes_Linux.tar.gz vendor/README.md

# Copy this overlay's profiles/, prompts/, scripts/, docs/ into the checkout;
# retain the existing kol_merkle example and AUDIT_POLICY.md.
# DO NOT overwrite scripts/setup-auditgo.sh, existing main or target files.

python3 -m json.tool profiles/gen-type-go-pr6/target.lock.json >/dev/null
bash -n scripts/replay-gen-type-go-pr6.sh
sh -n scripts/setup-auditgo.sh
sha256sum -c OVERLAY_SHA256SUMS.txt

git diff --stat --cached
git status --short
# After reviewing the staged changes:
git add profiles/gen-type-go-pr6 prompts/gen-type-go-pr6.md \
  scripts/replay-gen-type-go-pr6.sh docs/GEN_TYPE_GO_QWEN_PROFILE.md \
  docs/QWEN_PR6_AUDIT_HANDOFF.md
# git commit ... and git push -u origin ... (requires explicit authorized write access)
```

`OVERLAY_SHA256SUMS.txt` indexes only the imported AuditGo bootstrap files. The new profile is separate; do not rewrite the bootstrap hash index without changing its stated scope.

## In the Qwen coding session

The session must clone the *new audit-harness branch* (not the original bootstrap branch, and not Gen Type-Go as its workspace root). Then run:

```bash
echo QWEN_SHELL_OK
sh scripts/setup-auditgo.sh
sh scripts/auditgo-smoke.sh
bash scripts/replay-gen-type-go-pr6.sh
```

Read `prompts/gen-type-go-pr6.md` after bootstrap and complete the selected AuditGo ledger. The replay script runs source, digest, and Go checks and generates an unassessed seed; it **does not** independently evaluate the whole 383-item canonical inventory or grant a final passing disposition.

## Evidence/limits

- Output is ignored via the existing `audit-artifacts/` gitignore pattern.
- Clone, test, fuzz, and findings are confined to an isolated target copy. Do not execute unknown target code where private credentials or elevated host access are present.
- Missing shell, Go, Python 3, AuditGo, target SHA, or check execution is an explicit blocker/incomplete result.
- Do not post to PR #6 without independent publication authority; no approvals, merges, releases, deployment, or protected-authority promotion are part of this profile.
- If PR #6 head changes, **do not auto-retarget**. A human should update the lock and record a new checkpoint with fresh validation.
