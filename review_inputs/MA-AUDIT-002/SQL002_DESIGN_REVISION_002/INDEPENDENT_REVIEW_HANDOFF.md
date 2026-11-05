# SQL-002 v0.2 independent review handoff

**Disposition: DESIGN_V02_READY_FOR_INDEPENDENT_REVIEW.** Ready means a conditional design package is available to challenge, not that SQL-002 passes, authority is ratified or software exists. This assistant is now the repairing architect and must not authenticate its own independent acceptance. Read-only helper review identified adjacent consistency gaps; its feedback is design assistance, not independent approval.

Frozen revision: `4d654e8ee11c12e8d10ef27b90ea2e84cf43c55e`. Preserve v0.1 and original adversarial reports. All changes are in this isolated directory. No production code, live operations, target/fixture/model tests, signing keys, policy activation, selector/ledger changes, N/A transitions or gates were authorized or executed. Documentary parsing/hash checks only.

## Read order and authority

1. Verify SHA256SUMS.txt from this directory; verify original input pins in ROOT_CAUSE_RECONCILIATION.md.
2. Read ROOT_CAUSE_RECONCILIATION.md, SQL002_ARCHITECTURE_V02.md and TRANSACTION_STATE_MACHINE.json together.
3. Read AUTHORIZATION_CONTRACT_V02.json, especially exact field rules and normative chain/storage/evidence refinements.
4. Read CRASH_RECOVERY_MATRIX.md and REVIEW_FINDING_CLOSURE_MATRIX.md; each is a proposed evidence specification, never test authorization.
5. Read OPEN_OWNER_DECISIONS.md before interpreting guarantees as effective choices.

The architectural change is explicit: authoritative current, membership, sequence/epoch, authorization/job consumption, quota charge and COMMIT outbox obligation share one independently protected ACS transaction. Files are sealed and durable before registration/commit. Local CURRENT and external receipt delivery are recoverable projections, not jointly atomic resources. No offline CURRENT trust fallback. This change needs owner scope review and a qualified coordinator; provider independence and transaction guarantees are assumptions requiring proof.

## Closure status

F01 HIGH, F02 HIGH, F03 HIGH, F04 HIGH, F05 HIGH, F06 HIGH, F07 HIGH, F08 HIGH, F09 HIGH, F10 MEDIUM and F11 MEDIUM are each DESIGN_ADDRESSED. Each remains EVIDENCE_PROVEN=false and OWNER_RATIFIED=false. No original severity is reduced. Exact per-finding invariants, corrections, proof obligations, negative cases, expected outcomes and remaining limitations are in REVIEW_FINDING_CLOSURE_MATRIX.md.

The design includes 13 transition specifications, every required transition field, every durable cut and revocation/concurrency schedules. The conditional induction proof is architecture section 9, not a mechanized proof or executed model. Independent reviewers should attack its assumptions, not count document fields as assurance.

## Focused adversarial challenges

- Can actual ACS provider make all current/membership/ID/quota/policy/event changes one crash-durable transaction while preventing snapshot rollback and split brain? Never accept a separate second anchor update hidden behind a provider label.
- Can storage-prepared retention be invalidated by delayed GC? Normative refinement binds grants to transaction/registration/object incarnation and prohibits path reuse until old grants are durably retired. Plan a stale cleanup replay against a newly prepared identical manifest; it must reject without touching new data.
- Can signatures from individually valid foreign chains be mixed? Normative chain joins require identical dataset, principal, job, admission/delegation, recipe/schema, manifest, policy/epoch and audience plus interval containment. Plan cross-chain substitutions even when every signature is valid; all mismatches must deny.
- Can key rotation change a COMMIT envelope after log delivery and break ACK binding? T07 fixes one envelope durably before export; payload-chain and envelope digests are distinct. Plan crash/rotation between envelope persistence, log append and ACK; retries must use identical envelope bytes.
- Does independently admitted job policy genuinely constrain content, rather than simply laundering agent requests through a trusted identity? Exact recipes and trusted triggers are still unconfigured and default DENY.
- Are old pinned reads, new-read rollback exclusion, compromise quarantine and in-flight delivery boundaries compatible? Final delivery authorization is an explicit linearization point, not a promise to recall bytes after later revocation.
- Can a committed generation lose its exact content before ACK or later restore? v0.2 forbids committed GC; any missing content causes unavailable/quarantine without erasing the committed fact. Storage/controller/ACS durability remains unproven.
- Can a receipt exporter or operator bypass authority? Evidence keys cannot authorize catalogs/policy; administrative/restore/GC paths share ACS rules, and host-root compromise is honestly bounded.

## Unproven guarantees and open decisions

No implementation enforcement, isolation, canonicalization interoperability, deterministic derivation, clock/freshness behavior, coordinator durability/rollback resistance, log non-equivocation, power-cut persistence, source-to-running-process binding or authenticated independent acceptance is proven. Linux/ext4/local PLP is only a proposed qualification profile with exact patch/device/firmware parameters unset. ACS, log and identity/policy wire adapters require independently reviewed normative provider contracts before activation. The design explicitly cannot operate with those unset.

Owner must separately decide consistency/publication/authentication formal choices; SCOPE-01..05; ACS/provider/online availability; machine principals/datasets/admission recipes and budgets; clock/TTL; key custody/root recovery; exact storage profiles; quarantine/in-flight delivery; indefinite retention/capacity; audit-store/ACK contract and independent review authority. Required power-loss durability and no agent authority are not negotiable defaults. No owner values are selected here.

## Documentary validation and manifest interpretation

The input request and original review/package manifests were verified before work. Each generated/edited document received immediate syntax/structural checks. JSON parse, transition-field coverage, eleven finding status rows, 34 mandatory inventory and exact input hashes are checked as documentation integrity only. SHA256SUMS.txt pins the eight deliverable files; it deliberately does not hash itself. Its digest is supplied separately with the final response. Original manifests and review bytes are rechecked after creation.

All 34 mandatory objectives remain essential/open, including mandatory SQL-002 and BLD-005; SQL-004 generic constraints and SQL-005 explicit native booleans remain. 296 selected, eleven SQL whole checks UNKNOWN, zero counterpart equivalence certifications, no N/A or gate promotion. Gate 1 PENDING / NOT_READY_FOR_ACCEPTANCE; Gate 2 BLOCKED; Gate 3 NOT_AUTHORIZED.
