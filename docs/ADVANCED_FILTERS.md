# 🔐 Advanced Filters Guide

This guide explains the advanced filtering system in DockerShield, which provides granular, regex-based control over Docker operations.

## 🎯 Overview

Advanced filters (`DKRPRX__*` environment variables) provide fine-grained control over Docker operations beyond simple ACL enable/disable flags. They allow you to:

- **Whitelist/blacklist** specific images, volumes, networks by pattern
- **Enforce security policies** (deny privileged containers, block :latest tag, etc.)
- **Implement tenant isolation** in multi-tenant environments
- **Override ACL restrictions** when specific conditions are met

## 🚀 NEW: Filter Priority System

**Important:** Advanced filters now have **HIGHER PRIORITY** than basic ACL rules.

### How It Works

```
Request Flow:
┌─────────────────────────────────────────┐
│  1. Advanced Filter Middleware          │ ← Runs FIRST
│     - Checks DKRPRX__ rules             │
│     - If authorized → marks request     │
│     - If denied → blocks immediately    │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│  2. ACL Middleware                       │ ← Runs SECOND
│     - Checks if already authorized      │
│     - If yes → allows (bypass ACL)      │
│     - If no → apply ACL rules           │
└─────────────────────────────────────────┘
```

### Practical Example

**Scenario:** You want to allow pulling images ONLY from your private registry, even with `IMAGES=0`

**Before (❌ Not possible):**
```bash
export IMAGES=0  # Blocks ALL image operations
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry.company.com/.*"
# Result: All image pulls are blocked, filter never evaluated
```

**After (✅ Works!):**
```bash
export IMAGES=0  # Basic ACL blocks images
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry.company.com/.*"
# Result: Only images from registry.company.com can be pulled
#         ACL is bypassed for authorized operations
```

## 📚 Configuration Methods

### 1. Environment Variables (Recommended)

Environment variables have the highest priority and override everything else.

**Format:** `DKRPRX__SECTION__PARAMETER=value`

### 2. JSON Configuration File

Set `FILTERS_CONFIG=/path/to/filters.json` to load from a JSON file.

### 3. Default Security Filters

Built-in defaults (can be disabled with `DKRPRX__DISABLE_DEFAULTS=true`)

## 🔧 Available Filters

### Volume Filters

Control which volumes can be created and which host paths can be mounted.

```bash
# Allow only specific volume names (regex patterns)
export DKRPRX__VOLUMES__ALLOWED_NAMES="^data-.*,^app-.*,^logs-.*"

# Block sensitive system paths
export DKRPRX__VOLUMES__DENIED_PATHS="^/etc/.*,^/root/.*,^/sys/.*,^/proc/.*,^/var/run/.*"

# Allow only specific host paths
export DKRPRX__VOLUMES__ALLOWED_PATHS="^/data/.*,^/mnt/storage/.*"

# Restrict volume drivers
export DKRPRX__VOLUMES__ALLOWED_DRIVERS="local,nfs"
```

**Example: Block Docker socket mounting**
```bash
export VOLUMES=1
export DKRPRX__VOLUMES__DENIED_PATHS="^/var/run/docker\\.sock$,^/run/docker\\.sock$"
```

### Container Filters

Control which containers can be created based on images, names, and security settings.

```bash
# Allow only images from private registry with semantic versioning
export DKRPRX__CONTAINERS__ALLOWED_IMAGES="^registry.company.com/.*:v[0-9]+\\.[0-9]+\\.[0-9]+$"

# Block images with dangerous tags
export DKRPRX__CONTAINERS__DENIED_IMAGES=".*:(latest|dev|test)$"

# Enforce container naming convention
export DKRPRX__CONTAINERS__ALLOWED_NAMES="^(prod|staging|dev)-.*"

# Require specific labels
export DKRPRX__CONTAINERS__REQUIRE_LABELS="env=production,team=backend,cost-center=IT-001"

# Security restrictions
export DKRPRX__CONTAINERS__DENY_PRIVILEGED="true"
export DKRPRX__CONTAINERS__DENY_HOST_NETWORK="true"
```

**Example: Allow containers even with CONTAINERS=0**
```bash
export CONTAINERS=0  # Disable container creation by default
export DKRPRX__CONTAINERS__ALLOWED_IMAGES="^registry.company.com/approved/.*"
# Only containers from registry.company.com/approved/* can be created
```

### Image Filters

Control which images can be pulled or built.

```bash
# Allow only specific registries
export DKRPRX__IMAGES__ALLOWED_REPOS="^(docker\\.io/library|registry\\.company\\.com)/.*"

# Block suspicious registries
export DKRPRX__IMAGES__DENIED_REPOS=".*\\.(cn|ru|suspicious)/"

# Enforce semantic versioning
export DKRPRX__IMAGES__ALLOWED_TAGS="^v[0-9]+\\.[0-9]+\\.[0-9]+$"

# Block dangerous tags
export DKRPRX__IMAGES__DENIED_TAGS="^(latest|dev|test|alpha|beta|rc).*"
```

