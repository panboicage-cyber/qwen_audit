#!/bin/sh
# AuditGo v0.1.3 bootstrap for isolated Qwen coding workspaces.
# Installs a pinned local binary; never edits the audit target or Git configuration.
set -eu
umask 077

VERSION=0.1.3
ARCHIVE_NAME=AuditGo_v0.1.3_Crush_Hermes_Linux.tar.gz
EXPECTED_ARCHIVE=8c8625be74e54b70867196c88312d61b78b28f5737c81b7245b8339ca2e57640
EXPECTED_BINARY=b6c55aa3822c4dc6ae045c6a1f697a8f75903902e9e22ba957fa2e11f3a51a1f
EXPECTED_REGISTRY=385eaacb922eefa14231c59cc8c5f5cdfd9633938cc2daa5005045316af6451b

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_dir=$(CDPATH= cd -- "$script_dir/.." && pwd)
case "$(uname -s):$(uname -m)" in
  Linux:x86_64|Linux:amd64) ;;
  *) echo 'AUDITGO_SETUP_UNSUPPORTED_PLATFORM (requires Linux AMD64)' >&2; exit 2 ;;
esac
for tool in sha256sum tar mktemp grep; do
  command -v "$tool" >/dev/null 2>&1 || { echo "AUDITGO_SETUP_MISSING_TOOL: $tool" >&2; exit 2; }
done

archive=${AUDITGO_ARCHIVE:-"$repo_dir/vendor/$ARCHIVE_NAME"}
if [ ! -f "$archive" ] && [ -f "$repo_dir/$ARCHIVE_NAME" ] && [ -z "${AUDITGO_ARCHIVE:-}" ]; then
  archive="$repo_dir/$ARCHIVE_NAME"
fi
scratch=''
stage=''
cleanup() {
  if [ -n "$scratch" ] && [ -d "$scratch" ]; then rm -rf -- "$scratch"; fi
  if [ -n "$stage" ] && [ -d "$stage" ]; then rm -rf -- "$stage"; fi
}
trap cleanup EXIT HUP INT TERM

if [ ! -f "$archive" ] && [ -n "${AUDITGO_ARCHIVE_URL:-}" ]; then
  case "$AUDITGO_ARCHIVE_URL" in https://*) ;; *) echo 'AUDITGO_SETUP_INVALID_URL: HTTPS required' >&2; exit 2;; esac
  command -v curl >/dev/null 2>&1 || { echo 'AUDITGO_SETUP_MISSING_TOOL: curl' >&2; exit 2; }
  scratch=$(mktemp -d)
  archive="$scratch/$ARCHIVE_NAME"
  curl --fail --location --silent --show-error --max-time 120 --output "$archive" "$AUDITGO_ARCHIVE_URL"
fi
if [ ! -f "$archive" ]; then
  echo "AUDITGO_ARCHIVE_MISSING: $archive" >&2
  echo 'Place the pinned release archive under vendor/ or set AUDITGO_ARCHIVE to an existing file.' >&2
  echo 'Alternatively set AUDITGO_ARCHIVE_URL to an HTTPS URL hosting those exact bytes.' >&2
  exit 2
fi
actual_archive=$(sha256sum "$archive" | cut -d ' ' -f 1)
if [ "$actual_archive" != "$EXPECTED_ARCHIVE" ]; then
  echo "AUDITGO_ARCHIVE_HASH_MISMATCH: expected $EXPECTED_ARCHIVE, got $actual_archive" >&2
  exit 3
fi

install_base="$HOME/.local/share/auditgo"
install_dir="$install_base/auditgo-v$VERSION"
bin="$install_dir/dist/auditgo-linux-amd64"
mkdir -p "$install_base" "$HOME/.local/bin"
if [ ! -e "$bin" ]; then
  [ ! -e "$install_dir" ] || { echo 'AUDITGO_INSTALL_CONFLICT: existing incomplete install' >&2; exit 4; }
  stage=$(mktemp -d "$install_base/.stage.XXXXXXXX")
  tar -xzf "$archive" -C "$stage"
  staged_bin="$stage/auditgo-v$VERSION/dist/auditgo-linux-amd64"
  [ -f "$staged_bin" ] || { echo 'AUDITGO_PACKAGE_INVALID: binary missing' >&2; exit 3; }
  staged_hash=$(sha256sum "$staged_bin" | cut -d ' ' -f 1)
  [ "$staged_hash" = "$EXPECTED_BINARY" ] || { echo 'AUDITGO_BINARY_HASH_MISMATCH' >&2; exit 3; }
  chmod 755 "$staged_bin"
  mv -- "$stage/auditgo-v$VERSION" "$install_dir"
  rm -r -- "$stage"
fi
actual_binary=$(sha256sum "$bin" | cut -d ' ' -f 1)
[ "$actual_binary" = "$EXPECTED_BINARY" ] || { echo 'AUDITGO_INSTALLED_BINARY_HASH_MISMATCH' >&2; exit 3; }

doctor=$($bin doctor)
printf '%s\n' "$doctor"
printf '%s\n' "$doctor" | grep -Fx "AuditGo $VERSION" >/dev/null || { echo 'AUDITGO_DOCTOR_VERSION_MISMATCH' >&2; exit 3; }
printf '%s\n' "$doctor" | grep -Fx "Registry SHA-256: $EXPECTED_REGISTRY" >/dev/null || { echo 'AUDITGO_DOCTOR_REGISTRY_MISMATCH' >&2; exit 3; }

link="$HOME/.local/bin/auditgo"
if [ -e "$link" ] && [ ! -L "$link" ]; then
  echo "AUDITGO_SYMLINK_CONFLICT: $link exists and is not a symlink" >&2; exit 4
fi
ln -sfn "$bin" "$link"
printf 'AUDITGO_SETUP_PASS %s\n' "$link"
printf 'AUDITGO_ARCHIVE_SHA256 %s\n' "$EXPECTED_ARCHIVE"
printf 'AUDITGO_BINARY_SHA256 %s\n' "$EXPECTED_BINARY"
