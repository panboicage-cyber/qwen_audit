# Root-cause reconciliation v0.2

Architect-authored diagnostic synthesis, not independent acceptance. Inputs were read in this conversation and freshly verified before revision. Original files are preserved. Request SHA-256: `e6e4923f440e69416e5dc0f1e5a195e3107b5a0ca83bde52508c5410a3eb31ab`.

## Verified inputs

Paths below are relative to the parent evidence directory. Review manifest verification: 3/3 OK. Original package SHA256SUMS: 4/4 OK; shalist: 5/5 OK. Documentary checks are not target/fixture execution.

| Input | SHA-256 |
|---|---|
| SQL002_ADVERSARIAL_REVIEW_001.md | bbf2648e11d76f89ff23815bbeedbb324cbd170d52f95f49795d2128ff5cd9fd |
| SQL002_FAILURE_MODE_MATRIX_001.md | d130ddfcf0a9c237267ccf04e18facea5913f4f837126321a1a445bc540ee5f8 |
| SQL002_REVISION_RECOMMENDATIONS_001.md | 955c410e3732bc4f4e3d9b17cec1937f837fd4261a15101b7977aa6bf0777eed |
| MA_AUDIT_002_SQL002_SCOPED_MACHINE_AUTHORITY_V01/SQL002_SCOPED_MACHINE_AUTHORITY_DESIGN.md | d9dc09e40e10fe010d1182ef9adb06a6bcf6f6ca83c647b7c922edc0a6cfe606 |
| MA_AUDIT_002_SQL002_SCOPED_MACHINE_AUTHORITY_V01/SQL002_MACHINE_AUTHORITY_CONTRACT.json | eab761ced95ead1b673822f5d856d368ed20356aa38b5aa406d9e0257c56f66e |
| MA_AUDIT_002_SQL002_SCOPED_MACHINE_AUTHORITY_V01/SQL002_ACCEPTANCE_MATRIX.md | 1bb60c464bf8d7df8281417e2ef3a710628c7212206d604371e323cb66f91721 |
| MA_AUDIT_002_SQL002_SCOPED_MACHINE_AUTHORITY_V01/CRUSH_DESIGN_REVIEW_HANDOFF.md | 242224a8637dfdaa513f6c531e7e91b3d58a63280dc0098711f9823073cb9408 |
| MA_AUDIT_002_SQL002_SCOPED_MACHINE_AUTHORITY_V01/SHA256SUMS.txt | 6e2dd26716944f1ed98ec2daad81d7967e44d9fb7e42a29510e2ddbdb598afae |
| MA_AUDIT_002_SQL002_SCOPED_MACHINE_AUTHORITY_V01/shalist-0001.txt | 0b005ce21d92d9832b897eab40d7bc95b94df5454fadfeafa2dcbab7cfd89965 |

## Shared invariants, not isolated patches

1. **One authority for irreversible history (F01/F02/F03/F08/F10).** A filesystem pointer, external anchor, nonce database and audit outbox cannot be declared one atomic commit. Moving authoritative current, epoch, membership, consumed IDs, quota charge and receipt obligation into one ACS transaction removes the undefined cross-resource atomicity claim. Durable content is a prerequisite with explicit protected retention, while local pointer and external log delivery are recoverable projections. Revocation is not a cache refresh problem: it must share commit order.
2. **Authorization is a decision, not possession of a signer endpoint (F04/F05/F09).** Principal authentication plus path/size limits can still launder agent choice through a trusted worker. Independent admitted-job policy, exact immutable candidate identity, distinct typed key uses and administrative separation are required together. Root provisioning and recovery cannot be delegated to the same normal signer.
3. **Bytes and history have different lifetimes (F06/F07/F08).** Safe opens do not freeze bytes; signatures do not prove committed membership; old coherent readers do not authorize snapshot rollback. Sealed publisher-owned files plus ACS reader pins, final quarantine checks and no committed-data GC separate these concerns. Indefinite retention deliberately exposes capacity costs instead of inventing an owner retention horizon.
4. **Evidence is a separate claim domain (F02/F10/F11).** Signed receipt payloads can faithfully report an unverified assertion. Durable PREPARED/COMMIT state prevents lost obligations but does not prove flush honesty or runtime source identity. Independent log delivery, externally observed power-loss outcomes and BLD-005 evidence remain distinct.

