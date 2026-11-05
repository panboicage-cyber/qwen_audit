#!/bin/sh
# SQL-002 v0.2 design evidence admission only. No target execution, Git writes, or review promotion.
set -eu
umask 077
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
BUNDLE="$ROOT/review_inputs/MA-AUDIT-002/SQL002_DESIGN_REVISION_002"
EXPECTED=1252a4debd76cf7824a9d1cb17d5503feb090ea162e025d11d73ef68dcbf8529
command -v sha256sum >/dev/null 2>&1 || { echo 'SQL002_BOOTSTRAP_BLOCKED: sha256sum missing' >&2; exit 2; }
[ -d "$BUNDLE" ] || { echo 'SQL002_BOOTSTRAP_BLOCKED: design directory missing' >&2; exit 2; }
[ -f "$BUNDLE/SHA256SUMS.txt" ] && [ ! -L "$BUNDLE/SHA256SUMS.txt" ] || { echo 'SQL002_BOOTSTRAP_BLOCKED: manifest missing or symlink' >&2; exit 3; }
ACTUAL=$(sha256sum "$BUNDLE/SHA256SUMS.txt")
ACTUAL=${ACTUAL%% *}
[ "$ACTUAL" = "$EXPECTED" ] || { echo "SQL002_BOOTSTRAP_BLOCKED: manifest mismatch $ACTUAL" >&2; exit 3; }
# The manifest is pinned above; prohibit missing, symlinked, or unexpected artifact names.
for file in AUTHORIZATION_CONTRACT_V02.json CRASH_RECOVERY_MATRIX.md INDEPENDENT_REVIEW_HANDOFF.md OPEN_OWNER_DECISIONS.md REVIEW_FINDING_CLOSURE_MATRIX.md ROOT_CAUSE_RECONCILIATION.md SQL002_ARCHITECTURE_V02.md TRANSACTION_STATE_MACHINE.json; do
  [ -f "$BUNDLE/$file" ] && [ ! -L "$BUNDLE/$file" ] || { echo "SQL002_BOOTSTRAP_BLOCKED: missing/symlink $file" >&2; exit 3; }
done
( cd "$BUNDLE" && sha256sum -c SHA256SUMS.txt >/dev/null ) || { echo 'SQL002_BOOTSTRAP_BLOCKED: artifact digest failure' >&2; exit 3; }
# JSON parsing only when Python is already available; exact pinned hashes provide identity.
if command -v python3 >/dev/null 2>&1; then
  python3 -m json.tool "$BUNDLE/AUTHORIZATION_CONTRACT_V02.json" >/dev/null || { echo 'SQL002_BOOTSTRAP_BLOCKED: authorization JSON invalid' >&2; exit 3; }
  python3 -m json.tool "$BUNDLE/TRANSACTION_STATE_MACHINE.json" >/dev/null || { echo 'SQL002_BOOTSTRAP_BLOCKED: state machine JSON invalid' >&2; exit 3; }
fi
BUNDLE_V01="$ROOT/review_inputs/MA-AUDIT-002/SQL002_MACHINE_AUTHORITY_V01"
[ -f "$BUNDLE_V01/SHA256SUMS.txt" ] && [ ! -L "$BUNDLE_V01/SHA256SUMS.txt" ] || { echo 'SQL002_BOOTSTRAP_BLOCKED: v0.1 manifest missing/symlink' >&2; exit 3; }
ACTUAL_V01=$(sha256sum "$BUNDLE_V01/SHA256SUMS.txt")
ACTUAL_V01=${ACTUAL_V01%% *}
[ "$ACTUAL_V01" = 6e2dd26716944f1ed98ec2daad81d7967e44d9fb7e42a29510e2ddbdb598afae ] || { echo 'SQL002_BOOTSTRAP_BLOCKED: v0.1 manifest mismatch' >&2; exit 3; }
for file in SQL002_SCOPED_MACHINE_AUTHORITY_DESIGN.md SQL002_MACHINE_AUTHORITY_CONTRACT.json SQL002_ACCEPTANCE_MATRIX.md CRUSH_DESIGN_REVIEW_HANDOFF.md; do
  [ -f "$BUNDLE_V01/$file" ] && [ ! -L "$BUNDLE_V01/$file" ] || { echo "SQL002_BOOTSTRAP_BLOCKED: missing/symlink v0.1 $file" >&2; exit 3; }
done
( cd "$BUNDLE_V01" && sha256sum -c SHA256SUMS.txt >/dev/null ) || { echo 'SQL002_BOOTSTRAP_BLOCKED: v0.1 artifact digest failure' >&2; exit 3; }
printf '%s\n' 'SQL002_V01_SOURCE_DESIGN_READY (4/4 hashes; no review or execution)'
printf '%s\n' 'SQL002_V02_DESIGN_EVIDENCE_READY (8/8 hashes; review NOT executed; no owner ratification or gate change)'
printf 'SQL002_V02_HANDOFF %s\n' "$ROOT/docs/SQL002_QWEN_INDEPENDENT_REVIEW.md"
