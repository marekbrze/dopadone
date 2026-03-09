# CI/CD Pipeline

This document describes the Continuous Integration and Continuous Delivery pipeline for Dopadone.

## Overview

Dopadone uses GitHub Actions for automated building, testing, and releasing. The CI/CD pipeline ensures that every release is:

- **Reproducible**: Same source code always produces the same binaries
- **Cross-platform**: Builds for Linux, macOS, and Windows
- **Versioned**: Each binary contains version metadata
- **Verified**: SHA256 checksums for all releases

## Workflows

### Release Workflow (`.github/workflows/release.yml`)

The release workflow is the main CI/CD pipeline that handles building and publishing releases.

#### Triggers

The workflow runs when:

1. **A version tag is pushed** (e.g., `v1.0.0`, `v1.1.0-beta.1`)
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **Manual dispatch** via GitHub Actions UI
   - Useful for testing the workflow
   - Requires entering a tag name manually

#### Build Matrix

The workflow builds binaries for all supported platforms in parallel:

| OS       | Architecture | Runner          | Binary Extension |
|----------|--------------|-----------------|------------------|
| Linux    | amd64        | ubuntu-latest   | (none)           |
| macOS    | amd64        | macos-13        | (none)           |
| macOS    | arm64        | macos-latest    | (none)           |
| Windows  | amd64        | windows-latest  | .exe             |

#### Build Process

For each platform, the workflow:

1. **Checks out the code** with full git history
2. **Sets up Go 1.21** with build caching
3. **Extracts version information**:
   - Version: from git tag (e.g., `v1.0.0`)
   - Git commit: current commit SHA
   - Build date: UTC timestamp

4. **Builds the binary** with:
   ```bash
   CGO_ENABLED=0 \
   GOOS=$os GOARCH=$arch \
   go build -trimpath \
     -ldflags "-s -w \
       -X github.com/marekbrze/dopadone/internal/version.Version=$version \
       -X github.com/marekbrze/dopadone/internal/version.GitCommit=$commit \
       -X github.com/marekbrze/dopadone/internal/version.BuildDate=$date" \
     -o dopa-$os-$arch \
     ./cmd/dopa
   ```

   Build flags explained:
   - `CGO_ENABLED=0`: Pure Go build (no C dependencies)
   - `-trimpath`: Remove local paths for reproducible builds
   - `-s -w`: Strip debug info and symbol table (smaller binaries)
   - `-X`: Inject version variables at compile time

5. **Creates distribution archives**:
   - **Unix** (Linux/macOS): `tar.gz` format
   - **Windows**: `zip` format

6. **Generates SHA256 checksums**:
   - **Linux**: Uses `sha256sum`
   - **macOS**: Uses `shasum -a 256`
   - **Windows**: Uses PowerShell `Get-FileHash`

7. **Uploads artifacts** for the release job

#### Release Creation

After all builds complete, the release job:

1. **Downloads all build artifacts** from the build matrix
2. **Prepares release assets** by collecting all archives and checksums
3. **Determines release type**:
   - **Stable release**: Tag without hyphen (e.g., `v1.0.0`)
   - **Pre-release**: Tag with hyphen (e.g., `v1.0.0-beta.1`)

4. **Creates GitHub Release** with:
   - Release title: "Release v1.0.0"
   - Auto-generated release notes from commits
   - All binary archives (`.tar.gz` and `.zip`)
   - All checksum files (`.sha256`)
   - Pre-release flag (if applicable)

### CI Workflow (`.github/workflows/ci.yml`)

The CI workflow runs automated quality checks on every push and pull request to ensure code quality before releases.

#### Triggers

The workflow runs on:

1. **Push to main branch** - Every commit to main triggers CI
2. **Pull requests to main** - All PRs must pass CI before merging

#### Jobs

The CI workflow runs three jobs in parallel:

##### 1. Test & Coverage Job

