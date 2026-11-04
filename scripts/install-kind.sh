#!/bin/sh
set -eu
DEST=${1:-.tools/bin}
mkdir -p "$DEST"
case "$(uname -m)" in x86_64) ARCH=amd64 ;; aarch64|arm64) ARCH=arm64 ;; *) exit 1 ;; esac
TEMP=$(mktemp -d)
trap 'rm -rf "$TEMP"' EXIT HUP INT TERM
URL="https://kind.sigs.k8s.io/dl/v0.33.0/kind-linux-$ARCH"
curl --connect-timeout 15 --max-time 180 -fsSL "$URL" -o "$TEMP/kind-linux-$ARCH"
curl --connect-timeout 15 --max-time 180 -fsSL "$URL.sha256sum" -o "$TEMP/checksum"
(cd "$TEMP"; sha256sum -c checksum)
install -m 0755 "$TEMP/kind-linux-$ARCH" "$DEST/kind"
