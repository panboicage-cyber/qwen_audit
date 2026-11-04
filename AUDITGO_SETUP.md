# AuditGo for Qwen — workspace quick start

This repository is the audit harness. It is **not** the source tree under review.

In the coder.qwen.ai environment settings, enable code execution and enter this **Setup script** (provided the repository has been checked out):

```sh
sh scripts/setup-auditgo.sh
```

No environment variables or API keys are required: the pinned v0.1.3 Linux AMD64 archive is in `vendor/`. The setup refuses the wrong platform and mismatched hashes. If the workspace shell is unavailable or the setup fails, stop and report a blocker; never claim installation or tests passed.

After setup, test shell access and the installed binary:

```sh
echo QWEN_SHELL_OK
$HOME/.local/bin/auditgo doctor
sh scripts/auditgo-smoke.sh
```

`AUDITGO_SMOKE_PASS` proves only planning and prompt generation. For an actual review, follow `AUDIT_POLICY.md` and the separate pinned target checkout instructions. Use real source/test evidence, generate an audit receipt, and report `AUDIT_INCOMPLETE` whenever required evidence is missing. Do not approve, merge, release, or publish an audit result without separate authority.

Full details: [docs/AUDITGO-QWEN.md](docs/AUDITGO-QWEN.md).