**Purpose**: Run the test suite with race detection and generate coverage reports

**Steps**:
1. Check out code with full git history
2. Set up Go 1.21 with module caching
3. Download dependencies (`go mod download`)
4. Run tests with race detector and coverage:
   ```bash
   go test ./... -v -race -coverprofile=coverage.out -covermode=atomic
   ```
5. Upload coverage report as artifact (30-day retention)
6. Check coverage threshold (currently 20%, target: 70%)

**Coverage Threshold**:
- Current minimum: 20% (avoids breaking builds initially)
- Current coverage: ~28.4%
- Target: Gradually increase to 70% as test coverage improves

##### 2. Lint Job

**Purpose**: Ensure code quality and consistency

**Steps**:
1. Check out code with full git history
2. Set up Go 1.21 with module caching
3. Download dependencies
4. Run `go vet ./...` for basic static analysis
5. Run golangci-lint with comprehensive checks

**Enabled Linters** (configured in `.golangci.yml`):
- `gofmt` - Code formatting
- `goimports` - Import statement organization
- `govet` - Go vet checks
- `errcheck` - Error handling verification
- `staticcheck` - Advanced static analysis
- `ineffassign` - Detect ineffective assignments
- `typecheck` - Type checking
- `gosimple` - Code simplification suggestions
- `goconst` - Detect repeated strings that could be constants
- `gocyclo` - Cyclomatic complexity (threshold: 20)
- `dupl` - Code clone detection (threshold: 150)

##### 3. Build Job

**Purpose**: Verify code compiles successfully

**Dependencies**: Runs only after Test & Lint jobs pass

**Steps**:
1. Check out code with full git history
2. Set up Go 1.21 with module caching
3. Download dependencies
4. Build all packages: `go build -v ./...`
5. Verify dependencies: `go mod verify`

#### Workflow Benefits

- **Quality Gate**: Catches issues early in development
- **Race Detection**: Identifies concurrency bugs
- **Coverage Tracking**: Monitors test coverage over time
- **Code Consistency**: Enforces coding standards
- **Build Verification**: Ensures code always compiles

#### Local Testing

Run the same checks locally:

```bash
# Run tests with race detector and coverage
go test ./... -v -race -coverprofile=coverage.out -covermode=atomic

# View coverage
go tool cover -func=coverage.out

# Run go vet
go vet ./...

# Run golangci-lint (install first: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
golangci-lint run

# Verify build
go build -v ./...
```

## Version Injection

Version information is injected at compile time using Go's `-ldflags -X` flag. This allows the binary to report its version without external files.

### Version Package

The version package lives at `internal/version/version.go`:

```go
package version

var (
    Version   = "dev"      // Set via -ldflags
    GitCommit = "unknown"  // Set via -ldflags
    BuildDate = "unknown"  // Set via -ldflags
)

func Full() string {
    return fmt.Sprintf("%s (commit: %s, built: %s)", 
        Version, GitCommit, BuildDate)
}
```

### Usage in Binary

The version information is accessible via:

```bash
# Short version
dopa version
# Output: v1.0.0

# Full version details
dopa version --all
# Output: v1.0.0 (commit: abc123def, built: 2026-03-08T10:30:00Z)
```

## Archive Naming Convention

Archives follow a consistent naming pattern:

```
dopa-{os}-{arch}.{ext}
```

Examples:
- `dopa-linux-amd64.tar.gz`
- `dopa-darwin-amd64.tar.gz`
- `dopa-darwin-arm64.tar.gz`
- `dopa-windows-amd64.zip`

Checksum files append `.sha256`:
- `dopa-linux-amd64.tar.gz.sha256`
- `dopa-windows-amd64.zip.sha256`

## Local Development

### Building Locally

The Makefile provides targets for local builds:

```bash
# Build for current platform
make build

# Build with version info
VERSION=v1.0.0 make build-versioned

# Cross-compile for all platforms
make build-all

# Create distribution archives
make dist
```

