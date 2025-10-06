# Configuration Conflict Resolution - Docker Proxy

This document explains how docker-proxy handles conflicts between different configuration sources.

## 🔄 Configuration Priority Order

docker-proxy uses **3 configuration sources** with the following priority order:

```
1. Environment variables (DKRPRX__*)    ← HIGHEST PRIORITY
2. JSON file (FILTERS_CONFIG)           ← MEDIUM PRIORITY
3. Default security filters             ← LOWEST PRIORITY
```

## 📊 Merge Process (config/config.go)

```go
// Execution order in Load()
1. jsonFilters := loadAdvancedFilters(filtersPath)      // Load JSON
2. envFilters := LoadFiltersFromEnv()                   // Load ENV (priority)
3. mergedFilters := MergeFilters(jsonFilters, envFilters) // Merge
4. mergedFilters = ApplyDefaults(mergedFilters)         // Apply defaults
```

## 🎯 Conflict Examples and Resolutions

### Example 1: ENV vs JSON Conflict

**JSON Configuration** (`filters.json`):
```json
{
  "volumes": {
    "denied_paths": ["/data"]
  }
}
```

**Environment variable**:
```bash
DKRPRX__VOLUMES__DENIED_PATHS=/var/log,/tmp
```

**Result**: ENV wins ✅
```
denied_paths = ["/var/log", "/tmp"]  # JSON ignored
```

### Example 2: Conflict with Default Filters

**Default filter** (defaults.go):
```go
DeniedPaths: []string{
    "^/var/run/docker\\.sock$",
    "^/run/docker\\.sock$",
}
```

**JSON Configuration**:
```json
{
  "volumes": {
    "allowed_paths": ["/app"]
  }
}
```

**Result**: Merge (no direct conflict)
```
denied_paths  = ["/var/run/docker.sock", "/run/docker.sock"]  # Defaults
allowed_paths = ["/app"]                                       # JSON
```

### Example 3: Disabling Default Filters

**With DKRPRX__DISABLE_DEFAULTS**:
```bash
DKRPRX__DISABLE_DEFAULTS=true
```

**Result**: Default filters are **completely ignored** ⚠️
```go
if !CanOverrideDefaults() {
    mergedFilters = ApplyDefaults(mergedFilters)
}
// If DISABLE_DEFAULTS=true, ApplyDefaults() is NOT called
```

## 🔍 Detailed Merge Logic

### MergeFilters() - config/env_filters.go:259

```go
func MergeFilters(jsonFilter, envFilter *filters.AdvancedFilter) {
    // For each section (Volumes, Containers, Networks, Images):

    if envFilter.Volumes != nil {
        result.Volumes = envFilter.Volumes  // ENV replaces EVERYTHING
    } else {
        result.Volumes = jsonFilter.Volumes  // JSON if no ENV
    }
}
```

**Key point**: Merging is done **per entire section**, not per individual field.

### Concrete Example

**JSON**:
```json
{
  "volumes": {
    "denied_paths": ["/data"],
    "allowed_paths": ["/app"]
  }
}
```

**ENV**:
```bash
DKRPRX__VOLUMES__DENIED_PATHS=/tmp
# Note: allowed_paths is NOT defined in ENV
```

**Result**: ⚠️ **Entire Volumes section from JSON is replaced**
```
denied_paths  = ["/tmp"]     # From ENV
allowed_paths = []           # LOST! (was only in JSON)
```

## 🛡️ Default Security Rules

Default filters are **always applied** unless `DKRPRX__DISABLE_DEFAULTS=true`.

### ApplyDefaults() - config/defaults.go:47

```go
func ApplyDefaults(filter *filters.AdvancedFilter) {
    defaults := GetDefaultFilters()

    // Apply defaults ONLY if section is nil
    if filter.Volumes == nil {
        filter.Volumes = defaults.Volumes
    }
    // So: if you define EVEN ONE field in Volumes,
    // the defaults for Volumes are NOT applied!
}
```

### Security Defaults

**Volumes** (defaults.go:23):
```go
DeniedPaths: []string{
    `^/var/run/docker\.sock$`,  // Block Docker socket
    `^/run/docker\.sock$`,
}
```

**Containers** (defaults.go:15):
```go
DeniedNames: []string{
    `^docker-proxy$`,    // Protect the proxy itself
    `^/docker-proxy$`,
}
```

**Networks**:
```go
DeniedNames: []string{
    `^proxy-network$`,   // If PROXY_NETWORK_NAME is defined
}
```

## 📋 Conflicts in Filters (allowed vs denied)

### Evaluation Order (pkg/filters/advanced.go)

