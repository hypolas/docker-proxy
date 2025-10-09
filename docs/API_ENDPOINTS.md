# Docker API Endpoints Reference

This document explains what each environment variable allows in terms of Docker API access.

## 🔑 Access Control Model

When you set an environment variable like `CONTAINERS=1`, you enable access to **ALL endpoints** starting with `/containers`.

### Default Behavior

- **GET, HEAD**: Always allowed for enabled endpoints (read-only)
- **POST, PUT, PATCH, DELETE**: Require explicit permission via `POST=1`, `PUT=1`, `DELETE=1`

## 📋 Complete Endpoint Mapping

### `CONTAINERS=1` - Container Management

Enables access to **all** `/containers/*` endpoints:

#### Read Operations (GET/HEAD always allowed)
- `GET /containers/json` - List containers
- `GET /containers/{id}/json` - Inspect container
- `GET /containers/{id}/logs` - Get container logs
- `GET /containers/{id}/stats` - Get container stats (streaming)
- `GET /containers/{id}/top` - List processes in container
- `GET /containers/{id}/changes` - Get filesystem changes
- `GET /containers/{id}/export` - Export container filesystem
- `GET /containers/{id}/archive` - Get files/folders from container
- `HEAD /containers/{id}/archive` - Check if path exists in container

#### Write Operations (require `POST=1`, `DELETE=1`, `PUT=1`)
- `POST /containers/create` - Create container (needs `POST=1`)
- `POST /containers/{id}/start` - Start container (needs `POST=1`)
- `POST /containers/{id}/stop` - Stop container (needs `POST=1`)
- `POST /containers/{id}/restart` - Restart container (needs `POST=1`)
- `POST /containers/{id}/kill` - Kill container (needs `POST=1`)
- `POST /containers/{id}/pause` - Pause container (needs `POST=1`)
- `POST /containers/{id}/unpause` - Unpause container (needs `POST=1`)
- `POST /containers/{id}/wait` - Wait for container to stop (needs `POST=1`)
- `POST /containers/{id}/resize` - Resize container TTY (needs `POST=1`)
- `POST /containers/{id}/attach` - Attach to container (needs `POST=1`)
- `POST /containers/{id}/exec` - Create exec instance (needs `POST=1`)
- `POST /containers/{id}/rename` - Rename container (needs `POST=1`)
- `POST /containers/{id}/update` - Update container resources (needs `POST=1`)
- `POST /containers/prune` - Prune unused containers (needs `POST=1`)
- `PUT /containers/{id}/archive` - Extract files to container (needs `PUT=1`)
- `DELETE /containers/{id}` - Remove container (needs `DELETE=1`)

**Total**: ~25+ endpoints

---

### `IMAGES=1` - Image Management

Enables access to **all** `/images/*` endpoints:

#### Read Operations
- `GET /images/json` - List images
- `GET /images/{name}/json` - Inspect image
- `GET /images/{name}/history` - Get image history
- `GET /images/search` - Search images
- `GET /images/{name}/get` - Export image (tarball)
- `GET /images/get` - Export multiple images

#### Write Operations
- `POST /images/create` - Pull image (needs `POST=1`)
- `POST /images/{name}/push` - Push image (needs `POST=1`)
- `POST /images/{name}/tag` - Tag image (needs `POST=1`)
- `POST /images/load` - Import images (needs `POST=1`)
- `POST /images/prune` - Prune unused images (needs `POST=1`)
- `POST /build` - Build image from Dockerfile (needs `POST=1` + `BUILD=1`)
- `DELETE /images/{name}` - Remove image (needs `DELETE=1`)

**Total**: ~15+ endpoints

---

### `NETWORKS=1` - Network Management

Enables access to **all** `/networks/*` endpoints:

#### Read Operations
- `GET /networks` - List networks
- `GET /networks/{id}` - Inspect network

#### Write Operations
- `POST /networks/create` - Create network (needs `POST=1`)
- `POST /networks/{id}/connect` - Connect container to network (needs `POST=1`)
- `POST /networks/{id}/disconnect` - Disconnect container (needs `POST=1`)
- `POST /networks/prune` - Prune unused networks (needs `POST=1`)
- `DELETE /networks/{id}` - Remove network (needs `DELETE=1`)

**Total**: ~7 endpoints

---

