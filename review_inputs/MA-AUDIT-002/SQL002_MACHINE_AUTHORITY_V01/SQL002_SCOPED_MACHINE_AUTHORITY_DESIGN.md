# MA-AUDIT-002 / SQL-002 — Scoped Machine Authorization & Durable Generations

**Document class:** Design proposal; NOT an authorization, implementation, audit result, gate transition, or approval.  
**Version:** `SQL002-SCOPED-MACHINE-AUTHORITY-v0.1`  
**Frozen audit subject:** `4d654e8ee11c12e8d10ef27b90ea2e84cf43c55e` (unchanged).  
**Owner direction:** Some trusted automated processes may authorize catalog replacements, without blanket agent authority. Power-loss durability is required.

## 1. Four independent SQL-002 contract positions

| Contract ID | Proposed option | Current authority state | Guarantee to test |
|---|---|---|---|
| `SQL002-CONSISTENCY` | `A`: coherent immutable generation; complete older generation permitted | `PROPOSED_NOT_SELECTED` | Every successful read binds to one generation and validates its complete artifact set. No freshness/linearizability claim. |
| `SQL002-PUBLICATION` | `A`: one atomic switch to complete generation | `PROPOSED_NOT_SELECTED` | Externally visible active pointer never refers to a partially written or unauthorized generation. |
| `SQL002-DURABILITY` | `A`: power-loss durable acknowledgement and recovery | `OWNER_REQUIREMENT_CONFIRMED_DESIGN_PENDING_RATIFICATION` | An acknowledged publication survives power loss on an approved, evidenced filesystem/storage profile; recovery yields an authorized, intact generation. |
| `SQL002-AUTHENTICATION` | `A`: cryptographically verified exact-generation authorization, including delegated machine roles | `PROPOSED_NOT_SELECTED` | Reader/publisher verify an owner-root-authorized signer, exact generation identity, policy scope, active epoch, and rollback controls. |

**No scope ratification or Gate 1 acceptance is implied by these design positions.** Historical frozen source has separate `catalog`/sidecar reads and separate context/sidecar publication; it is not certified to meet this future contract. None of the 34 mandatory audit checks is waived. The SCOPE-01..05 A/A/A/A/B choices remain proposals.

## 2. Trust zones and capabilities

| Zone | Can do | Must not do |
|---|---|---|
| Candidate producer (AI agent, build tool, worker) | Generate proposed catalog/context bytes; compute hashes; submit to intake | Sign, modify policy, publish active pointer, revoke keys, impersonate signer |
| Authenticated machine workload | Request authorization under a pinned delegated capability | Request arbitrary hashes to sign; choose its own identity, scope, trust root or policy revision |
| Independent policy/signing service | Verify workload identity and immutable delegation; issue scoped, signed one-use authorization | Accept free-form agent policy, delegate more rights than owner root, expose private-key material |
| Restricted publisher | Verify receipt/signature/epoch, safely read prepared bytes, serialize publication, enforce durable commit | Sign/expand policy or bypass revocation; write unrelated directories |
| Reader/verifier | Verify immutable generation, exact paths/length/digests/signature/authority epoch; reject unauthorized data | Trust matching local hashes as ownership authentication |
| Owner / offline root | Create/revoke delegated policy and recovery authority, rotate root via separate process | Be impersonated by a documentation field, chat transcript, or AI-issued approval string |
| Independent evidence reviewer | Inspect revision-bound receipts, tests, failure and scope claims | Be treated as publisher or signer merely for reviewing |

Private keys should be non-exportable or in an isolated signer/HSM/TPM/KMS-backed service with distinct access policies. Hardware protection **alone is insufficient** if an agent can invoke an unrestricted signing API. Do not place machine signing credentials, root keys or policy-admin credentials in the agent's filesystem, environment, OpenCode tools or general-purpose shell. Enforce this with OS/service permissions, workload authentication, network boundaries and signed policy, not instructions to the model.

## 3. Delegation constraints (all must be checked)

Delegation is owner-root-signed (or approved via separately specified owner-root threshold). It binds:

- `trust_domain`, `workspace_id`, exact `dataset_id` and explicit path allowlist (catalog + exact associated sidecar/context set; no unchecked globs, `..`, symlinks, mount escapes or arbitrary destinations).
- Named workload principal, key ID, signing-service audience and role; permitted operation `publish_catalog_generation`, not `sign_anything`.
- Maximum bytes, maximum file count, rate/quota limit, permitted schema and policy versions, and permitted generation ancestry.
- `not_before`, `expires_at`, one-use authorization ID, exact canonical manifest digest, current authority epoch and policy digest.
- Separate recovery escalation and delegation management forbidden to automated normal signers.

