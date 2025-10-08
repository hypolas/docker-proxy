# 🚀 Release Guide

This guide explains how releases are automated when you push a Git tag.

## 📦 What Happens on Tag Push

When you create and push a tag (e.g., `v1.0.0`), the workflow **automatically**:

1. ✅ **Builds Docker images** for 3 platforms (amd64, arm64, armv7)
2. ✅ **Pushes to Docker Hub** with multiple tags
3. ✅ **Compiles binaries** for 6 platforms
4. ✅ **Creates GitHub Release** with:
   - Release notes
   - Downloadable binaries
   - SHA256 checksums
   - Links to documentation

## 🏗️ Artifacts Created

### Docker Images (Published to Docker Hub)

**Platforms:**
- `linux/amd64` (x86_64)
- `linux/arm64` (ARM 64-bit - Raspberry Pi 4, Apple M1/M2, AWS Graviton)
- `linux/arm/v7` (ARM 32-bit - Raspberry Pi 2/3)

**Tags for `v1.2.3`:**
```
hypolas/dockershield:1.2.3
hypolas/dockershield:1.2
hypolas/dockershield:1
hypolas/dockershield:latest
```

### Binary Releases (Published to GitHub Releases)

**Platforms:**
- `dockershield-linux-amd64` - Linux x86_64
- `dockershield-linux-arm64` - Linux ARM 64-bit
- `dockershield-linux-armv7` - Linux ARM 32-bit (v7)
- `dockershield-darwin-amd64` - macOS Intel
- `dockershield-darwin-arm64` - macOS Apple Silicon (M1/M2/M3)
- `dockershield-windows-amd64.exe` - Windows x86_64
- `checksums.txt` - SHA256 checksums for all binaries

## 🎯 Creating a Release

### Step 1: Prepare the Release

```bash
# Make sure you're on main branch with latest changes
git checkout main
git pull origin main

# Run tests (if you have them)
go test ./...

# Build locally to verify
go build -o dockershield ./cmd/dockershield
./dockershield --version
```

### Step 2: Create and Push Tag

```bash
# Create annotated tag with detailed message
git tag -a v1.0.0 -m "Release v1.0.0 - Initial Stable Release

## 🎉 Features
- Secure Docker socket proxy with ACL
- Advanced regex-based filtering
- Multi-platform support
- CI/CD optimized

## 🔒 Security
- Default socket protection
- Self-protection mechanisms
- Zero-trust architecture

## 📚 Documentation
- Complete environment variable reference
- CI/CD integration examples
- Security best practices
"

# Push the tag (triggers the workflow)
git push origin v1.0.0
```

### Step 3: Monitor the Workflow

```bash
# View workflow runs
gh run list --workflow=docker-publish.yml

# Watch latest run
gh run watch

# Or view in browser
https://github.com/hypolas/dockershield/actions
```

### Step 4: Verify the Release

After ~10-15 minutes:

**Docker Hub:**
```bash
docker pull hypolas/dockershield:1.0.0
docker run --rm hypolas/dockershield:1.0.0 --version
```

**GitHub Release:**
```bash
# View in browser
https://github.com/hypolas/dockershield/releases/tag/v1.0.0

# Or via CLI
gh release view v1.0.0

# Download binary
gh release download v1.0.0 -p "dockershield-linux-amd64"
chmod +x dockershield-linux-amd64
./dockershield-linux-amd64 --version
```

## 📝 Release Types

### Stable Release

```bash
git tag -a v1.0.0 -m "Stable release v1.0.0"
```

- ✅ Creates Docker tags: `1.0.0`, `1.0`, `1`, `latest`
- ✅ Creates GitHub Release (not draft, not prerelease)
- ✅ Includes all binaries

### Pre-release (Beta, RC, Alpha)

```bash
git tag -a v1.0.0-beta -m "Beta release v1.0.0-beta"
# or
git tag -a v1.0.0-rc1 -m "Release candidate 1"
```

- ✅ Creates Docker tag: `1.0.0-beta` or `1.0.0-rc1`
- ❌ Does NOT update `latest` tag
- ✅ Creates GitHub Release marked as **prerelease**
- ✅ Includes all binaries

## 🔧 Binary Build Details

### Build Configuration

Binaries are built with:
```bash
-ldflags="-s -w
  -X main.Version=$TAG
  -X main.BuildDate=$DATE
  -X main.GitCommit=$SHA"
```

**Flags:**
- `-s` - Strip symbol table
- `-w` - Strip DWARF debug info
- Result: Smaller binaries (~30% reduction)

### Version Information

