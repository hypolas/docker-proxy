# 🏷️ Tag Format Guide

## ✅ Recommended Tag Formats

### Standard Releases (Preferred)

```bash
v1.0.0          # ✅ Perfect - Standard semver
v2.3.5          # ✅ Perfect - Standard semver
v1.0.0-beta     # ✅ Good - Pre-release
v1.0.0-rc1      # ✅ Good - Release candidate
v1.0.0-alpha    # ✅ Good - Alpha version
```

**Docker tags created:**
- `v1.0.0` → `1.0.0`, `1.0`, `1`, `latest`
- `v1.0.0-beta` → `1.0.0-beta`

### Non-Standard Formats (Supported but not recommended)

```bash
v0.1.0-beta.01  # ⚠️ Works but not standard semver
v1.0.0-beta.1   # ⚠️ Works but not standard semver
v1.0.0-rc.2     # ⚠️ Works but not standard semver
```

**Docker tags created:**
- `v0.1.0-beta.01` → `0.1.0-beta.01` (single tag, via fallback)

## ❌ Invalid Tag Formats

These will **NOT trigger** the workflow:

```bash
1.0.0           # ❌ Missing 'v' prefix
v1.0            # ❌ Missing patch version
release-1.0.0   # ❌ Wrong prefix
latest          # ❌ Not a version
test            # ❌ Not a version
```

## 🎯 Best Practices

### 1. Use Standard Semver

For maximum compatibility, use standard semantic versioning:

```bash
# ✅ Recommended
git tag -a v1.0.0 -m "Release v1.0.0"

# ⚠️ Avoid
git tag -a v1.0.0-beta.01 -m "Beta 01"

# ✅ Better
git tag -a v1.0.0-beta1 -m "Beta 1"
# or
git tag -a v1.0.0-beta-01 -m "Beta 01"
```

### 2. Pre-release Versions

Use hyphens instead of dots in pre-release identifiers:

```bash
# ✅ Good
v1.0.0-beta1
v1.0.0-beta-1
v1.0.0-rc1
v1.0.0-alpha

# ⚠️ Works but not standard
v1.0.0-beta.1
v1.0.0-rc.2
```

### 3. Version Numbering

Follow semantic versioning rules:

- **MAJOR**: Breaking changes (v1.0.0 → v2.0.0)
- **MINOR**: New features, backward compatible (v1.0.0 → v1.1.0)
- **PATCH**: Bug fixes (v1.0.0 → v1.0.1)

```bash
# Initial release
git tag v1.0.0

# Bug fix
git tag v1.0.1

# New feature
git tag v1.1.0

# Breaking change
git tag v2.0.0
```

## 📊 Tag Conversion Examples

### Standard Tags

| Git Tag | Docker Tags |
|---------|-------------|
| `v1.0.0` | `1.0.0`, `1.0`, `1`, `latest` |
| `v1.2.3` | `1.2.3`, `1.2`, `1`, `latest` |
| `v2.0.0` | `2.0.0`, `2.0`, `2`, `latest` |

### Pre-release Tags (Standard)

| Git Tag | Docker Tags |
|---------|-------------|
| `v1.0.0-beta` | `1.0.0-beta` |
| `v1.0.0-beta1` | `1.0.0-beta1` |
| `v1.0.0-rc1` | `1.0.0-rc1` |
| `v1.0.0-alpha` | `1.0.0-alpha` |

### Pre-release Tags (Non-standard)

| Git Tag | Docker Tags |
|---------|-------------|
| `v0.1.0-beta.01` | `0.1.0-beta.01` (fallback) |
| `v1.0.0-rc.2` | `1.0.0-rc.2` (fallback) |

**Note:** Non-standard formats only create a single Docker tag (no major, major.minor, or latest).

## 🔄 Migration from Non-Standard Tags

If you've been using `v0.1.0-beta.01` format and want to switch:

```bash
# Current (non-standard)
v0.1.0-beta.01
v0.1.0-beta.02
v0.1.0-beta.03

# Recommended (standard)
v0.1.0-beta1
v0.1.0-beta2
v0.1.0-beta3

# Or with hyphens
v0.1.0-beta-01
v0.1.0-beta-02
v0.1.0-beta-03
```

**For your next tag:**

```bash
# Instead of
git tag v0.1.0-beta.04  # ⚠️

# Use
git tag v0.1.0-beta4    # ✅
# or
git tag v0.1.0-beta-04  # ✅
```

## 🧪 Testing Your Tag

Before pushing, verify the tag format:

```bash
# Create tag locally
git tag v1.0.0

# Check if it matches semver
# Should print version components
echo "v1.0.0" | grep -P '^v\d+\.\d+\.\d+(-[a-zA-Z0-9-]+)?$'

# If exit code is 0, format is valid
echo $?  # Should be 0
```

## 📚 Semver Resources

- [Semantic Versioning 2.0.0](https://semver.org/)
- [Regex for Semver](https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string)
- [NPM Semver Parser](https://github.com/npm/node-semver)

## 🚀 Quick Reference

### Creating a Tag

```bash
# 1. Standard release
git tag -a v1.0.0 -m "Release v1.0.0

- Added feature X
- Fixed bug Y
"

# 2. Push to trigger workflow
git push origin v1.0.0

# 3. Check workflow
# https://github.com/hypolas/docker-proxy/actions

# 4. Verify Docker Hub
# https://hub.docker.com/r/hypolas/proxy-docker/tags
```

### Deleting a Wrong Tag

```bash
# Delete locally
git tag -d v1.0.0

# Delete on remote
git push origin :refs/tags/v1.0.0

# Create correct tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## ❓ FAQ

### Q: Why does my tag `v0.1.0-beta.01` only create one Docker tag?

**A:** The dot in `-beta.01` makes it non-standard semver. The workflow uses a fallback that creates a single tag. Use `-beta01` or `-beta-01` instead for full tag generation.

### Q: Can I use `v1.0` instead of `v1.0.0`?

**A:** No, all three version components (major.minor.patch) are required.

### Q: What about build metadata like `v1.0.0+20130313144700`?

**A:** Semver supports it, but Docker tags don't allow `+`. Avoid build metadata in Git tags if you want Docker tags.

### Q: Should I use `v` prefix?

**A:** **Yes!** The workflow requires `v` prefix. Tags without `v` won't trigger the build.

---

**Need help?** See [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
