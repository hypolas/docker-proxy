# 🤖 GitHub Actions Workflows

This directory contains automated workflows for the docker-proxy project.

## 📋 Available Workflows

### 🐳 Docker Publish (`docker-publish.yml`)

Automatically builds and publishes multi-platform Docker images to Docker Hub when you create version tags.

**Triggers:**
- ✅ Git tags matching `v*.*.*` (e.g., `v1.0.0`, `v2.3.1`)
- ✅ Pre-release tags `v*.*.*-*` (e.g., `v1.0.0-beta`, `v1.0.0-rc1`)
- ❌ Branch pushes (disabled)
- ❌ Pull requests (disabled)

**Platforms built:**
- `linux/amd64` (x86_64)
- `linux/arm64` (ARM 64-bit)
- `linux/arm/v7` (ARM 32-bit)

**Docker tags created:**

For tag `v1.2.3`:
```
hypolas/proxy-docker:1.2.3
hypolas/proxy-docker:1.2
hypolas/proxy-docker:1
hypolas/proxy-docker:latest
```

**Setup required:** See [DOCKER_HUB_SETUP.md](DOCKER_HUB_SETUP.md)

## 🚀 Quick Start

### 1. Configure secrets (one-time setup)

```bash
# Go to repository settings
https://github.com/hypolas/docker-proxy/settings/secrets/actions

# Add these secrets:
# - DOCKERHUB_USERNAME: hypolas
# - DOCKERHUB_TOKEN: <your-docker-hub-access-token>
```

See detailed instructions in [DOCKER_HUB_SETUP.md](DOCKER_HUB_SETUP.md)

### 2. Create a release

```bash
# Tag and push
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Watch the build
https://github.com/hypolas/docker-proxy/actions
```

### 3. Verify on Docker Hub

```bash
# Check Docker Hub
https://hub.docker.com/r/hypolas/proxy-docker/tags

# Pull and test
docker pull hypolas/proxy-docker:latest
docker run -d hypolas/proxy-docker:latest
```

## 📚 Documentation

- **[DOCKER_HUB_SETUP.md](DOCKER_HUB_SETUP.md)** - Complete setup guide
- **[Workflow file](workflows/docker-publish.yml)** - Workflow source code

## 🔍 Workflow Details

### Build Process

1. **Checkout** - Clone repository
2. **Setup QEMU** - Enable multi-platform emulation
3. **Setup Buildx** - Configure Docker buildx
4. **Extract metadata** - Generate tags and labels from Git tag
5. **Login** - Authenticate with Docker Hub
6. **Build & Push** - Build for all platforms and push to registry
7. **Update description** - Sync README.md to Docker Hub
8. **Summary** - Generate build report

### Build Time

- **First build**: ~10-15 minutes (no cache)
- **Subsequent builds**: ~5-7 minutes (with cache)

### Build Cache

The workflow uses GitHub Actions cache to speed up builds:
- Cache is stored in GitHub Actions
- Shared across workflow runs
- Automatically cleaned up after 7 days of inactivity

## 🐛 Troubleshooting

### Workflow didn't trigger

**Problem:** Pushed code but no workflow run appeared

**Solution:** The workflow only triggers on **tags**, not branch pushes.

```bash
# ❌ This won't trigger the workflow
git push origin main

# ✅ This will trigger it
git tag v1.0.0
git push origin v1.0.0
```

### Authentication failed

**Problem:** `Error: unauthorized: authentication required`

**Solutions:**
1. Check secrets are set: `Settings → Secrets → Actions`
2. Verify `DOCKERHUB_USERNAME` is `hypolas`
3. Regenerate Docker Hub access token
4. Ensure token has `Read, Write, Delete` permissions

### Build failed on specific platform

**Problem:** Build succeeds on `amd64` but fails on `arm64`

**Solutions:**
1. Check if dependencies support ARM
2. Test locally with QEMU:
   ```bash
   docker buildx build --platform linux/arm64 .
   ```
3. Update Dockerfile for cross-platform compatibility

## 🔐 Security

### Secrets Management

- ✅ Secrets are encrypted by GitHub
- ✅ Secrets are never exposed in logs
- ✅ Use Docker Hub access tokens (not passwords)
- ✅ Rotate tokens every 6-12 months

### Permissions

The workflow has minimal permissions:
```yaml
permissions:
  contents: read    # Read repository code
  packages: write   # Push Docker images
```

## 📊 Monitoring

### View workflow runs:

```bash
# CLI
gh run list --workflow=docker-publish.yml
gh run view <run-id> --log

# Web UI
https://github.com/hypolas/docker-proxy/actions
```

### Check Docker Hub stats:

```bash
# Pulls
https://hub.docker.com/r/hypolas/proxy-docker

# Tags
https://hub.docker.com/r/hypolas/proxy-docker/tags
```

## 🎯 Best Practices

1. **Semantic Versioning**: Use `vMAJOR.MINOR.PATCH` format
2. **Changelog**: Document changes in each release
3. **Testing**: Test images before tagging
4. **Pre-releases**: Use `-beta`, `-rc1` for testing
5. **Release Notes**: Create GitHub releases with notes

### Example release workflow:

```bash
# 1. Test locally
docker build -t test .
docker run -d test
# ... test the image ...

# 2. Update version/changelog
vim CHANGELOG.md
git commit -am "Prepare v1.0.0"

# 3. Create tag with detailed message
git tag -a v1.0.0 -m "Release v1.0.0

Changes:
- Added feature X
- Fixed bug Y
- Updated dependencies
"

# 4. Push tag (triggers workflow)
git push origin v1.0.0

# 5. Wait for workflow completion
gh run watch

# 6. Create GitHub release
gh release create v1.0.0 \
  --title "v1.0.0" \
  --notes-file CHANGELOG.md

# 7. Announce
# - Update documentation
# - Post on social media
# - Notify users
```

## 🔄 Future Workflows

Potential additions:
- 🧪 **Tests** - Run unit/integration tests on PR
- 🔍 **Linting** - Go linting and formatting checks
- 🔒 **Security** - Vulnerability scanning with Trivy
- 📦 **Release** - Automated GitHub releases
- 📊 **Benchmarks** - Performance regression testing

---

**Need help?** Check [DOCKER_HUB_SETUP.md](DOCKER_HUB_SETUP.md) or open an issue.
