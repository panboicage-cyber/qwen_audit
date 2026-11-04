#!/bin/sh
# One Qwen setup command: deterministic bootstrap only, NO target checkout/tests.
set -eu
umask 077
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
for exe in go sha256sum; do
  command -v "$exe" >/dev/null 2>&1 || { echo "MASTER_SETUP_MISSING_TOOL $exe" >&2; exit 2; }
done
cd "$ROOT"
sha256sum -c OVERLAY_SHA256SUMS.txt >/dev/null || { echo 'MASTER_SETUP_OVERLAY_INTEGRITY_BLOCKED' >&2; exit 3; }
sh scripts/setup-auditgo.sh
sh scripts/auditgo-smoke.sh
sh scripts/repo-memory-smoke.sh
go run ./tools/masteraudit/main.go -root "$ROOT" verify
go run ./tools/masteraudit/main.go -root "$ROOT" list
printf '%s\n' 'MASTER_AUDIT_WORKSPACE_READY (no target audit executed; no memory written)'
