# Independent adversarial design review — MA-AUDIT-002 / SQL-002 v0.2

- Reviewer: Qwen (independent technical challenger; not the designing architect, not owner)
- Date/time UTC: 2026-11-05
- Mode: READ-ONLY analytical review of a document-only proposal. No source changes to the frozen target, no target checkout/tests, no execution of `scripts/replay-*`, no key creation, no policy activation, no selector/ledger/gate mutation, no owner-decision approval.
- Intake contract: `docs/SQL002_QWEN_INDEPENDENT_REVIEW.md`
- Frozen audit subject (NOT modified or re-examined as code): `4d654e8ee11c12e8d10ef27b90ea2e84cf43c55e`
- Workspace branch: `qwen-code-2a49bf06-346e-4294-b8c0-0dbe8b7f1161` (harness branch, not target source)

## Disposition

**REVISION_REQUIRED** (design-level, conditional). The package is internally honest about its non-operative status and materially improves v0.1, but it contains at least two normative internal contradictions that make T05/T08 unreproducible as specified, several unsupported or circular guarantees, and an authority-core whose decisive properties are entirely delegated to an unselected, contract-less ACS provider. Under the package's own rules ("no implementation enforcement … proven", F01–F11 all `EVIDENCE_PROVEN=false`), this disposition is an independent technical challenge only — **not** an acceptance, PASS, owner ratification, or gate change. Gate 1 remains `PENDING / NOT_READY_FOR_ACCEPTANCE`; Gate 2 `BLOCKED`; Gate 3 `NOT_AUTHORIZED`. All 34 mandatory objectives remain essential/open (verified against `OPEN_OWNER_DECISIONS.md` §"Preserved mandatory inventory": count = 34, including SQL-002 and BLD-005).

## Evidence ledger (what was actually inspected)

| ID | Check | Class | Result | Evidence link |
|---|---|---|---|---|
| E01 | Bootstrap `sh scripts/setup-master-audit-workspace.sh` | EXECUTED | PASS; ended `MASTER_AUDIT_WORKSPACE_READY` | Script output incl. `MASTER_SETUP` overlay verify, `AUDITGO_SMOKE_PASS`, `REPO_MEMORY_SMOKE_PASS`, `MASTER_AUDIT_VERIFY_PASS workspace_sha256=71acab25… queue_sha256=ca9714b0…`, `SQL002_V02_DESIGN_EVIDENCE_READY (8/8 hashes)` |
| E02 | v0.2 manifest pin | EXECUTED | SHA-256 of `review_inputs/MA-AUDIT-002/SQL002_DESIGN_REVISION_002/SHA256SUMS.txt` = `1252a4debd76cf7824a9d1cb17d5503feb090ea162e025d11d73ef68dcbf8529` (matches intake) | local sha256sum |
| E03 | v0.2 artifact digests | EXECUTED | 8/8 OK via `sha256sum -c` in bundle dir | per-file digests recorded in `SHA256SUMS.txt` |
| E04 | v0.1 package pin + artifacts | EXECUTED | manifest = `6e2dd26716944f1ed98ec2daad81d7967e44d9fb7e42a29510e2ddbdb598afae` (matches intake); 4/4 OK | local sha256sum |
| E05 | JSON validity of both v0.2 JSON docs | EXECUTED | parse OK; transitions T01–T13 present (13); typed envelope set = delegation/admission/authorization/storage_prepared/commit_event/ack (6) | python3 json.load |
| E06 | Full-text read of all 8 v0.2 artifacts | OBSERVED | complete | files listed below with line citations |
| E07 | v0.1 design rename-as-commit-point baseline | OBSERVED | v0.1 §D:86 made `CURRENT` rename "the visibility commit point"; v0.2 removes it (ARCH §1, §5) | `SQL002_MACHINE_AUTHORITY_V01/SQL002_SCOPED_MACHINE_AUTHORITY_DESIGN.md:86` |
| E08 | Original adversarial reports (`SQL002_ADVERSARIAL_REVIEW_001.md`, `SQL002_FAILURE_MODE_MATRIX_001.md`, `SQL002_REVISION_RECOMMENDATIONS_001.md`) | UNVERIFIED — **upstream evidence unavailable** | Not bundled in this workspace; pinned expected digests recorded in intake. Astra High's F01–F11 findings and "Agree" corrections treated as untrusted claims; no original wording was invented or relied upon. Request the three pinned files before any claim requiring exact original text. | intake doc §Bootstrap |
| E09 | Frozen target revision content | NOT INSPECTED | Per intake, deliberately not fetched/mutated/tested | — |

