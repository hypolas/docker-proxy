# Configuration des Filtres via Variables d'Environnement

Les filtres avancés peuvent être configurés soit via un fichier JSON, soit via des variables d'environnement. **Les variables d'environnement sont prioritaires** sur le fichier JSON.

## 🔒 Sécurité par Défaut

**Par défaut, le montage du socket Docker est interdit** pour prévenir l'escalade de privilèges.

Chemins bloqués automatiquement :
- `/var/run/docker.sock`
- `/run/docker.sock`

Pour désactiver cette protection (⚠️ non recommandé) :
```bash
export DKRPRX__DISABLE_DEFAULTS="true"
```

Pour autoriser explicitement le socket Docker :
```bash
export DKRPRX__VOLUMES__ALLOWED_PATHS="^/var/run/docker\\.sock$"
```

## 📋 Format des Variables

Toutes les variables de filtres utilisent le préfixe `DKRPRX__` suivi de la structure hiérarchique avec double underscore `__`.

### Structure générale
```
DKRPRX__<SECTION>__<PARAMETER>="value"
```

## 🎯 Variables Disponibles

### Filtres de Volumes

| Variable | Type | Description | Exemple |
|----------|------|-------------|---------|
| `DKRPRX__VOLUMES__ALLOWED_NAMES` | Array | Noms de volumes autorisés (regex) | `"^data-.*,^app-.*"` |
| `DKRPRX__VOLUMES__DENIED_NAMES` | Array | Noms de volumes interdits (regex) | `"^system-.*,^tmp-.*"` |
| `DKRPRX__VOLUMES__ALLOWED_PATHS` | Array | Chemins host autorisés (regex) | `"^/data/.*,^/mnt/.*"` |
| `DKRPRX__VOLUMES__DENIED_PATHS` | Array | Chemins host interdits (regex) | `"^/etc/.*,^/root/.*"` |
| `DKRPRX__VOLUMES__ALLOWED_DRIVERS` | Array | Drivers de volumes autorisés | `"local,nfs"` |

### Filtres de Conteneurs

| Variable | Type | Description | Exemple |
|----------|------|-------------|---------|
| `DKRPRX__CONTAINERS__ALLOWED_IMAGES` | Array | Images autorisées (regex) | `"^myregistry.com/.*"` |
| `DKRPRX__CONTAINERS__DENIED_IMAGES` | Array | Images interdites (regex) | `".*:latest$"` |
| `DKRPRX__CONTAINERS__ALLOWED_NAMES` | Array | Noms de conteneurs autorisés (regex) | `"^prod-.*,^staging-.*"` |
| `DKRPRX__CONTAINERS__DENIED_NAMES` | Array | Noms de conteneurs interdits (regex) | `"^test-.*"` |
| `DKRPRX__CONTAINERS__REQUIRE_LABELS` | Map | Labels requis (key=value) | `"env=production,team=backend"` |
| `DKRPRX__CONTAINERS__DENY_PRIVILEGED` | Bool | Interdire mode privileged | `"true"` |
| `DKRPRX__CONTAINERS__DENY_HOST_NETWORK` | Bool | Interdire host network | `"true"` |

### Filtres de Réseaux

| Variable | Type | Description | Exemple |
|----------|------|-------------|---------|
| `DKRPRX__NETWORKS__ALLOWED_NAMES` | Array | Noms de réseaux autorisés (regex) | `"^app-.*"` |
| `DKRPRX__NETWORKS__DENIED_NAMES` | Array | Noms de réseaux interdits (regex) | `"^host$"` |
| `DKRPRX__NETWORKS__ALLOWED_DRIVERS` | Array | Drivers de réseaux autorisés | `"bridge,overlay"` |

### Filtres d'Images

| Variable | Type | Description | Exemple |
|----------|------|-------------|---------|
| `DKRPRX__IMAGES__ALLOWED_REPOS` | Array | Repos/registries autorisés (regex) | `"^docker.io/.*,^myregistry.com/.*"` |
| `DKRPRX__IMAGES__DENIED_REPOS` | Array | Repos/registries interdits (regex) | `".*untrusted.*"` |
| `DKRPRX__IMAGES__ALLOWED_TAGS` | Array | Tags autorisés (regex) | `"^v[0-9]+\\.[0-9]+\\.[0-9]+$"` |
| `DKRPRX__IMAGES__DENIED_TAGS` | Array | Tags interdits (regex) | `"^latest$,^dev$"` |

## 📝 Format des Valeurs

### Arrays (Tableaux)
Séparés par virgule `,`, pipe `|` ou point-virgule `;` :

```bash
# Avec virgule (recommandé)
export DKRPRX__VOLUMES__ALLOWED_NAMES="^data-.*,^app-.*,^logs-.*"

# Avec pipe
export DKRPRX__VOLUMES__ALLOWED_NAMES="^data-.*|^app-.*|^logs-.*"

# Avec point-virgule
export DKRPRX__VOLUMES__ALLOWED_NAMES="^data-.*;^app-.*;^logs-.*"
```

### Maps (Clé-Valeur)
Format `key=value` séparés par virgule :

```bash
export DKRPRX__CONTAINERS__REQUIRE_LABELS="env=production,team=backend,app=web"
```

### Booléens
Valeurs acceptées : `true`, `false`, `1`, `0`, `yes`, `no`, `on`, `off` (case-insensitive)

```bash
export DKRPRX__CONTAINERS__DENY_PRIVILEGED="true"
export DKRPRX__CONTAINERS__DENY_HOST_NETWORK="1"
```

