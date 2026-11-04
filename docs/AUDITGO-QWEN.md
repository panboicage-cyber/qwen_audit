# AuditGo v0.1.3 — coder.qwen.ai workspace integration

This is an **opt-in repository-side bootstrap**, not a claim that coder.qwen.ai has loaded this tool or can run a custom MCP server. The Qwen browser environment is remote: your laptop's `/home/sticky-ricky/...` paths are not accessible there. This repository is the **audit harness**, not the code under review.

## Code execution / environment settings

- Enable code execution to use the Linux binary through Shell (not proof of authorized testing).
- **Setup script** in the Qwen environment: `sh scripts/setup-auditgo.sh` (only if the setup runs after repository checkout). If it runs before checkout, first navigate to a checked-out copy of this repository.
- **Environment variables:** no secrets required. Optional `AUDITGO_ARCHIVE=/absolute/path/to/AuditGo_v0.1.3_Crush_Hermes_Linux.tar.gz`. Or `AUDITGO_ARCHIVE_URL=https://...` to retrieve exactly the pinned bytes. If neither is configured, store the archive at `vendor/AuditGo_v0.1.3_Crush_Hermes_Linux.tar.gz` or at the repository root.
- `PATH`: not required. Invoke `$HOME/.local/bin/auditgo`. The installer uses no Go, Node, Python, or MCP configuration.

## Artifact prerequisite and provenance

Obtain `AuditGo_v0.1.3_Crush_Hermes_Linux.tar.gz` from the original AuditGo v0.1.3 package, *not* from an arbitrary third-party URL. This deployment overlay includes the pinned archive in `vendor/`. Installation blocks if it is missing in the checked-out GitHub repository or does not match the expected hash. Expected SHA-256 of entire TAR.GZ:

`8c8625be74e54b70867196c88312d61b78b28f5737c81b7245b8339ca2e57640`

Expected Linux AMD64 binary SHA-256:

`b6c55aa3822c4dc6ae045c6a1f697a8f75903902e9e22ba957fa2e11f3a51a1f`

Expected 383-check registry SHA-256:

`385eaacb922eefa14231c59cc8c5f5cdfd9633938cc2daa5005045316af6451b`

The overlay becomes self-contained only after the archive is committed to GitHub and the Qwen workspace successfully runs the installer. A successful Qwen setup script is not evidence that any audit target was inspected.

## Preflight and first audit

```sh
echo QWEN_SHELL_OK
sh scripts/setup-auditgo.sh
sh scripts/auditgo-smoke.sh
$HOME/.local/bin/auditgo doctor
```

The smoke test creates temporary selector, plan and mission files only. To prepare a real audit, obtain the *target repository's* exact SHA, then run `auditgo init`, `auditgo plan`, `auditgo render`; assess the selected checks against the target using actual source and evidence, and use `auditgo validate`/`closeout` on the completed return. Technical risk must be distinguished from completeness, publication and authority.

## Failure behavior

| Status | Meaning |
|---|---|
| `AUDITGO_SETUP_PASS` | Pinned binary installed and doctor verified; no audit performed |
| `AUDITGO_ARCHIVE_MISSING` | No verified release archive provided |
| `AUDITGO_ARCHIVE_HASH_MISMATCH` | Wrong release bytes, refuse to extract |
| `AUDITGO_SETUP_UNSUPPORTED_PLATFORM` | Non-Linux or non-AMD64 Qwen workspace |
| `AUDITGO_SMOKE_PASS` | Local command orchestration works, not audit completion |

Do not grant MCP or agent mutation authority automatically. If coder.qwen.ai's browser product supports tool registration, validate it separately; shell execution does not demonstrate MCP registration.

## Existing target workflow

`prompts/kol-merkle-pr4.md` and `scripts/setup-qwen-target.sh` remain separate from this installer. Do not silently switch targets or call the workspace SHA the target SHA. Publication to another GitHub PR requires a separately authorized GitHub identity and permission.