### `VOLUMES=1` - Volume Management

Enables access to **all** `/volumes/*` endpoints:

#### Read Operations
- `GET /volumes` - List volumes
- `GET /volumes/{name}` - Inspect volume

#### Write Operations
- `POST /volumes/create` - Create volume (needs `POST=1`)
- `POST /volumes/prune` - Prune unused volumes (needs `POST=1`)
- `DELETE /volumes/{name}` - Remove volume (needs `DELETE=1`)

**Total**: ~5 endpoints

---

### `EXEC=1` - Exec Management

Enables access to **all** `/exec/*` endpoints:

#### Read Operations
- `GET /exec/{id}/json` - Inspect exec instance

#### Write Operations
- `POST /exec/{id}/start` - Start exec instance (needs `POST=1`)
- `POST /exec/{id}/resize` - Resize exec TTY (needs `POST=1`)

**Note**: Creating exec requires `CONTAINERS=1` + `POST=1` (via `/containers/{id}/exec`)

**Total**: ~3 endpoints

---

### `SERVICES=1` - Swarm Services

Enables access to **all** `/services/*` endpoints (Swarm mode):

#### Read Operations
- `GET /services` - List services
- `GET /services/{id}` - Inspect service
- `GET /services/{id}/logs` - Get service logs

#### Write Operations
- `POST /services/create` - Create service (needs `POST=1`)
- `POST /services/{id}/update` - Update service (needs `POST=1`)
- `DELETE /services/{id}` - Remove service (needs `DELETE=1`)

**Total**: ~6 endpoints

---

### `TASKS=1` - Swarm Tasks

Enables access to **all** `/tasks/*` endpoints (Swarm mode):

#### Read Operations
- `GET /tasks` - List tasks
- `GET /tasks/{id}` - Inspect task
- `GET /tasks/{id}/logs` - Get task logs

**Total**: ~3 endpoints

---

### `NODES=1` - Swarm Nodes

Enables access to **all** `/nodes/*` endpoints (Swarm mode):

#### Read Operations
- `GET /nodes` - List nodes
- `GET /nodes/{id}` - Inspect node

#### Write Operations
- `POST /nodes/{id}/update` - Update node (needs `POST=1`)
- `DELETE /nodes/{id}` - Remove node (needs `DELETE=1`)

**Total**: ~4 endpoints

---

### `SWARM=1` - Swarm Management

Enables access to **all** `/swarm/*` endpoints:

#### Read Operations
- `GET /swarm` - Inspect swarm

#### Write Operations
- `POST /swarm/init` - Initialize swarm (needs `POST=1`)
- `POST /swarm/join` - Join swarm (needs `POST=1`)
- `POST /swarm/leave` - Leave swarm (needs `POST=1`)
- `POST /swarm/update` - Update swarm (needs `POST=1`)
- `GET /swarm/unlockkey` - Get unlock key

**Total**: ~6 endpoints

---

### `SECRETS=1` - Docker Secrets (Swarm)

Enables access to **all** `/secrets/*` endpoints:

#### Read Operations
- `GET /secrets` - List secrets
- `GET /secrets/{id}` - Inspect secret

#### Write Operations
- `POST /secrets/create` - Create secret (needs `POST=1`)
- `POST /secrets/{id}/update` - Update secret (needs `POST=1`)
- `DELETE /secrets/{id}` - Remove secret (needs `DELETE=1`)

**Total**: ~5 endpoints

---

### `CONFIGS=1` - Docker Configs (Swarm)

Enables access to **all** `/configs/*` endpoints:

#### Read Operations
- `GET /configs` - List configs
- `GET /configs/{id}` - Inspect config

#### Write Operations
- `POST /configs/create` - Create config (needs `POST=1`)
- `POST /configs/{id}/update` - Update config (needs `POST=1`)
- `DELETE /configs/{id}` - Remove config (needs `DELETE=1`)

**Total**: ~5 endpoints

---

### `PLUGINS=1` - Plugin Management

Enables access to **all** `/plugins/*` endpoints:

#### Read Operations
- `GET /plugins` - List plugins
- `GET /plugins/{name}/json` - Inspect plugin
- `GET /plugins/privileges` - Get plugin privileges

