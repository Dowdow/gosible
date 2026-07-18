#!/usr/bin/env bash
# Regenerates demo.gif using the demo/ fixtures instead of a real homelab.
# Requires vhs (https://github.com/charmbracelet/vhs) to be installed.
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

echo "Building gosible..."
go build -o build/gosible ./cmd

echo "Building the demo fake SSH server..."
go build -o build/fakesshd ./demo/server

echo "Starting the demo fake SSH server..."
build/fakesshd &
server_pid=$!
trap 'kill "$server_pid" 2>/dev/null || true' EXIT

# Give the server a moment to bind before vhs starts typing commands.
sleep 1

echo "Recording demo/demo.gif with vhs..."
PATH="$repo_root/build:$PATH" vhs demo/demo.tape

echo "Done: demo/demo.gif regenerated."
