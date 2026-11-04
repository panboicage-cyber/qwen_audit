# Qwen audit recovery — 2026-11-03

## What failed

The Qwen log showed that the Shell tool produced no verifiable output even for `echo hello`. The environment was cloned from `panboicage-cyber/qwen_audit` (`main` SHA `2e9d23dc27bb37abab60b80e12f2a2530b7106f9`), which currently contains only `qwen_audit_starter.zip`. This is **not** the target code repository. The setup command `test -f go.mod` and the SHA check for `c64502...` would fail when run in that checkout.

## Recovery sequence

1. Start a fresh Qwen coding session, or reset its shell tool. Run `echo QWEN_SHELL_OK` and demand actual stdout with exit status 0. If this fails, stop and return `AUDIT_INCOMPLETE`. Repeated Ctrl-C/Ctrl-D attempts and file-create probes are not substitutes for execution.
2. Ensure this repository's ZIP is extracted, or upload the unpacked scaffold files to `main`. ZIP contents are not inspectable as ordinary text by Qwen's Edit tool.
3. From a working shell run `bash scripts/setup-qwen-target.sh` (a copy of this scaffold must be present in the checkout). It clones **kol_merkle** into `/tmp/kol-merkle-pr4-target` and verifies its pinned SHA; it does not confuse the audit workspace SHA with the target SHA.
4. In that target checkout run `go test -count=1 ./...`, `go test -race -count=1 ./...`, and `go vet ./...`. Review and adversarially test actual Go parser, SQL transaction, FTS5 and Merkle code. Python XML parser tests alone do not test the Go implementation.
5. Confirm the Qwen GitHub App is installed with the required PR-comment permissions on **donCannoli-burns/kol_merkle**, not just `panboicage-cyber/qwen_audit`. If it cannot post, return a completed comment draft for manual publication. Do not falsely claim publication or use a second account to simulate independent review.

The code under audit must remain unchanged. No `APPROVE`, merge, release, or deployment is authorized.