## Finding positions and preserved diagnostic evidence

D/C/A denote original architecture/contract/acceptance matrix in the original package. Exact original findings remain in the untouched report; this table preserves ID, severity and source references. No finding is EVIDENCE_PROVEN or OWNER_RATIFIED.

| ID | Severity | Original source references | Architect position / linked correction |
|---|---|---|---|
| F01 | HIGH | D:58,81-88,96-99; C:152,160-179,184; A:11-15,20-21 | Agree: unified ACS commit and content preparation retention, I1-I3 |
| F02 | HIGH | D:64-74,83-88,103-107; C:160-171,181-187; A:12,24 | Agree: precommit intent and atomic COMMIT outbox, I7 |
| F03 | HIGH | D:43,81-86,94-98,115; C:155-156,173-179; A:19 | Agree: serial revocation/commit, I4 |
| F04 | HIGH | D:23-43,103,111-114; C:60-88,103-118,147-158; A:16-18 | Agree: independent job admission and no agent credentials, I5 |
| F05 | HIGH | D:35-58,64-74; C:103-145; A:17,22 | Agree: exact signed envelope fields, decimal strings and acyclic hash graph |
| F06 | HIGH | D:37,58,77,82-85; C:162-166; A:7-9 | Agree: sealed own inodes, no aliases, pins/no committed GC, I8 |
| F07 | HIGH | D:12,15,27,58,77,94-99,115; C:26,155-179; A:7,19-21 | Agree: new admission vs pinned read vs final quarantine distinction, I6 |
| F08 | HIGH | D:81-90; C:160-171,189-195; A:10,12-15 | Agree: explicit uncertain persistence, qualified platform and ACK conditions, I9-I10 |
| F09 | HIGH | D:28,31,41-43,82,94-99,113-117; C:50-58,83-94,154,173-179,212-225; A:9,15-21 | Agree: TCB/admin/restore paths must obey ACS; no local bypass |
| F10 | MEDIUM | C:83-88,181-187; D:26,103-107; A:23-24 | Agree: publication evidence key distinct from authorization key; independent log |
| F11 | MEDIUM | A:7-25,27; D:86-88,121; C:156,176,179,190,212-225 | Agree: admissible state sets, independent oracle and bounded claims |

## Adjacent failures found while reconstructing

ACS availability is now essential; do not silently restore offline CURRENT trust. Prepared data cannot be reclaimed on an elapsed lease while ACS commit outcome is unknown. A committed generation may be quarantined between commit and ACK; preserve commit fact and withhold ACK, never label it aborted. Signed authorization and receipt files must be outside their own inventory. ACK evidence cannot require recursively acknowledging its own remote ACK event. Multiple datasets sharing one delegation need atomic aggregate quotas. A schema-valid catalog/context must satisfy independent derivation semantics, not just signed hash agreement. Host-root access defeats same-host ACL separation; keep external coordinator/root/audit administration in separate failure domains. These constraints are integrated into the architecture and transitions.

## Unresolved interpretation, not suppressed disagreement

No substantive original finding is rejected. Two demands require qualified interpretation: no algorithm can give a partitioned client certain commit knowledge, so durable ACS state is binary while caller knowledge may be OUTCOME_UNAVAILABLE; and revocation cannot recall bytes already authorized for in-flight delivery. Both are explicit proposed semantics requiring owner and independent review. The design can be reviewed without selecting provider/TTL/role parameters, but cannot be implemented or activated with them unset. A Linux/ext4 profile name is a qualification target, not tested support. Provider identity/policy/log adapters still require normative wire contracts and evidence; no automatic trust of arbitrary receipts.

Frozen revision `4d654e8ee11c12e8d10ef27b90ea2e84cf43c55e`, 296 selected objectives, all 34 mandatory essential/open and all gates are preserved. This synthesis changes documentation only.
