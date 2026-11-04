#!/bin/sh
# Explicit opt-in download of pinned upstream binaries with SHA-256 verification.
set -eu
DEST=${1:-.tools/bin}
mkdir -p "$DEST"
DEST=$(cd "$DEST" && pwd)
case "$(uname -s)" in Linux) OS=linux ;; *) echo 'Only Linux is supported by this installer' >&2; exit 1 ;; esac
case "$(uname -m)" in x86_64) ARCH=amd64 ;; aarch64|arm64) ARCH=arm64 ;; *) echo 'Unsupported architecture' >&2; exit 1 ;; esac
TOFU_VERSION=${TOFU_VERSION:-1.13.5}
TERRAFORM_VERSION=${TERRAFORM_VERSION:-1.17.1}
KUBECTL_VERSION=${KUBECTL_VERSION:-1.37.2}
TEMP=$(mktemp -d)
trap 'rm -rf "$TEMP"' EXIT HUP INT TERM
zip_tool() {
 name=$1; version=$2; base=$3
 archive="${name}_${version}_${OS}_${ARCH}.zip"
 curl --connect-timeout 15 --max-time 180 --fail --silent --show-error --location "$base/$archive" -o "$TEMP/$archive"
 curl --connect-timeout 15 --max-time 180 --fail --silent --show-error --location "$base/${name}_${version}_SHA256SUMS" -o "$TEMP/sums"
 (cd "$TEMP"; awk -v file="$archive" '$2==file || $2=="*"file {print}' sums > selected; test -s selected; sha256sum -c selected)
 unzip -qo "$TEMP/$archive" "$name" -d "$TEMP"
 install -m 0755 "$TEMP/$name" "$DEST/$name"
}
if [ "${SKIP_TOFU:-0}" != 1 ]; then
TOFU_SOURCE_SHA256=${TOFU_SOURCE_SHA256:-14475a9be680579788273bc004a246e7b8b58aee2a83921687e0c2657033d38c}
command -v go >/dev/null
curl --connect-timeout 15 --max-time 180 -fsSL "https://github.com/opentofu/opentofu/archive/refs/tags/v$TOFU_VERSION.tar.gz" -o "$TEMP/tofu.tar.gz"
(cd "$TEMP"; printf '%s  tofu.tar.gz\n' "$TOFU_SOURCE_SHA256" | sha256sum -c -)
mkdir "$TEMP/source"
tar -xzf "$TEMP/tofu.tar.gz" -C "$TEMP/source" --strip-components=1
(cd "$TEMP/source"; CGO_ENABLED=0 go build -trimpath -ldflags '-s -w -X github.com/opentofu/opentofu/version.dev=no' -o "$DEST/tofu" ./cmd/tofu)
fi
zip_tool terraform "$TERRAFORM_VERSION" "https://releases.hashicorp.com/terraform/$TERRAFORM_VERSION"
base="https://dl.k8s.io/release/v$KUBECTL_VERSION/bin/$OS/$ARCH/kubectl"
curl --connect-timeout 15 --max-time 180 --fail --silent --show-error --location "$base" -o "$TEMP/kubectl"
curl --connect-timeout 15 --max-time 180 --fail --silent --show-error --location "$base.sha256" -o "$TEMP/kubectl.sha256"
(cd "$TEMP"; printf '%s  kubectl\n' "$(cat kubectl.sha256)" | sha256sum -c -)
install -m 0755 "$TEMP/kubectl" "$DEST/kubectl"
echo "Tools installed in $DEST"
