# Crash/recovery and authority-cut matrix v0.2

Specification only. No fault injection, model execution, fixture or power-cut tests are authorized or executed. T01-T13 refer to TRANSACTION_STATE_MACHINE.json. Actual membership is COMMITTED iff the durable ACS commit exists; an unreachable ACS produces OUTCOME_UNAVAILABLE, never a guessed abort or a second commit.

## Cuts at every durable transition

Each row requires later independent testing immediately before and immediately after the named durability boundary, plus error return and lost-response schedules. Repeating an operation must use the same immutable transaction digest/ID.

| Cut | Persistent state allowed | Recovery / expected outcome | Independent oracle and remaining prerequisite |
|---|---|---|---|
| T01 admission transaction | Absent or complete admitted job/event | Query ID; never accept producer assertions as admission | ACS history and authenticated recipe/trigger; real admission provider not selected |
| T02 signer durable record / envelope response | No usable authorization or same stored signed envelope | Return same payload/signature identity; conflicting ID rejects | Independent canonical-byte/signature vectors; no keys exist |
| T03 each payload write, close, fsync | Partial or durable unreachable files | No ACS commit; validate or discard only after authoritative cleanup decision | Exact bytes/length and syscall evidence; no reader path into staging |
| T03 manifest fsync | Payloads may be durable, metadata incomplete | No preparation registration until complete metadata verified | Acyclic inventory oracle; missing/extra artifact rejects |
| T03 authorization fsync | Manifest durable, authorization incomplete | No preparation registration; no local hash-only acceptance | Exact envelope and manifest digest oracle |
| T03 PREPARED intent fsync | Intent missing/complete independently of other staging files | Only complete verified/synced set qualifies | Intent binds authorization/generation/profile/fence; unknown observations stay unknown |
| T03 stage-directory fsync | Entries may disappear until synced | No visibility or ACK; reconstruct complete set | Directory durability on qualified ext4 profile required |
| T03 no-overwrite directory rename | Old staging or complete installed directory | Verify exact same generation/transaction; never replace collision | Duplicate generation/name, mount and alias negatives |
| T03 fsync destination parent then source parent/created ancestors | Namespace durability partly uncertain | Mark unhealthy on error; no T04; retry only with complete revalidation | Include EIO, ENOSPC, delayed writeback and uncertain close errors |
| T04 ACS PREPARED commit / response | Registration absent or wholly durable | Query exact transaction; do not GC on timeout; protect objects | Coordinator history and storage-retention reservation; no expiring deletion lease |
| T05 before ACS transaction | Prepared data only | Parent/epoch/expiry/IDs revalidated; no reserved authority | Concurrent revoke and same-parent schedules |
| T05 during ACS atomic durable transaction | PREPARED or complete COMMITTED with all I2 fields | Full tuple/nonce/quota/event agree; no partially consumed authorization | Provider crash atomicity is required, not established by this design |
| T05 after commit, before response | COMMITTED and durable outbox, client uncertain | Query same ID; report COMMITTED_UNACKNOWLEDGED; never new-ID retry | External ACK/commit observation independent of failed host |
| T06 CURRENT temp fsync | Old cache or temp present | Readers query ACS and use committed generation; discard/reconcile temp | CURRENT cannot admit data without ACS |
| T06 CURRENT rename before parent fsync | New cache visible but possibly not durable | No logical commit change; sync error quarantines local service and withholds ACK | Independent reader sees only ACS selected immutable generation, or fails |
| T06 cache parent fsync / evidence recording | Durable cache, reconciliation evidence possibly absent | Re-read/reconcile; cache lag or missing event cannot erase commit | Verify current or later committed successor; no stale cache downgrade |
| T07 evidence signing and fixed-envelope durability before log append | COMMIT payload durable; canonical envelope absent or durably fixed | If absent, sign and durably fix before export; if fixed, retransmit exact bytes across restart/key rotation, never re-sign | Unsigned payload chain digest distinct from fixed envelope digest; ACK/log bind same envelope |
| T07 log durable append before receipt response | Log entry exists but exporter uncertain | Resubmit event ID; same digest returns same inclusion proof | Protected log/witness non-equivocation evidence, not hash chain alone |
| T07 ACS delivery-reference transaction | Log proof may exist before local reference | Re-obtain proof and record idempotently | No ACK without authenticated independent COMMIT inclusion |
| T08 ACK_READY durable transaction | COMMITTED or ACK_READY with proof references | Issue/retrieve same ACK identity only after prerequisites; no claim client got token | ACK capture outside failed machine; verify quarantine at decision |
| T08 response or later ACK-event export | ACK_READY durable; client or remote ACK event may be missing | Token retry; outbox eventually exports ACK event without recursive ACK requirement | Durable ACK fact and actual receipt by client are distinct |
| T09 policy/revoke/quarantine atomic update | Entire old or new epoch/state | T05/T10/T11 ordered consistently; no mixed registry | Owner authenticated update, protected epoch and provider ordering |
| T10 reader pin creation / open | Pin absent/present, data never GCed if committed | Obtain fresh admission on restart; open exact digest | Reader-vs-publication schedules; no old signed but uncommitted candidate |
| T11 final delivery authorization then release | Permission absent/recorded, pin may remain after crash | Revalidate for a new delivery; idempotent release | Revoke ordered before check denies; later revoke cannot recall in-flight bytes |
| T12 ABORTED commit then uncommitted deletion fsync | Abort recorded; orphan may remain | Commit forbidden; delete only exact grant and never committed objects | Abort-vs-commit race and same-ID reuse rejection |
| T13 restored files sync / reconciliation | Committed fact fixed; local data partial or exact restored | Fail closed until full seal/verification; resume projection/log steps | Restore exact authenticated bytes; recovery events preserve historical unknowns |

