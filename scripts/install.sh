#!/bin/bash

echo "NOT FINISHED"
exit 1

set -e

TMP_DIR="$(mktemp -d)"
# shellcheck disable=SC2064 # intentionally expands here
trap "rm -rf \"$TMP_DIR\"" EXIT INT TERM

# /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

OS="$(uname -s)"
ARCH="$(uname -m)"
test "$ARCH" = "aarch64" && ARCH="arm64"

TAR_FILE="awssso_${OS}_${ARCH}.tar.gz"

LATEST_VERSION=$(curl -sSL https://api.github.com/repos/webdestroya/aws-sso/releases/latest \
  | grep "tag_name.*" \
  | cut -d : -f 2,3 \
  | cut -d \" -f 2)

# dlurl=$(curl -sSL https://api.github.com/repos/webdestroya/aws-sso/releases/latest \
#   | grep "browser_download_url.*awssso_${OS}_${ARCH}.tar.gz\"" \
#   | cut -d : -f 2,3 \
#   | tr -d \" \
#   | xargs)


(
  cd "$TMP_DIR"
  echo "Downloading aws-sso version ${LATEST_VERSION}..."

  curl -sfLO "https://github.com/webdestroya/aws-sso/releases/download/${LATEST_VERSION}/${TAR_FILE}"
  curl -sfLO "https://github.com/webdestroya/aws-sso/releases/download/${LATEST_VERSION}/checksums.txt"

  echo "Verifying checksums..."
  sha256sum --ignore-missing --quiet --check checksums.txt

  # 

  echo "[${dlurl}]"

)