Binaries include embedded version info:
```go
// In your code
var (
    Version   string // Injected at build time
    BuildDate string // Injected at build time
    GitCommit string // Injected at build time
)
```

Users can check version:
```bash
./dockershield --version
# Output:
# dockershield v1.0.0
# Built: 2025-01-10T12:34:56Z
# Commit: a1b2c3d
```

### Checksums

SHA256 checksums are automatically generated:
```bash
# Download and verify
wget https://github.com/hypolas/dockershield/releases/download/v1.0.0/dockershield-linux-amd64
wget https://github.com/hypolas/dockershield/releases/download/v1.0.0/checksums.txt

# Verify
sha256sum -c checksums.txt --ignore-missing
# Should output:
# dockershield-linux-amd64: OK
```

## 📊 Release Notes

### Automatic Content

The workflow automatically includes:
- Installation instructions (Docker + Binary)
- Supported platforms list
- Links to documentation
- Full changelog link
- Contact information

### Customizing Release Notes

To add custom release notes, edit the tag message:

```bash
git tag -a v1.0.0 -m "Release v1.0.0

## 🎉 New Features
- Feature X: Description
- Feature Y: Description

## 🐛 Bug Fixes
- Fixed issue #123
- Fixed memory leak in proxy handler

## ⚠️ Breaking Changes
- Changed default port from 2375 to 2376
- Removed deprecated ENABLE_OLD_API flag

## 📚 Documentation
- Updated README with new examples
- Added troubleshooting guide

## 🙏 Contributors
- @user1 - Fixed bug #123
- @user2 - Added feature X
"
```

The workflow will keep this message and append standard installation instructions.

## 🐛 Troubleshooting

### Binary build failed

**Error:** `go: cannot find main module`

**Solution:** Ensure you have `go.mod` in repository root.

### No binaries in release

**Error:** `fail_on_unmatched_files: true` failed

**Solution:** Check that all binaries were built successfully. Look at "Build binaries" step logs.

### Checksums don't match

**Cause:** Binary was modified after checksum generation.

**Solution:** Don't manually modify binaries. Re-run workflow if needed.

### Release is marked as draft

**Cause:** You manually created a draft release with same tag.

**Solution:** Delete the draft release, then re-push the tag:
```bash
gh release delete v1.0.0 --yes
git push origin v1.0.0 --force
```

## 📈 Release Workflow Timeline

Typical workflow execution time:

```
Step                          Time
────────────────────────────  ──────
Checkout & Setup              1 min
Build Docker (3 platforms)    8 min
Push to Docker Hub            1 min
Build Binaries (6 platforms)  2 min
Create GitHub Release         1 min
────────────────────────────  ──────
Total                         ~13 min
```

## 🎯 Best Practices

### 1. Semantic Versioning

```bash
v1.0.0  # Initial release
v1.0.1  # Bug fix (patch)
v1.1.0  # New feature (minor)
v2.0.0  # Breaking change (major)
```

### 2. Test Before Release

```bash
# Run tests
go test ./...

# Build and test locally
go build ./cmd/dockershield
./dockershield

# Test Docker build
docker build -t test .
docker run --rm test
```

### 3. Pre-release for Testing

```bash
# Release beta for testing
git tag v1.0.0-beta1
git push origin v1.0.0-beta1

# After testing passes
git tag v1.0.0
git push origin v1.0.0
```

### 4. Detailed Tag Messages

Include:
- Summary of changes
- New features
- Bug fixes
- Breaking changes
- Migration guide (if needed)

### 5. Keep Changelog

Maintain a `CHANGELOG.md`:
```markdown
# Changelog

## [1.0.0] - 2025-01-10
### Added
- Initial stable release
- Docker socket proxy with ACL
- Advanced filtering

### Changed
- Improved performance

### Fixed
- Fixed bug #123
```

## 🔄 Updating a Release

### Option 1: Create new patch version (recommended)

```bash
# Fix the issue
git commit -m "Fix bug"
git push origin main

# Create new version
git tag v1.0.1
git push origin v1.0.1
```

### Option 2: Delete and recreate (not recommended)

```bash
# Delete release
gh release delete v1.0.0 --yes

# Delete tag locally and remotely
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# Fix and recreate
git tag v1.0.0
git push origin v1.0.0
```

## 📞 Support

- **Workflow issues:** See [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
- **Tag format:** See [TAG_FORMAT.md](TAG_FORMAT.md)
- **Docker Hub setup:** See [DOCKER_HUB_SETUP.md](DOCKER_HUB_SETUP.md)

---

**Ready to create your first release?** 🚀

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```
