#!/bin/bash
set -e

echo "🔧 Setting up development environment..."

# Update package manager
echo "📦 Updating package manager..."
sudo apt update

# Install tmux
echo "📦 Installing tmux and pipx"
sudo apt install -y tmux pipx
echo "Installing sql linter"
pipx install sqlfluff

# Install golangci-lint
echo "📦 Installing golangci-lint v2.7.2..."
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.7.2

echo "✅ Development environment setup complete!"
