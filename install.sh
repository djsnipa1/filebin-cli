#!/bin/sh
set -e

# filebin - install script
# Downloads the latest filebin CLI binary (linux-amd64) from GitHub releases
# and installs it to a bin directory.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/djsnipa1/filebin-cli/main/install.sh | sh
#
# Environment:
#   FILEBIN_BIN_DIR - install directory (default: $HOME/.local/bin)

REPO=djsnipa1/filebin-cli
BINARY=filebin-linux-amd64
APP=filebin
BIN_DIR="${FILEBIN_BIN_DIR:-$HOME/.local/bin}"

# 1) Resolve the latest release tag
log() { printf '\033[36m%s\033[0m\n' "$*"; }
fail() { printf '\033[31mError: %s\033[0m\n' "$*" >&2; exit 1; }

log "Resolving latest release for $REPO..."
TAG=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n1)
[ -z "$TAG" ] && fail "unable to find the latest release for $REPO"

# 2) Download the linux-amd64 binary
URL="https://github.com/$REPO/releases/download/$TAG/$BINARY"
log "Downloading $BINARY from $TAG..."
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

if ! curl -fsSL "$URL" -o "$TMPDIR/$APP"; then
  fail "failed to download $URL"
fi

# 3) Install
mkdir -p "$BIN_DIR"
install -m 0755 "$TMPDIR/$APP" "$BIN_DIR/$APP"

printf 'Installed \033[1m%s\033[0m (%s) to \033[1m%s/%s\033[0m\n' \
  "$APP" "$TAG" "$BIN_DIR" "$APP"

case ":$PATH:" in
  *":$BIN_DIR:"*) ;;
  *) printf 'Note: %s is not on your PATH. Add it with:\n  export PATH="%s:$PATH"\n' "$BIN_DIR" "$BIN_DIR" ;;
esac
