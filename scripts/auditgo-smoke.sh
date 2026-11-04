#!/bin/sh
# Read-only AuditGo orchestration smoke test. Outputs only under a disposable temp dir.
set -eu
umask 077
bin=${AUDITGO_BIN:-"$HOME/.local/bin/auditgo"}
[ -x "$bin" ] || { echo 'AUDITGO_SMOKE_MISSING_BINARY' >&2; exit 2; }
work=$(mktemp -d)
trap 'rm -rf -- "$work"' EXIT HUP INT TERM
"$bin" doctor
"$bin" init -out "$work/audit.json" -target qwen-auditgo-smoke -claim 'The Go implementation preserves its contract' -type IMPLEMENTATION -module GO_NATIVE_BINARY
"$bin" plan -config "$work/audit.json" > "$work/plan.json"
"$bin" render -config "$work/audit.json" -out "$work/mission.md"
[ -s "$work/plan.json" ] && [ -s "$work/mission.md" ] || { echo 'AUDITGO_SMOKE_OUTPUT_MISSING' >&2; exit 3; }
grep -q '"source_revision": "NOT_PROVIDED"' "$work/plan.json" || { echo 'AUDITGO_SMOKE_UNEXPECTED_PLAN' >&2; exit 3; }
printf 'AUDITGO_SMOKE_PASS (planning only; no repository audit executed)\n'
