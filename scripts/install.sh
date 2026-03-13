#!/usr/bin/env bash
#
# Install script for dopa
# Usage: curl -sSL https://raw.githubusercontent.com/marekbrze/dopadone/main/scripts/install.sh | bash
#
# Flags:
#   --dry-run    Simulate platform detection and download URL without installing
#   --yes        Skip confirmation prompts (auto-confirm upgrades)
#   --no-verify  Skip installation verification
#   --help       Show this help message
#
# Environment variables:
#   INSTALL_DIR  Installation directory (default: /usr/local/bin)
#
# Exit codes:
#   0 - Success
#   1 - Error (missing dependencies, download failure, etc.)
#   2 - User cancelled (declined upgrade confirmation)

set -e

REPO="marekbrze/dopadone"
PROJECT_NAME="dopadone"
BINARY_NAME="dopa"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

DRY_RUN=false
AUTO_YES=false
VERIFY=true

usage() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS]

Install dopa CLI tool.

Options:
    --dry-run     Simulate installation without actually downloading or installing
    --yes         Skip confirmation prompts (useful for automation)
    --no-verify   Skip installation verification (dopa version check)
    --help        Show this help message

Environment Variables:
    INSTALL_DIR   Installation directory (default: /usr/local/bin)

Examples:
    # Standard installation
    curl -sSL https://raw.githubusercontent.com/marekbrze/dopadone/main/scripts/install.sh | bash

    # Dry run to see what would be installed
    curl -sSL https://raw.githubusercontent.com/marekbrze/dopadone/main/scripts/install.sh | bash -s -- --dry-run

    # Install to custom directory
    INSTALL_DIR=~/.local/bin curl -sSL https://raw.githubusercontent.com/marekbrze/dopadone/main/scripts/install.sh | bash

    # Unattended installation (skip prompts)
    curl -sSL https://raw.githubusercontent.com/marekbrze/dopadone/main/scripts/install.sh | bash -s -- --yes
EOF
}

check_dependencies() {
    local missing=""
    
    if ! command -v curl > /dev/null 2>&1; then
        missing="$missing curl"
    fi
    
    if ! command -v tar > /dev/null 2>&1; then
        missing="$missing tar"
    fi
    
    if [ -n "$missing" ]; then
        echo "Error: Missing required dependencies:$missing" >&2
        echo "" >&2
        echo "Please install the missing tools:" >&2
        
        for dep in $missing; do
            case "$dep" in
                curl)
                    echo "  - curl: Usually available via package manager (apt, brew, etc.)" >&2
                    ;;
                tar)
                    echo "  - tar: Usually pre-installed on most systems" >&2
                    ;;
            esac
        done
        
        exit 1
    fi
}

detect_platform() {
    local os
    local arch
    
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)
    
    case "$arch" in
        x86_64|amd64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        *)
            echo "Error: Unsupported architecture: $arch" >&2
            exit 1
            ;;
    esac
    
    case "$os" in
        darwin)
            os="darwin"
            ;;
        linux)
            os="linux"
            ;;
        *)
            echo "Error: Unsupported OS: $os" >&2
            exit 1
            ;;
    esac
    
    echo "${os}-${arch}"
}

get_latest_version() {
    local version
    version=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$version" ]; then
        echo "Error: Could not determine latest version" >&2
        exit 1
    fi
    
    echo "$version"
}

get_download_url() {
    local version="$1"
    local platform="$2"
    local ext="tar.gz"
    
    if [ "$(echo "$platform" | cut -d'-' -f1)" = "windows" ]; then
        ext="zip"
    fi
    
    echo "https://github.com/${REPO}/releases/download/${version}/${PROJECT_NAME}-${platform}.${ext}"
}

download_and_extract() {
    local version="$1"
    local platform="$2"
    local ext="tar.gz"
    local tmp_dir
    local archive_name
    local extracted_binary
    
    if [ "$(echo "$platform" | cut -d'-' -f1)" = "windows" ]; then
        ext="zip"
    fi
    
    local url
    url=$(get_download_url "$version" "$platform")
    archive_name="${PROJECT_NAME}-${platform}.${ext}"
    
    tmp_dir=$(mktemp -d)
    trap 'rm -rf "$tmp_dir"' EXIT
    
    echo "Downloading ${BINARY_NAME} ${version} for ${platform}..." >&2
    
    if ! curl -sSL -f -o "${tmp_dir}/${archive_name}" "$url"; then
        echo "Error: Failed to download from $url" >&2
        exit 1
    fi
    
    cd "$tmp_dir"
    
    if [ "$ext" = "zip" ]; then
        if ! unzip -o "$archive_name"; then
            echo "Error: Failed to extract archive" >&2
            exit 1
        fi
        extracted_binary="${BINARY_NAME}.exe"
    else
        if ! tar xzf "$archive_name" 2>/dev/null; then
            echo "Error: Failed to extract archive" >&2
            exit 1
        fi
        extracted_binary="${BINARY_NAME}"
    fi
    
    if [ ! -f "$extracted_binary" ]; then
        echo "Error: Expected binary '$extracted_binary' not found in archive" >&2
        echo "Archive contents:" >&2
        ls -la >&2
        exit 1
    fi
    
    echo "$tmp_dir/$extracted_binary"
}

