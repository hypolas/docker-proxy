# 🚀 Quick Start Deployment Guide

This guide helps you deploy dockershield quickly for different use cases.

## 📦 Prerequisites

- Docker installed
- Docker Compose (optional, for easy deployment)
- Access to Docker socket (usually `/var/run/docker.sock`)

## 🎯 Use Case 1: Basic Read-Only Access (Safest)

Perfect for monitoring tools, dashboards, CI/CD pipelines that only need to read Docker state.

```bash
docker run -d \
  --name dockershield \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -p 2375:2375 \
  -e CONTAINERS=1 \
  -e IMAGES=1 \
  -e NETWORKS=1 \
  -e VOLUMES=1 \
  -e INFO=1 \
  hypolas/proxy-docker:latest
```

**What this does**:
- ✅ Allows reading container, image, network, and volume information
- ❌ Blocks all write operations (POST, PUT, DELETE)
- 🔒 Security by default: Blocks Docker socket mounting

**Test it**:
```bash
export DOCKER_HOST=tcp://localhost:2375
docker ps        # ✅ Works
docker images    # ✅ Works
docker rm XXX    # ❌ Forbidden
```

## 🎯 Use Case 2: CI/CD Build & Deploy

For CI/CD pipelines that need to build images and manage containers.

```bash
docker run -d \
  --name dockershield \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -p 2375:2375 \
  -e BUILD=1 \
  -e IMAGES=1 \
  -e CONTAINERS=1 \
  -e NETWORKS=1 \
  -e POST=1 \
  -e DELETE=1 \
  -e PROXY_CONTAINER_NAME=dockershield \
  hypolas/proxy-docker:latest
```

**What this does**:
- ✅ Allows building images
- ✅ Allows creating/removing containers
- ✅ Allows creating/removing networks
- 🔒 Protects the proxy container itself from manipulation
- 🔒 Blocks Docker socket mounting by default

**GitHub Actions example**:
```yaml
services:
  docker:
    image: your-registry/dockershield:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - /tmp:/tmp
    env:
      BUILD: 1
      IMAGES: 1
      CONTAINERS: 1
      POST: 1
      DELETE: 1
      LISTEN_SOCKET: unix:///tmp/dockershield.sock
```

## 🎯 Use Case 3: Advanced Filtering (Production)

For production environments with strict security requirements.

**docker-compose.yml**:
```yaml
version: '3.8'

services:
  dockershield:
    image: your-registry/dockershield:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - /tmp:/tmp
      - ./filters.json:/etc/dockershield/filters.json:ro
    environment:
      # Basic access
      CONTAINERS: 1
      IMAGES: 1
      NETWORKS: 1
      POST: 1
      DELETE: 1

      # Advanced filters
      FILTERS_CONFIG: /etc/dockershield/filters.json
      LISTEN_SOCKET: unix:///tmp/dockershield.sock

      # Self-protection
      PROXY_CONTAINER_NAME: dockershield
      PROXY_NETWORK_NAME: proxy_network

      # Additional env-based filters
      DKRPRX__IMAGES__ALLOWED_REGISTRIES: "docker.io,ghcr.io"
      DKRPRX__IMAGES__DENIED_TAGS: "latest,master"
      DKRPRX__VOLUMES__DENIED_PATHS: "/var/run/docker.sock,/run/docker.sock,/etc,/root"

    networks:
      - proxy_network

networks:
  proxy_network:
    name: proxy_network
```

**filters.json**:
```json
{
  "images": {
    "allowed_registries": ["docker.io", "ghcr.io", "gcr.io"],
    "denied_tags": ["^latest$", "^master$"],
    "allowed_architectures": ["amd64", "arm64"]
  },
  "volumes": {
    "denied_paths": [
      "^/var/run/docker\\.sock$",
      "^/run/docker\\.sock$",
      "^/etc/.*",
      "^/root/.*"
    ],
    "allowed_drivers": ["local"]
  },
  "containers": {
    "denied_images": [".*:latest$"],
    "required_labels": {
      "environment": "production"
    }
  },
  "networks": {
    "allowed_drivers": ["bridge", "overlay"],
    "denied_subnets": ["10.0.0.0/8"]
  }
}
```

**Deploy**:
```bash
docker-compose up -d
```

## 🎯 Use Case 4: Multi-Tenant Environment

For hosting providers or multi-tenant platforms.

```bash
docker run -d \
  --name dockershield-tenant1 \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v /tmp:/tmp \
  -e LISTEN_SOCKET=unix:///tmp/dockershield.sock \
  -e CONTAINERS=1 \
  -e IMAGES=1 \
  -e NETWORKS=1 \
  -e POST=1 \
  -e DELETE=1 \
  -e DKRPRX__CONTAINERS__ALLOWED_NAMES="^tenant1-.*" \
  -e DKRPRX__NETWORKS__ALLOWED_NAMES="^tenant1-.*" \
  -e DKRPRX__IMAGES__ALLOWED_REGISTRIES="registry.tenant1.com" \
  -e PROXY_CONTAINER_NAME=dockershield-tenant1 \
  hypolas/proxy-docker:latest
```

**What this does**:
- 🔐 Isolates tenant resources by naming patterns
- 🔐 Restricts to tenant-specific registry
- 🔐 Prevents cross-tenant access

## 🎯 Use Case 5: Unix Socket Mode

For local development or sidecar patterns.

```bash
docker run -d \
  --name dockershield \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v /tmp:/tmp \
  -e LISTEN_SOCKET=unix:///tmp/dockershield.sock \
  -e SOCKET_PERMS=0666 \
  -e CONTAINERS=1 \
  -e IMAGES=1 \
  -e INFO=1 \
  hypolas/proxy-docker:latest
```

