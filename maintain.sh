#!/bin/bash
# go-jutsu maintenance helper

set -e

CMD=${1:-help}

case $CMD in
  check)
    echo "Checking all examples compile..."
    for pattern_dir in */; do
      # Skip hidden directories and non-pattern directories
      if [[ "$pattern_dir" == .* ]] || [[ ! -f "${pattern_dir}example.go" ]]; then
        continue
      fi
      echo "  Checking ${pattern_dir}example.go..."
      (cd "$pattern_dir" && go build -o /dev/null example.go)
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
    for pattern_dir in */; do
      if [[ "$pattern_dir" == .* ]] || [[ ! -f "${pattern_dir}example.go" ]]; then
        continue
      fi
      echo ""
      echo "=== Running ${pattern_dir}example.go ==="
      (cd "$pattern_dir" && timeout 2s go run example.go) || echo "Timed out or errored"
    done
    ;;
    
  validate)
    echo "Validating repository structure..."
    
    # Check if all patterns have both pattern.md and example.go
    echo "Checking pattern structure..."
    for pattern_dir in */; do
      # Skip hidden directories and special directories
      if [[ "$pattern_dir" == .* ]] || [[ "$pattern_dir" == ".github/" ]]; then
        continue
      fi
      
      # Skip if not a pattern directory (must have either pattern.md or example.go)
      if [[ ! -f "${pattern_dir}pattern.md" ]] && [[ ! -f "${pattern_dir}example.go" ]]; then
        continue
      fi
      
      if [[ ! -f "${pattern_dir}pattern.md" ]]; then
        echo "⚠ Warning: ${pattern_dir} missing pattern.md"
      fi
      
      if [[ ! -f "${pattern_dir}example.go" ]]; then
        echo "⚠ Warning: ${pattern_dir} missing example.go"
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
