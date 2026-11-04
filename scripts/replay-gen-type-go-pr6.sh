#!/usr/bin/env bash
# Pinned, fail-closed project profile for an isolated Qwen coding environment.
# This script does not publish findings, approve PRs, or mutate upstream GitHub.
set -Eeuo pipefail
umask 077

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)"
LOCK="$ROOT/profiles/gen-type-go-pr6/target.lock.json"
OUT_BASE="${QWEN_AUDIT_OUTPUT_DIR:-$ROOT/audit-artifacts/gen-type-go-pr6}"
AUDITGO="${AUDITGO_BIN:-$HOME/.local/bin/auditgo}"
mkdir -p "$OUT_BASE"
OUT="$(mktemp -d "$OUT_BASE/run.XXXXXXXX")"

for exe in git go python3 sha256sum cmp grep mktemp date; do
  command -v "$exe" >/dev/null 2>&1 || { echo "BLOCKED: missing $exe" >&2; exit 2; }
done
[ -x "$AUDITGO" ] || { echo "BLOCKED: AuditGo binary absent; run sh scripts/setup-auditgo.sh" >&2; exit 2; }
[ -f "$LOCK" ] || { echo "BLOCKED: lock absent" >&2; exit 2; }

# Load and strictly type-check the immutable target lock in one operation.
eval "$(python3 - "$LOCK" <<'PY'
import json, re, shlex, sys
from pathlib import Path
p=Path(sys.argv[1]); v=json.loads(p.read_text())
if v['schema'] != 'qwen-audit-target-lock-v1': raise SystemExit('invalid lock schema')
for k in ('source_head','source_tree','base_main','original_v04_head'):
    if not re.fullmatch('[0-9a-f]{40}',v[k]): raise SystemExit('bad revision '+k)
for k in ('upstream_mutations_authorized','merge_authorized','release_authorized','new_authority_authorized'):
    if type(v[k]) is not bool or v[k]: raise SystemExit('unsafe lock '+k)
for k in ('target_repo','source_head','source_tree','base_main','original_v04_head'):
    print(k.upper()+'='+shlex.quote(v[k]))
PY
)"

export GIT_TERMINAL_PROMPT=0 GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off GOFLAGS=-mod=readonly
if [ "${QWEN_AUDIT_TARGET_DIR:-}" ]; then
  TARGET="$QWEN_AUDIT_TARGET_DIR"
  if [ -e "$TARGET" ]; then echo 'BLOCKED: target directory already exists; use a fresh directory' >&2; exit 2; fi
  mkdir -p "$(dirname -- "$TARGET")"
else
  TARGET="$(mktemp -d "${TMPDIR:-/tmp}/qwen-pr6-target.XXXXXXXX")"
  rmdir "$TARGET"
fi
printf 'TARGET_DIR %s\n' "$TARGET" > "$OUT/preflight.txt"
# Preserve the target checkout and all logs as audit evidence; do not clean automatically.
git clone --no-checkout "$TARGET_REPO" "$TARGET" > "$OUT/clone.log" 2>&1 || { echo 'BLOCKED: target clone failed' >&2; exit 3; }
git -C "$TARGET" checkout --detach "$SOURCE_HEAD" >> "$OUT/clone.log" 2>&1 || { echo 'BLOCKED: target commit unavailable' >&2; exit 3; }

python3 - "$LOCK" "$TARGET" > "$OUT/source-integrity.txt" <<'PY'
import json, pathlib, subprocess, sys, re
v=json.loads(pathlib.Path(sys.argv[1]).read_text()); root=pathlib.Path(sys.argv[2]);
def git(*args):
    return subprocess.check_output(['git','-C',str(root),*args],text=True).strip()
assert git('remote','get-url','origin')==v['target_repo'], 'remote mismatch'
assert git('rev-parse','HEAD')==v['source_head'], 'HEAD mismatch'
assert git('rev-parse','HEAD^{tree}')==v['source_tree'], 'tree mismatch'
assert git('show','-s','--format=%P','HEAD').split()==v['parents'], 'parents mismatch'
assert git('status','--porcelain')=='', 'checkout dirty'
entries=git('ls-files','--cached','-z').split('\0')[:-1]
assert len(entries)==v['tracked_files_expected'], 'tracked count mismatch'
manifest=root/'SHA256SUMS'; rows=manifest.read_text().splitlines()
assert len(rows)==v['manifest_entries_expected'], 'manifest count mismatch'
parsed=[]
for row in rows:
    m=re.fullmatch(r'([a-f0-9]{64})  (.+)',row)
    assert m, 'malformed manifest entry'
    parsed.append(m.group(2))
assert parsed==sorted(set(entries)-{'SHA256SUMS'}), 'manifest incomplete or unsorted'
paths=git('diff','--name-only',v['original_v04_head'],v['source_head']).splitlines()
assert paths==sorted(v['reconciliation_paths']), 'unexpected reconciliation diff'
for p in entries:
    if p.endswith('.go') or p in ('go.mod','.github/workflows/ci.yml'):
        assert git('rev-parse',v['original_v04_head']+':'+p)==git('rev-parse',v['source_head']+':'+p), 'Go/CI changed: '+p
