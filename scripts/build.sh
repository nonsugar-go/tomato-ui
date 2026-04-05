#!/usr/bin/bash
set -e

echo "Building binaries..."

mkdir -p bin

GOOS=windows GOARCH=amd64 go build -o bin/utmconv.exe ./cmd/utmconv
GOOS=windows GOARCH=amd64 go build -o bin/vm-tui.exe ./cmd/vm-tui

echo "Done!"
