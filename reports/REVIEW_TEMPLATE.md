# Independent audit report — <repo> PR #<number>

- Reviewer identity: <actual GitHub username/bot>
- Target SHA: <40-hex>
- Compared against base SHA: <40-hex>
- Date/time UTC: <ISO 8601>
- Toolchain and sandbox: <facts>
- Scope: <all PR files or explicit subset>

## Evidence ledger

| ID | Check | Class (EXECUTED / OBSERVED / REPORTED / UNVERIFIED) | Result | Evidence link / log |
|---|---|---|---|---|
| E01 | Exact-head and clean checkout |  |  |  |
| E02 | Diff and changed-file inventory |  |  |  |
| E03 | `go test -count=1 ./...` |  |  |  |
| E04 | `go test -race -count=1 ./...` |  |  |  |
| E05 | `go vet ./...` |  |  |  |
| E06 | Malformed XML and rollback adversarial cases |  |  |  |
| E07 | Merkle / SQLite / FTS5 invariant review |  |  |  |

## Findings

For each finding: ID, severity, file/line, expected versus actual behavior, reproduction inputs, proof/log, affected contracts, and proposed diagnostic direction. Distinguish a real defect from a hypothesis.

## Coverage and exclusions

- Independently executed:
- Observed only through GitHub:
- Not evaluated / unsupported:
- Deferred contract areas:

## Disposition

`AUDIT_PASS` / `AUDIT_FAIL` / `AUDIT_INCOMPLETE`

Rationale, unresolved questions, and next review gate. **This technical assessment is not an independent human approval or permission to merge/release.**
