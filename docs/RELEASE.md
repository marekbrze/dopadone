# Release Process

This document describes how to release new versions of dopa.

## Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Incompatible API changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible

## Release Checklist

### 1. Prepare the Release

```bash
# Ensure you're on main and up to date
git checkout main
git pull origin main

# Run all tests
make test

# Run linting
make lint

# Test installation script locally
./scripts/install.sh --dry-run
./scripts/test-install.sh
```

### 2. Create a Version Tag

```bash
# For a patch release (bug fixes)
git tag -a v1.0.1 -m "Release v1.0.1"

# For a minor release (new features)
git tag -a v1.1.0 -m "Release v1.1.0"

# For a major release (breaking changes)
git tag -a v2.0.0 -m "Release v2.0.0"
```

### 3. Push the Tag

```bash
git push origin v1.0.1
```

### 4. Automated Release Process

Once the tag is pushed, GitHub Actions will automatically:

1. **Build binaries** for all platforms:
   - Linux (amd64)
   - macOS (amd64, arm64)
   - Windows (amd64)

2. **Inject version information** into each binary:
   - Version number from the git tag
   - Git commit SHA
   - Build timestamp

3. **Create distribution archives**:
   - `.tar.gz` for Linux and macOS
   - `.zip` for Windows

4. **Generate SHA256 checksums** for all archives

5. **Create a GitHub Release** with:
   - All binary archives
   - SHA256 checksum files
   - Auto-generated release notes from commits
   - Pre-release marker (if tag contains a hyphen)

The entire process is defined in `.github/workflows/release.yml` and typically completes in 5-10 minutes.

### 5. Verify the Release

1. Check the [Releases page](https://github.com/marekbrze/dopadone/releases)
2. Download and test the binaries
3. Verify the version command shows correct info:

```bash
./dopa version --all
```

## Testing the Release Workflow

You can test the release workflow without creating an actual release:

### Option 1: Manual Workflow Dispatch

1. Go to GitHub Actions → Release workflow
2. Click "Run workflow"
3. Enter a test tag (e.g., `v0.0.1-test`)
4. The workflow will run but won't create a release unless the tag exists

### Option 2: Create a Test Tag

```bash
# Create a test pre-release tag
git tag -a v0.0.1-test -m "Test release workflow"
git push origin v0.0.1-test

# After testing, delete the tag
git tag -d v0.0.1-test
git push origin :refs/tags/v0.0.1-test
```

## Pre-release Versions

For pre-release versions, use a hyphen in the tag:

```bash
# Alpha release
git tag -a v1.1.0-alpha.1 -m "Alpha release v1.1.0-alpha.1"

# Beta release
git tag -a v1.1.0-beta.1 -m "Beta release v1.1.0-beta.1"

# Release candidate
git tag -a v1.1.0-rc.1 -m "Release candidate v1.1.0-rc.1"
```

Pre-release versions will be automatically marked as "pre-release" on GitHub (the workflow detects the hyphen in the tag name).

## Manual Build (for testing)

```bash
# Build for current platform
make build

# Cross-compile for all platforms
make build-all

# Build with version info
VERSION=v1.0.0 make build-versioned
```

## Rollback

If a release has critical issues:

1. Go to [Releases](https://github.com/marekbrze/dopadone/releases)
2. Find the problematic release
3. Click "Delete" to remove it
4. Delete the tag locally and remotely:

```bash
git tag -d v1.0.1
git push origin :refs/tags/v1.0.1
```

---

# User Upgrade Instructions

## Checking Your Version

```bash
dopa version --all
```

## Automatic Upgrade (Recommended)

The `upgrade` command handles everything automatically:
- Downloads the latest release for your platform
- Replaces the binary
- Runs database migrations

```bash
dopa upgrade
```

To skip migrations during upgrade:

```bash
dopa upgrade --skip-migrate
```

## Manual Migrations

Database migrations are embedded in the binary. To run them manually:

```bash
# Check migration status
dopa migrate status

# Apply pending migrations
dopa migrate up

# Rollback last migration
dopa migrate down

# Reset database (rollback all, then apply all)
dopa migrate reset
```

## Manual Installation

### Option 1: Quick Install (Linux/macOS)

The automated installation script handles everything:

```bash
# Install latest version
curl -sSL https://raw.githubusercontent.com/marekbrze/dopadone/main/scripts/install.sh | sh
```

**Script Features**:
- Automatic platform detection (Linux, macOS Intel, macOS ARM)
- Dependency checking (curl, tar, unzip)
- Installation verification (`dopa version`)
- Upgrade support (detects and replaces existing installation)
- Dry-run mode for testing

**Advanced Usage**:

```bash
# Download and inspect before running
curl -sSL https://raw.githubusercontent.com/marekbrze/dopadone/main/scripts/install.sh -o install.sh
chmod +x install.sh

# Test what would be installed (dry-run)
./install.sh --dry-run

# Install with custom directory
INSTALL_DIR=$HOME/bin ./install.sh

# Upgrade without prompts (CI/automation)
./install.sh --yes

# Skip installation verification
./install.sh --no-verify

# Show help
./install.sh --help
```

**Testing the Script**:

The repository includes a comprehensive test suite for the installation script:

```bash
# Run all tests
./scripts/test-install.sh

# Test specific functionality
./scripts/test-install.sh dry_run
./scripts/test-install.sh platform_detection
```

### Option 2: Manual Download

1. Go to [Releases](https://github.com/marekbrze/dopadone/releases/latest)
2. Download the archive for your platform:
   - Linux: `dopa-linux-amd64.tar.gz`
   - macOS (Intel): `dopa-darwin-amd64.tar.gz`
   - macOS (Apple Silicon): `dopa-darwin-arm64.tar.gz`
   - Windows: `dopa-windows-amd64.zip`

3. Extract and install:

**Linux/macOS:**
```bash
tar xzf dopa-linux-amd64.tar.gz
sudo mv dopa /usr/local/bin/
```

**Windows (PowerShell):**
```powershell
Expand-Archive dopa-windows-amd64.zip
Move-Item dopa.exe -Destination "$env:USERPROFILE\bin\"
```

4. Run migrations:
```bash
dopa migrate up
```

5. Verify:
```bash
dopa version
```

## Specific Version Installation

To install a specific version:

```bash
VERSION=v1.0.0
curl -sSL https://github.com/marekbrze/dopadone/releases/download/${VERSION}/dopa-linux-amd64.tar.gz | tar xz
sudo mv dopa /usr/local/bin/
dopa migrate up
```