### Testing the Workflow Locally

While you can't run GitHub Actions locally, you can test the build commands:

```bash
# Test Linux build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -trimpath \
    -ldflags "-s -w -X github.com/marekbrze/dopadone/internal/version.Version=v0.0.1-test" \
    -o dopa-linux-amd64 \
    ./cmd/dopa

# Test archive creation
tar -czvf dopa-linux-amd64.tar.gz dopa-linux-amd64

# Test checksum generation
sha256sum dopa-linux-amd64.tar.gz > dopa-linux-amd64.tar.gz.sha256
```

## Troubleshooting

### Build Fails

**Problem**: Build fails with "package not found" errors

**Solution**: Ensure `go.mod` is up to date:
```bash
go mod tidy
go mod verify
```

### Version Not Injected

**Problem**: Binary shows `version: dev` instead of actual version

**Solution**: Verify ldflags path matches your module:
```bash
# Check module path in go.mod
head -1 go.mod

# Ensure ldflags uses correct path
-X github.com/marekbrze/dopadone/internal/version.Version=$version
```

### Checksum Mismatch

**Problem**: Downloaded binary checksum doesn't match published checksum

**Solution**: 
1. Re-download the binary (might be corrupted)
2. Verify you're using the correct checksum file for your platform
3. Use the correct checksum command for your OS:
   ```bash
   # Linux
   sha256sum -c dopa-linux-amd64.tar.gz.sha256
   
   # macOS
   shasum -a 256 -c dopa-darwin-amd64.tar.gz.sha256
   
   # Windows (PowerShell)
   Get-FileHash -Algorithm SHA256 dopa-windows-amd64.zip
   ```

### Release Not Created

**Problem**: Tag pushed but no release appears

**Solution**:
1. Check GitHub Actions tab for workflow status
2. Verify tag follows `v*` pattern (e.g., `v1.0.0`, not `1.0.0`)
3. Check workflow logs for errors
4. Ensure `GITHUB_TOKEN` has correct permissions (set in workflow)

## Security Considerations

### Reproducible Builds

The workflow uses:
- **Fixed Go version**: `1.21` (not `latest`)
- **Specific runner versions**: Pinned to major versions
- **`-trimpath` flag**: Removes local paths from binaries
- **`CGO_ENABLED=0`**: Pure Go build (no system dependencies)

This ensures the same source code always produces identical binaries.

### Checksum Verification

All releases include SHA256 checksums. Users should verify downloads:

```bash
# Linux
sha256sum -c dopa-linux-amd64.tar.gz.sha256

# macOS
shasum -a 256 -c dopa-darwin-amd64.tar.gz.sha256

# Windows (PowerShell)
$expectedHash = (Get-Content dopa-windows-amd64.zip.sha256).Split()[0]
$actualHash = (Get-FileHash dopa-windows-amd64.zip).Hash.ToLower()
if ($expectedHash -eq $actualHash) { "OK" } else { "MISMATCH" }
```

### Artifact Retention

Build artifacts are retained for only 1 day to reduce storage costs. The GitHub Release stores the final assets permanently.

## Future Improvements

Potential enhancements to the CI/CD pipeline:

- [x] Add automated testing before release (implemented in CI workflow)
- [ ] Sign binaries with GPG/code signing certificates
- [ ] Generate SBOM (Software Bill of Materials)
- [ ] Add Docker image builds
- [ ] Implement automated rollback on critical issues
- [ ] Add performance benchmarking in CI
- [ ] Generate and publish API documentation
- [ ] Increase coverage threshold to 70%
- [ ] Add code coverage visualization (e.g., codecov integration)
- [ ] Add security vulnerability scanning

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Build Modes](https://pkg.go.dev/cmd/go#hdr-Build_modes)
- [Semantic Versioning](https://semver.org/)
- [Reproducible Builds](https://reproducible-builds.org/)
