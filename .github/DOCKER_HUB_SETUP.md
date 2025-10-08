# 🐳 Docker Hub Publishing Setup

This document explains how to configure GitHub Actions to automatically publish Docker images to Docker Hub.

## 📋 Prerequisites

1. **Docker Hub account**: https://hub.docker.com
2. **GitHub repository**: https://github.com/hypolas/dockershield
3. **Repository admin access** to configure secrets

## 🔑 Step 1: Create Docker Hub Access Token

1. Log in to Docker Hub: https://hub.docker.com
2. Click on your **username** → **Account Settings**
3. Go to **Security** tab
4. Click **New Access Token**
5. Configure:
   - **Description**: `GitHub Actions - dockershield`
   - **Access permissions**: `Read, Write, Delete`
6. Click **Generate**
7. **⚠️ IMPORTANT**: Copy the token immediately (you won't see it again!)

## 🔐 Step 2: Configure GitHub Secrets

Go to your GitHub repository:
```
https://github.com/hypolas/dockershield/settings/secrets/actions
```

### Add the following secrets:

#### 1. `DOCKERHUB_USERNAME`
- **Value**: `hypolas` (your Docker Hub username)
- Click **New repository secret**
- Name: `DOCKERHUB_USERNAME`
- Secret: `hypolas`
- Click **Add secret**

#### 2. `DOCKERHUB_TOKEN`
- **Value**: The access token you just created
- Click **New repository secret**
- Name: `DOCKERHUB_TOKEN`
- Secret: `<paste-your-token-here>`
- Click **Add secret**

## 🏷️ Step 3: Create and Push Tags

The workflow **only runs on tags**, not on branch pushes.

### Create a new release:

```bash
# Make sure you're on the main branch with latest changes
git checkout main
git pull origin main

# Create and push a version tag
git tag -a v1.0.0 -m "Release v1.0.0 - Initial stable release"
git push origin v1.0.0
```

### Tag naming conventions:

| Tag Format | Example | Description |
|------------|---------|-------------|
| `v*.*.*` | `v1.0.0` | Stable release |
| `v*.*.*-*` | `v1.0.0-beta` | Pre-release (beta, rc, alpha) |

### Docker Hub tags generated:

**⚠️ Important:** The `v` prefix is automatically removed from Docker tags.

| Git Tag | Docker Tags Created |
|---------|---------------------|
| `v1.2.3` | `1.2.3`, `1.2`, `1`, `latest` |
| `v2.0.0` | `2.0.0`, `2.0`, `2`, `latest` |
| `v1.5.0-beta` | `1.5.0-beta` (pre-release, no latest) |

**Example for Git tag `v1.2.3`:**
- ✅ `hypolas/dockershield:1.2.3` ← Full version (v removed)
- ✅ `hypolas/dockershield:1.2` ← Major.minor
- ✅ `hypolas/dockershield:1` ← Major only
- ✅ `hypolas/dockershield:latest` ← If stable release

## 🚀 Step 4: Verify the Workflow

After pushing a tag:

1. Go to **Actions** tab: https://github.com/hypolas/dockershield/actions
2. You should see **"Build and Publish Docker Image"** running
3. Wait for completion (~5-10 minutes for multi-platform build)
4. Check Docker Hub: https://hub.docker.com/r/hypolas/dockershield/tags

## 🔍 Step 5: Test the Published Image

```bash
# Pull the image
docker pull hypolas/dockershield:latest

# Or specific version
docker pull hypolas/dockershield:1.0.0

# Run it
docker run -d \
  --name dockershield \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -p 2375:2375 \
  -e CONTAINERS=1 \
  -e IMAGES=1 \
  hypolas/dockershield:latest

# Test
export DOCKER_HOST=tcp://localhost:2375
docker ps
```

## 🏗️ Multi-Platform Support

The workflow builds for multiple architectures:
- ✅ **linux/amd64** - x86_64 (Intel/AMD)
- ✅ **linux/arm64** - ARM 64-bit (Raspberry Pi 4, Apple M1/M2, AWS Graviton)
- ✅ **linux/arm/v7** - ARM 32-bit (Raspberry Pi 2/3)

Users can pull the appropriate image for their platform automatically:
```bash
# Docker automatically selects the right architecture
docker pull hypolas/dockershield:latest
```

## 📝 Workflow Features

### ✅ What the workflow does:

1. **Triggers only on tags** (not branches)
2. **Multi-platform builds** (amd64, arm64, arm/v7)
3. **Automatic version tagging** from Git tags
4. **Docker Hub description sync** from README.md
5. **Build cache** for faster builds
6. **Release summary** in GitHub Actions UI

### ❌ What triggers are disabled:

- ❌ Branch pushes (including `main`)
- ❌ Pull requests
- ❌ Manual workflow dispatch

This ensures **only tagged releases** generate Docker images.

## 🔄 Release Workflow

### For a new release:

```bash
# 1. Update version in code/docs (if needed)
# 2. Commit changes
git add .
git commit -m "Prepare release v1.0.0"
git push origin main

# 3. Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0

- Added feature X
- Fixed bug Y
- Improved documentation
"
git push origin v1.0.0

# 4. GitHub Actions automatically:
#    - Builds multi-platform images
#    - Pushes to Docker Hub
#    - Updates Docker Hub description

# 5. Create GitHub release (optional but recommended)
gh release create v1.0.0 \
  --title "v1.0.0 - Initial Release" \
  --notes "See CHANGELOG.md for details"
```

## 🐛 Troubleshooting

### ❌ "Error: failed to solve: failed to read dockerfile"
- **Cause**: `Dockerfile` not found in repository root
- **Solution**: Ensure `Dockerfile` exists at `/home/laslite/git/dockershield/Dockerfile`

### ❌ "Error: unauthorized: authentication required"
- **Cause**: Invalid or missing Docker Hub credentials
- **Solution**:
  1. Verify secrets are set correctly
  2. Regenerate Docker Hub token
  3. Update `DOCKERHUB_TOKEN` secret

### ❌ "Error: denied: requested access to the resource is denied"
- **Cause**: Token doesn't have write permissions
- **Solution**: Create new token with `Read, Write, Delete` permissions

### ⚠️ Workflow didn't trigger
- **Cause**: Pushed to branch instead of tag
- **Solution**: Create and push a tag: `git tag v1.0.0 && git push origin v1.0.0`

### ⏱️ Build takes too long (>15 minutes)
- **Cause**: Multi-platform builds + cold cache
- **Solution**: Normal for first build, subsequent builds use cache (~5 min)

## 📊 Monitoring

### Check build status:
```bash
# View recent workflow runs
gh run list --workflow=docker-publish.yml

# View specific run details
gh run view <run-id>

# View logs
gh run view <run-id> --log
```

### Docker Hub statistics:
- **Pulls**: https://hub.docker.com/r/hypolas/dockershield
- **Tags**: https://hub.docker.com/r/hypolas/dockershield/tags
- **Builds**: Managed by GitHub Actions (not Docker Hub automated builds)

## 🔒 Security Best Practices

1. ✅ **Use Access Tokens** (not passwords)
2. ✅ **Limit token scope** to specific repository
3. ✅ **Rotate tokens** every 6-12 months
4. ✅ **Use repository secrets** (not organization secrets)
5. ✅ **Enable 2FA** on Docker Hub account
6. ✅ **Review action logs** regularly

## 📚 Additional Resources

- [Docker Hub Access Tokens](https://docs.docker.com/security/for-developers/access-tokens/)
- [GitHub Actions Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Docker Buildx](https://docs.docker.com/build/buildx/)
- [Multi-platform Images](https://docs.docker.com/build/building/multi-platform/)

## 🎯 Next Steps

After setup:
1. ✅ Test the workflow by creating a tag
2. ✅ Verify images appear on Docker Hub
3. ✅ Pull and test the published image
4. ✅ Update README.md with installation instructions
5. ✅ Create release notes for each version

---

**Questions?** Open an issue at https://github.com/hypolas/dockershield/issues