A request must fail closed on any mismatch, stale/expired policy, absent authority epoch, duplicate authorization ID, changed file inventory, invalid parent/sequence, or unavailable required revocation source. Multiple automated signers can have different datasets and bounded budgets; access is not transitive. High-impact policy changes, trust-root rotation and rollback recovery require explicit owner-controlled authority and ideally multiple independent approvals.

## 4. Exact signed subject

Choose a reviewed signature scheme (illustratively Ed25519 where supported by the selected hardware/key provider). Sign **domain-separated canonical bytes** using a well-defined canonicalization contract such as RFC 8785 JCS for I-JSON. Do not use ad-hoc serializer ordering or ambiguous JSON numbers. Signatures must cover, at minimum:

```
contract_version, trust_domain, workspace_id, dataset_id,
generation_id, previous_generation_id, monotonic_sequence,
file_inventory[{relative_path, size_bytes, sha256}],
manifest_digest, policy_digest, authority_epoch,
issuer_key_id, workload_principal, authorization_id,
operation, not_before, expires_at
```

Use an independently provisioned public-key trust root outside the writable catalog. Validate schema, duplicate-key rejection, signature, signed path set, size limits, policy compatibility and digest matches. The authorization receipt is bound to **one candidate generation**, not an abstract command or whole directory. One-use signing authorization is consumed by the publisher's protected state. Generation contents need not expire merely because the short-lived authorization expires; acceptance of a previously committed generation after key revocation follows an explicit compromise/quarantine policy.

## 5. Immutable generation protocol and durability

**Layout (conceptual, exact paths to be source-reviewed):**

```
trusted-data/
  CURRENT                     # small atomic pointer to generation ID + digest
  generations/<generation-id>/
    catalog.json
    digest-sidecar
    dependent-context
    MANIFEST.json             # exact file inventory and parent sequence
    AUTHORIZATION.json        # signed, one-use delegated authorization
    PUBLICATION_INTENT.json   # prepared evidence/outbox for reconciliation
  trusted-publisher-state/    # protected, not agent-writable
```

An immutable generation directory contains the complete set. Reader opens `CURRENT` once, safely opens only that generation, verifies the exact signed inventory and never crosses into another generation mid-read. Old generations are retained for concurrent readers and rollback-recovery policy; garbage collection is separately guarded. A simple two-file rename is **not** a multi-file atomic transaction.

**Publisher algorithm (specification, not executed):**

1. Authenticate publisher and signer; acquire exclusive publication lock; re-read protected current sequence/parent, policy epoch and revocation view. Reject stale proposals via compare-and-swap semantics.
2. Open trusted root by retained directory descriptor; validate ownership, mount boundaries, ancestor stability and exact path allowlist. On unsupported filesystem/permissions, refuse publication. Never trust producer-supplied absolute paths.
3. Create a unique staging directory on the **same verified filesystem**. Write each file with exclusive creation and safe fd-relative opens. Verify exact byte length, schema and SHA-256 while reading already-opened file descriptors; no mutable path reopen window.
4. `fsync` each generation data file, then the manifest and authorization receipt; `fsync` the staging directory. Signer authorization must be validated against exactly these bytes and the current revocation/policy epoch immediately before commit.
5. Atomically rename staging directory to `generations/<id>`; `fsync` the `generations/` directory and any directory whose entry changed. Do not expose a pointer until the complete generation is durable.
6. Create `CURRENT.tmp` with new ID/digest/sequence; `fsync(CURRENT.tmp)`; atomically rename over `CURRENT` within its directory; `fsync` the containing directory. Only this is the visibility commit point. Use documented recovery ordering and check all errors.
7. Re-open `CURRENT`, verify authorized generation and protected anti-rollback floor; persist a publication evidence/outbox record. Return `ACKNOWLEDGED_DURABLE` **only after** required file/directory sync and recovery/receipt obligations succeed. If pointer may have changed but ACK failed, return `COMMITTED_UNACKNOWLEDGED` and reconcile by authorization ID on retry; never blindly republish.
8. On restart, scan only the active pointer + protected journal/outbox, validate full chain/epoch/contents, reconcile interrupted intent and append the missing audit event. Never silently choose a newer but unauthorized generation or accept torn content.

A rename can make pointer replacement atomic for concurrent readers but does not by itself make it durable. Linux `fsync(2)` documents the need to sync directory entries separately. A claimed guarantee is bounded to an explicitly approved local filesystem, kernel, mount and storage profile whose flush/power-cut behavior has been tested. Do not assert arbitrary NFS/FUSE/removable-drive/cloud-volume power-loss guarantees, or that a successful VM kill equals hardware power failure. The trust model must state what happens when a disk lies about flushes, or when a privileged attacker controls storage.

