#!/bin/bash
set -e

# Install Terraform Provider locally
echo "Installing oneuptime Terraform Provider..."

# Build the provider
echo "Building provider..."
go build -o terraform-provider-oneuptime

# Create plugin directory
OS_ARCH="$(go env GOOS)_$(go env GOARCH)"
PLUGIN_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/oneuptime/oneuptime/12.0.28/$OS_ARCH"
mkdir -p "$PLUGIN_DIR"

# Copy binary
echo "Installing provider to $PLUGIN_DIR"
cp terraform-provider-oneuptime "$PLUGIN_DIR/"

echo "✅ Provider installed successfully!"
echo "You can now use it in your Terraform configuration:"
echo ""
echo "terraform {"
echo "  required_providers {"
echo "    oneuptime = {"
echo "      source = \"oneuptime/oneuptime\""
echo "      version = \"12.0.28\""
echo "    }"
echo "  }"
echo "}"
