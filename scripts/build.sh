#!/usr/bin/env bash
set -euo pipefail

LINUX_BIN="recheck/bin/recheck-linux-amd64"

# Get currently installed recheck version via npm
CURRENT_VERSION=$(npm list recheck --depth=0 2>/dev/null | grep 'recheck@' | sed 's/.*recheck@//')
if [ -z "$CURRENT_VERSION" ]; then
  echo "Error: Could not determine recheck version. Run 'mise run recheck-update' first."
  exit 1
fi

# If all binaries exist, skip build
if [ -f "$LINUX_BIN" ] && \
   [ -f "gorecheck" ] && \
   [ -f "gorecheck-linux-amd64" ]; then
  echo "recheck v${CURRENT_VERSION}: all binaries exist. Skipping build."
  exit 0
fi

# Download linux binary if missing
if [ ! -f "$LINUX_BIN" ]; then
  echo "Downloading recheck-linux-amd64 v${CURRENT_VERSION}..."
  curl -L -o "$LINUX_BIN" "https://github.com/makenowjust-labs/recheck/releases/download/v${CURRENT_VERSION}/recheck-linux-x64"
  chmod +x "$LINUX_BIN"
fi

go build -o gorecheck main.go
GOOS=linux GOARCH=amd64 go build -o gorecheck-linux-amd64 main.go

echo "Built with recheck v${CURRENT_VERSION}"