**Example: Block :latest tag but allow everything else**
```bash
export IMAGES=1
export DKRPRX__IMAGES__DENIED_TAGS="^latest$"
# docker pull nginx:latest     ← ❌ Denied
# docker pull nginx:1.25.3     ← ✅ Allowed
```

### Network Filters

Control network creation.

```bash
# Enforce network naming convention
export DKRPRX__NETWORKS__ALLOWED_NAMES="^app-.*"

# Block host network mode
export DKRPRX__NETWORKS__DENIED_NAMES="^host$"

# Restrict network drivers
export DKRPRX__NETWORKS__ALLOWED_DRIVERS="bridge,overlay"
```

## 🎓 Use Cases

### Use Case 1: Enforce Private Registry (Override IMAGES=0)

**Goal:** Allow image pulls ONLY from private registry, block everything else

```bash
# Basic ACL: Disable images
export IMAGES=0

# Advanced filter: Allow only private registry
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry.company.com/.*"

# Result:
# docker pull nginx                          ← ❌ Blocked by ACL (no advanced filter match)
# docker pull registry.company.com/nginx     ← ✅ Allowed (advanced filter overrides ACL)
```

### Use Case 2: Block :latest Tag in CI/CD

**Goal:** Enforce semantic versioning, prevent `:latest` tag usage

```bash
export IMAGES=1
export CONTAINERS=1
export DKRPRX__IMAGES__DENIED_TAGS="^latest$"
export DKRPRX__CONTAINERS__DENIED_IMAGES=".*:latest$"

# Result:
# docker pull nginx:latest       ← ❌ Denied by image filter
# docker run nginx:latest        ← ❌ Denied by container filter
# docker pull nginx:1.25.3       ← ✅ Allowed
# docker run nginx:1.25.3        ← ✅ Allowed
```

### Use Case 3: Multi-Tenant Isolation

**Goal:** Each tenant can only access their own resources

```bash
export TENANT_ID="tenant-123"

# Prefix-based isolation
export DKRPRX__VOLUMES__ALLOWED_NAMES="^${TENANT_ID}-.*"
export DKRPRX__CONTAINERS__ALLOWED_NAMES="^${TENANT_ID}-.*"
export DKRPRX__NETWORKS__ALLOWED_NAMES="^${TENANT_ID}-.*"

# Enforce tenant label
export DKRPRX__CONTAINERS__REQUIRE_LABELS="tenant=${TENANT_ID}"

# Result:
# docker volume create tenant-123-data      ← ✅ Allowed
# docker volume create tenant-456-data      ← ❌ Denied
# docker network create tenant-123-net      ← ✅ Allowed
# docker network create tenant-456-net      ← ❌ Denied
```

### Use Case 4: Production Security Hardening

**Goal:** Maximum security for production environment

```bash
export CONTAINERS=1
export IMAGES=1
export VOLUMES=1

# Only versioned images from production registry
export DKRPRX__CONTAINERS__ALLOWED_IMAGES="^registry.prod.company.com/.*:v[0-9]+\\.[0-9]+\\.[0-9]+$"

# Block sensitive paths
export DKRPRX__VOLUMES__ALLOWED_PATHS="^/data/prod/.*"
export DKRPRX__VOLUMES__DENIED_PATHS="^/(etc|root|home|sys|proc|var/run)/.*"

# Security restrictions
export DKRPRX__CONTAINERS__DENY_PRIVILEGED="true"
export DKRPRX__CONTAINERS__DENY_HOST_NETWORK="true"

# Require production labels
export DKRPRX__CONTAINERS__REQUIRE_LABELS="env=production,approved=true,security-scan=passed"
```

### Use Case 5: CI/CD Flexibility with Security

**Goal:** Allow builds but restrict deployments

```bash
# Allow image operations
export IMAGES=1

# But block :latest tag
export DKRPRX__IMAGES__DENIED_TAGS="^latest$"

# Only allow images from CI registry
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry.ci.company.com/.*"

# Block sensitive volume mounts
export DKRPRX__VOLUMES__DENIED_PATHS="^/(etc|root|home|sys|proc)/.*"

# Result:
# docker build -t registry.ci.company.com/app:1.0.0 .   ← ✅ Allowed
# docker build -t registry.ci.company.com/app:latest .  ← ❌ Denied
# docker build -t external.io/app:1.0.0 .               ← ❌ Denied
```

## 🔄 Filter Evaluation Logic

### Deny Lists vs Allow Lists