check_existing_installation() {
    local install_path="${INSTALL_DIR}/${BINARY_NAME}"
    
    if [ -f "$install_path" ]; then
        return 0
    fi
    return 1
}

get_current_version() {
    local install_path="${INSTALL_DIR}/${BINARY_NAME}"
    
    if [ -f "$install_path" ] && [ -x "$install_path" ]; then
        "$install_path" version 2>/dev/null | head -1 | awk '{print $NF}' || echo "unknown"
    else
        echo "unknown"
    fi
}

prompt_upgrade() {
    local current_version="$1"
    local new_version="$2"
    
    echo ""
    echo "Existing installation detected at ${INSTALL_DIR}/${BINARY_NAME}"
    
    if [ "$current_version" != "unknown" ]; then
        echo "Current version: $current_version"
    fi
    echo "New version: $new_version"
    echo ""
    
    if [ "$AUTO_YES" = true ]; then
        echo "Auto-confirming upgrade (--yes flag provided)"
        return 0
    fi
    
    if [ ! -t 0 ]; then
        if [ -e /dev/tty ]; then
            read -r -p "Replace existing installation? [y/N] " response < /dev/tty
        else
            echo "Error: Cannot prompt in non-interactive mode. Use --yes flag." >&2
            return 1
        fi
    else
        read -r -p "Replace existing installation? [y/N] " response
    fi
    case "$response" in
        [yY][eE][sS]|[yY])
            return 0
            ;;
        *)
            echo "Upgrade cancelled."
            return 1
            ;;
    esac
}

backup_existing() {
    local install_path="${INSTALL_DIR}/${BINARY_NAME}"
    local backup_path
    backup_path="${install_path}.backup.$(date +%Y%m%d%H%M%S)"
    
    if [ -f "$install_path" ]; then
        echo "Backing up existing binary to ${backup_path}..."
        cp "$install_path" "$backup_path"
    fi
}

install_binary() {
    local binary_path="$1"
    local install_path="${INSTALL_DIR}/${BINARY_NAME}"
    
    if [ ! -d "$INSTALL_DIR" ]; then
        echo "Creating installation directory: ${INSTALL_DIR}"
        mkdir -p "$INSTALL_DIR"
    fi
    
    echo "Installing to ${install_path}..."
    
    if [ -w "$INSTALL_DIR" ]; then
        mv "$binary_path" "$install_path"
        chmod +x "$install_path"
    else
        echo "sudo required for installation to ${INSTALL_DIR}"
        sudo mv "$binary_path" "$install_path"
        sudo chmod +x "$install_path"
    fi
    
    echo "Installed ${BINARY_NAME} to ${install_path}"
}

verify_installation() {
    local install_path="${INSTALL_DIR}/${BINARY_NAME}"
    
    echo ""
    echo "Verifying installation..."
    
    if [ ! -f "$install_path" ]; then
        echo "Error: Binary not found at ${install_path}" >&2
        return 1
    fi
    
    if [ ! -x "$install_path" ]; then
        echo "Error: Binary is not executable" >&2
        return 1
    fi
    
    if ! "$install_path" version 2>&1; then
        echo "Error: Failed to run 'dopa version'" >&2
        return 1
    fi
    
    echo ""
    echo "Verification successful!"
    return 0
}

print_dry_run_info() {
    local version="$1"
    local platform="$2"
    local url
    
    url=$(get_download_url "$version" "$platform")
    
    echo "=== DRY RUN MODE ==="
    echo ""
    echo "Platform:      ${platform}"
    echo "Version:       ${version}"
    echo "Download URL:  ${url}"
    echo "Install to:    ${INSTALL_DIR}/${BINARY_NAME}"
    
    if check_existing_installation; then
        local current_version
        current_version=$(get_current_version)
        echo ""
        echo "Note: Existing installation detected (version: ${current_version})"
        echo "      This would be an upgrade operation."
    fi
    
    echo ""
    echo "No files were downloaded or modified."
}

main() {
    local platform
    local version
    local binary_path
    
    while [ $# -gt 0 ]; do
        case "$1" in
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            --yes)
                AUTO_YES=true
                shift
                ;;
            --no-verify)
                VERIFY=false
                shift
                ;;
            --help|-h)
                usage
                exit 0
                ;;
            *)
                echo "Error: Unknown option: $1" >&2
                echo "Use --help for usage information" >&2
                exit 1
                ;;
        esac
    done
    
    echo "Installing ${BINARY_NAME}..."
    echo ""
    
    check_dependencies
    
    platform=$(detect_platform)
    version=$(get_latest_version)
    
    echo "Latest version: ${version}"
    echo "Platform: ${platform}"
    
    if [ "$DRY_RUN" = true ]; then
        print_dry_run_info "$version" "$platform"
        exit 0
    fi
    
    if check_existing_installation; then
        local current_version
        current_version=$(get_current_version)
        
        if ! prompt_upgrade "$current_version" "$version"; then
            exit 2
        fi
        
        backup_existing
    fi
    
    binary_path=$(download_and_extract "$version" "$platform")
    
    trap - EXIT
    install_binary "$binary_path"
    
    if [ "$VERIFY" = true ]; then
        if ! verify_installation; then
            exit 1
        fi
    fi
    
    echo ""
    echo "Installation complete!"
    
    if [ "$VERIFY" = false ]; then
        echo "Run '${BINARY_NAME} version' to verify manually."
    fi
}

main "$@"
