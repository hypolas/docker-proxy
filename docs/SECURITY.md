# Security - Docker Socket Proxy

This document describes the security protections implemented in the proxy and recommended best practices.

## 🛡️ Default Protections

The proxy implements several security layers enabled **by default** to prevent privilege escalation.

### 1. Docker Socket Protection

**Automatically blocks mounting of the Docker socket** to prevent an unsecured container from taking full control of Docker.

Paths blocked by default:
- `/var/run/docker.sock`
- `/run/docker.sock`

**Example of blocked attempt:**
```bash
docker run -v /var/run/docker.sock:/var/run/docker.sock alpine
# ❌ Denied: "Volume creation denied by advanced filter"
# Reason: "host path is denied: /var/run/docker.sock"
```

### 2. Proxy Container Self-Protection

**The docker-proxy container protects itself** against any manipulation via the API it exposes.

Automatic protection:
- ❌ Cannot stop the proxy container
- ❌ Cannot restart the proxy container
- ❌ Cannot delete the proxy container
- ❌ Cannot modify the proxy container
- ❌ Cannot execute commands in the proxy

**Example of blocked attempt:**
```bash
# Attempt to stop the proxy via the API it exposes
docker stop docker-proxy
# ❌ Denied: "Container operation denied by advanced filter"
# Reason: "container name is denied: docker-proxy"
```

### 3. Proxy Network Protection

If the proxy uses a dedicated network, it is also protected.

**Example:**
```bash
docker network rm docker-proxy
# ❌ Denied: "Network operation denied by advanced filter"
# Reason: "network name is denied: docker-proxy"
```

## ⚙️ Protection Configuration

### Environment Variables

```bash
# Name of container to protect (default: docker-proxy)
export PROXY_CONTAINER_NAME="docker-proxy"

# Name of network to protect (optional)
export PROXY_NETWORK_NAME="docker-proxy-network"
```

### Docker Compose

```yaml
services:
  docker-proxy:
    container_name: docker-proxy
    environment:
      - PROXY_CONTAINER_NAME=docker-proxy
      - PROXY_NETWORK_NAME=docker-proxy
    networks:
      - docker-proxy

networks:
  docker-proxy:
    name: docker-proxy
```

## 🔓 Disabling Protections

### ⚠️ Complete Disabling (NOT RECOMMENDED)

```bash
export DKRPRX__DISABLE_DEFAULTS="true"
```

**Consequences:**
- Docker socket can be mounted freely
- Proxy container can be manipulated
- Proxy network can be removed
- **High risk of privilege escalation**

### ✅ Selective Disabling (Recommended)

To allow only the Docker socket:

```bash
# Via env var (overrides defaults for volumes only)
export DKRPRX__VOLUMES__ALLOWED_PATHS="^/var/run/docker\\.sock$"
```

Or via JSON:
```json
{
  "volumes": {
    "allowed_paths": ["^/var/run/docker\\.sock$"]
  }
}
```

## 🚨 Blocked Attack Vectors

### 1. Escalation via Docker Socket

**Attack:**
```bash
# Attacker tries to mount Docker socket
docker run -v /var/run/docker.sock:/var/run/docker.sock \
  alpine sh -c "docker run --privileged --pid=host alpine nsenter -t 1 -m -u -i sh"
```

**Protection:**
- ✅ Socket mounting blocked by default
- ✅ Even if VOLUMES=1, path filter blocks it

### 2. Escalation via Proxy Manipulation

**Attack:**
```bash
# Attacker tries to stop proxy to bypass restrictions
docker stop docker-proxy
docker run -v /var/run/docker.sock:/var/run/docker.sock alpine
```

**Protection:**
- ✅ Proxy container manipulation blocked
- ✅ Proxy API refuses operations on itself

### 3. Escalation via Host Network

**Attack:**
```bash
# Attacker tries to get full network access
docker run --network=host alpine
```

**Protection (optional):**
```bash
export DKRPRX__CONTAINERS__DENY_HOST_NETWORK="true"
```

### 4. Escalation via Privileged Container

**Attack:**
```bash
# Attacker tries to launch privileged container
docker run --privileged alpine
```

**Protection (optional):**
```bash
export DKRPRX__CONTAINERS__DENY_PRIVILEGED="true"
```

## 📋 Security Checklist

### Production Deployment

- [ ] **Docker socket read-only**: `-v /var/run/docker.sock:/var/run/docker.sock:ro`
- [ ] **Default protections enabled**: Do not set `DKRPRX__DISABLE_DEFAULTS`
- [ ] **Container name configured**: `PROXY_CONTAINER_NAME=docker-proxy`
- [ ] **Dedicated network**: `PROXY_NETWORK_NAME=docker-proxy`
- [ ] **Read-only mode**: `POST=0`, `DELETE=0`, `PUT=0`
- [ ] **Minimal endpoints**: Only enable necessary endpoints
- [ ] **No public exposure**: NEVER expose on the Internet
- [ ] **Advanced filters**: Configure filters adapted to your use case

### Recommended Advanced Filters

```bash
# Forbid privileged containers
export DKRPRX__CONTAINERS__DENY_PRIVILEGED="true"

# Forbid host network
export DKRPRX__CONTAINERS__DENY_HOST_NETWORK="true"

# Forbid mounting sensitive directories
export DKRPRX__VOLUMES__DENIED_PATHS="^/etc/.*,^/root/.*,^/sys/.*,^/proc/.*"

# Allow only images from private registry
export DKRPRX__CONTAINERS__ALLOWED_IMAGES="^registry.private.com/.*"

# Forbid :latest tag
export DKRPRX__IMAGES__DENIED_TAGS="^latest$"
```

## 🔍 Audit and Monitoring

### Security Logs

The proxy logs all blocked attempts:

```json
{
  "level": "warn",
  "msg": "Volume creation denied",
  "reason": "host path is denied: /var/run/docker.sock",
  "path": "/v1.43/volumes/create",
  "method": "POST"
}
```

### Recommended Monitoring

Monitor these events in logs:
- `"denied"` - Blocked attempts
- `"forbidden"` - Refused access
- `"privileged"` - Privileged container attempts
- `"docker.sock"` - Socket mounting attempts

## 📚 References

- [ADVANCED_FILTERS.md](ADVANCED_FILTERS.md) - Complete filter documentation
- [ENV_FILTERS.md](ENV_FILTERS.md) - Configuration via environment variables
- [Docker Security Best Practices](https://docs.docker.com/engine/security/)
- [CIS Docker Benchmark](https://www.cisecurity.org/benchmark/docker)

## 🆘 Support

If you discover a security vulnerability, please report it responsibly via GitHub Issues marking it as "security".

Do NOT publicly disclose critical vulnerabilities before a fix is available.
