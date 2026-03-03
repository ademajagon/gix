#!/usr/bin/env bash
# build-deb.sh — builds a .deb package for gix
# Usage: ./build-deb.sh <version> <arch> <binary_path>
# Example: ./build-deb.sh v0.3.0 amd64 ./dist/gix-linux-amd64

set -euo pipefail

VERSION="${1:?version required}"
ARCH="${2:?arch required (amd64|arm64)}"
BINARY="${3:?binary path required}"

DEB_VERSION="${VERSION#v}"
PACKAGE_NAME="gix_${DEB_VERSION}_${ARCH}"
STAGING="${TMPDIR:-/tmp}/${PACKAGE_NAME}"

echo "Building ${PACKAGE_NAME}.deb"

mkdir -p "${STAGING}/DEBIAN"
mkdir -p "${STAGING}/usr/bin"

install -m 0755 "${BINARY}" "${STAGING}/usr/bin/gix"

cat > "${STAGING}/DEBIAN/control" <<EOF
Package: gix
Version: ${DEB_VERSION}
Architecture: ${ARCH}
Maintainer: Agon Ademaj <agon@ademajagon.com>
Description: AI powered git commit assistant
 gix generates conventional commit messages from your staged diff
 and splits large changes into small atomic commits using AI.
Homepage: https://github.com/ademajagon/gix
Section: utils
Priority: optional
EOF

dpkg-deb --build --root-owner-group "${STAGING}" "${PACKAGE_NAME}.deb"

echo "Built ${PACKAGE_NAME}.deb"