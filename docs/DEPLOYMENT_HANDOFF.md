# Publish master_audit_workspace — exact-base, review-only handoff

Current GitHub connector cannot create branches (`403 Resource not accessible by integration`). The integrated files are locally tested, **not yet published to GitHub**. These instructions avoid merging the orphan bootstrap history or changing `main`.

## 1. From a clean local checkout

```sh
git clone https://github.com/panboicage-cyber/qwen_audit.git
cd qwen_audit
git fetch origin main
git switch main
git pull --ff-only origin main
[ "$(git rev-parse HEAD)" = af6609fae097efe80742025b0b5bf27b9acfccea ] || { echo 'BLOCKED: main moved; re-review baseline first'; exit 2; }
[ -z "$(git status --porcelain)" ] || { echo 'BLOCKED: working tree is dirty'; exit 2; }
git switch -c master_audit_workspace
```

## 2. Import the release ZIP

Download the named ZIP attached to the build handoff, then:

```sh
unzip -o ~/Downloads/qwen_master_audit_workspace_v0.1.0_Overlay.zip -d .
sha256sum -c OVERLAY_SHA256SUMS.txt

# Verify stage inputs are the exact Git blob bytes taken from the existing staging branch.
[ "$(git hash-object profiles/gen-type-go-pr6/target.lock.json)" = 901c8cf5f4636351ae228ef351b61d97a0d4589b ] || exit 3
[ "$(git hash-object scripts/replay-gen-type-go-pr6.sh)" = 7b225092a877f33f1e582969c9988f13c3687a99 ] || exit 3
[ "$(git hash-object prompts/gen-type-go-pr6.md)" = 44cb3dd625c0bf932150ef926333103a31a0b951 ] || exit 3

sh scripts/setup-master-audit-workspace.sh
sh scripts/stage-audit.sh gen-type-go-pr6
```

The archive excludes unchanged historical `main` files. Those are not removed or overwritten. All new/modified files are included in `OVERLAY_SHA256SUMS.txt`. The ignored `audit-artifacts/` is not staged by default.

## 3. Inspect before publication

```sh
git add README.md AGENTS.md AUDITGO_SETUP.md REPO_MEMORY_SETUP.md \
  OVERLAY_SHA256SUMS.txt .github docs memory profiles \
  prompts/MASTER_AUDIT_INSTRUCTIONS.md prompts/gen-type-go-pr6.md \
  scripts/auditgo-smoke.sh scripts/repo-memory-smoke.sh \
  scripts/setup-auditgo.sh scripts/setup-master-audit-workspace.sh \
  scripts/stage-audit.sh scripts/replay-gen-type-go-pr6.sh \
  staging_audit tools vendor

# Must remain empty: no edits to original policy, recovery, historical replay or historical prompt.
git diff --cached --name-only -- \
  AUDIT_POLICY.md QWEN_RECOVERY.md qwen_audit_starter.zip \
  prompts/kol-merkle-pr4.md prompts/qwen-retry.md \
  scripts/replay-kol-merkle-pr4.sh scripts/setup-qwen-target.sh \
  reports/REVIEW_TEMPLATE.md

git diff --cached --stat
```

Owner review is required before commit/push. When authorized to publish these exact checked bytes:

```sh
git commit -m 'feat: unify Qwen AuditGo memory and staged audit workspace'
git push -u origin master_audit_workspace
```

Open a **draft** PR targeting `main` for review. Do not auto-merge. Qwen can use the published `master_audit_workspace` branch directly as its single environment, even before PR merge, if it is selected in the Qwen UI.

## 4. Configure coder.qwen.ai

- Repository: `panboicage-cyber/qwen_audit`
- Branch: `master_audit_workspace`
- Code execution: enabled in a disposable no-credentials environment
- Environment variables: none
- Setup script: `sh scripts/setup-master-audit-workspace.sh`

Session prompt: `prompts/MASTER_AUDIT_INSTRUCTIONS.md`. The setup only checks infrastructure and staged metadata; no target replay is invoked. A successful startup is not evidence of a completed source audit or owner approval.
