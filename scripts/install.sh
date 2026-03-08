#!/bin/bash
# Install script for dopa
# Usage: curl -sSL https://raw.githubusercontent.com/marekbrze/dopadone/main/scripts/install.sh | sh

set -e

REPO="marekbrze/dopadone"
BINARY_NAME="dopa"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo "Unsupported architecture: $ARCH" >&2
            exit 1
            ;;
    esac
    
    if [ "$OS" = "darwin" ]; then
        OS="darwin"
    elif [ "$OS" = "linux" ]; then
        OS="linux"
    else
        echo "Unsupported OS: $OS" >&2
        exit 1
    fi
    
    echo "${OS}-${ARCH}"
}

get_latest_version() {
    curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
}

download_binary() {
    local VERSION="$1"
    local PLATFORM="$2"
    local EXT="tar.gz"
    
    if [ "$(echo "$PLATFORM" | cut -d'-' -f1)" = "windows" ]; then
        EXT="zip"
    fi
    
    local URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}-${PLATFORM}.${EXT}"
    local TMP_DIR=$(mktemp -d)
    
    echo "Downloading ${BINARY_NAME} ${VERSION} for ${PLATFORM}..."
    curl -sSL -o "${TMP_DIR}/${BINARY_NAME}.${EXT}" "$URL"
    
    cd "$TMP_DIR"
    if [ "$EXT" = "zip" ]; then
        unzip -o "${BINARY_NAME}.${EXT}"
    else
        tar xzf "${BINARY_NAME}.${EXT}"
    fi
    
    echo "${TMP_DIR}/${BINARY_NAME}"
}

install_binary() {
    local BINARY_PATH="$1"
    
    echo "Installing to ${INSTALL_DIR}..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY_PATH" "${INSTALL_DIR}/${BINARY_NAME}"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        echo "sudo required for installation to ${INSTALL_DIR}"
        sudo mv "$BINARY_PATH" "${INSTALL_DIR}/${BINARY_NAME}"
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    echo "Installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"
}

main() {
    echo "Installing ${BINARY_NAME}..."
    
    PLATFORM=$(detect_platform)
    VERSION=$(get_latest_version)
    
    echo "Latest version: ${VERSION}"
    echo "Platform: ${PLATFORM}"
    
    BINARY_PATH=$(download_binary "$VERSION" "$PLATFORM")
    install_binary "$BINARY_PATH"
    
    echo ""
    echo "Installation complete!"
    echo "Run '${BINARY_NAME} version' to verify."
}

main "$@"