## 💡 Exemples Complets

### Exemple 1 : Sécurité Production

```bash
# Volumes : uniquement chemins sécurisés
export DKRPRX__VOLUMES__ALLOWED_PATHS="^/data/prod/.*"
export DKRPRX__VOLUMES__DENIED_PATHS="^/etc/.*,^/root/.*,^/home/.*"

# Conteneurs : registry privé uniquement, pas de privileged
export DKRPRX__CONTAINERS__ALLOWED_IMAGES="^registry.prod.com/.*"
export DKRPRX__CONTAINERS__DENY_PRIVILEGED="true"
export DKRPRX__CONTAINERS__DENY_HOST_NETWORK="true"
export DKRPRX__CONTAINERS__REQUIRE_LABELS="env=production,approved=true"

# Images : tags versionnés uniquement
export DKRPRX__IMAGES__DENIED_TAGS="^latest$,^dev$,^test$"
export DKRPRX__IMAGES__ALLOWED_TAGS="^v[0-9]+\\.[0-9]+\\.[0-9]+$"

# Lancer le proxy
./docker-proxy
```

### Exemple 2 : Environnement Multi-tenant

```bash
# Tenant ID
TENANT_ID="client123"

# Volumes : préfixe par tenant
export DKRPRX__VOLUMES__ALLOWED_NAMES="^${TENANT_ID}-.*"

# Conteneurs : préfixe par tenant
export DKRPRX__CONTAINERS__ALLOWED_NAMES="^${TENANT_ID}-.*"
export DKRPRX__CONTAINERS__REQUIRE_LABELS="tenant=${TENANT_ID}"

# Réseaux : préfixe par tenant
export DKRPRX__NETWORKS__ALLOWED_NAMES="^${TENANT_ID}-.*"

./docker-proxy
```

### Exemple 3 : Bloquer Volumes Sensibles

```bash
# Interdire montage de répertoires système
export DKRPRX__VOLUMES__DENIED_PATHS="^/etc/.*,^/root/.*,^/var/run/.*,^/sys/.*,^/proc/.*"

# Autoriser uniquement volumes nommés ou /data
export DKRPRX__VOLUMES__ALLOWED_PATHS="^/data/.*"
export DKRPRX__VOLUMES__ALLOWED_NAMES=".*"

./docker-proxy
```

### Exemple 4 : Interdire Tag :latest

```bash
# Refuser :latest et tags non-versionnés
export DKRPRX__IMAGES__DENIED_TAGS="^latest$"
export DKRPRX__CONTAINERS__DENIED_IMAGES=".*:latest$"

./docker-proxy
```

## 🔄 Combinaison JSON + Variables d'Environnement

Les variables d'environnement **écrasent** le fichier JSON :

**filters.json :**
```json
{
  "volumes": {
    "allowed_names": ["^data-.*"]
  },
  "containers": {
    "deny_privileged": false
  }
}
```

**Variables d'environnement :**
```bash
export FILTERS_CONFIG=./filters.json
export DKRPRX__CONTAINERS__DENY_PRIVILEGED="true"
```

**Résultat :**
- Volumes : règles du JSON (`allowed_names: ^data-.*`)
- Containers : règle de l'env var (`deny_privileged: true`) ← **prioritaire**

## 🐳 Docker Compose

```yaml
version: '3.8'

services:
  docker-proxy:
    build: .
    environment:
      # Configuration de base
      - LISTEN_ADDR=:2375
      - CONTAINERS=1
      - VOLUMES=1
      - POST=1

      # Filtres avancés via env vars
      - DKRPRX__VOLUMES__DENIED_PATHS=^/etc/.*,^/root/.*
      - DKRPRX__CONTAINERS__DENY_PRIVILEGED=true
      - DKRPRX__CONTAINERS__ALLOWED_IMAGES=^myregistry.com/.*
      - DKRPRX__IMAGES__DENIED_TAGS=^latest$$
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    ports:
      - "2375:2375"
```

## 🔍 Vérification de Configuration

Au démarrage, le proxy affiche les filtres actifs dans les logs :

```bash
INFO Detected Docker API version: v1.43
INFO Advanced filters loaded from environment variables
INFO   - Volumes: denied_paths=[^/etc/.* ^/root/.*]
INFO   - Containers: deny_privileged=true
```

## ⚙️ Ordre de Priorité

1. **Variables d'environnement** (`DKRPRX__*`) ← **Plus haute priorité**
2. **Fichier JSON** (`FILTERS_CONFIG`)
3. **Aucun filtre** (toutes opérations autorisées)

## 🚀 Migration depuis JSON

Pour migrer d'une configuration JSON vers env vars :

**Avant (JSON) :**
```json
{
  "volumes": {
    "denied_paths": ["^/etc/.*", "^/root/.*"]
  }
}
```

**Après (Env vars) :**
```bash
export DKRPRX__VOLUMES__DENIED_PATHS="^/etc/.*,^/root/.*"
```

## 💡 Bonnes Pratiques

1. **Utilisez les env vars pour** :
   - Configuration dynamique (CI/CD)
   - Secrets (via orchestrateurs)
   - Override rapide en production

2. **Utilisez le JSON pour** :
   - Configuration complexe
   - Partage entre environnements
   - Versioning dans Git

3. **Combinez les deux** :
   - JSON pour config de base
   - Env vars pour override par environnement

## 🔗 Voir Aussi

- [ADVANCED_FILTERS.md](ADVANCED_FILTERS.md) - Documentation complète des filtres
- [filters.example.json](filters.example.json) - Exemple de configuration JSON
