# Docker Socket Proxy - Sécurité CI/CD & Contrôle d'Accès Granulaire

> **Proxy Docker sécurisé pour CI/CD, DevOps et environnements multi-tenants** avec système de filtrage avancé et granulaire.

Un proxy de socket Docker professionnel avec **filtrage regex avancé**, conçu spécifiquement pour **sécuriser les pipelines CI/CD** (GitHub Actions, GitLab CI, Jenkins, CircleCI, etc.) et les **environnements cloud-native**. Inspiré de [Tecnativa/docker-socket-proxy](https://github.com/Tecnativa/docker-socket-proxy), implémenté en Go haute performance avec Gin et Resty.

## 🎯 Cas d'Usage Principaux

### 🔧 CI/CD & DevOps
**Idéal pour sécuriser vos pipelines CI/CD** en exposant uniquement les fonctionnalités Docker nécessaires :
- ✅ **GitHub Actions, GitLab CI, Jenkins** : Limitez les actions Docker autorisées
- ✅ **Docker-in-Docker (DinD)** sécurisé : Contrôlez build, push, run
- ✅ **Registry privé obligatoire** : Forcez l'utilisation de vos registries internes
- ✅ **Interdiction de :latest** : Imposez le versioning sémantique
- ✅ **Audit complet** : Logs structurés de toutes les opérations

### ☁️ Plateformes Cloud & Multi-tenant
- **Kubernetes, Docker Swarm, Nomad** : Isolation entre namespaces/tenants
- **PaaS & Container-as-a-Service** : Contrôle granulaire par client
- **Environnements partagés** : Sécurité et isolation stricte

### 🏢 Entreprise & Production
- **Zero-trust architecture** : Principe du moindre privilège appliqué
- **Compliance & Audit** : Traçabilité complète des opérations
- **Sécurité multi-couches** : Protection contre l'escalade de privilèges

---

## ✨ Fonctionnalités Principales

### 🚀 Pourquoi ce proxy ?

Ce proxy Docker offre une **sécurité avancée pour vos environnements Docker** grâce à :
- **Policy enforcement** et **RBAC** pour Docker API
- **Admission control** pour conteneurs avec filtres regex
- **Zero-trust architecture** appliquée au socket Docker
- **Multi-tenant isolation** avec contrôle granulaire par namespace
- **Container escape prevention** via restrictions de volumes et privilèges
- **Audit trail** complet de toutes les opérations Docker

Idéal pour **cloud-native security**, **Kubernetes**, **Docker Swarm**, **PaaS**, et **CI/CD pipelines** (GitHub Actions, GitLab CI, Jenkins, CircleCI, Azure DevOps).

### 🎯 Contrôle d'Accès Granulaire
- **ACL par endpoint** : Activez uniquement les endpoints Docker nécessaires
- **Filtres avancés avec regex** : Contrôle précis sur volumes, conteneurs, images, réseaux
- **Filtrage de contenu** : Inspectez et validez les requêtes avant de les transmettre
- **Configuration flexible** : JSON ou variables d'environnement (prioritaires)

### 🔐 Sécurité Multi-Couches
- **Protection du socket Docker** : Bloqué par défaut pour éviter l'escalade de privilèges
- **Auto-protection** : Le proxy se protège lui-même contre toute manipulation
- **Mode lecture seule** : Désactiver POST/DELETE/PUT par défaut
- **Détection automatique** : Version de l'API Docker auto-détectée

### 🛠️ Exemples de Filtres Avancés
```bash
# Interdire montage de répertoires sensibles
export DKRPRX__VOLUMES__DENIED_PATHS="^/etc/.*,^/root/.*,^/home/.*"

# Autoriser uniquement images d'un registry privé
export DKRPRX__CONTAINERS__ALLOWED_IMAGES="^registry.company.com/.*"

# Interdire tag :latest
export DKRPRX__IMAGES__DENIED_TAGS="^latest$"

# Interdire conteneurs privilégiés
export DKRPRX__CONTAINERS__DENY_PRIVILEGED="true"

# Exiger des labels spécifiques
export DKRPRX__CONTAINERS__REQUIRE_LABELS="env=production,team=backend"
```

## ⚠️ Avertissement de Responsabilité

**CE LOGICIEL EST FOURNI "TEL QUEL", SANS GARANTIE D'AUCUNE SORTE.**

Le développeur décline toute responsabilité concernant :
- Les dommages directs ou indirects causés par l'utilisation de ce logiciel
- Les failles de sécurité ou vulnérabilités
- La perte de données ou l'interruption de service
- Toute utilisation malveillante ou inappropriée

**Vous utilisez ce proxy à vos propres risques.** Il est de votre responsabilité de :
- Configurer correctement les filtres de sécurité
- Tester la configuration dans un environnement de développement
- Auditer régulièrement les accès et les logs
- Ne JAMAIS exposer le proxy sur un réseau public

## 📜 Licence Double (GPL-3.0 + Commerciale)

Ce projet utilise un **modèle de licence double** (comme Qt, MySQL, GitLab) :

### 🆓 Option 1 : GPL-3.0 (GRATUIT)

**✅ Utilisation GRATUITE si vous :**
- Partagez vos modifications (open-source)
- Distribuez le code source
- Utilisez GPL-3.0 pour votre projet

**📋 Obligations GPL-3.0 :**
- Partager vos modifications sous GPL-3.0
- Fournir le code source aux utilisateurs
- Conserver les mentions de copyright
- Propager la licence GPL-3.0

**Parfait pour :**
- Projets open-source
- Recherche & éducation
- Usage personnel
- Entreprises qui partagent leur code

---

### 💼 Option 2 : Licence Commerciale (PAYANTE)

**Licence commerciale requise si vous :**
- ❌ Ne voulez PAS partager votre code source
- ❌ Utilisez dans un produit propriétaire/fermé
- ❌ Fournissez un SaaS sans partager le code
- ❌ Intégrez dans un système embarqué fermé

**✅ Avantages licence commerciale :**
- Pas d'obligation de partager le code
- Aucune contrainte GPL
- Support prioritaire inclus
- Termes personnalisables

**💰 Tarification indicative :**

| Licence | Employés | Prix/an | Support |
|---------|----------|---------|---------|
| **Startup** | < 10 | €500 | Email 48h |
| **SME** | < 100 | €2,000 | Email 24h |
| **Enterprise** | > 100 | €10,000 | Priority 4h |
| **OEM** | Sur mesure | Custom | Dédié |

---

### 📧 Obtenir une Licence Commerciale

**Contact :**
- 📧 Email : nicolas.hypolite@gmail.com
- 🌐 Web : https://github.com/hypolas/docker-proxy
- 📄 Template : Voir [LICENSE-COMMERCIAL](LICENSE-COMMERCIAL)

**Process :**
1. Contactez-nous avec vos besoins
2. Recevez un devis personnalisé
3. Signez l'accord de licence
4. Recevez votre licence immédiatement

---

### 🔍 Quelle licence choisir ?

| Si vous voulez... | Utilisez |
|-------------------|----------|
| Contribuer à l'open-source | GPL-3.0 (gratuit) |
| Garder votre code privé | Commerciale (payant) |
| Projet personnel/éducatif | GPL-3.0 (gratuit) |
| Produit SaaS commercial | Commerciale (payant) |
| Startup qui débute | GPL-3.0 puis commercial plus tard |

---

**📖 Détails complets :** Voir [LICENSE](LICENSE) pour tous les termes juridiques.

**⚖️ Compatibilité :** Compatible avec toutes nos dépendances (Docker SDK/Apache 2.0, Gin/MIT, etc.)

## 🚀 Installation

```bash
go mod download
go build -o docker-proxy ./cmd/docker-proxy
```

## 📋 Configuration

La configuration se fait via des variables d'environnement :

### Configuration de base

| Variable | Description | Défaut |
|----------|-------------|---------|
| `LISTEN_ADDR` | Adresse d'écoute TCP | `:2375` |
| `LISTEN_SOCKET` | Chemin vers le socket Unix pour écouter (optionnel, prioritaire sur LISTEN_ADDR) | - |
| `DOCKER_SOCKET` | Chemin vers le socket Docker | `unix:///var/run/docker.sock` |
| `LOG_LEVEL` | Niveau de log (debug, info, warn, error) | `info` |
| `API_VERSION` | Version de l'API Docker (auto-détectée si non définie) | Auto-détection |
| `SOCKET_PERMS` | Permissions du socket Unix (format octal) | `0666` |

### Contrôle d'accès aux endpoints

Par défaut **autorisés** (valeur: `1`) :
- `EVENTS` - Événements Docker
- `PING` - Healthcheck
- `VERSION` - Version de Docker

Par défaut **refusés** (valeur: `0`), doivent être activés explicitement :
- `AUTH` - Authentification
- `BUILD` - Construction d'images
- `COMMIT` - Commit de conteneurs
- `CONFIGS` - Configurations Swarm
- `CONTAINERS` - Gestion des conteneurs
- `DISTRIBUTION` - Distribution d'images
- `EXEC` - Exécution de commandes
- `IMAGES` - Gestion des images
- `INFO` - Informations système
- `NETWORKS` - Gestion des réseaux
- `NODES` - Nœuds Swarm
- `PLUGINS` - Plugins Docker
- `SECRETS` - Secrets Swarm
- `SERVICES` - Services Swarm
- `SESSION` - Sessions
- `SWARM` - Swarm mode
- `SYSTEM` - Système Docker
- `TASKS` - Tâches Swarm
- `VOLUMES` - Gestion des volumes

### Méthodes HTTP

- `GET`, `HEAD` : Toujours autorisées (lecture seule)
- `POST` : Défaut `0` (variable `POST=1` pour activer)
- `DELETE` : Défaut `0` (variable `DELETE=1` pour activer)
- `PUT`, `PATCH` : Défaut `0` (variable `PUT=1` pour activer)

## 💡 Exemples d'utilisation

### Mode lecture seule (défaut)

```bash
export CONTAINERS=1
export IMAGES=1
./docker-proxy
```

### Mode lecture/écriture

```bash
export CONTAINERS=1
export IMAGES=1
export POST=1
export DELETE=1
./docker-proxy
```

### Écoute sur Unix socket

```bash
# Écoute sur Unix socket au lieu de TCP
export LISTEN_SOCKET=/tmp/docker-proxy.sock
export CONTAINERS=1
export IMAGES=1
./docker-proxy

# Test avec curl
curl --unix-socket /tmp/docker-proxy.sock http://localhost/v1.41/containers/json
```

> ℹ️ `LISTEN_SOCKET` prend toujours le pas sur `LISTEN_ADDR`. Pensez à ajuster `SOCKET_PERMS` si le socket doit être partagé avec d’autres utilisateurs (ex. `export SOCKET_PERMS=0660`).

```mermaid
graph LR
    subgraph "CI/CD Runner"
        A[Pipeline Step\n(docker build/push)]
    end

    subgraph "Proxy Host"
        B[Docker Proxy\nLISTEN_SOCKET=/tmp/docker-proxy.sock]
        C[(Unix Socket\n/tmp/docker-proxy.sock)]
    end

    subgraph "Docker Host"
        D[(Docker Socket\n/var/run/docker.sock)]
        E[Docker Engine]
    end

    A -- "DOCKER_HOST=unix:///tmp/docker-proxy.sock" --> B
    B -- binds --> C
    C -- "proxy traffic" --> D
    D -- native socket --> E
```

### 🔧 Intégration CI/CD

#### GitHub Actions

```yaml
# .github/workflows/docker.yml
name: Docker Build
on: [push]

services:
  docker-proxy:
    image: your-registry/docker-proxy:latest
    env:
      CONTAINERS: 1
      IMAGES: 1
      BUILD: 1
      POST: 1
      # Forcer registry privé
      DKRPRX__CONTAINERS__ALLOWED_IMAGES: "^registry.company.com/.*"
      # Interdire :latest
      DKRPRX__IMAGES__DENIED_TAGS: "^latest$"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build Docker image
        env:
          DOCKER_HOST: tcp://docker-proxy:2375
        run: docker build -t registry.company.com/app:${{ github.sha }} .
```

#### GitLab CI

```yaml
# .gitlab-ci.yml
variables:
  DOCKER_HOST: tcp://docker-proxy:2375

services:
  - name: your-registry/docker-proxy:latest
    alias: docker-proxy
    variables:
      CONTAINERS: "1"
      IMAGES: "1"
      BUILD: "1"
      POST: "1"
      DKRPRX__CONTAINERS__ALLOWED_IMAGES: "^registry.company.com/.*"

build:
  script:
    - docker build -t registry.company.com/app:$CI_COMMIT_SHA .
    - docker push registry.company.com/app:$CI_COMMIT_SHA
```

#### Jenkins Pipeline

```groovy
pipeline {
  agent any
  environment {
    DOCKER_HOST = 'tcp://docker-proxy:2375'
  }
  stages {
    stage('Build') {
      steps {
        sh 'docker build -t registry.company.com/app:${BUILD_NUMBER} .'
      }
    }
  }
}
```

### Avec Docker Compose

```yaml
services:
  docker-proxy:
    build: .
    ports:
      - "2375:2375"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      - CONTAINERS=1
      - IMAGES=1
      - NETWORKS=1
      - VOLUMES=1
      - POST=0
      - DELETE=0
      - LOG_LEVEL=info
```

### Avec Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o docker-proxy ./cmd/docker-proxy

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/docker-proxy .

EXPOSE 2375

CMD ["./docker-proxy"]
```

## 🏗️ Architecture

```
.
├── cmd/
│   └── docker-proxy/       # Point d'entrée de l'application
│       └── main.go
├── config/                 # Configuration et chargement des env vars
│   └── config.go
├── internal/
│   ├── middleware/         # Middlewares Gin
│   │   ├── acl.go         # Contrôle d'accès
│   │   └── logging.go     # Logging structuré
│   └── proxy/             # Handler de proxy
│       └── handler.go
└── pkg/
    └── rules/             # Moteur de règles d'accès
        └── matcher.go
```

## 🔐 Système de Filtrage Avancé

**Le point fort de ce proxy : un système de filtrage extrêmement granulaire et puissant.**

### 🎯 Contrôle Précis avec Regex

Contrairement aux proxies basiques, ce proxy permet de **contrôler finement chaque aspect** des opérations Docker via des patterns regex :

#### 📦 Filtres de Volumes
```bash
# Autoriser uniquement volumes nommés spécifiques
export DKRPRX__VOLUMES__ALLOWED_NAMES="^data-.*,^app-.*,^logs-.*"

# Interdire montage de répertoires système sensibles
export DKRPRX__VOLUMES__DENIED_PATHS="^/etc/.*,^/root/.*,^/sys/.*,^/proc/.*,^/var/run/.*"

# Autoriser uniquement bind mounts dans /data
export DKRPRX__VOLUMES__ALLOWED_PATHS="^/data/.*"

# Restreindre aux drivers locaux
export DKRPRX__VOLUMES__ALLOWED_DRIVERS="local"
```

#### 🐳 Filtres de Conteneurs
```bash
# Autoriser uniquement images d'un registry privé avec version sémantique
export DKRPRX__CONTAINERS__ALLOWED_IMAGES="^registry.company.com/.*:v[0-9]+\.[0-9]+\.[0-9]+$"

# Interdire toute image avec tag :latest ou :dev
export DKRPRX__CONTAINERS__DENIED_IMAGES=".*:(latest|dev|test)$"

# Exiger des noms de conteneurs préfixés par l'environnement
export DKRPRX__CONTAINERS__ALLOWED_NAMES="^(prod|staging|dev)-.*"

# Exiger des labels obligatoires
export DKRPRX__CONTAINERS__REQUIRE_LABELS="env=production,team=backend,cost-center=IT-001"

# Interdire conteneurs privilégiés et host network
export DKRPRX__CONTAINERS__DENY_PRIVILEGED="true"
export DKRPRX__CONTAINERS__DENY_HOST_NETWORK="true"
```

#### 🖼️ Filtres d'Images
```bash
# Autoriser uniquement registries approuvés
export DKRPRX__IMAGES__ALLOWED_REPOS="^(docker\.io/library|registry\.company\.com)/.*"

# Interdire registries non sécurisés
export DKRPRX__IMAGES__DENIED_REPOS=".*\.(cn|ru|suspicious)/"

# Autoriser uniquement tags versionnés (semver)
export DKRPRX__IMAGES__ALLOWED_TAGS="^v[0-9]+\.[0-9]+\.[0-9]+$"

# Interdire tags de développement
export DKRPRX__IMAGES__DENIED_TAGS="^(latest|dev|test|alpha|beta|rc).*"
```

#### 🌐 Filtres de Réseaux
```bash
# Autoriser uniquement réseaux applicatifs
export DKRPRX__NETWORKS__ALLOWED_NAMES="^app-.*"

# Interdire réseau host (sécurité)
export DKRPRX__NETWORKS__DENIED_NAMES="^host$"

# Restreindre aux drivers bridge et overlay
export DKRPRX__NETWORKS__ALLOWED_DRIVERS="bridge,overlay"
```

### 📋 Configuration via JSON (Alternative)

Pour des configurations complexes, utilisez JSON :

```json
{
  "volumes": {
    "allowed_names": ["^data-.*", "^app-.*"],
    "denied_paths": ["^/etc/.*", "^/root/.*", "^/sys/.*", "^/proc/.*"],
    "allowed_paths": ["^/data/.*", "^/mnt/volumes/.*"]
  },
  "containers": {
    "allowed_images": ["^registry.company.com/.*:v[0-9]+\\.[0-9]+\\.[0-9]+$"],
    "denied_images": [".*:(latest|dev)$"],
    "require_labels": {
      "env": "production",
      "approved": "true"
    },
    "deny_privileged": true,
    "deny_host_network": true
  },
  "networks": {
    "allowed_names": ["^app-.*"],
    "allowed_drivers": ["bridge", "overlay"]
  },
  "images": {
    "allowed_repos": ["^registry.company.com/.*"],
    "denied_tags": ["^latest$"]
  }
}
```

```bash
export FILTERS_CONFIG=./filters.json
```

**Note :** Les variables d'environnement sont **prioritaires** sur le JSON.

### 💡 Cas d'Usage Avancés

#### Multi-tenant avec Isolation
```bash
TENANT_ID="client-123"
export DKRPRX__VOLUMES__ALLOWED_NAMES="^${TENANT_ID}-.*"
export DKRPRX__CONTAINERS__ALLOWED_NAMES="^${TENANT_ID}-.*"
export DKRPRX__CONTAINERS__REQUIRE_LABELS="tenant=${TENANT_ID}"
export DKRPRX__NETWORKS__ALLOWED_NAMES="^${TENANT_ID}-.*"
```

#### Environnement de Production Strict
```bash
# Uniquement images versionnées d'un registry privé
export DKRPRX__CONTAINERS__ALLOWED_IMAGES="^registry.prod.company.com/.*:v[0-9]+\.[0-9]+\.[0-9]+$"

# Uniquement montages dans /data/prod
export DKRPRX__VOLUMES__ALLOWED_PATHS="^/data/prod/.*"

# Labels obligatoires
export DKRPRX__CONTAINERS__REQUIRE_LABELS="env=production,approved=true,security-scan=passed"

# Sécurité renforcée
export DKRPRX__CONTAINERS__DENY_PRIVILEGED="true"
export DKRPRX__CONTAINERS__DENY_HOST_NETWORK="true"
```

#### CI/CD avec Restrictions
```bash
# Autoriser build mais pas latest
export DKRPRX__IMAGES__DENIED_TAGS="^latest$"

# Autoriser uniquement registry CI
export DKRPRX__IMAGES__ALLOWED_REPOS="^registry.ci.company.com/.*"

# Interdire volumes sensibles
export DKRPRX__VOLUMES__DENIED_PATHS="^/(etc|root|home|sys|proc)/.*"
```

### 📚 Documentation Complète

- **[ADVANCED_FILTERS.md](ADVANCED_FILTERS.md)** - Guide complet des filtres avancés avec exemples détaillés
- **[ENV_FILTERS.md](ENV_FILTERS.md)** - Configuration via variables d'environnement
- **[SECURITY.md](SECURITY.md)** - Guide de sécurité et vecteurs d'attaque bloqués

**Le système de filtrage permet un contrôle aussi précis que nécessaire pour votre environnement !**

## 🔒 Sécurité par Défaut

Le proxy applique **plusieurs protections par défaut** pour éviter l'escalade de privilèges :

### 🛡️ Protection du Socket Docker
Les chemins suivants sont **bloqués par défaut** :
- `/var/run/docker.sock`
- `/run/docker.sock`

### 🛡️ Protection du Conteneur Proxy
Le conteneur `docker-proxy` lui-même est **protégé contre toute manipulation** :
- ❌ Impossible de stopper/redémarrer le conteneur proxy
- ❌ Impossible de modifier le conteneur proxy
- ❌ Impossible de supprimer le conteneur proxy

### 🛡️ Protection du Réseau Proxy
Si le proxy utilise un réseau dédié, celui-ci est également protégé.

### ⚙️ Configuration
```bash
# Nom du conteneur proxy (défaut: docker-proxy)
export PROXY_CONTAINER_NAME="docker-proxy"

# Nom du réseau proxy (optionnel)
export PROXY_NETWORK_NAME="docker-proxy-network"
```

### 🔓 Désactivation (non recommandé)
Pour désactiver toutes les protections :
```bash
export DKRPRX__DISABLE_DEFAULTS="true"
```

Pour autoriser explicitement le socket Docker :
```bash
export DKRPRX__VOLUMES__ALLOWED_PATHS="^/var/run/docker\\.sock$"
```

## ⚠️ Avertissements de sécurité

- **N'exposez JAMAIS ce proxy sur un réseau public**
- **Ne montez JAMAIS le socket Docker dans un conteneur non sécurisé**
- Activez uniquement les endpoints nécessaires
- Utilisez le mode lecture seule quand possible
- Montez le socket Docker en lecture seule quand possible (`:ro`)
- Utilisez les filtres avancés pour un contrôle granulaire

## 🧪 Tests

Pour tester le proxy :

### Mode TCP (défaut)

```bash
# Démarrer le proxy avec containers en lecture seule
export CONTAINERS=1
./docker-proxy

# Dans un autre terminal
curl http://localhost:2375/v1.41/containers/json

# Tester un endpoint refusé
curl http://localhost:2375/v1.41/images/json  # 403 Forbidden
```

### Mode Unix socket

```bash
# Démarrer le proxy sur Unix socket
export LISTEN_SOCKET=/tmp/docker-proxy.sock
export CONTAINERS=1
./docker-proxy

# Dans un autre terminal
curl --unix-socket /tmp/docker-proxy.sock http://localhost/v1.41/containers/json

# Ou avec Docker CLI
export DOCKER_HOST=unix:///tmp/docker-proxy.sock
docker ps
```

## 🔧 Intégration CI/CD

Ce proxy est **spécifiquement conçu pour sécuriser les pipelines CI/CD**. Voir [CICD_EXAMPLES.md](CICD_EXAMPLES.md) pour des exemples détaillés :

- **GitHub Actions** - Sécuriser Docker dans les workflows
- **GitLab CI** - Contrôle d'accès granulaire avec services
- **Jenkins** - Pipeline sécurisé avec Docker proxy
- **Azure DevOps** - Intégration avec ACR
- **CircleCI** - Build sécurisé avec remote Docker

**Cas d'usage typiques :**
- Forcer l'utilisation d'un registry privé uniquement
- Interdire les tags `:latest`, `:dev`, `:test`
- Bloquer les conteneurs privileged en CI
- Audit complet des opérations Docker
- Isolation multi-tenant dans les runners partagés

## 🔐 Sécurité

Voir [SECURITY.md](SECURITY.md) pour :
- Protections par défaut implémentées
- Vecteurs d'attaque bloqués
- Checklist de sécurité pour la production
- Bonnes pratiques de déploiement

## 📝 License

**Dual License: GPL-3.0 (Free) OR Commercial**

- 🆓 **FREE** for open-source projects (GPL-3.0)
- 💼 **Commercial license** for proprietary/closed-source use
- 📋 See [LICENSE](LICENSE) for complete terms
- 💰 Pricing: Startup €500/year, SME €2k/year, Enterprise €10k/year

Contact nicolas.hypolite@gmail.com for commercial licensing.

---

## 🏆 Avantages Techniques

### Architecture & Performance
Développé en **Go (Golang)** haute performance avec **Gin framework** et **Resty HTTP client**, ce proxy offre une **API wrapper** robuste autour du **Docker Engine API** et **Docker SDK**. Support natif des **Unix sockets** et **TCP sockets** pour une intégration flexible.

### Sécurité Avancée
Implémente les principes de **zero-trust**, **least privilege**, et **defense in depth** pour prévenir :
- **Privilege escalation prevention** : Blocage des conteneurs privileged et host network
- **Container escape prevention** : Restrictions strictes sur volumes et bind mounts
- **Socket injection prevention** : Protection automatique du socket Docker

### Filtrage & Contrôle
- **Regex-based filtering** : Patterns avancés pour images, volumes, réseaux
- **Network policy** et **volume restriction** granulaires
- **Image policy** et **tag policy** personnalisables
- **Label enforcement** pour conformité organisationnelle
- **Registry whitelist** : Forcer l'utilisation de registries approuvés

### Intégrations CI/CD
Compatibilité native avec :
- **GitHub Actions**, **GitLab CI/CD**, **Jenkins**, **CircleCI**, **Azure DevOps**
- **Travis CI**, **Drone CI**, **Bamboo**, **TeamCity**
- **Docker-in-Docker (DinD) security** améliorée
- **Kubernetes admission controller** via webhook

### Orchestrateurs & Plateformes
Support pour :
- **Kubernetes** (pods, deployments, namespaces)
- **Docker Swarm** (services, stacks, secrets)
- **HashiCorp Nomad** (jobs, tasks)
- **Rancher**, **Portainer**, **OpenShift**
- Toute plateforme utilisant **Docker API**

---

**Développé pour les équipes DevOps, SRE et Security qui recherchent un contrôle granulaire sur Docker dans des environnements multi-tenants, cloud-native et CI/CD.**