Read artifacts: `INDEPENDENT_REVIEW_HANDOFF.md`, `ROOT_CAUSE_RECONCILIATION.md`, `SQL002_ARCHITECTURE_V02.md` (ARCH), `TRANSACTION_STATE_MACHINE.json` (TSM), `AUTHORIZATION_CONTRACT_V02.json` (AC), `CRASH_RECOVERY_MATRIX.md` (CRM), `REVIEW_FINDING_CLOSURE_MATRIX.md` (FCM), `OPEN_OWNER_DECISIONS.md` (OOD). Abbreviations below cite these.

## Independent findings

Severities are reviewer-assigned for this design review and do not modify the preserved original F01–F11 severities.

### HIGH

**QV2-01 — Hidden second anchor inside the single-anchor mandate (contradiction; authority bypass).**
ARCH §5: "ACS Commit executes serializably **for K and policy/delegation quota state**". TSM T05 `synchronization`: "Single crash-durable linearizable ACS transaction **including policy/delegation/quota rows**; no external second anchor update". But AC `chain_validation.policy_join` requires "policy_digest and authority_epoch equal across delegation, admission and authorization **and current ACS policy**", and ARCH §2 states ACS durably holds "never-reused policy versions" and root/delegated key state. Revocation/policy updates (T09) mutate epoch/policy/key state that T05 reads. If policy/root/delegation state is *not* inside the same serialized transaction domain as K's membership (as the narrower ARCH §5 scope implies), then epoch/policy consumption is exactly the undefined cross-resource atomicity the redesign claimed to eliminate — a hidden second anchor behind the "one ACS transaction" label, which the handoff itself forbids ("Never accept a separate second anchor update hidden behind a provider label", IRH:25). Either the architecture text understates the required transaction scope or the state machine overstates coverage. Proof obligation: a normative statement that T05, T09, quota charge, consumed-ID insertion, and policy/epoch/key-state validation all serialize in one ACS linearization order, with a revoke-vs-commit race schedule per data center. Affected: F01/F03 closure claims.

**QV2-02 — COMMIT outbox event schema omits the immutable request digest, breaking durable idempotency (crash-recovery weakness, contradicts ARCH §2).**
ARCH §2 defines transaction identity as "`(K, authorization_id)` **and immutable request digest**; conflicting reuse always rejects", and T05/T04 `retry` rules ("Same transaction returns same result"; "Identical request returns same reservation; changed request rejects") require ACS to compare a stored request digest on replay. Yet AC `types.commit_event.exact_required_fields` (line 71) contains no request-digest field, and no other ACS-resident type carries it. After an ACS restart or outbox replay, a same-`(K, authorization_id)` retry with a *different* payload cannot be mechanically rejected from the committed event alone; "conflicting reuse always rejects" is unimplementable as specified. Counterexample schedule: publisher crashes after T05 durable commit but before response; attacker/bug retries same authorization_id with altered manifest bytes; spec provides no durable basis to reject. Fix obligation: bind the canonical commit-request digest into the COMMIT event (or another durable ACS record) and into T04's registration record schema.

**QV2-03 — Storage PREPARED evidence is not join-bound to admission/job, enabling cross-job laundering through the honest-storage assumption (authority bypass).**
AC `types.storage_prepared` fields include `authorization_id`, `generation_id`, digests, `fence`, `sealed`, `sync_complete` — but **no `admission_id`, no `job_id`, no `policy_digest`/`authority_epoch`**. AC `artifact_join` binds prepared evidence only to "exact authorization envelope digest and all represented artifact fields". Consequence: a storage principal (or a captured/stale-but-reissued fence context) can sign `SQL002_STORAGE_PREPARED` for a generation digest under an authorization belonging to job A while the actual staged bytes originate from a non-admitted source B, provided the digest matches; T05 checks admission validity of the *authorization*, not of the *bytes' provenance path*. This partially reopens F04/F06 ("authorization is a decision, not possession of a signer endpoint") at the storage step, and interacts with ARCH §4's explicit caveat that "this assumes honest qualified storage-service assertions; signatures cannot independently prove successful flush." Proof obligation: bind admission_id/job_id (and epoch/policy) into the prepared envelope and the chain joins, plus a negative test where valid storage evidence for a foreign admission must deny.

