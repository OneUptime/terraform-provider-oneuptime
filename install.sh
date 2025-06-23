#!/bin/bash
set -e

# Install Terraform Provider locally
echo "Installing oneuptime Terraform Provider..."

# Build the provider
echo "Building provider..."
go build -o terraform-provider-oneuptime

# Create plugin directory
PLUGIN_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/oneuptime/oneuptime/1.0.0/darwin_amd64"
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
echo "      source = "oneuptime/oneuptime""
echo "      version = "1.0.0""
echo "    }"
echo "  }"
echo "}"
