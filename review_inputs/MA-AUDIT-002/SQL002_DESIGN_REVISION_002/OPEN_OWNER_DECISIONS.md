# Open owner decisions v0.2

This document records no owner selection. Confirmed constraints are power-loss durable acknowledgement, scoped trusted machine authorization and no agent signing/publication/policy/revocation authority. These are not reopened as optional controls. All companion contract selections remain null; durability is required but formal contract ratification is pending.

| Decision | Proposed design position, NOT ratified | Consequence / evidence needed before activation |
|---|---|---|
| O01 SQL002-CONSISTENCY | Coherent current generation admitted by ACS, complete older pinned read permitted; no universal freshness promise | Owner accepts online admission/final delivery, quarantine and in-flight delivery semantics; assess actual reader implementation |
| O02 SQL002-PUBLICATION | ACS atomic current/member/nonce/epoch/quota/event commit, not CURRENT rename | Material change from local file-only scope; owner accepts external coordinator TCB, availability and deployment costs; independent transaction/rollback proof |
| O03 SQL002-AUTHENTICATION | Root-delegated signer plus independent admitted job and exact generation manifest | Name allowed trusted principals, datasets, job recipes/content rules and trust domains; no default agent-mediated admission |
| O04 Durability recording | Required power-loss protection under qualified explicit profile | Formal ratification must not weaken required power-loss property; storage/controller/kernel/ACS/audit-store evidence required |
| O05 SCOPE-01..05 | Existing A/A/A/A/B proposals unchanged | This expanded future service architecture may exceed source-only assurance; independent scope review needed, no retrospective change to F |
| O06 ACS provider and trust anchor | One logical crash-durable serial transaction system, independent of data snapshots and agent/admin credentials | Choose authenticated provider/wire adapter, protected checkpoint/witness bootstrap and recovery model; do not replace with unsigned local DB backup |
| O07 Machine policies and budgets | Exact dataset/path/operation/schema, independently verified recipe/trigger, aggregate quota and expiry | Owner sets max bytes/rate/window/TTL/resource bounds; missing values deny; no unrestricted defaults |
| O08 Time and identity | Trusted bounded ACS clock, challenge-bound fresh sessions and audience-bound workload/job proof | Specify clock/skew/rollback alarm, identity provider/adapters and key-use isolation; no producer clock or identity claims |
| O09 Key provider/root recovery | Ed25519 v0.2 contract, separate admin/authorization/storage/receipt/log roles, out-of-band root-compromise recovery | Select capable provider, root custody/threshold/rotation channels; different algorithm requires reviewed contract version, no fallback |
| O10 Supported storage | PROFILE-LINUX-EXT4-LOCAL-PLP-01 qualification target only | Pin exact Linux 6.6 patch/config, ext4 mount, NVMe/controller/firmware/cache/PLP evidence and independent ACS/log profiles; nothing currently supported |
| O11 Quarantine/history | Compromise defaults to all generations signed by key since activation; normal rotation preserves eligible historical membership | Owner reviews availability/compromise interval and final-delivery ordering; committed history never erased; no automatic grandfathered precommit authorization |
| O12 Retention/capacity | No committed-generation, membership or required-evidence GC in v0.2 | Owner accepts storage cost/capacity denial; future finite retention needs separate contract/revision, not an undocumented cleanup job |
| O13 Receipt/ACK contract | Independent COMMIT inclusion receipt required for ACK_READY; ACK event export retryable without recursion | Choose independent append-only store/witness/retention/export wire protocol; accept that log outage withholds ACK after commit |
| O14 Privileged/administrative boundaries | Storage host root trusted for availability/flush assertions; independent ACS/root/log administration; restore through same coordinator | Name TCB, service identities, operator duties and excluded media/privileged compromise; prove no legacy consumer falsely asserts v0.2 assurance |
| O15 Independent review/evidence authority | Architect cannot certify own repair or approve SQL-002 | Authenticated independent reviewer and exact artifact-binding receipt needed; future implementation/tests require separate scoped authorization |

Unselected parameters do not prevent review of the conditional design; they block activation and evidence claims. If owner rejects online ACS, indefinite retention or delivery semantics, the transaction model must be revised explicitly rather than falling back to v0.1 assumptions. No actual authority receipt or policy is created here.

## Preserved mandatory inventory

All 34 remain essential/open: AGT-001, AGT-003, AGT-004, AGT-008, AGT-036, AGT-037, AGT-043; BLD-002, BLD-005; CORE-001, CORE-002, CORE-005, CORE-006, CORE-014, CORE-019, CORE-022, CORE-023, CORE-024; GNB-012, GNB-031, GNB-041, GNB-062, GNB-109; HSH-001, HSH-030, HSH-031; LED-001; PRC-004, PRC-005, PRC-010; SQL-002; TSP-001, TSP-004, TSP-009.

296 selected objectives remain inventoried; eleven SQL whole checks UNKNOWN; zero certified counterparts; SQL-004 generic constraints and SQL-005 explicit booleans retained. Frozen revision `4d654e8ee11c12e8d10ef27b90ea2e84cf43c55e`. Gate 1 PENDING / NOT_READY_FOR_ACCEPTANCE, Gate 2 BLOCKED, Gate 3 NOT_AUTHORIZED. Owner choice values remain unset.