Filters are evaluated in this order:

1. **Denied list** checked first
   - If matches → ❌ **DENIED**
   - If no match → continue to step 2

2. **Allowed list** checked second
   - If no allowed list configured → ✅ **ALLOWED**
   - If allowed list exists:
     - Matches → ✅ **ALLOWED**
     - No match → ❌ **DENIED**

### Example

```bash
export DKRPRX__IMAGES__DENIED_TAGS="^(latest|dev)$"
export DKRPRX__IMAGES__ALLOWED_TAGS="^v[0-9]+\\.[0-9]+\\.[0-9]+$"

# Evaluation:
# nginx:latest    → Matches denied list → ❌ DENIED
# nginx:dev       → Matches denied list → ❌ DENIED
# nginx:v1.0.0    → Not denied, matches allowed → ✅ ALLOWED
# nginx:1.0.0     → Not denied, but doesn't match allowed → ❌ DENIED
# nginx:stable    → Not denied, but doesn't match allowed → ❌ DENIED
```

## 📋 Regex Pattern Examples

### Exact Match
```bash
# Block exact tag name
export DKRPRX__IMAGES__DENIED_TAGS="^latest$"
```

### Prefix Match
```bash
# Allow only company registry
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry\\.company\\.com/.*"
```

### Multiple Alternatives
```bash
# Block multiple tags
export DKRPRX__IMAGES__DENIED_TAGS="^(latest|dev|test|master|main)$"
```

### Semantic Versioning
```bash
# Match v1.2.3 format
export DKRPRX__IMAGES__ALLOWED_TAGS="^v[0-9]+\\.[0-9]+\\.[0-9]+$"
```

### Path Blocking
```bash
# Block system directories
export DKRPRX__VOLUMES__DENIED_PATHS="^/(etc|root|home|sys|proc|var/run)/.*"
```

## 🐛 Troubleshooting

### Filter Not Working

**Problem:** Advanced filter seems to be ignored

**Solution:** Check filter evaluation order
```bash
# Enable debug logging
export LOG_LEVEL=debug

# Check logs
docker logs dockershield | grep -i "filter"
```

### ACL Still Blocking Despite Filter

**Problem:** `IMAGES=0` blocks even with `DKRPRX__IMAGES__ALLOWED_REPOS` set

**Solution:** This is the NEW behavior! Advanced filters now override ACL.

✅ **Current behavior (after update):**
```bash
export IMAGES=0
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry.company.com/.*"
# docker pull registry.company.com/nginx → ✅ ALLOWED (filter overrides ACL)
```

### Regex Not Matching

**Problem:** Pattern doesn't match expected strings

**Solution:** Escape special characters
```bash
# Wrong: Will not work
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry.company.com/.*"

# Correct: Escape dots
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry\\.company\\.com/.*"
```

### Testing Filters

Use curl to test filter behavior:
```bash
# Test image pull
curl -X POST "http://localhost:2375/v1.41/images/create?fromImage=nginx&tag=latest"

# Expected responses:
# ✅ Success: HTTP 200
# ❌ Denied: HTTP 403 with "Image operation denied by advanced filter"
```

## 📖 Related Documentation

- [README.md](../README.md) - Main documentation
- [SECURITY.md](SECURITY.md) - Security best practices
- [CICD_EXAMPLES.md](CICD_EXAMPLES.md) - CI/CD integration examples
- [CONFLICT_RESOLUTION.md](CONFLICT_RESOLUTION.md) - Filter priority and merging

## 💡 Best Practices

1. **Start with deny-first approach**
   - Set restrictive ACL (IMAGES=0, CONTAINERS=0)
   - Use advanced filters to allow specific operations

2. **Use semantic versioning enforcement**
   ```bash
   export DKRPRX__IMAGES__DENIED_TAGS="^latest$"
   export DKRPRX__IMAGES__ALLOWED_TAGS="^v[0-9]+\\.[0-9]+\\.[0-9]+$"
   ```

3. **Implement tenant isolation**
   ```bash
   export DKRPRX__VOLUMES__ALLOWED_NAMES="^${TENANT_ID}-.*"
   export DKRPRX__CONTAINERS__ALLOWED_NAMES="^${TENANT_ID}-.*"
   ```

4. **Test in development first**
   - Enable `LOG_LEVEL=debug`
   - Monitor logs for denied operations
   - Adjust filters accordingly

5. **Document your filters**
   - Keep a config file with comments
   - Share with team members
   - Version control your filter configurations

## 🔗 Quick Links

- [Docker API Reference](https://docs.docker.com/engine/api/)
- [Regex101](https://regex101.com/) - Test your regex patterns
- [GitHub Issues](https://github.com/hypolas/dockershield/issues) - Report bugs or request features
