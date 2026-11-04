#!/usr/bin/env bash
# Execute ONLY in an isolated sandbox with no production credentials.
set -euo pipefail

TARGET_URL='https://github.com/donCannoli-burns/kol_merkle.git'
TARGET_SHA='c64502cb5db3efd9316991ce1c0756d97729a6c3'
OUTPUT_DIR="${AUDIT_OUTPUT_DIR:-${PWD}/audit-artifacts}"
mkdir -p "$OUTPUT_DIR"
OUTPUT_DIR="$(cd "$OUTPUT_DIR" && pwd)"
LOG="$OUTPUT_DIR/kol-merkle-pr4-replay.log"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

export GOTOOLCHAIN=local
export GOWORK=off
export GOFLAGS=-mod=readonly
export GIT_TERMINAL_PROMPT=0

run() { echo "> $*" | tee -a "$LOG"; "$@" 2>&1 | tee -a "$LOG"; }
: > "$LOG"
run go version
run git clone --no-checkout "$TARGET_URL" "$WORK/target"
run git -C "$WORK/target" checkout --detach "$TARGET_SHA"
ACTUAL_SHA="$(git -C "$WORK/target" rev-parse HEAD)"
[[ "$ACTUAL_SHA" == "$TARGET_SHA" ]] || { echo 'FATAL: target SHA mismatch' | tee -a "$LOG"; exit 1; }
echo "Verified target SHA: $ACTUAL_SHA" | tee -a "$LOG"
cd "$WORK/target"
run go mod download
run go mod verify
run go test -count=1 ./...
run go test -race -count=1 ./...
run go vet ./...
run git diff --exit-code
[[ -z "$(git status --porcelain)" ]] || { echo 'FATAL: target checkout dirty' | tee -a "$LOG"; exit 1; }
echo 'All requested replay commands completed successfully.' | tee -a "$LOG"
sha256sum "$LOG" | tee "$OUTPUT_DIR/kol-merkle-pr4-replay.sha256"
