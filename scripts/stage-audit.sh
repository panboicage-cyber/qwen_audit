#!/bin/sh
# Derive a review handoff into ignored artifact storage; never executes target scripts.
set -eu
umask 077
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
PROFILE=${1:-}
case "$PROFILE" in
  gen-type-go-pr6) ;;
  *) echo 'STAGING_BLOCKED: unsupported, disabled or unspecified profile' >&2; exit 2 ;;
esac
command -v go >/dev/null 2>&1 || { echo 'STAGING_BLOCKED: go unavailable' >&2; exit 2; }
command -v sha256sum >/dev/null 2>&1 || { echo 'STAGING_BLOCKED: sha256sum unavailable' >&2; exit 2; }
cd "$ROOT"
# Verify the entire frozen baseline before writing derived data.
go run ./tools/masteraudit/main.go -root "$ROOT" verify
STAGE="$ROOT/audit-artifacts/staging/$PROFILE"
mkdir -p "$STAGE"
# The output is NOT a repo memory record or an audit return.
TMP=$(mktemp "$STAGE/context.XXXXXXXX")
trap 'rm -f "$TMP"' EXIT HUP INT TERM
go run ./tools/masteraudit/main.go -root "$ROOT" -profile "$PROFILE" context > "$TMP"
mv -f -- "$TMP" "$STAGE/context.json"
(cd "$STAGE" && sha256sum context.json > CONTEXT_SHA256SUMS.txt)
printf 'MASTER_AUDIT_STAGE_READY (NOT_EXECUTED) %s\n' "$STAGE/context.json"