#### Write Operations
- `POST /plugins/pull` - Install plugin (needs `POST=1`)
- `POST /plugins/{name}/enable` - Enable plugin (needs `POST=1`)
- `POST /plugins/{name}/disable` - Disable plugin (needs `POST=1`)
- `POST /plugins/{name}/upgrade` - Upgrade plugin (needs `POST=1`)
- `POST /plugins/{name}/push` - Push plugin (needs `POST=1`)
- `POST /plugins/{name}/set` - Configure plugin (needs `POST=1`)
- `POST /plugins/create` - Create plugin (needs `POST=1`)
- `DELETE /plugins/{name}` - Remove plugin (needs `DELETE=1`)

**Total**: ~11 endpoints

---

### Other Endpoints

#### `BUILD=1` - Build Images
- `POST /build` - Build image from Dockerfile (also needs `POST=1`)

#### `COMMIT=1` - Commit Containers
- `POST /commit` - Create image from container (needs `POST=1`)

#### `INFO=1` - System Info
- `GET /info` - Get system information

#### `VERSION=1` - Docker Version
- `GET /version` - Get Docker version (always allowed by default)

#### `PING=1` - Health Check
- `GET /_ping` - Ping daemon (always allowed by default)
- `HEAD /_ping` - Ping daemon HEAD

#### `EVENTS=1` - System Events
- `GET /events` - Stream events (always allowed by default)

#### `AUTH=1` - Registry Authentication
- `POST /auth` - Check registry authentication (needs `POST=1`)

#### `DISTRIBUTION=1` - Image Distribution
- `GET /distribution/{name}/json` - Get image distribution info

#### `SYSTEM=1` - System Operations
- `GET /system/df` - Get data usage information
- `POST /system/prune` - Prune unused data (needs `POST=1`)

#### `SESSION=1` - Experimental Session
- `POST /session` - Create session (needs `POST=1`)

---

## 🎯 Common Configuration Examples

### Example 1: Read-Only Container Monitoring
```bash
CONTAINERS=1
IMAGES=1
# GET/HEAD only - no modifications possible
```

**Allows**:
- List and inspect containers
- List and inspect images
- Get logs, stats, etc.

**Blocks**:
- Create, start, stop containers
- Pull, push, delete images

---

### Example 2: Full Container Management
```bash
CONTAINERS=1
IMAGES=1
POST=1
DELETE=1
```

**Allows**:
- Everything from Example 1, PLUS:
- Create, start, stop, remove containers
- Pull, remove images
- Restart, kill, pause containers

**Blocks**:
- PUT operations (like uploading files to containers)

---

### Example 3: CI/CD Pipeline
```bash
CONTAINERS=1
IMAGES=1
NETWORKS=1
VOLUMES=1
BUILD=1
POST=1
DELETE=1
PUT=1
```

**Allows**:
- Full container lifecycle
- Build images
- Create/remove networks and volumes
- Upload files to containers

---

### Example 4: Swarm Orchestration
```bash
SERVICES=1
TASKS=1
NODES=1
SWARM=1
SECRETS=1
CONFIGS=1
POST=1
DELETE=1
PUT=1
```

**Allows**:
- Full Swarm management
- Service deployment
- Secret/config management

---

## 🔒 Security Best Practices

### Principle of Least Privilege

1. **Start minimal**: Only enable what you need
   ```bash
   CONTAINERS=1  # Just listing
   ```

2. **Add write access carefully**:
   ```bash
   CONTAINERS=1
   POST=1  # Now can start/stop
   ```

3. **Avoid wildcards**: Don't enable everything
   ```bash
   # ❌ BAD: Too permissive
   CONTAINERS=1 IMAGES=1 VOLUMES=1 NETWORKS=1 POST=1 DELETE=1 PUT=1

   # ✅ GOOD: Only what's needed
   CONTAINERS=1
   POST=1
   ```

### Read-Only Mode

For monitoring/observability tools:
```bash
CONTAINERS=1
IMAGES=1
INFO=1
EVENTS=1
# NO POST/DELETE/PUT = Read-only
```

### Audit Logging

Enable logging to track API calls:
```bash
LOG_LEVEL=debug
```

---

## 📚 Reference

- **Source**: `pkg/rules/matcher.go:63`
- **HTTP Methods**: `pkg/rules/matcher.go:32`
- **Path Matching**: Regex-based, prefix matching

---

**Version**: 1.0
**Last Updated**: 2025-10-06