**QV2-04 — Reader session/pin receipts have no cryptographic contract and no freshness binding; I6's anti-replay guarantee is asserted, not designed (unsupported guarantee).**
ARCH §7 and TSM T10 require "challenge-bound responses, trusted service identity, monotonic checkpoint" and forbid "copied signed responses" establishing current membership; FCM's required challenge asks whether "new reader admission can replay a cached ACS response". But AC defines **no session receipt, pin receipt, or read-response envelope type at all** (only delegation/admission/authorization/storage_prepared/commit_event/ack exist; `pin_id` appears solely as an opaque ID in `normative_encoding.opaque_ids`). There is therefore no normative way to verify what T10 calls "Pin receipt" (TSM `audit_evidence`) or to enforce challenge binding, audience, nonce, or monotonic-checkpoint inclusion in a read response. The I6 rollback-protection guarantee currently rests on prose, not a signed contract — an evidence gap indistinguishable (to an implementer) from the v0.1 problems F07/F09 identified. Obligation: normative read-admission/pin envelope types with session-nonce and checkpoint fields, plus replay negatives in `scope_limits.test_vectors`.

**QV2-05 — Indefinite retention collides with abort tombstones and "terminal" ABORTED state; GC grant retirement is unbounded (architectural contradiction / availability honesty gap).**
ARCH §7 mandates "no GC of committed generations, committed authorizations, membership or required audit evidence"; ARCH §8a requires "T03/T04 must not reuse a digest path until **every previous deletion grant is durably retired**" and T12 creates a "Terminal abort tombstone retaining used request identity". Nothing bounds tombstone/grant-retirement history, and every aborted same-digest attempt permanently blocks re-preparation of that digest path unless retirement itself is a defined, provable, bounded operation — which the package never specifies (no retirement transition exists in T01–T13; CRM has no cut row for grant retirement). Counterexample schedule: repeated aborted preparations of the same generation digest → permanent inability to prepare that content even though no committed object exists; capacity denial without any committed fact justifying it. Owner consequence (O12) is understated: exhaustion affects uncommitted paths too. Obligation: define retirement semantics (transition, durability cuts, oracle) or explicitly accept per-digest poisoning by aborted attempts as an owner-visible cost.

### MEDIUM

**QV2-06 — `COMMITTED_UNACKNOWLEDGED` is a fourth publication state missing from the normative state enumeration.** Used in ARCH §5/§6 and CRM T05-row and FCM F02, but absent from TSM `states` and `orthogonal_states`. Since the package stakes its correctness story on exhaustive transition/state coverage ("The design includes 13 transition specifications, every required transition field", IRH:21), an ad-hoc state name outside the enumeration is a normative inconsistency; add it (likely as `delivery=PENDING ∧ knowledge` combination) or replace the usage.

**QV2-07 — Manifest role-path ordering is internally inconsistent.** AC `manifest.role_paths` enumerates `catalog.json`, `digest-sidecar`, `dependent-context`, while `file_order` is `["catalog.json","dependent-context","digest-sidecar"]`. Neither is marked normative for JCS array ordering of `files`, so two implementers can compute different `manifest_digest`s for identical bytes — violating the deterministic-derivation requirement ARCH §3/IRH:36 flags as unproven, and making `test_vectors` ambiguous. Pin one ordering and state that `files` must serialize in `file_order`.

