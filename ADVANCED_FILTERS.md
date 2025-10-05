# Filtres Avancés - Guide d'utilisation

Les filtres avancés permettent un contrôle granulaire sur les opérations Docker au-delà des simples autorisations d'endpoints.

## 🔒 Sécurité par Défaut

**Par défaut, le montage du socket Docker est interdit** pour éviter l'escalade de privilèges.

Chemins bloqués automatiquement :
- `/var/run/docker.sock`
- `/run/docker.sock`

Pour désactiver cette protection (⚠️ non recommandé) :
```bash
export DKRPRX__DISABLE_DEFAULTS="true"
```

## 📋 Configuration

Les filtres avancés se configurent via un fichier JSON et la variable d'environnement `FILTERS_CONFIG`.

```bash
export FILTERS_CONFIG=/path/to/filters.json
./docker-proxy
```

## 🎯 Types de filtres disponibles

### 1. Filtres de Volumes

Contrôlez quels volumes peuvent être créés ou montés :

```json
{
  "volumes": {
    "allowed_names": ["^data-.*", "^app-.*"],
    "denied_names": ["^system-.*"],
    "allowed_paths": ["^/data/.*", "^/mnt/volumes/.*"],
    "denied_paths": ["^/etc/.*", "^/root/.*"],
    "allowed_drivers": ["local", "nfs"]
  }
}
```

**Paramètres :**
- `allowed_names` : Patterns regex des noms de volumes autorisés
- `denied_names` : Patterns regex des noms de volumes interdits
- `allowed_paths` : Patterns regex des chemins host autorisés pour les bind mounts
- `denied_paths` : Patterns regex des chemins host interdits
- `allowed_drivers` : Liste des drivers de volumes autorisés

**Exemples :**
```bash
# ✅ Autorisé : nom commence par "data-"
docker volume create data-mysql

# ❌ Refusé : nom commence par "system-"
docker volume create system-config

# ✅ Autorisé : chemin dans /data/
docker run -v /data/mysql:/var/lib/mysql mysql

# ❌ Refusé : chemin dans /etc/
docker run -v /etc/passwd:/etc/passwd ubuntu
```

### 2. Filtres de Conteneurs

Contrôlez quels conteneurs peuvent être créés :

```json
{
  "containers": {
    "allowed_images": ["^docker.io/.*", "^myregistry.com/.*"],
    "denied_images": [".*:latest$"],
    "allowed_names": ["^prod-.*", "^staging-.*"],
    "denied_names": ["^test-.*"],
    "require_labels": {
      "env": "production",
      "team": "backend"
    },
    "deny_privileged": true,
    "deny_host_network": true
  }
}
```

**Paramètres :**
- `allowed_images` : Patterns regex des images autorisées
- `denied_images` : Patterns regex des images interdites
- `allowed_names` : Patterns regex des noms de conteneurs autorisés
- `denied_names` : Patterns regex des noms de conteneurs interdits
- `require_labels` : Labels obligatoires (clé-valeur exacte)
- `deny_privileged` : Interdire les conteneurs privilégiés
- `deny_host_network` : Interdire le mode réseau host

**Exemples :**
```bash
# ✅ Autorisé : image du registry autorisé
docker run myregistry.com/app:v1.0.0

# ❌ Refusé : tag :latest interdit
docker run nginx:latest

# ❌ Refusé : mode privileged interdit
docker run --privileged ubuntu
```

### 3. Filtres de Réseaux

Contrôlez quels réseaux peuvent être créés :

```json
{
  "networks": {
    "allowed_names": ["^app-.*", "^service-.*"],
    "denied_names": ["^host$"],
    "allowed_drivers": ["bridge", "overlay"]
  }
}
```

**Paramètres :**
- `allowed_names` : Patterns regex des noms de réseaux autorisés
- `denied_names` : Patterns regex des noms de réseaux interdits
- `allowed_drivers` : Liste des drivers de réseaux autorisés

**Exemples :**
```bash
# ✅ Autorisé : nom commence par "app-"
docker network create app-backend

# ❌ Refusé : driver non autorisé
docker network create --driver macvlan my-net
```

### 4. Filtres d'Images

Contrôlez quelles images peuvent être pull/build :

```json
{
  "images": {
    "allowed_repos": ["^docker.io/library/.*", "^myregistry.com/.*"],
    "denied_repos": [".*untrusted.*"],
    "allowed_tags": ["^v[0-9]+\\.[0-9]+\\.[0-9]+$", "^stable$"],
    "denied_tags": ["^latest$", "^dev$"]
  }
}
```

