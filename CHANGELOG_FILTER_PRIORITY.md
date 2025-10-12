# Changelog: Filter Priority System

## 🎯 Overview

This update introduces a **new filter priority system** where advanced filters (`DKRPRX__*` environment variables) can **override basic ACL restrictions**.

## 🚀 What Changed

### Before (Old Behavior)

```
Request Flow:
1. ACL Middleware (first)
   - IMAGES=0 → ALL /images/* requests blocked
   - Advanced filters never evaluated
2. Advanced Filter Middleware (second)
   - Never reached if ACL blocks
```

**Problem:** If `IMAGES=0`, you could NOT use `DKRPRX__IMAGES__ALLOWED_REPOS` to allow specific registries.

### After (New Behavior)

```
Request Flow:
1. Advanced Filter Middleware (first)
   - Checks DKRPRX__ rules
   - If authorized → marks request in context
   - If denied → blocks immediately
2. ACL Middleware (second)
   - Checks if already authorized by advanced filter
   - If yes → allows (bypass ACL)
   - If no → apply normal ACL rules
```

**Benefit:** You can now set restrictive ACLs and use advanced filters to allow specific operations!

## 📝 Examples

### Example 1: Allow Only Private Registry

```bash
# Disable all image operations by default
export IMAGES=0

# But allow pulls from private registry
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry\.company\.com/.*"

# Result:
docker pull nginx                          # ❌ Blocked (no advanced filter match)
docker pull registry.company.com/nginx     # ✅ Allowed (advanced filter overrides ACL)
```

### Example 2: Selective Container Creation

```bash
# Disable container creation by default
export CONTAINERS=0

# But allow specific approved images
export DKRPRX__CONTAINERS__ALLOWED_IMAGES="^registry\.company\.com/approved/.*"

# Result:
docker run nginx                                    # ❌ Blocked
docker run registry.company.com/approved/app        # ✅ Allowed
docker run registry.company.com/unapproved/app     # ❌ Blocked
```

### Example 3: Default Deny, Explicit Allow

```bash
# Deny everything by default
export IMAGES=0
export CONTAINERS=0
export VOLUMES=0
export NETWORKS=0

# Allow only specific operations via advanced filters
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry\.prod\.company\.com/.*"
export DKRPRX__CONTAINERS__ALLOWED_IMAGES="^registry\.prod\.company\.com/.*:v[0-9]+\\.[0-9]+\\.[0-9]+$"
export DKRPRX__VOLUMES__ALLOWED_NAMES="^prod-.*"
export DKRPRX__NETWORKS__ALLOWED_NAMES="^prod-.*"

# This creates a zero-trust environment where only explicitly allowed operations succeed!
```

## 🔧 Technical Changes

### Modified Files

1. **internal/middleware/advanced_filter.go**
   - Added `handled` flag to track which operations are processed by advanced filters
   - Set `c.Set("advanced_filter_authorized", true)` when operation is authorized
   - This marks the request as authorized in the Gin context

2. **internal/middleware/acl.go**
   - Added check for `advanced_filter_authorized` context value
   - If found and true, bypass ACL check
   - Otherwise, apply normal ACL rules

3. **cmd/dockershield/main.go**
   - Inverted middleware order
   - `AdvancedFilterMiddleware` now runs **before** `ACLMiddleware`
   - Added comment explaining the priority

4. **internal/middleware/acl_test.go**
   - Added `TestACLMiddlewareWithAdvancedFilterOverride`
   - Tests 3 scenarios:
     - Without advanced filter → blocked by ACL
     - With advanced filter → ACL bypassed
     - Advanced filter doesn't affect other endpoints

### New Files

1. **docs/ADVANCED_FILTERS.md**
   - Complete guide to advanced filters
   - Explains new priority system
   - Includes use cases and examples
   - Regex pattern examples
   - Troubleshooting guide

## ✅ Backward Compatibility

**This change is FULLY backward compatible:**

- If no advanced filters are configured, behavior is identical to before
- Existing configurations will work exactly the same
- Only adds new capability when `DKRPRX__*` variables are used

## 🧪 Testing

All existing tests pass, plus new tests added:

```bash
go test ./... -v
# PASS: 100% of tests
```

New test: `TestACLMiddlewareWithAdvancedFilterOverride`
- Verifies ACL blocking without advanced filter
- Verifies ACL bypass with advanced filter
- Verifies advanced filter doesn't affect unrelated endpoints

## 📖 Documentation Updates

1. **README.md**
   - Added "🆕 NEW: Advanced Filters Override ACL" section
   - Added bullet point in features list
   - Link to full advanced filters documentation

2. **docs/ADVANCED_FILTERS.md** (NEW)
   - Complete advanced filters guide
   - Filter priority explanation
   - Use cases and examples
   - Best practices
   - Troubleshooting

## 🎓 Use Cases

### Security Hardening
Set restrictive defaults, allow specific operations:
```bash
export IMAGES=0
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry\.company\.com/.*"
```

### Multi-Tenant Isolation
Each tenant gets their own filtered access:
```bash
export IMAGES=0
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry\.tenant-${TENANT_ID}\.com/.*"
```

### CI/CD Control
Allow builds, but enforce versioning:
```bash
export IMAGES=1
export DKRPRX__IMAGES__DENIED_TAGS="^latest$"
export DKRPRX__IMAGES__ALLOWED_TAGS="^v[0-9]+\\.[0-9]+\\.[0-9]+$"
```

## 🚦 Migration Guide

### If You're Not Using Advanced Filters
**No action needed!** Your configuration works as before.

### If You're Using Advanced Filters
**Good news!** Your filters now have MORE power:

**Before:**
```bash
export IMAGES=1  # Had to enable
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry\.company\.com/.*"
```

**After (More Secure):**
```bash
export IMAGES=0  # Can disable
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry\.company\.com/.*"
# Advanced filter overrides the ACL!
```

## 💡 Best Practices

1. **Use Default-Deny Approach**
   ```bash
   export IMAGES=0
   export CONTAINERS=0
   export DKRPRX__IMAGES__ALLOWED_REPOS="^approved-registry\\.com/.*"
   ```

2. **Combine with HTTP Method Restrictions**
   ```bash
   export POST=0  # Disable creation by default
   export DKRPRX__CONTAINERS__ALLOWED_IMAGES="^registry\\.company\\.com/.*"
   ```

3. **Layer Your Security**
   - ACL for broad categories
   - Advanced filters for fine-grained control
   - Both work together!

## 🐛 Known Issues

None at this time.

## 📞 Support

- Documentation: [docs/ADVANCED_FILTERS.md](docs/ADVANCED_FILTERS.md)
- Issues: https://github.com/hypolas/dockershield/issues
- Email: nicolas.hypolite@gmail.com

---

**Version:** 1.1.0 (unreleased)
**Date:** 2025-10-12
**Author:** DockerShield Team
