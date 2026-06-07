#!/usr/bin/bash
set -e

echo "Building binaries..."

mkdir -p bin

GOOS=windows GOARCH=amd64 go build -o bin/utmconv.exe ./cmd/utmconv
GOOS=windows GOARCH=amd64 go build -o bin/cp-sh-pkg.exe ./cmd/cp-sh-pkg
GOOS=windows GOARCH=amd64 go build -o bin/cp-dump.exe ./cmd/cp-dump
GOOS=windows GOARCH=amd64 go build -o bin/vm-tui.exe ./cmd/vm-tui

echo "Done!"
