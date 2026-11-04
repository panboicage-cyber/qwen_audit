#!/bin/sh
# One-way validation of curated context; intentionally never saves memory.
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
command -v go >/dev/null || { echo 'REPO_MEMORY_BLOCKED: Go unavailable' >&2; exit 1; }
(cd "$ROOT" && go run ./tools/repomemory/main.go verify)
(cd "$ROOT" && go run ./tools/repomemory/main.go list >/dev/null)
(cd "$ROOT" && go run ./tools/repomemory/main.go context >/dev/null)
printf '%s\n' 'REPO_MEMORY_SMOKE_PASS (curated read-only data; no capture or audit execution)'
