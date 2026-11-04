# Qwen — fixed Gen Type-Go PR #6 independent audit (AuditGo-backed)

This file is a **task instruction**, not evidence and not permission to mutate GitHub.

1. Read `AUDIT_POLICY.md`, `AGENTS.md`, `AUDITGO_SETUP.md`, and `profiles/gen-type-go-pr6/target.lock.json`. This workspace is an audit harness, not the target repository. The exact target is `es-3581100/Gen_Type-go` PR #6 at `0af4707d66bb30aef09fe80316c408568e256e53`, reconciled against `main` at `7b071c7c0ed5037b107f4ae9a263f3c06853c46f`.
2. Shell preflight: run `echo QWEN_SHELL_OK`, capture stdout and exit code. If shell execution is unavailable, stop as `AUDIT_INCOMPLETE`. Run `sh scripts/setup-auditgo.sh` and `sh scripts/auditgo-smoke.sh` only in an isolated environment and verify output. The bootstrap is on the separate `feat/auditgo-qwen-bootstrap-v0.1.3` source branch and must have been incorporated into this workspace.
3. Run `bash scripts/replay-gen-type-go-pr6.sh` in a disposable Qwen workspace. The script fails closed on source SHA, tree, parent, manifest identity, Go/CI code preservation, PR #5 governance import, registry identity, and test environment; it writes evidence under ignored `audit-artifacts/`. It does **not** claim independent audit completion.
4. Study the generated `auditgo-mission.md`, `auditgo-plan.json`, and `auditgo-unassessed-return.json` **and all raw logs**. Investigate the entire audit selection, not just canned commands. Mark all relevant items with actual evidence; do not mislabel generated seed records as passing or call the AuditGo validator proof of semantic truth.
5. Adversarially test strict JSON rejection, escaped duplicate keys, invalid UTF-8, type/null/schema errors, cycles and edges, map iteration and permutation determinism, graph-root/child/ancestor/descendant order, malformed input and CLI exit codes, 2 MiB input/4 MiB output limits, short output writes, and untrusted IDs/claims. Existing `phasea.DecodeStrict` remains the only graph-validity oracle. A bad diagnostic hint is not necessarily an accepted-invalid defect.
6. Reproduce actual evidence in unmodified source or clearly label a temporary instrumented copy. Keep generated fuzz corpora outside upstream. Inspect 14-file PR diff against reconciled `main`, and separately the seven-file reconciliation delta against original `f8c6ba94...`.
7. Produce `QWEN_PR6_INDEPENDENT_AUDIT.md`, completed AuditGo machine return (not the unassessed seed), and `QWEN_PR6_REPRO_STEPS.md`. Run `auditgo validate` with `-evidence-root` after populating real evidence and record validator output and report hash. If full coverage isn't obtained, return `AUDIT_INCOMPLETE`, with precise blocked items.
8. Report actual observed GitHub identity and permissions. Do **not** automatically comment to the PR, submit `APPROVE`/`REQUEST_CHANGES`, merge, release, alter source, or promote authority. The authority gate is `NEGATIVE_DISPOSITION_ORACLE_OPEN__RATIFICATION_PENDING__PHASE1_BLOCKED`.

### Fixed audit decisions

- AuditGo selected inventory is generated from exact lock `target.lock.json`. False selectors are not silently enabled.
- `AUDIT_PASS`: fully supported bounded scope and no confirmed blocker.
- `AUDIT_FAIL`: reproduced blocker with source/line, input, command, actual exit, expected-vs-observed, and evidence.
- `AUDIT_INCOMPLETE`: missing essential execution, provenance, or coverage.
- An AuditGo seed, successful CI or green replay script is **not** an approval.

Final answer should show source identity, archive/checksum coverage, Go/test/race/vet/build/fuzz execution counts, unexpected source changes, findings, unreached tests, actual artifact locations and hashes, and the independent disposition. Do not invent results or report that Qwen executed the script until it has actually run.
