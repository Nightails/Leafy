#!/usr/bin/env bash

set -euo pipefail

VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo dev)"
BINARY_NAME="leafy"

echo "Building ${BINARY_NAME} version ${VERSION}..."

go build \
  -ldflags "-X main.version=${VERSION}" \
  -o "${BINARY_NAME}" \
  .

echo "Done: ./${BINARY_NAME}"