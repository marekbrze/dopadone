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

Once the tag is pushed, GitHub Actions will:

1. Build binaries for all platforms:
   - Linux (amd64)
   - macOS (amd64, arm64)
   - Windows (amd64)

2. Generate release notes from commits

3. Create a GitHub Release with:
   - All binary archives
   - Auto-generated changelog
   - Installation instructions

### 5. Verify the Release

1. Check the [Releases page](https://github.com/example/dopa/releases)
2. Download and test the binaries
3. Verify the version command shows correct info:

```bash
./dopa version --all
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

Pre-release versions will be marked as "pre-release" on GitHub.

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

1. Go to [Releases](https://github.com/example/dopa/releases)
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

```bash
curl -sSL https://raw.githubusercontent.com/example/dopa/main/scripts/install.sh | sh
```

### Option 2: Manual Download

1. Go to [Releases](https://github.com/example/dopa/releases/latest)
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
curl -sSL https://github.com/example/dopa/releases/download/${VERSION}/dopa-linux-amd64.tar.gz | tar xz
sudo mv dopa /usr/local/bin/
dopa migrate up
```
