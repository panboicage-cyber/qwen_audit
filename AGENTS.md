# Qwen audit workspace rules

This repository is an audit harness, NOT the source repository being audited. Do not confuse its commit SHA with any target PR commit. Read `AUDIT_POLICY.md` and `QWEN_RECOVERY.md` before running tests.

## AuditGo v0.1.3

The optional AuditGo CLI is installed only after `scripts/setup-auditgo.sh` prints `AUDITGO_SETUP_PASS`. Check with `$HOME/.local/bin/auditgo doctor` and `sh scripts/auditgo-smoke.sh`. Shell access is not a claim of test execution: record stdout and exit code. `auditgo init/plan/render` generates a mission; it does not inspect or verify the target on its own.

The setup is pinned to exact archive and binary SHA-256 values. If its archive is absent or the hash fails, report `AUDITGO_SETUP_BLOCKED` and do not substitute an unverified release. The deployment overlay contains `vendor/AuditGo_v0.1.3_Crush_Hermes_Linux.tar.gz`; confirm it was actually committed and fetched before claiming installation. Do not run arbitrary downloaded scripts, install global packages, or weaken the hash check.

Perform audits against a separate explicitly pinned target checkout. Maintain `EXECUTED`, `OBSERVED`, `REPORTED`, `UNVERIFIED` evidence distinctions and preserve `AUDIT_PASS`, `AUDIT_FAIL`, `AUDIT_INCOMPLETE` from repository policy. The AuditGo internal risk/completeness states are supporting evidence and do not override these outcomes. A generated or validated evidence receipt is not independent human approval.

Return a short operational summary, then exact artifact paths, SHA-256 of materialized files, verified and unverified coverage, and next gates. Do not commit generated review data unless separately authorized. Do not approve a PR, merge, publish, or deploy by inference. Never expose credentials.
