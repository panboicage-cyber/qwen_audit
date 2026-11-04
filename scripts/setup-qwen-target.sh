#!/usr/bin/env bash
# Run only after an actual test command confirms the Qwen shell is functional.
# This script prepares a SEPARATE, detached checkout of the audited project.
set -Eeuo pipefail

TARGET_REPO='https://github.com/donCannoli-burns/kol_merkle.git'
TARGET_SHA='c64502cb5db3efd9316991ce1c0756d97729a6c3'
TARGET_DIR="${QWEN_AUDIT_TARGET_DIR:-/tmp/kol-merkle-pr4-target}"

printf '%s\n' 'QWEN_SHELL_PREFLIGHT_OK'
command -v git >/dev/null || { echo 'ERROR: git unavailable'; exit 1; }
command -v go >/dev/null || { echo 'ERROR: Go unavailable'; exit 1; }
command -v sha256sum >/dev/null || { echo 'ERROR: sha256sum unavailable'; exit 1; }
git --version
go version

if [ -e "$TARGET_DIR" ] && [ ! -d "$TARGET_DIR/.git" ]; then
  echo 'ERROR: target path exists but is not a Git checkout'
  exit 1
fi
if [ ! -d "$TARGET_DIR/.git" ]; then
  git clone --no-checkout "$TARGET_REPO" "$TARGET_DIR"
else
  CURRENT_REMOTE="$(git -C "$TARGET_DIR" remote get-url origin)"
  [ "$CURRENT_REMOTE" = "$TARGET_REPO" ] || {
    echo "ERROR: unexpected target remote: $CURRENT_REMOTE"
    exit 1
  }
  git -C "$TARGET_DIR" fetch origin
fi

git -C "$TARGET_DIR" checkout --detach "$TARGET_SHA"
ACTUAL_SHA="$(git -C "$TARGET_DIR" rev-parse HEAD)"
[ "$ACTUAL_SHA" = "$TARGET_SHA" ] || {
  echo "ERROR: target SHA mismatch, got $ACTUAL_SHA"
  exit 1
}
[ -f "$TARGET_DIR/go.mod" ] || {
  echo 'ERROR: target go.mod missing'
  exit 1
}

(cd "$TARGET_DIR" && go mod download && go mod verify)
printf 'QWEN_TARGET_READY %s %s\n' "$ACTUAL_SHA" "$TARGET_DIR"
printf '%s\n' 'Do not treat setup success as an audit PASS. Run and record all tests separately.'
