# 🏷️ Tagging Guide - Git Tags vs Docker Tags

This document explains how Git tags are converted to Docker Hub tags.

## 📋 Quick Reference

| Git Tag | Docker Tags Created | Description |
|---------|---------------------|-------------|
| `v1.0.0` | `1.0.0`, `1.0`, `1`, `latest` | Stable release |
| `v2.3.5` | `2.3.5`, `2.3`, `2`, `latest` | Stable release |
| `v1.0.0-beta` | `1.0.0-beta` | Pre-release (no latest) |
| `v1.0.0-rc1` | `1.0.0-rc1` | Release candidate |
| `v1.0.0-alpha` | `1.0.0-alpha` | Alpha version |

## 🔄 Tag Conversion Rules

### ✅ The `v` Prefix is Removed

The GitHub Action automatically removes the `v` prefix when creating Docker tags:

```
Git Tag:      v1.2.3
             ↓ (v removed)
Docker Tag:   1.2.3
```

### 📊 Semantic Versioning Tags

For a stable release like `v1.2.3`, multiple Docker tags are created:

```bash
git tag v1.2.3
git push origin v1.2.3

# Creates Docker tags:
# ├─ hypolas/proxy-docker:1.2.3   (full version)
# ├─ hypolas/proxy-docker:1.2     (major.minor)
# ├─ hypolas/proxy-docker:1       (major)
# └─ hypolas/proxy-docker:latest  (latest stable)
```

### 🧪 Pre-release Tags

Pre-release tags (beta, rc, alpha) only create a single Docker tag:

```bash
git tag v1.0.0-beta
git push origin v1.0.0-beta

# Creates only:
# └─ hypolas/proxy-docker:1.0.0-beta
#
# Note: Does NOT update 'latest' tag
```

## 📝 Examples

### Example 1: First Stable Release

```bash
# Create tag
git tag -a v1.0.0 -m "First stable release"
git push origin v1.0.0

# GitHub Actions creates:
docker pull hypolas/proxy-docker:1.0.0
docker pull hypolas/proxy-docker:1.0
docker pull hypolas/proxy-docker:1
docker pull hypolas/proxy-docker:latest   # ← Points to 1.0.0
```

### Example 2: Patch Release

```bash
# Create patch tag
git tag -a v1.0.1 -m "Bug fix release"
git push origin v1.0.1

# GitHub Actions creates:
docker pull hypolas/proxy-docker:1.0.1
docker pull hypolas/proxy-docker:1.0       # ← Now points to 1.0.1
docker pull hypolas/proxy-docker:1         # ← Still points to 1.0.1
docker pull hypolas/proxy-docker:latest    # ← Now points to 1.0.1
```

### Example 3: Minor Version

```bash
# Create minor version tag
git tag -a v1.1.0 -m "New features"
git push origin v1.1.0

# GitHub Actions creates:
docker pull hypolas/proxy-docker:1.1.0
docker pull hypolas/proxy-docker:1.1       # ← Points to 1.1.0
docker pull hypolas/proxy-docker:1         # ← Now points to 1.1.0
docker pull hypolas/proxy-docker:latest    # ← Now points to 1.1.0
```

### Example 4: Major Version

```bash
# Create major version tag
git tag -a v2.0.0 -m "Major release with breaking changes"
git push origin v2.0.0

# GitHub Actions creates:
docker pull hypolas/proxy-docker:2.0.0
docker pull hypolas/proxy-docker:2.0       # ← Points to 2.0.0
docker pull hypolas/proxy-docker:2         # ← Points to 2.0.0
docker pull hypolas/proxy-docker:latest    # ← Now points to 2.0.0

# Note: Tag '1' still points to latest 1.x version
docker pull hypolas/proxy-docker:1         # ← Still points to 1.1.0
```

### Example 5: Beta Release

```bash
# Create beta tag
git tag -a v1.2.0-beta -m "Beta testing"
git push origin v1.2.0-beta

# GitHub Actions creates:
docker pull hypolas/proxy-docker:1.2.0-beta

# Note: Does NOT create 1.2, 1, or latest tags
# latest still points to previous stable (e.g., 1.1.0)
```

### Example 6: Release Candidate