for p in ('CHANGELOG.md','docs/BUILD_STATUS.md','docs/V0.3_PROMOTION_RECONCILIATION.md'):
    assert git('rev-parse',v['base_main']+':'+p)==git('rev-parse',v['source_head']+':'+p), 'governance import mismatch '+p
print('SOURCE_IDENTITY_PASS')
print('MANIFEST_PATH_COVERAGE_PASS')
print('RECONCILIATION_BLOBS_PASS')
print('ACTUAL_HEAD='+v['source_head'])
print('ACTUAL_TREE='+v['source_tree'])
PY
(cd "$TARGET" && sha256sum -c SHA256SUMS) > "$OUT/checksums.log" 2>&1 || { echo 'FAIL: manifest bytes mismatch' >&2; exit 4; }

"$AUDITGO" doctor > "$OUT/auditgo-doctor.log" 2>&1
python3 - "$LOCK" "$OUT/auditgo-doctor.log" <<'PY'
import json,pathlib,sys
v=json.loads(pathlib.Path(sys.argv[1]).read_text())
s=pathlib.Path(sys.argv[2]).read_text()
assert 'AuditGo '+v['auditgo']['version'] in s, 'AuditGo version mismatch'
assert 'Registry SHA-256: '+v['auditgo']['registry_sha256'] in s, 'AuditGo registry mismatch'
PY

# AuditGo init is deterministic with these exact options and immutable lock metadata.
python3 - "$LOCK" "$AUDITGO" "$OUT" <<'PY'
import json,subprocess,sys,pathlib
v=json.loads(pathlib.Path(sys.argv[1]).read_text()); bin=sys.argv[2];out=pathlib.Path(sys.argv[3]);
cmd=[bin,'init','-out',str(out/'auditgo-selector.json'),'-target',v['profile_id'],'-description','Pinned read-only Gen Type-Go PR #6 independent audit','-revision',v['source_head']]
for claim in v['claims']:cmd+=['-claim',claim]
for a in v['audit_types']:cmd+=['-type',a]
for m in v['modules']:cmd+=['-module',m]
subprocess.run(cmd,check=True)
PY
"$AUDITGO" plan -config "$OUT/auditgo-selector.json" > "$OUT/auditgo-plan.json"
"$AUDITGO" render -config "$OUT/auditgo-selector.json" -out "$OUT/auditgo-mission.md"
"$AUDITGO" seed -config "$OUT/auditgo-selector.json" -out "$OUT/auditgo-unassessed-return.json"

: > "$OUT/commands.tsv"
run_check() {
  local name="$1"; shift
  local exitcode=0
  (cd "$TARGET" && "$@") > "$OUT/${name}.log" 2>&1 || exitcode=$?
  printf '%s\t%s\t%s\n' "$name" "$exitcode" "$OUT/${name}.log" >> "$OUT/commands.tsv"
  if [ "$exitcode" -eq 0 ]; then printf 'PASS %s\n' "$name"; else printf 'FAIL %s (exit %d)\n' "$name" "$exitcode"; fi
}
run_check go_version go version
run_check gofmt bash -c 'test -z "$(gofmt -l cmd internal)"'
run_check go_mod_verify go mod verify
run_check go_test go test -count=1 ./...
run_check go_vet go vet ./...
run_check go_race go test -race -count=1 ./...
run_check go_build go build ./...
run_check inspect go run ./cmd/gen-type-go inspect --id c testdata/stats_forest.json
run_check diagnose_valid go run ./cmd/gen-type-go diagnose testdata/stats_forest.json
run_check fuzz_inspect go test ./internal/graphinspect -run '^$' -fuzz '^FuzzInspectPermutation$' -fuzztime=30s
run_check fuzz_diagnose go test ./internal/graphdiagnostics -run '^$' -fuzz '^FuzzDiagnoseDeterminism$' -fuzztime=30s

# Not a complete independent security audit: selected checklist items remain UNASSESSED.
# Qwen must interpret logs, investigate adverse behavior, complete the seeded ledger,
# and run AuditGo validate only after evidence has been assessed.
python3 - "$OUT/commands.tsv" "$OUT" <<'PY'
import json,pathlib,sys
rows=[]
for s in pathlib.Path(sys.argv[1]).read_text().splitlines():
    name,code,log=s.split('\t');rows.append(dict(check=name,exit_code=int(code),log=log,passed=int(code)==0))
out=pathlib.Path(sys.argv[2]);obj={'schema':'qwen-pr6-execution-receipt-v1','commands':rows,'all_commands_passed':all(x['passed'] for x in rows),'complete_security_audit':False,'formal_approval':False,'upstream_mutation':False}
(out/'execution-receipt.json').write_text(json.dumps(obj,indent=2)+'\n')
PY
(cd "$OUT" && sha256sum source-integrity.txt checksums.log auditgo-selector.json auditgo-plan.json auditgo-mission.md auditgo-unassessed-return.json commands.tsv execution-receipt.json > EVIDENCE_SHA256SUMS.txt)
printf 'QWEN_PR6_REPLAY_COMPLETE (not a final audit verdict)\nARTIFACTS %s\nTARGET %s\n' "$OUT" "$TARGET"
