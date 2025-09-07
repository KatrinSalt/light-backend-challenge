#!/bin/bash

# Test script for the invoice approval workflow system
# This script builds and runs the comprehensive test suite

set -e

echo "🧪 Invoice Approval Workflow - Test Runner"
echo "=========================================="

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "📁 Project root: $PROJECT_ROOT"

# Change to project root
cd "$PROJECT_ROOT"

echo "🔨 Building test binary..."
go build -o bin/test ./cmd/test

echo "🚀 Running comprehensive test suite..."
echo ""

# Run the tests
./bin/test

echo ""
echo "✅ Test execution completed!"

# Clean up
echo "🧹 Cleaning up..."
rm -f bin/test

echo "🎉 All done!"