```bash
# Create RC tag
git tag -a v2.0.0-rc1 -m "Release candidate 1"
git push origin v2.0.0-rc1

# GitHub Actions creates:
docker pull hypolas/proxy-docker:2.0.0-rc1

# After testing, create stable release:
git tag -a v2.0.0 -m "Stable release"
git push origin v2.0.0

# Now creates full set of tags:
docker pull hypolas/proxy-docker:2.0.0
docker pull hypolas/proxy-docker:2.0
docker pull hypolas/proxy-docker:2
docker pull hypolas/proxy-docker:latest    # ← Updates to 2.0.0
```

## 🎯 Best Practices

### 1. Use Semantic Versioning

Follow [SemVer](https://semver.org/) format: `vMAJOR.MINOR.PATCH`

```bash
v1.0.0   # Initial release
v1.0.1   # Bug fix
v1.1.0   # New feature (backward compatible)
v2.0.0   # Breaking changes
```

### 2. Test with Pre-releases

Before stable releases, use pre-release tags:

```bash
v1.0.0-alpha    # Early testing
v1.0.0-beta     # Feature complete, testing
v1.0.0-rc1      # Release candidate
v1.0.0          # Stable release
```

### 3. Annotated Tags (Recommended)

Use annotated tags with messages:

```bash
# ✅ Good - Annotated tag with message
git tag -a v1.0.0 -m "Release v1.0.0

- Added feature X
- Fixed bug Y
- Updated dependencies
"

# ❌ Avoid - Lightweight tag
git tag v1.0.0
```

### 4. Tag Naming Patterns

Valid patterns:
- ✅ `v1.0.0` - Stable release
- ✅ `v1.0.0-beta` - Beta release
- ✅ `v1.0.0-rc1` - Release candidate
- ✅ `v2.0.0-alpha.1` - Alpha with iteration

Invalid patterns:
- ❌ `1.0.0` - Missing `v` prefix (won't trigger workflow)
- ❌ `release-1.0.0` - Wrong format
- ❌ `v1.0` - Missing patch version

## 🔍 Verifying Tags

### Check Git tags

```bash
# List all tags
git tag -l

# Show tag details
git show v1.0.0

# List remote tags
git ls-remote --tags origin
```

### Check Docker Hub tags

```bash
# View on Docker Hub
https://hub.docker.com/r/hypolas/proxy-docker/tags

# Or via CLI
docker pull hypolas/proxy-docker:latest
docker images hypolas/proxy-docker
```

### Check GitHub Actions

```bash
# View workflow runs
https://github.com/hypolas/docker-proxy/actions

# Or via CLI
gh run list --workflow=docker-publish.yml
gh run view <run-id>
```

## 🐛 Troubleshooting

### Tag exists but Docker image not created

**Check:**
1. Tag format is correct (`v*.*.*`)
2. Workflow ran successfully in Actions tab
3. Docker Hub credentials are valid

### Wrong Docker tag created

**Cause:** Incorrect Git tag format

**Solution:**
```bash
# Delete wrong tag
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# Create correct tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Latest tag not updated

**Cause:** Pre-release tag or not on default branch

**Solution:** Only stable releases on the default branch update `latest`

```bash
# This updates 'latest'
git checkout main
git tag v1.0.0
git push origin v1.0.0

# This does NOT update 'latest'
git tag v1.0.0-beta
git push origin v1.0.0-beta
```

## 📚 Reference

### Workflow Configuration

The tag conversion is configured in `.github/workflows/docker-publish.yml`:

```yaml
tags: |
  # Git tag v1.2.3 → Docker tag 1.2.3 (removes 'v' prefix)
  type=semver,pattern={{version}}
  # Git tag v1.2.3 → Docker tag 1.2 (major.minor)
  type=semver,pattern={{major}}.{{minor}}
  # Git tag v1.2.3 → Docker tag 1 (major only)
  type=semver,pattern={{major}}
  # Always tag as 'latest' for stable releases
  type=raw,value=latest,enable={{is_default_branch}}
```

### Related Documentation

- [DOCKER_HUB_SETUP.md](DOCKER_HUB_SETUP.md) - Complete setup guide
- [README.md](.github/README.md) - Workflow documentation
- [Semantic Versioning](https://semver.org/) - Official SemVer specification

---

**Questions?** See [DOCKER_HUB_SETUP.md](DOCKER_HUB_SETUP.md) or open an issue.