## Authority and concurrency schedules

| Schedule | Expected result |
|---|---|
| Two valid different candidates share one current parent | At most one T05 commits; loser receives stale-parent rejection. No shared ID can bind two payloads. |
| Same authorization sent concurrently to multiple publishers | ACS returns one immutable result; no duplicate quota charge or second membership. Stale storage fence rejects. |
| Revoke before T02 or between T02/T04/T05 | Signer denies when aware; regardless, ACS T05 denies if revoke serialized first. Signature validity is insufficient. |
| Revoke concurrent with T05 | ACS order defines winner. Earlier commit remains committed, subject to quarantine; later commit with old epoch rejects. |
| Policy/root rotation concurrent with preparation or restart | Prepared bytes retain no authority reservation; current version/epoch checked at T05. Old snapshot cannot reinstall trust. |
| Expiry during slow staging / clock step backward | T05 enforces current trusted clock; expired or unhealthy clock rejects, never producer clock fallback. |
| Delegation quotas across keys/workers/datasets | ACS atomic shared bucket prevents overspend; exact retry is not charged twice. |
| Agent uses workload proxy with valid schema/digests | No independently admitted job/trigger means DENY; path/size scope alone cannot authorize. |
| Alias mutation after verification | Enforced publisher-owned sealing prevents untrusted writer; if integrity changes detected, deny/quarantine, never mixed successful delivery. |
| New publication while reader pinned | Reader may deliver complete old pinned member after final eligibility check; no file crossing or arbitrary stale new admission. |
| Quarantine between reader pin and final check | Deny delivery; after final authorization, in-flight delivery is an explicit bounded irreversibility limit. |
| ACS outage or restore below protected checkpoint | No new publication/admission/delivery/ACK, no local CURRENT fallback; preserve data, query after trusted recovery. |
| GC/restore/admin attempts to bypass coordinator | No committed GC in v0.2; prepared cleanup needs terminal abort and exact grant; new content requires ordinary transaction. |
| Remote log outage after commit | Commit/read eligibility can persist; ACK withheld, durable outbox retries, no fake log receipt. |

## Failure classes and evidence boundaries

Process interruption may leave OS cache intact; it cannot establish power-loss durability. Kernel crash loses unsynced memory; completed durable boundaries only are assumed. Machine power loss requires pinned qualified hardware/cache/flush and ACS/audit-store profiles. Media destruction, dishonest flush, privileged compromise or correlated loss outside provisioned storage guarantees cause explicit outage/integrity incident, not silent rollback or a claimed pass. The proposed Linux 6.6/ext4/local PLP profile has not been qualified. Exact patch, device/controller firmware, mount options, clock, ACS provider and independent log are open owner/evidence decisions.

Proof plan requires an oracle outside the failed host to record durable ACK tokens and recovered ACS membership, plus independent verification of exact manifest, payload, authorization, prepared evidence and event chain. Signed service assertions cannot be their own oracle. The state-machine induction in architecture section 9 is conditional reasoning only; no model checker or target run occurred.

Frozen `4d654e8ee11c12e8d10ef27b90ea2e84cf43c55e` and all 34 mandatory objectives remain essential/open. Gate states unchanged.
