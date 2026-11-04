# Frozen audit handoff: Gen Type-Go PR #6

**Target:** https://github.com/es-3581100/Gen_Type-go/pull/6

**Exact HEAD:** `0af4707d66bb30aef09fe80316c408568e256e53`

**Required tree:** `121189c9a8beacbddf3545ceeec79afbd8e39ea9`

**Reconciled main:** `7b071c7c0ed5037b107f4ae9a263f3c06853c46f`

**Pre-reconciliation v0.4:** `f8c6ba9461c01d165313a6bb8b1076d1cea2d576`

Existing GitHub evidence: PR #6 ready for review, unmerged. Exact-head push CI 41770732901 and PR CI 41770739176 both successful. Non-author technical comment PR6-REC-02 is not `APPROVED` and no owner merge decision has been given.

**Test targets:** `inspect --id NODE_ID FILE`, `diagnose FILE`; reuse `phasea.DecodeStrict`; canonical hashing; 2 MiB input, 4 MiB output; forest traversal, explicit edges, deterministic single best-effort diagnostic, rejected-input exit code 2.

**Threat model:** local untrusted graph JSON and source-provided metadata (IDs and content are not authenticated). No network ingestion, DB, issuer, policy activation, or production auth service. Historical proof must remain inert. The authority contract is blocked: `NEGATIVE_DISPOSITION_ORACLE_OPEN__RATIFICATION_PENDING__PHASE1_BLOCKED`.

**AuditGo:** v0.1.3, registry SHA-256 `385eaacb922eefa14231c59cc8c5f5cdfd9633938cc2daa5005045316af6451b`. Runner uses `init`, `plan`, `render`, and `seed` but must **not** claim seed validation or risk disposition until actual findings and evidence exist. The selected Go/canonical/build/compound modules are a bounded profile; no other inventory is silently treated as audited.

**Return:** exact observed source and tree hash; proof of 69-file manifest coverage and content hashes; deterministic test/fuzz logs, source and runtime evidence; findings with reproducible cases; complete/partial coverage ledger; SHA-256 receipts. Distinguish actual executed, observed, reported, and unverified claims. Final technical `AUDIT_PASS`, `AUDIT_FAIL`, or `AUDIT_INCOMPLETE`; no GitHub approval/merge/release without separate owner authority.
