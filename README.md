# Qwen Audit

Evidence-first, read-only code-review workspace for Qwen Coder AI. This repository contains review instructions, independent replay helpers, and report templates; it is **not** a mirror of projects under audit.

## First assignment: kol_merkle PR #4

- PR: https://github.com/donCannoli-burns/kol_merkle/pull/4
- Frozen candidate: `c64502cb5db3efd9316991ce1c0756d97729a6c3`
- Review target: all seven changed files / six commits, including the XML repair and surrounding invariants
- Verdicts: `AUDIT_PASS`, `AUDIT_FAIL`, `AUDIT_INCOMPLETE`

### Workflow

1. Read [AUDIT_POLICY.md](AUDIT_POLICY.md).
2. Provide Qwen [prompts/kol-merkle-pr4.md](prompts/kol-merkle-pr4.md).
3. Reproduce tests in a **sandboxed** isolated checkout; [scripts/replay-kol-merkle-pr4.sh](scripts/replay-kol-merkle-pr4.sh) is an optional helper.
4. Complete [reports/REVIEW_TEMPLATE.md](reports/REVIEW_TEMPLATE.md).
5. Publish the final review as a **PR conversation comment** using the actual submitting GitHub identity. Report execution gaps rather than inventing outcomes.

### Boundaries

No automatic GitHub `APPROVE`/`REQUEST_CHANGES`, merge, releases, deployments, or edits to source being audited. Audit PASS is a technical finding, **not** independent human approval or owner authorization. Use separate bot credentials if distinct GitHub attribution is required.

This starter does not install Qwen, configure a GitHub App or API keys, or start an automated workflow. Never commit credentials here.

## Troubleshooting Qwen's first audit run

The initial Qwen run reached an audit workspace whose only tracked file was the starter ZIP, and the Shell tool returned no verified output. See [QWEN_RECOVERY.md](QWEN_RECOVERY.md), [scripts/setup-qwen-target.sh](scripts/setup-qwen-target.sh) and [prompts/qwen-retry.md](prompts/qwen-retry.md). **A working shell is the first gate**; setup against the audit-workspace HEAD is not a valid test of the target PR SHA.