**QV2-08 — "Independent audit store" identity lacks a typed key use; independence is asserted, not enforced.** ARCH §3 says "Independent audit store identity is separate again", but AC `key_usage` enum instances are OWNER_DELEGATION_ADMIN / JOB_ADMISSION / CATALOG_AUTHORIZATION / STORAGE_EVIDENCE_ONLY / ACS_EVIDENCE_ONLY — none covers log/witness/export protocol messages, and `scope_limits.receipt_schemas` defers the log wire contract entirely. Cross-type substitution tests (FCM F10) cannot cover log-side impersonation until an AUDIT_STORE/WITNESS key use and inclusion-receipt schema exist. Keep F10 at MEDIUM severity retained, not closed.

**QV2-09 — Key rotation vs. historical verification is contradictory across documents.** OOD O11 and AC `policy_enforcement.key_rotation` assert normal rotation "preserves eligible historical membership"/"does not retroactively invalidate eligible committed data", while ARCH §7's compromise default quarantines "all generations authorized by that key since its activation". Rotation that retires a key leaves historical envelopes verifiable only if old public keys remain pinned and trusted for verification; no transition or field governs verifier trust-window retention for rotated keys (T09 changes key state; readers verify "signed inventory" with which key set?). Specify post-rotation verification key-set semantics distinct from revocation.

**QV2-10 — Signer idempotency record durability is outside the qualified power-loss model.** T02 `synchronization`: "Durable signer record before returning envelope". The signer is a distinct service from ACS and the data host, yet ARCH §8 qualifies durability only for data-host ext4 profile, ACS, and audit store. If the signer's "durable" record loses a write on power cut, `crash_restart` ("Recover same envelope or fail closed") has no defined recovery source, risking either double signature of divergent payloads under one `authorization_id` (violating T02 `retry`: "Same authorization ID and payload only") or unavailability. Add the signer platform to the qualification profile obligations or move idempotency state into ACS.

**QV2-11 — Single-writer serialization cost and dataset-scoped lock semantics are unquantified (SCOPE-05 risk).** T05 serializes per-K "and policy/delegation quota state"; combined with QV2-01's ambiguity, worst case is effectively one global serial writer for policy state. OOD O05 acknowledges scope risk but the package offers no throughput/fairness analysis; a shared quota bucket across datasets (CRM schedule "Delegation quotas across keys/workers/datasets") can serialize unrelated datasets through the same rows. Owner should see this cost stated before accepting O02/O06.

### LOW / completeness notes

**QV2-12 — Genesis bootstrap multi-document consistency untested.** TSM `genesis` and AC `manifest.genesis_previous` require an owner-provisioned sequence-zero tuple under an ACS transaction; there is no genesis envelope type nor a CRM cut row for partial genesis provisioning (root created, dataset not yet). Low severity because fail-closed defaults likely deny, but specify.

**QV2-13 — `audience_join` self-referential configuration dependency.** Audience equality is defined against "independently configured ACS commit-service identity" whose distribution mechanism is only described for root updates (ARCH §3). Until the workload-authenticator's audience provisioning is normatively specified, audience confusion between a rogue ACS stand-in and the real one is only prevented by operational care, not by signed fields.

## Assessment of Astra High-derived claims (untrusted-input handling)

- The eleven findings F01–F11 and their severities (9 HIGH, 2 MEDIUM) originate from upstream reports **not present in this workspace** (E08). Their text was not verified; only the architect's restatements in ROOT_CAUSE_RECONCILIATION.md were read. Accordingly, no F-status here is treated as independently verified, and none is re-certified: DESIGN_ADDRESSED ≠ EVIDENCE_PROVEN, exactly as the package itself insists.
- The "shared invariants" synthesis (RCR §"Shared invariants") is directionally sound — consolidating current/membership/nonce/quota/outbox into one transaction genuinely removes v0.1's rename-as-commit-point error (confirmed against v0.1 D:86 in E07). However, per QV2-01/02/03, the consolidation is asserted at prose level while the normative schemas still leave anchors (request digest, policy state scope, admission↔storage binding) unjoined.
- The conditional induction proof (ARCH §9) honestly lists its dependencies; this review confirms those dependencies are load-bearing rather than incidental: linearizable crash-durable ACS, honest qualified storage, and monotonic bootstrap are precisely the three assumptions with zero provider contracts in the package ("no approved provider exists here", ARCH §8). The design therefore currently certifies nothing, which the documents acknowledge — the residual problem is that the closure matrix language could tempt an owner into reading DESIGN_ADDRESSED as architectural sufficiency; QV2-01..05 show it is not yet sufficient even at design level.

