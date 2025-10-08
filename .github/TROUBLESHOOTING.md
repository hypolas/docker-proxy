# 🐛 GitHub Actions Troubleshooting

Common issues and solutions for the Docker publish workflow.

## ❌ Error: "tag is needed when pushing to registry"

### Full Error Message
```
ERROR: failed to build: tag is needed when pushing to registry
```

### Possible Causes

#### 1. **No tags generated** (most common)

The metadata action couldn't generate tags from your Git ref.

**Check:**
```bash
# Verify your tag format
git tag -l

# Should see tags like:
# v1.0.0
# v2.1.3
# NOT: 1.0.0 (missing 'v')
# NOT: release-1.0.0 (wrong format)
```

**Solution:**
```bash
# Create tag with 'v' prefix
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

#### 2. **Wrong trigger**

Workflow only runs on tags, not branches.

**Check GitHub Actions logs:**
```
Git ref: refs/tags/v1.0.0     ✅ Correct (tag)
Git ref: refs/heads/main      ❌ Wrong (branch)
```

**Solution:**
```bash
# Don't push to branches expecting build
git push origin main          # ❌ Won't trigger

# Create and push a tag
git tag v1.0.0
git push origin v1.0.0        # ✅ Will trigger
```

#### 3. **Missing Dockerfile**

Dockerfile not found in repository root.

**Check workflow logs:**
```
❌ Error: Dockerfile not found in repository root
```

**Solution:**
```bash
# Verify Dockerfile exists
ls -la Dockerfile

# If missing, create it
# Then commit and push
git add Dockerfile
git commit -m "Add Dockerfile"
git push origin main

# Then create tag
git tag v1.0.0
git push origin v1.0.0
```

#### 4. **Invalid semver tag**

Tag doesn't match semantic versioning pattern.

**Invalid examples:**
```bash
git tag 1.0.0           # ❌ Missing 'v' prefix
git tag v1.0            # ❌ Missing patch version
git tag release-1.0.0   # ❌ Wrong prefix
git tag latest          # ❌ Not semver
```

**Valid examples:**
```bash
git tag v1.0.0          # ✅ Stable release
git tag v1.2.3          # ✅ Stable release
git tag v2.0.0-beta     # ✅ Pre-release
git tag v1.0.0-rc1      # ✅ Release candidate
```

**Solution:**
```bash
# Delete wrong tag
git tag -d wrong-tag
git push origin :refs/tags/wrong-tag

# Create correct tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Debug Steps

1. **Check workflow run details:**
   ```
   https://github.com/hypolas/dockershield/actions
   ```

2. **Look for "Debug - Show generated tags" step:**
   ```
   Generated tags:
   hypolas/proxy-docker:1.0.0
   hypolas/proxy-docker:1.0
   hypolas/proxy-docker:1
   hypolas/proxy-docker:latest
   ```

   If empty → metadata action failed

3. **Check Git ref info:**
   ```
   Git ref: refs/tags/v1.0.0
   Git ref name: v1.0.0
   Git ref type: tag
   ```

4. **Verify tag locally:**
   ```bash
   git show v1.0.0
   ```

### Quick Fix Checklist

- [ ] Tag has `v` prefix (e.g., `v1.0.0`)
- [ ] Tag is semantic version (`vMAJOR.MINOR.PATCH`)
- [ ] Tag is pushed to GitHub (`git push origin v1.0.0`)
- [ ] Dockerfile exists in repository root
- [ ] Workflow triggered on tag (not branch)
- [ ] GitHub secrets are configured (DOCKERHUB_USERNAME, DOCKERHUB_TOKEN)

## ❌ Error: "unauthorized: authentication required"

### Full Error Message
```
ERROR: failed to solve: failed to do request: Head "https://registry-1.docker.io/v2/hypolas/proxy-docker/manifests/1.0.0": unauthorized: authentication required
```

### Causes

1. **Missing or invalid Docker Hub credentials**

**Check:**
```
Settings → Secrets and variables → Actions
```

Should have:
- `DOCKERHUB_USERNAME` = `hypolas`
- `DOCKERHUB_TOKEN` = `<your-access-token>`

