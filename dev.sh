#!/bin/bash
# Development helper script

set -e

CMD="go run ./cmd/projectdb --db projectdb.db"

case "$1" in
    seed)
        echo "Seeding database with unique tasks..."
        rm -f projectdb.db
        ./scripts/seed-test-data.sh projectdb.db
        ;;
    
    tui)
        echo "Starting TUI..."
        make run
        ;;
    
    test)
        echo "Running tests..."
        make test
        ;;
    
    *)
        echo "Usage: $0 {seed|tui|test}"
        exit 1
        ;;
esac