**Usage**:
```bash
export DOCKER_HOST=unix:///tmp/dockershield.sock
docker ps
```

## 🔧 Configuration Cheat Sheet

### Common Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LISTEN_ADDR` | `:2375` | TCP address to listen on |
| `LISTEN_SOCKET` | - | Unix socket path (overrides LISTEN_ADDR) |
| `DOCKER_SOCKET` | `unix:///var/run/docker.sock` | Docker socket to proxy to |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `SOCKET_PERMS` | `0666` | Unix socket permissions (octal) |

### Access Control (Default: Denied)

| Variable | Description |
|----------|-------------|
| `CONTAINERS=1` | Allow container operations |
| `IMAGES=1` | Allow image operations |
| `VOLUMES=1` | Allow volume operations |
| `NETWORKS=1` | Allow network operations |
| `BUILD=1` | Allow image building |
| `EXEC=1` | Allow container exec |
| `POST=1` | Allow POST requests (create operations) |
| `DELETE=1` | Allow DELETE requests (remove operations) |
| `PUT=1` | Allow PUT/PATCH requests (update operations) |

### Security Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PROXY_CONTAINER_NAME` | `dockershield` | Name of proxy container (for self-protection) |
| `PROXY_NETWORK_NAME` | - | Name of proxy network (for protection) |
| `DKRPRX__DISABLE_DEFAULTS` | `false` | Disable security defaults |

### Advanced Filters

See [ENV_FILTERS.md](ENV_FILTERS.md) for complete reference.

**Examples**:
```bash
# Restrict image registries
DKRPRX__IMAGES__ALLOWED_REGISTRIES=docker.io,ghcr.io

# Block specific tags
DKRPRX__IMAGES__DENIED_TAGS=latest,master,dev

# Block volume paths
DKRPRX__VOLUMES__DENIED_PATHS=/etc,/root,/var/run/docker.sock

# Require labels
DKRPRX__CONTAINERS__REQUIRED_LABELS=environment=production,team=devops

# Restrict networks
DKRPRX__NETWORKS__ALLOWED_DRIVERS=bridge,overlay
```

## 🧪 Testing Your Deployment

### 1. Check Proxy Health
```bash
curl http://localhost:2375/_ping
# Expected: OK

curl http://localhost:2375/version
# Expected: Docker version info
```

### 2. Test Access Control
```bash
export DOCKER_HOST=tcp://localhost:2375

# Should work (if CONTAINERS=1)
docker ps

# Should fail (if POST=0)
docker run hello-world
# Expected: 403 Forbidden
```

### 3. Test Advanced Filters
```bash
# Should fail (blocked socket)
docker run -v /var/run/docker.sock:/var/run/docker.sock alpine
# Expected: 403 Forbidden: Volume path /var/run/docker.sock is denied

# Should fail (if denying :latest)
docker pull nginx:latest
# Expected: 403 Forbidden: Image tag latest is denied
```

### 4. Check Logs
```bash
docker logs dockershield
```

## 🆘 Troubleshooting

### Issue: "Permission denied" accessing Docker socket

**Solution**: Ensure socket is mounted with read permissions:
```bash
-v /var/run/docker.sock:/var/run/docker.sock:ro
```

### Issue: "403 Forbidden" on all requests

**Cause**: Access rules too restrictive.

**Solution**: Enable required endpoints:
```bash
-e CONTAINERS=1 -e IMAGES=1 -e POST=1
```

### Issue: Advanced filters not working

**Check**:
1. Env vars are properly formatted: `DKRPRX__SECTION__PARAMETER`
2. Array separators: use comma, pipe, or semicolon
3. Regex patterns are valid (test with `grep -P 'pattern'`)

**Debug**:
```bash
docker logs dockershield | grep -i filter
```

### Issue: Security defaults blocking legitimate operations

**Temporary override** (use with caution):
```bash
-e DKRPRX__DISABLE_DEFAULTS=true
```

**Better solution**: Customize filters to allow specific operations:
```bash
-e DKRPRX__VOLUMES__ALLOWED_PATHS=/my/safe/path
```

## 🔒 Security Best Practices

1. ✅ **Always mount Docker socket as read-only** (`:ro`)
2. ✅ **Use Unix sockets when possible** (better isolation)
3. ✅ **Enable only required endpoints** (principle of least privilege)
4. ✅ **Use advanced filters in production** (defense in depth)
5. ✅ **Set PROXY_CONTAINER_NAME** (prevent self-manipulation)
6. ✅ **Monitor proxy logs** (detect suspicious activity)
7. ✅ **Use network isolation** (dedicated network for proxy)
8. ✅ **Regularly update** (security patches)

## 📚 Next Steps

- Read [SECURITY.md](SECURITY.md) for security guidelines
- Explore [ADVANCED_FILTERS.md](ADVANCED_FILTERS.md) for complex filtering
- Check [CICD_EXAMPLES.md](CICD_EXAMPLES.md) for CI/CD integration
- Review [ENV_FILTERS.md](ENV_FILTERS.md) for complete env var reference

## 📞 Support

- **Issues**: https://github.com/hypolas/dockershield/issues
- **Email**: nicolas.hypolite@gmail.com
- **Documentation**: https://github.com/hypolas/dockershield

---

**Pro Tip**: Start with the most restrictive configuration and gradually enable features as needed. It's easier to add permissions than to remove them after a security incident! 🛡️