```go
func CheckVolumeMount(volumeName, hostPath, driver string) (bool, string) {
    // 1. Check denied BEFORE allowed
    if matchesDeniedList(vf.DeniedPaths, hostPath) {
        return false, "denied"  // BLOCKED
    }

    // 2. If allowed_paths is defined, check inclusion
    if len(vf.AllowedPaths) > 0 {
        if !matchesAllowedList(vf.AllowedPaths, hostPath) {
            return false, "not in allowed"  // BLOCKED
        }
    }

    // 3. Everything is OK
    return true, ""
}
```

### Example: Both Allowed AND Denied defined

```json
{
  "volumes": {
    "denied_paths": ["/var"],
    "allowed_paths": ["/app", "/data"]
  }
}
```

**Test**: Mount `/var/log`
```
1. Check denied: /var/log matches ^/var → BLOCKED ❌
2. allowed_paths is not even checked
```

**Test**: Mount `/tmp`
```
1. Check denied: /tmp doesn't match /var → OK
2. Check allowed: /tmp doesn't match /app or /data → BLOCKED ❌
```

**Test**: Mount `/app/config`
```
1. Check denied: /app/config doesn't match /var → OK
2. Check allowed: /app/config matches /app → ALLOWED ✅
```

## ⚙️ Practical Use Cases

### Case 1: Maximum Security (Defaults Enabled)

```bash
# No DKRPRX__ variables defined
# Defaults protect against:
# - Docker socket mounting
# - Proxy container manipulation
# - Proxy network manipulation
```

### Case 2: Custom Configuration (JSON)

```json
{
  "volumes": {
    "allowed_paths": ["/app", "/data"]
  }
}
```

```bash
FILTERS_CONFIG=/etc/docker-proxy/filters.json
# Volumes defaults are REPLACED by JSON
# ⚠️ Docker socket is NO LONGER protected!
```

### Case 3: Override with ENV

```bash
# JSON defines rules
FILTERS_CONFIG=/etc/filters.json

# ENV overrides a specific section
DKRPRX__VOLUMES__DENIED_PATHS=/var/run/docker.sock,/tmp

# Result: ENV completely replaces the Volumes section from JSON
```

### Case 4: Complete Defaults Disabling

```bash
DKRPRX__DISABLE_DEFAULTS=true

# ⚠️ DANGEROUS: No default protection anymore
# You must define ALL rules yourself
```

## 🔧 Recommendations

### ✅ Best Practices

1. **Keep defaults enabled** in production
2. **Use JSON for base configuration**
3. **Use ENV for temporary overrides**
4. **Always define `denied` before `allowed`**

### ❌ To Avoid

1. **Disabling defaults** without valid reason
2. **Mixing JSON and ENV for the same section** (ENV overwrites everything)
3. **Forgetting that merge is per section, not per field**

## 🧪 Testing Configuration

### Test Script

```bash
# Display final configuration
docker run --rm \
  -e CONTAINERS=1 \
  -e DKRPRX__VOLUMES__DENIED_PATHS=/var/log \
  -e LOG_LEVEL=debug \
  hypolas/proxy-docker:latest

# Check logs to see applied configuration
```

### Production Verification

```bash
# See active rules in logs
docker logs docker-proxy | grep "Access Rules Configuration"
docker logs docker-proxy | grep "Advanced Filters"
```

## 📚 Reference Files

| File | Responsibility |
|------|---------------|
| `config/config.go:58` | Entry point `Load()` - orchestration |
| `config/env_filters.go:259` | `MergeFilters()` - ENV + JSON merge |
| `config/defaults.go:47` | `ApplyDefaults()` - apply defaults |
| `pkg/filters/advanced.go` | Filter evaluation logic |

## 🆘 Debugging Conflicts

### Enable Debug Logs

```bash
LOG_LEVEL=debug
```

### Logs to Monitor

```
Access Rules Configuration:
  Granted endpoints: [...]
  Allowed methods: [...]

Advanced Filters:
  Volumes: {...}
  Containers: {...}
```

## 🔄 Merge Behavior Summary

| Scenario | Behavior | Example |
|----------|----------|---------|
| **ENV defined** | ENV replaces entire section | `DKRPRX__VOLUMES__*` → JSON Volumes ignored |
| **Only JSON defined** | JSON used, defaults applied to other sections | JSON Volumes → Default Containers |
| **Only defaults** | All defaults applied | No config → Full protection |
| **DISABLE_DEFAULTS=true** | No defaults, only user config | Requires complete manual config |

## ⚡ Quick Reference

### Priority Chain
```
ENV > JSON > Defaults
```

### Section-Level Merge
```
If ANY field in section is set via ENV → Entire section from JSON is ignored
```

### Evaluation Order
```
denied (checked first) → allowed (if defined) → default allow
```

### Default Override
```
Define ANY field in section → Section defaults are NOT applied
```

---

**Version**: 1.0
**Last Updated**: 2025-10-06
**Maintainer**: Nicolas HYPOLITE (nicolas.hypolite@gmail.com)
