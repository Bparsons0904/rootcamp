#!/bin/bash
# Simple watch script for TUI development

while true; do
    # Build the app
    go build -o ./tmp/main ./cmd/rootcamp

    # Run it (will block until exit)
    ./tmp/main

    # Wait a moment before allowing rebuild
    sleep 1
done