**Solution:**
1. Go to Docker Hub: https://hub.docker.com
2. Account Settings → Security → New Access Token
3. Create token with `Read, Write, Delete` permissions
4. Copy token
5. Add to GitHub secrets

2. **Token expired or revoked**

**Solution:**
Regenerate Docker Hub token and update GitHub secret.

3. **Wrong username**

**Check:**
```yaml
# Should match your Docker Hub username
username: ${{ secrets.DOCKERHUB_USERNAME }}
```

**Solution:**
Verify secret value is exactly `hypolas`

## ❌ Error: "denied: requested access to the resource is denied"

### Causes

1. **Token lacks write permissions**

**Solution:**
Create new token with `Read, Write, Delete` permissions.

2. **Repository doesn't exist on Docker Hub**

**Solution:**
First push creates the repository automatically. Ensure you have permissions.

3. **Username mismatch**

**Check:**
- Docker Hub username: `hypolas`
- GitHub secret DOCKERHUB_USERNAME: `hypolas`
- Image name in workflow: `hypolas/proxy-docker`

Must all match.

## ❌ Error: "no builder instance found"

### Full Error Message
```
ERROR: no builder instance found
```

### Solution

This is automatically handled by the workflow. If you see this:

1. Check if `docker/setup-buildx-action` step succeeded
2. Verify Docker is available in the runner

## ❌ Build succeeds but image not on Docker Hub

### Causes

1. **`push: false` in build step**

**Check:**
```yaml
- name: Build and push Docker image
  uses: docker/build-push-action@v5
  with:
    push: true  # ← Must be true
```

2. **Workflow completed with warnings**

**Check workflow logs** for warnings.

3. **Docker Hub rate limit**

**Solution:**
Wait a few minutes and retry.

## ⚠️ Warning: "provenance: true should be used"

### Message
```
provenance: true should be used with multi-platform builds
```

### Solution (Optional)

Add provenance attestations:

```yaml
- name: Build and push Docker image
  uses: docker/build-push-action@v5
  with:
    push: true
    provenance: true  # Add this
    sbom: true        # Optional: Software Bill of Materials
```

## 🐌 Build taking too long (>15 minutes)

### Normal Times

- **First build**: 10-15 minutes (no cache, 3 platforms)
- **Subsequent builds**: 5-7 minutes (with cache)

### Solutions

1. **Use cache** (already configured):
   ```yaml
   cache-from: type=gha
   cache-to: type=gha,mode=max
   ```

2. **Reduce platforms** (if needed):
   ```yaml
   env:
     PLATFORMS: linux/amd64  # Single platform for testing
   ```

3. **Check if stuck**:
   - View real-time logs in Actions tab
   - Look for hanging dependencies

## 🔍 Debugging Tips

### 1. Test locally

```bash
# Install buildx
docker buildx install

# Create builder
docker buildx create --name multiarch --use

# Build for single platform
docker buildx build --platform linux/amd64 -t test .

# Build for all platforms
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  -t hypolas/proxy-docker:test \
  .
```

### 2. Check action versions

Ensure you're using latest stable versions:
```yaml
uses: actions/checkout@v4           # ✅
uses: docker/setup-qemu-action@v3   # ✅
uses: docker/setup-buildx-action@v3 # ✅
uses: docker/login-action@v3        # ✅
uses: docker/build-push-action@v5   # ✅
uses: docker/metadata-action@v5     # ✅
```

### 3. Enable debug logging

```bash
# Re-run workflow with debug logging
gh run rerun <run-id> --debug
```

### 4. Check runner logs

Look for system issues in the runner environment.

## 📞 Getting Help

If issue persists:

1. **Check existing issues:**
   ```
   https://github.com/hypolas/dockershield/issues
   ```

2. **Create new issue with:**
   - Workflow run URL
   - Full error message
   - Git tag used
   - Steps already tried

3. **Relevant logs:**
   - Copy "Debug - Show generated tags" output
   - Copy "Show Git ref info" output
   - Copy error message

## 📚 Reference

- [Docker Build Push Action](https://github.com/docker/build-push-action)
- [Docker Metadata Action](https://github.com/docker/metadata-action)
- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Docker Hub](https://hub.docker.com/r/hypolas/proxy-docker)

---

**Still stuck?** Open an issue with full details and we'll help! 🚀
