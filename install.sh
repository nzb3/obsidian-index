#!/bin/bash

# obsidian-index installation script
# This script downloads and installs obsidian-index on macOS

set -e

# Configuration
BINARY_NAME="obsidian-index" 
INSTALL_DIR="/usr/local/bin"
REPO="nzb3/obsidian-index"
LATEST_RELEASE_URL="https://api.github.com/repos/$REPO/releases/latest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to detect macOS architecture
detect_arch() {
    local arch=$(uname -m)
    case $arch in
        x86_64)
            echo "amd64"
            ;;
        arm64)
            echo "arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
}

# Function to get latest release info
get_latest_release() {
    local response=$(curl -s "$LATEST_RELEASE_URL")
    local version=$(echo "$response" | grep '"tag_name"' | cut -d'"' -f4)
    local download_url=$(echo "$response" | grep '"browser_download_url"' | grep "darwin-$(detect_arch)" | cut -d'"' -f4)
    
    if [ -z "$version" ] || [ -z "$download_url" ]; then
        print_error "Failed to get latest release information"
        exit 1
    fi
    
    echo "$version|$download_url"
}

# Function to download and install
install_binary() {
    local release_info=$1
    local version=$(echo "$release_info" | cut -d'|' -f1)
    local download_url=$(echo "$release_info" | cut -d'|' -f2)
    
    print_status "Installing $BINARY_NAME version $version..."
    
    # Create temporary directory
    local temp_dir=$(mktemp -d)
    local archive_path="$temp_dir/$BINARY_NAME.tar.gz"
    
    # Download the binary
    print_status "Downloading from $download_url..."
    if ! curl -L -o "$archive_path" "$download_url"; then
        print_error "Failed to download binary"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # Extract the binary
    print_status "Extracting binary..."
    tar -xzf "$archive_path" -C "$temp_dir"
    
    # Find the extracted binary
    local binary_path=$(find "$temp_dir" -name "$BINARY_NAME" -type f | head -1)
    
    if [ -z "$binary_path" ]; then
        print_error "Binary not found in archive"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # Make binary executable
    chmod +x "$binary_path"
    
    # Install to system directory
    print_status "Installing to $INSTALL_DIR..."
    if ! sudo cp "$binary_path" "$INSTALL_DIR/"; then
        print_error "Failed to install binary to $INSTALL_DIR"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # Clean up
    rm -rf "$temp_dir"
    
    print_status "Installation completed successfully!"
    print_status "You can now use '$BINARY_NAME' command from anywhere."
}

# Function to check if already installed
check_existing() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local current_version=$($BINARY_NAME --version 2>/dev/null || echo "unknown")
        print_warning "$BINARY_NAME is already installed (version: $current_version)"
        read -p "Do you want to update it? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_status "Installation cancelled."
            exit 0
        fi
    fi
}

# Function to verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version=$($BINARY_NAME --version 2>/dev/null || echo "unknown")
        print_status "Installation verified! Version: $version"
        print_status "Try running: $BINARY_NAME --help"
    else
        print_error "Installation failed - binary not found in PATH"
        exit 1
    fi
}

# Main installation flow
main() {
    print_status "obsidian-index installer for macOS"
    print_status "=================================="
    
    # Check if running on macOS
    if [[ "$OSTYPE" != "darwin"* ]]; then
        print_error "This installer is for macOS only"
        exit 1
    fi
    
    # Check for required tools
    if ! command -v curl >/dev/null 2>&1; then
        print_error "curl is required but not installed"
        exit 1
    fi
    
    if ! command -v sudo >/dev/null 2>&1; then
        print_error "sudo is required but not available"
        exit 1
    fi
    
    # Check if already installed
    check_existing
    
    # Get latest release and install
    local release_info
    if ! release_info=$(get_latest_release); then
        exit 1
    fi
    
    install_binary "$release_info"
    verify_installation
}

# Run main function
main "$@"