**Paramètres :**
- `allowed_repos` : Patterns regex des repos/registries autorisés
- `denied_repos` : Patterns regex des repos/registries interdits
- `allowed_tags` : Patterns regex des tags autorisés
- `denied_tags` : Patterns regex des tags interdits

**Exemples :**
```bash
# ✅ Autorisé : registry autorisé avec tag semver
docker pull myregistry.com/app:v1.2.3

# ❌ Refusé : tag :latest interdit
docker pull nginx:latest

# ❌ Refusé : registry non autorisé
docker pull untrustedregistry.com/malware:v1
```

## 📝 Exemples complets

### Exemple 1 : Environnement de production strict

```json
{
  "volumes": {
    "allowed_paths": ["^/data/prod/.*"],
    "denied_paths": ["^/.*"],
    "allowed_drivers": ["local"]
  },
  "containers": {
    "allowed_images": ["^registry.prod.com/.*"],
    "require_labels": {
      "env": "production",
      "approved": "true"
    },
    "deny_privileged": true,
    "deny_host_network": true
  },
  "images": {
    "allowed_repos": ["^registry.prod.com/.*"],
    "allowed_tags": ["^v[0-9]+\\.[0-9]+\\.[0-9]+$"]
  }
}
```

### Exemple 2 : Environnement de développement permissif

```json
{
  "volumes": {
    "denied_paths": ["^/etc/.*", "^/root/.*", "^/var/run/.*"]
  },
  "containers": {
    "denied_images": [".*malware.*"],
    "deny_privileged": false
  }
}
```

### Exemple 3 : Multi-tenant avec isolation

```json
{
  "volumes": {
    "allowed_names": ["^tenant-${TENANT_ID}-.*"]
  },
  "containers": {
    "allowed_names": ["^${TENANT_ID}-.*"],
    "require_labels": {
      "tenant": "${TENANT_ID}"
    }
  },
  "networks": {
    "allowed_names": ["^${TENANT_ID}-.*"]
  }
}
```

## 🔧 Patterns Regex utiles

| Pattern | Description | Exemple |
|---------|-------------|---------|
| `^prod-.*` | Commence par "prod-" | prod-web, prod-db |
| `.*-staging$` | Se termine par "-staging" | app-staging |
| `^v[0-9]+\.[0-9]+\.[0-9]+$` | Version semver | v1.2.3 |
| `^docker\.io/.*` | Registry Docker Hub | docker.io/nginx |
| `^/data/.*` | Chemin dans /data | /data/mysql |
| `^(prod\|staging)-.*` | prod- OU staging- | prod-web, staging-api |

## ⚙️ Utilisation

```bash
# 1. Créer votre fichier de filtres
cat > filters.json << 'EOF'
{
  "volumes": {
    "allowed_names": ["^data-.*"],
    "denied_paths": ["^/etc/.*", "^/root/.*"]
  },
  "containers": {
    "deny_privileged": true
  }
}
EOF

# 2. Configurer et démarrer le proxy
export FILTERS_CONFIG=./filters.json
export CONTAINERS=1
export VOLUMES=1
export POST=1
./docker-proxy

# 3. Les filtres sont appliqués automatiquement
docker volume create data-mysql      # ✅ Autorisé
docker volume create system-config   # ❌ Refusé
docker run --privileged ubuntu       # ❌ Refusé
```

## 🔍 Ordre d'évaluation

Les filtres sont évalués dans cet ordre :

1. **Liste noire (denied_*)** : Si un pattern correspond, refus immédiat
2. **Liste blanche (allowed_*)** : Si définie, doit correspondre
3. **Contraintes spécifiques** : Labels requis, privileged, etc.

Si aucun filtre ne s'applique, l'opération est **autorisée par défaut**.

## 🚨 Messages d'erreur

Lorsqu'une opération est refusée, vous recevrez un message explicite :

```json
{
  "message": "Volume creation denied by advanced filter",
  "reason": "host path is denied: /etc/passwd"
}
```

## 💡 Bonnes pratiques

1. **Commencez permissif** : Activez les filtres progressivement
2. **Testez en dev** : Validez vos patterns avant la production
3. **Loggez tout** : Surveillez les logs pour comprendre les refus
4. **Utilisez des patterns précis** : Évitez `.*` qui matche tout
5. **Documentez vos règles** : Ajoutez des commentaires dans le JSON
6. **Versionnez vos filtres** : Git pour tracer les changements

## 🔗 Intégration avec les ACL basiques

Les filtres avancés fonctionnent **en complément** des ACL basiques :

1. ACL basique vérifie l'accès à l'endpoint (ex: `VOLUMES=1`)
2. Filtre avancé vérifie le contenu de la requête (ex: chemin autorisé)

Les deux doivent passer pour que l'opération soit autorisée.
