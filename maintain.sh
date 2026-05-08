#!/bin/bash
# go-jutsu maintenance helper

set -e

CMD=${1:-help}

case $CMD in
  check)
    echo "Checking all examples compile..."
    for example in examples/*.go; do
      if [ -f "$example" ]; then
        echo "  Checking $example..."
        go build -o /dev/null "$example"
      fi
    done
    echo "✓ All examples compile"
    ;;
    
  format)
    echo "Formatting all Go code..."
    gofmt -w .
    echo "✓ Code formatted"
    ;;
    
  run-all)
    echo "Running all examples..."
    for example in examples/*.go; do
      if [ -f "$example" ]; then
        echo ""
        echo "=== Running $example ==="
        timeout 2s go run "$example" || echo "Timed out or errored"
      fi
    done
    ;;
    
  validate)
    echo "Validating repository structure..."
    
    # Check if all patterns have examples
    echo "Checking patterns have examples..."
    for pattern in patterns/*.md; do
      if [ "$pattern" != "patterns/README.md" ]; then
        basename=$(basename "$pattern" .md)
        if [ ! -f "examples/$basename.go" ]; then
          echo "⚠ Warning: Pattern $basename has no example"
        fi
      fi
    done
    
    echo "✓ Validation complete"
    ;;
    
  help|*)
    echo "go-jutsu maintenance helper"
    echo ""
    echo "Usage: ./maintain.sh [command]"
    echo ""
    echo "Commands:"
    echo "  check      - Verify all examples compile"
    echo "  format     - Run gofmt on all code"
    echo "  run-all    - Run all examples (2s timeout each)"
    echo "  validate   - Check repo structure consistency"
    echo "  help       - Show this help"
    ;;
esac
