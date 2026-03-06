#!/bin/bash
# Generate release notes from git commits
# Usage: ./scripts/generate-changelog.sh <version>

set -e

VERSION="${1:-}"
REPO="${GITHUB_REPOSITORY:-example/dopadone}"

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>" >&2
    exit 1
fi

# Get previous tag
PREV_TAG=$(git describe --tags --abbrev=0 HEAD^ 2>/dev/null || echo "")

echo "## Release ${VERSION}"
echo ""

if [ -n "$PREV_TAG" ]; then
    echo "### Changes since ${PREV_TAG}"
    echo ""
    
    # Generate commit list
    git log --pretty=format:"- %s (%h)" "$PREV_TAG"..HEAD 2>/dev/null | head -50
    echo ""
    echo ""
fi

echo "### Installation"
echo ""
echo "Download the appropriate binary for your platform:"
echo ""
echo "| Platform | Architecture | Download |"
echo "|----------|--------------|----------|"
echo "| Linux | amd64 | \`dopa-linux-amd64.tar.gz\` |"
echo "| macOS | amd64 | \`dopa-darwin-amd64.tar.gz\` |"
echo "| macOS | arm64 (M1/M2) | \`dopa-darwin-arm64.tar.gz\` |"
echo "| Windows | amd64 | \`dopa-darwin-amd64.zip\` |"
echo ""
echo "### Quick Install (Linux/macOS)"
echo ""
echo '```bash'
echo "curl -sSL https://github.com/${REPO}/releases/download/${VERSION}/dopa-linux-amd64.tar.gz | tar xz"
echo "sudo mv dopa /usr/local/bin/"
echo '```'
echo ""
echo "### Verification"
echo ""
echo "After installation, verify with:"
echo '```bash'
echo "dopa version --all"
echo '```'