## Contradiction / assumption / proof-obligation register (condensed)

| # | Type | Location(s) | Required resolution |
|---|---|---|---|
| QV2-01 | Contradiction | ARCH §5 vs TSM T05 vs AC policy_join + ARCH §2 | One-transaction scope statement incl. policy/epoch/key rows; revoke-race schedules |
| QV2-02 | Gap/crash-weakness | ARCH §2 vs AC commit_event vs TSM T05 retry | Durable request-digest binding; replay-with-different-payload negative test |
| QV2-03 | Authority bypass | AC storage_prepared/artifact_join vs F04/F06 intent | Admission/job binding in prepared envelope; foreign-admission negative test |
| QV2-04 | Unsupported guarantee | ARCH §7, TSM T10/T11 vs AC (no types) | Normative session/pin/read-response envelopes; cached-response replay tests |
| QV2-05 | Contradiction | ARCH §7 vs §8a vs TSM T12 tombstone | Retirement transition + durability cuts, or declared per-digest poisoning cost |
| QV2-06 | Inconsistency | ARCH §5/§6, CRM, FCM vs TSM states | Enumerate COMMITTED_UNACKNOWLEDGED normatively |
| QV2-07 | Ambiguity | AC manifest.role_paths vs file_order | Pin canonical files-array order |
| QV2-08 | Evidence gap | ARCH §3 vs AC key usages vs scope_limits | Audit-store/witness key use + receipt schema |
| QV2-09 | Cross-doc contradiction | OOD O11 / AC key_rotation vs ARCH §7 compromise default | Post-rotation verifier key-set semantics |
| QV2-10 | Crash-weakness | TSM T02 vs ARCH §8 profiles | Signer durability qualification or ACS-held idempotency |
| QV2-11 | Assumption | TSM T05 sync clause; OOD O05 | Throughput/serialization analysis for owner |
| QV2-12 | Gap | TSM genesis / AC genesis_previous | Genesis envelope + partial-provisioning cut |
| QV2-13 | Assumption | AC audience_join | Audience provisioning contract |

## What was NOT done / explicit gaps

- Did not fetch, mutate, build, or test frozen target `4d654e8e…`; did not run `scripts/replay-*`; created no keys; executed no audits; changed no PR states; approved nothing; advanced no MA-AUDIT-002 gates; wrote no RepoMemory capture (setup confirmed `capture=DISABLED`).
- Upstream original review reports unavailable (E08) — any claim needing their exact wording is marked insufficient evidence; the three pinned files are requested, not assumed.
- No live ACS, provider, hardware PLP, clock, or cryptographic-canonicalization interoperability was tested; none was claimed. This report is not an authenticated independent sign-off, does not confer EVIDENCE_PROVEN or OWNER_RATIFIED on any finding, and does not pass SQL-002.

## Mandatory-objective preservation check

Independently recounted from OOD §"Preserved mandatory inventory": 34 identifiers listed (AGT×7, BLD×2 incl. BLD-005, CORE×10, GNB×5, HSH×3, LED×1, PRC×3, SQL-002×1, TSP×3) — consistent with intake, ARCH §10, RCR closing line, CRM closing line, FCM closing line, and setup output (`mandatory_selected_essential: 34` in both JSON pins). All remain **essential/open**; eleven SQL whole-checks UNKNOWN; zero counterpart equivalences; zero N/A transitions; 296 selected objectives unchanged.

## Recommended next actions (owner/architect, none authorized by this review)

1. Issue v0.3 revisions resolving QV2-01..05 (HIGH) at normative-schema level, not prose.
2. Supply the three pinned upstream reports so F01–F11 original text can be verified before any closure discussion.
3. Produce normative provider contracts (ACS, audit store/witness, identity) with the added envelope types from QV2-04/08 before any activation conversation.
4. Present O02/O05/O12 costs to the owner including QV2-05 and QV2-11 consequences.
