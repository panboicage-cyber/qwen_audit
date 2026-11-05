# SQL-002 — Acceptance/Evidence Matrix (Design Only)

No test is authorized or claimed executed by this matrix. Tests must be performed on separately approved disposable fixtures/platforms with exact revision-locked drivers and receipts. All 34 selected mandatory objectives remain essential/open.

| Group | Case | Expected result / evidence | Cross-check |
|---|---|---|---|
| CONSISTENCY | Reader starts before pointer swap and continues after | Sees complete old generation or complete new generation; never mixed | SQL-002, CORE-009 |
| CONSISTENCY | Concurrent file truncation/replacement attempt in staged/published data | Deny mutation or reject changed bytes; never accept mismatched object/manifest | HSH-023, GNB-028 |
| CONSISTENCY | Single path/parent symlink or mount replacement | Open fails closed via verified directory descriptors/mount policy | GNB-027/028, PRC-003 |
| PUBLICATION | Fault after each individual write and fsync | CURRENT remains prior authorized generation until complete commit | SQL-002, CORE-003 |
| PUBLICATION | Two publishers with same parent / duplicated nonce | Exactly one generation accepted and serialized; other rejected or reconciled | SQL-002, CORE-005 |
| PUBLICATION | Committed pointer but ACK/remote receipt fails | `COMMITTED_UNACKNOWLEDGED`, retry recovers same authorization without duplicate commit | SQL-002, PRC-004 |
| DURABILITY | SIGKILL/process crash at every recorded boundary | Restart selects complete authorized previous or committed generation; no torn active | CORE-003, PRC-004 |
| DURABILITY | Abrupt host power cut at every critical fsync/rename stage on pinned storage profile | Every `ACKNOWLEDGED_DURABLE` generation survives; no corruption/rollback | GNB-029, PRC-004 |
| DURABILITY | Restore old directory, newer CURRENT, partial backup; lost/lying flush | Reject invalid restore; explicitly report unsupported storage guarantees | SQL-010, PRC-010 |
| AUTHENTICATION | Agent without delegation submits coherent hash pair | DENY despite matching hashes | SQL-007, CORE-005 |
| AUTHENTICATION | Machine signer requests other workspace/path/operation/size/expiry | DENY; signed denial trace without private key leakage | SQL-002, AGT-043 |
| AUTHENTICATION | Compromised worker with access to signing endpoint asks for arbitrary digest | DENY because policy service evaluates exact authorized manifest & principal | CORE-005/006, HSH-030 |
| AUTHENTICATION | Revoke key/epoch while request staged, before commit | DENY unless an explicitly defined safe grandfathering contract is separately authorized (recommended: none) | SQL-002, LED-001 |
| AUTHENTICATION | Replay old authorization/nonce/parent after new generation | DENY | SQL-007, HSH-031 |
| AUTHENTICATION | Restore old signed policy + CURRENT together | Reject using separately protected epoch/sequence anchor | SQL-007, PRC-010 |
| AUTHENTICATION | Cryptographic verification: missing signature, wrong key, canonicalization ambiguity, duplicate JSON fields | DENY | HSH-001, BLD-005 |
| RECEIPTS | Independent auditor examines one approved and one denied update | Reconstruct exact workload, delegation/policy hash, manifest, publisher result, source/runtime identity, revocation epoch | CORE-006, LED-001 |
| RECEIPTS | Crash after pointer swap before external logging | Durable outbox allows recovery and append without falsifying original event time | PRC-004, CORE-006 |
| PROVENANCE | Deployed verifier/publisher binary differs from audited source | Flag unresolved/reject assurance claim, never treat checksum of source archive as binary linkage | BLD-005 |

**Criteria for objective PASS are not defined solely by green test suite count.** Reviewer requires observed command, environment, toolchain, commit/source and binary provenance, preserved stdout/stderr/exit, exact test input/expected/actual, supported filesystem profile, and authenticated acceptance of residual limits.

**Unresolved owner choices:** consistency A; publication A; authentication A are recommendations, not ratifications. Durability A is the explicit design requirement; contract recording and implementation remain pending.
