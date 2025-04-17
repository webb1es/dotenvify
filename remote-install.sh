#!/bin/bash

# Remote installer for DotEnvify
# Usage: curl -s https://raw.githubusercontent.com/webb1es/dotenvify/main/remote-install.sh | bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to show colored messages
show_message() {
  local type=$1
  local message=$2

  if [ "$type" == "success" ]; then
    echo -e "${GREEN}✓ $message${NC}"
  elif [ "$type" == "error" ]; then
    echo -e "${RED}✗ $message${NC}"
  elif [ "$type" == "warning" ]; then
    echo -e "${YELLOW}⚠ $message${NC}"
  else
    echo "$message"
  fi
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
  show_message "error" "Go is not installed. Please install Go first."
  echo "Visit https://golang.org/doc/install for installation instructions."
  exit 1
fi

# Create temporary directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# GitHub repository details
REPO_OWNER="webb1es"  # Replace with your GitHub username
REPO_NAME="dotenvify"
BRANCH="main"  # or "master" depending on your default branch

show_message "info" "Downloading DotEnvify from GitHub..."

# Create directory structure
mkdir -p plugins/azure

# Download the Go source files
show_message "info" "Downloading DotEnvify source files..."
curl -s "https://raw.githubusercontent.com/$REPO_OWNER/$REPO_NAME/$BRANCH/dotenvify.go" -o dotenvify.go
curl -s "https://raw.githubusercontent.com/$REPO_OWNER/$REPO_NAME/$BRANCH/plugins/azure/azure.go" -o plugins/azure/azure.go
curl -s "https://raw.githubusercontent.com/$REPO_OWNER/$REPO_NAME/$BRANCH/go.mod" -o go.mod

if [ ! -f "dotenvify.go" ] || [ ! -f "plugins/azure/azure.go" ]; then
  show_message "error" "Failed to download source code from GitHub."
  exit 1
fi

# Compile
show_message "info" "Compiling DotEnvify..."
go build

if [ ! -f "dotenvify" ]; then
  show_message "error" "Compilation failed."
  exit 1
fi

# Install
INSTALL_DIR="/usr/local/bin"
sudo mv dotenvify "$INSTALL_DIR/"

# Verify installation
if [ -x "$INSTALL_DIR/dotenvify" ]; then
  show_message "success" "DotEnvify has been installed successfully!"
  echo "Usage: dotenvify source_file [output_file]"
else
  show_message "error" "Installation failed."
fi

# Clean up
cd - > /dev/null
rm -rf "$TEMP_DIR"
