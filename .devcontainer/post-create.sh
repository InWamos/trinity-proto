#!/bin/bash
set -e

echo "🔧 Setting up development environment..."

# Update package manager
echo "📦 Updating package manager..."
sudo apt update

# Install Go 1.25.7
echo "📦 Installing Go 1.25.7"
GO_VERSION="1.25.7"
case "$(uname -m)" in
  x86_64)   GO_ARCH="amd64" ;;
  aarch64)  GO_ARCH="arm64" ;;
  armv7l)   GO_ARCH="armv6l" ;;
  *)        echo "Unsupported architecture: $(uname -m)"; exit 1 ;;
esac
cd /tmp
wget -q https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-${GO_ARCH}.tar.gz
rm go${GO_VERSION}.linux-${GO_ARCH}.tar.gz
echo "✅ Go ${GO_VERSION} installed"

# Install tmux
echo "📦 Installing tmux and pipx"
sudo apt install -y tmux pipx

echo "Installing sql linters"
pipx install sqlfluff
pipx install squawk-cli

# Install golangci-lint
echo "📦 Installing golangci-lint"
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.9.0

# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Set up protoc architecture (different naming than Go)
case "$(uname -m)" in
  x86_64)   PB_ARCH="amd64" ;;
  aarch64)  PB_ARCH="aarch64" ;;
  armv7l)   PB_ARCH="armv7" ;;
  *)        echo "Unsupported architecture for protobuf"; exit 1 ;;
esac

PB_REL="https://github.com/protocolbuffers/protobuf/releases"
curl -LO $PB_REL/download/v33.5/protoc-33.5-linux-${PB_ARCH}.zip
unzip protoc-33.5-linux-${PB_ARCH}.zip -d $HOME/.local
export PATH="$PATH:$HOME/.local/bin"

echo "✅ Development environment setup complete!"
