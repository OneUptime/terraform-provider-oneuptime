#!/bin/bash
set -e

echo "Running tests for oneuptime Terraform Provider..."

# Run unit tests
echo "Running unit tests..."
go test ./... -v

# Run acceptance tests if TF_ACC is set
if [ "$TF_ACC" = "1" ]; then
    echo "Running acceptance tests..."
    TF_ACC=1 go test ./... -v -timeout 120m
fi

echo "✅ Tests completed successfully!"
