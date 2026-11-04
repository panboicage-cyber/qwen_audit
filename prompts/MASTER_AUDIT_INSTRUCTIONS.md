# Qwen — master audit workspace startup and handoff

**First** run `echo QWEN_SHELL_OK`, then `sh scripts/setup-master-audit-workspace.sh` and record exact exit status/stdout. Stop on any failure. These commands do not audit an upstream repository.

**Read-only preparation** (no target execution):

```sh
go run ./tools/masteraudit/main.go list
sh scripts/stage-audit.sh gen-type-go-pr6
cat audit-artifacts/staging/gen-type-go-pr6/context.json
go run ./tools/repomemory/main.go context
```

Read `AUDIT_POLICY.md` and `profiles/gen-type-go-pr6/target.lock.json`; the Gen Type-Go PR6 HEAD is exactly `0af4707d66bb30aef09fe80316c408568e256e53`. The historical `kol-merkle-pr4` queue entry is inactive and not interchangeable. Treat the staged context as `NOT_EXECUTED` and `source_verified=false` until independently established.

Only after separate explicit instruction, run the staged `scripts/replay-gen-type-go-pr6.sh` in an isolated sandbox with no credentials. Preserve logs and hashes. Never substitute a branch tip for a pinned commit or treat replay success as independent audit completion.

When actual source evidence exists, produce an evidence-return JSON, a short natural-language report, and SHA-256 sidecars. Distinguish reviewed evidence from inherited claims. A green setup/CI, selected AuditGo checks, a generated seed, or a checked hash does not imply owner approval, source correctness, or permission to publish or merge.
