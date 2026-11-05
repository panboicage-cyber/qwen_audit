# Crush / independent reviewer handoff — SQL002-SCOPED-MACHINE-AUTHORITY-v0.1

Read `SQL002_SCOPED_MACHINE_AUTHORITY_DESIGN.md`, `SQL002_MACHINE_AUTHORITY_CONTRACT.json`, and `SQL002_ACCEPTANCE_MATRIX.md` as **proposed design inputs**, not verified deployed controls.

**Frozen audit subject**: `4d654e8ee11c12e8d10ef27b90ea2e84cf43c55e`.

**Authorized scope from user:** design scoped signing for some trusted automated processes, revocation, receipts and power-loss durability; do not grant blanket AI agent authority. The only definitive technical choice is durable power-loss protection. The other contract options are recommendations. No implementation, fixture execution, policy install, source repair, selector/ledger/gate mutation, or authority issuance was authorized.

**Task:** inspect this design against the frozen source and owner-scope packet. Look for unintended signing escalation, confused deputies, mutable-parent TOCTOU, incomplete fsync/rename sequencing, compromised-root/epoch rollback, stale revocation, duped receipts, crash windows, compromised storage and approval laundering. Propose corrections and explicit evidence prerequisites; do not quietly choose unresolved owner parameters. Verify JSON consistency, hash inputs, and report exact findings and missing evidence.

Keep the existing selected mandatory 34 essential/open and all 11 SQL whole-check statuses UNKNOWN, with `BLD-005`, SQL-004 generic constraints and SQL-005 explicit booleans retained. Zero whole-check N/A promotions and zero certified counterpart replacements. Gate 1 PENDING, Gate 2 BLOCKED, Gate 3 NOT_AUTHORIZED.