## 6. Revocation, rollback and recovery

- Owner-root or threshold-signed **authority registry** lists delegated role keys, explicit bounds, policy versions, key expiry and an increasing `authority_epoch`.
- Revoke by publishing a newer authenticated registry/epoch removing the key or scope; signer's policy service stops authorizing immediately after learning the change. The publisher must obtain a fresh authoritative revocation view before commit and fail closed if freshness cannot be established. Readers must apply the specified compromised-generation quarantine policy rather than blindly invalidating all historical valid data.
- Protect the highest accepted generation sequence **and** authority epoch outside attacker-writable snapshots (TPM-backed monotonic state, isolated service or externally anchored transparency/checkpoint). A local `previous_generation_id` hash chain alone cannot stop complete filesystem rollback.
- Every authorization has expiry, nonce/idempotency ID, scoped principal and exact parent. Reject replay, duplicate commit, stale parent, cross-dataset and old-epoch signatures.
- Key compromise: revoke delegation, rotate to a new key and authority epoch, quarantine affected generations until reattested; root-key compromise invokes out-of-band recovery. Revocation cannot undo data an already compromised authorized signer successfully committed before revocation.
- Restore must verify the external anti-rollback anchor and signatures; a copied old `CURRENT` plus copied old registry must not count as trustworthy recovery.

## 7. Auditable receipts and authority evidence

**Authorization receipt** (signing service): authorization ID; workload principal + authenticated identity evidence reference; policy digest/version and delegation key ID; exact manifest digest; permitted dataset/operation/parent/sequence; epoch; timestamp/expiry; decision ALLOW or DENY; signature. Denial receipts should avoid leaking secrets.

**Publication receipt** (publisher): authorization receipt digest, expected/observed parent, committed generation ID, pointer digest, manifest digest, filesystem profile, durable write/sync step outcomes, commit state (`PREPARED`, `COMMITTED_UNACKNOWLEDGED`, `ACKNOWLEDGED_DURABLE`, `REJECTED`, or `RECOVERED`), execution identity, and separate publisher signature. A signed claim of fsync success is evidence of a reported operation, not independent power-failure proof.

Send immutable/event-hash-linked copies to a **separate append-only log or protected audit store**, with an exportable evidence bundle. Maintain a durable local outbox for crashes between pointer publication and external log acknowledgement, with recovery/reconciliation. Record the exact tested binary/source/toolchain and environment in test receipts. Local checksums cannot establish authentic reviewer identity or prove that a process truly executed.

## 8. Failure modes and prohibitions

- Agent can submit a correct hash pair but cannot authorize it.
- Machine worker cannot increase its own allowlist, rate, signer scope, authority epoch or policy version.
- Trust-root, signing and publication credentials must not be collocated within an unrestricted agent execution boundary.
- No automatic root-key delegation, policy rewrite, arbitrary rollback, broad wildcard paths, or self-issued signed approval receipt.
- If signer/revocation state is unreachable or stale, **deny new publication**; retaining an already verified existing generation is a separate availability policy.
- A forged or replayed signed receipt, modified generation, incorrect parent, wrong workspace, unknown signer, expired policy or multiple concurrent commits must fail closed.
- Stage/recovery paths do not bypass signed validation or directory bounds.

## 9. Evidence and acceptance plan

All results must be objective-bound and independently reviewed; none is yet claimed executed. The companion `SQL002_ACCEPTANCE_MATRIX.md` defines representative negative/positive and interruption tests. Required evidence includes source-level enforcement, signed policy fixture vectors, signature/canonicalization vectors, concurrency schedule, fault injection at every write/rename/fsync/receipt boundary, actual power-loss/restore tests on the exact supported storage profile, hostile credentials/replay/revocation tests, readback, and source-to-running-process linkage for BLD-005. Stop work without separate scoped execution approval.

**Gate state remains:** Gate 1 PENDING / NOT_READY_FOR_ACCEPTANCE; Gate 2 BLOCKED; Gate 3 NOT_AUTHORIZED. This is a future design and may show the frozen implementation fails desired guarantees; it cannot retrofit them into frozen source.

## 10. External design references (informative; no claim of compliance)

- Linux fsync(2): https://www.man7.org/linux/man-pages/man2/fsync.2.html
- Linux rename(2): https://man7.org/linux/man-pages/man2/rename.2.html
- The Update Framework delegated-role, revocation, rollback model: https://theupdateframework.github.io/specification/v1.0.36/
- RFC 8785 JCS canonical JSON: https://www.rfc-editor.org/rfc/rfc8785.html
- NIST SP 800-57 Part 1 key management: https://www.nist.gov/publications/recommendation-key-management-part-1-general-1
