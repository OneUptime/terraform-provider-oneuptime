#!/bin/bash
set -e

echo "Building oneuptime Terraform Provider..."

# Clean previous builds
rm -f terraform-provider-oneuptime

# Build for current platform
go build -o terraform-provider-oneuptime

echo "✅ Build completed successfully!"
echo "Binary: terraform-provider-oneuptime"
