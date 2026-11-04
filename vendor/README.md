# AuditGo archive placement

This folder is an optional location for the **pinned v0.1.3 Linux AMD64 archive**.

The release archive `AuditGo_v0.1.3_Crush_Hermes_Linux.tar.gz` is included in this deployment overlay. Commit the archive along with this README to make the repository self-contained. Alternatively configure `AUDITGO_ARCHIVE` / `AUDITGO_ARCHIVE_URL`. The installer verifies its SHA-256 before unpacking and checks the installed binary SHA-256 and registry identity.

See [`docs/AUDITGO-QWEN.md`](../docs/AUDITGO-QWEN.md) for expected hashes and recovery steps. Never substitute unpinned downloads or disable the integrity check. Do not commit secrets or audit-target artifacts here.
