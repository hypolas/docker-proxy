# Configurer Docker pour écouter sur TCP 2375

## ⚠️ Avertissement de Sécurité

**ATTENTION**: Exposer Docker sur TCP sans TLS est dangereux !
- N'utilisez ceci que sur `localhost` ou réseau privé
- Pour production, utilisez TLS (port 2376)
- Ou mieux : utilisez docker-proxy comme protection

## Méthode 1 : systemd (Ubuntu/Debian/CentOS)

### Étape 1 : Créer le fichier de configuration

```bash
sudo mkdir -p /etc/systemd/system/docker.service.d
sudo nano /etc/systemd/system/docker.service.d/override.conf
```

### Étape 2 : Ajouter la configuration

**Pour écouter sur localhost uniquement (SÉCURISÉ):**
```ini
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd -H fd:// -H tcp://127.0.0.1:2375
```

**Pour écouter sur toutes les interfaces (DANGEREUX):**
```ini
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd -H fd:// -H tcp://0.0.0.0:2375
```

### Étape 3 : Recharger et redémarrer

```bash
sudo systemctl daemon-reload
sudo systemctl restart docker
```

### Étape 4 : Vérifier

```bash
# Vérifier que Docker écoute
sudsudo netstat -tlnp | grep 2375o netstat -tlnp | grep 2375
# ou
sudo ss -tlnp | grep 2375

# Tester l'API
curl http://localhost:2375/version
```

## Méthode 2 : Fichier daemon.json

### Étape 1 : Éditer daemon.json

```bash
sudo nano /etc/docker/daemon.json
```

### Étape 2 : Ajouter la configuration

**Pour localhost uniquement (SÉCURISÉ):**
```json
{
  "hosts": ["unix:///var/run/docker.sock", "tcp://127.0.0.1:2375"]
}
```

**Pour toutes les interfaces (DANGEREUX):**
```json
{
  "hosts": ["unix:///var/run/docker.sock", "tcp://0.0.0.0:2375"]
}
```

### Étape 3 : Redémarrer Docker

```bash
sudo systemctl restart docker
```

### ⚠️ Conflit possible

Si vous avez une erreur du type "unable to configure the Docker daemon with file", c'est que systemd définit déjà `-H`. Utilisez plutôt la Méthode 1.

## Méthode 3 : Docker Desktop (Windows/Mac)

### Windows

1. Ouvrir Docker Desktop
2. Settings → General
3. Cocher "Expose daemon on tcp://localhost:2375 without TLS"
4. Apply & Restart

### Mac

1. Ouvrir Docker Desktop
2. Preferences → General
3. Cocher "Expose daemon on tcp://localhost:2375 without TLS"
4. Apply & Restart

## Méthode 4 : Temporaire (pour tests)

```bash
# Arrêter Docker
sudo systemctl stop docker

# Lancer manuellement
sudo dockerd -H unix:///var/run/docker.sock -H tcp://127.0.0.1:2375

# Dans un autre terminal, tester
export DOCKER_HOST=tcp://localhost:2375
docker ps
```

## Tester la configuration

### Avec curl

```bash
# Version
curl http://localhost:2375/version

# Lister les conteneurs
curl http://localhost:2375/v1.41/containers/json

# Info système
curl http://localhost:2375/info
```

### Avec Docker CLI

```bash
export DOCKER_HOST=tcp://localhost:2375
docker ps
docker version
```

### Avec docker-proxy (RECOMMANDÉ)

Au lieu d'exposer Docker directement, utilisez docker-proxy comme intermédiaire sécurisé :

```bash
# Docker écoute sur localhost:2375 (ou socket)
# docker-proxy écoute sur votre port désiré avec sécurité

export DOCKER_SOCKET=tcp://localhost:2375  # Votre Docker
export LISTEN_ADDR=:2376                   # Port sécurisé
export CONTAINERS=1
export IMAGES=1
./dockershield
```

## Sécuriser avec TLS (Production)

Pour production, utilisez TLS (port 2376) :

### Générer les certificats

```bash
# Créer CA
openssl genrsa -out ca-key.pem 4096
openssl req -new -x509 -days 365 -key ca-key.pem -sha256 -out ca.pem

# Créer certificat serveur
openssl genrsa -out server-key.pem 4096
openssl req -new -key server-key.pem -out server.csr
openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem

# Créer certificat client
openssl genrsa -out key.pem 4096
openssl req -new -key key.pem -out client.csr
openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out cert.pem
```

### Configurer Docker avec TLS

```bash
sudo nano /etc/systemd/system/docker.service.d/override.conf
```

```ini
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd \
  --tlsverify \
  --tlscacert=/etc/docker/certs/ca.pem \
  --tlscert=/etc/docker/certs/server-cert.pem \
  --tlskey=/etc/docker/certs/server-key.pem \
  -H fd:// \
  -H tcp://0.0.0.0:2376
```

```bash
sudo systemctl daemon-reload
sudo systemctl restart docker
```

### Utiliser avec TLS

```bash
docker --tlsverify \
  --tlscacert=ca.pem \
  --tlscert=cert.pem \
  --tlskey=key.pem \
  -H tcp://localhost:2376 \
  ps
```

## Pare-feu

### Autoriser le port localement

```bash
# UFW (Ubuntu)
sudo ufw allow from 127.0.0.1 to any port 2375

# Firewalld (CentOS/Fedora)
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="127.0.0.1" port protocol="tcp" port="2375" accept'
sudo firewall-cmd --reload

# iptables
sudo iptables -A INPUT -s 127.0.0.1 -p tcp --dport 2375 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 2375 -j DROP
```

## Problèmes courants

### Erreur : "unable to configure the Docker daemon"

**Cause** : Conflit entre daemon.json et systemd

**Solution** : Utilisez uniquement la Méthode 1 (systemd override)

### Erreur : "address already in use"

**Cause** : Le port 2375 est déjà utilisé

**Solution** :
```bash
# Trouver le processus
sudo lsof -i :2375

# Ou changer de port
ExecStart=/usr/bin/dockerd -H fd:// -H tcp://127.0.0.1:2380
```

### Docker ne démarre plus

**Solution** : Revenir à la config par défaut
```bash
sudo rm /etc/systemd/system/docker.service.d/override.conf
sudo systemctl daemon-reload
sudo systemctl restart docker
```

## Architecture recommandée

```
┌─────────────────────────────────────────────┐
│            Votre Application                │
│         (docker CLI, SDK, etc.)             │
└──────────────────┬──────────────────────────┘
                   │ tcp://localhost:2376
                   ↓
┌─────────────────────────────────────────────┐
│           docker-proxy (sécurisé)           │
│  - ACL rules                                │
│  - Advanced filters                         │
│  - Audit logs                               │
└──────────────────┬──────────────────────────┘
                   │ tcp://localhost:2375
                   ↓
┌─────────────────────────────────────────────┐
│         Docker Daemon (localhost)           │
│    Écoute uniquement sur 127.0.0.1:2375     │
└─────────────────────────────────────────────┘
```

## Résumé des commandes

### Configuration rapide (localhost seulement)

```bash
# 1. Créer la config
sudo mkdir -p /etc/systemd/system/docker.service.d
echo '[Service]
ExecStart=
ExecStart=/usr/bin/dockerd -H fd:// -H tcp://127.0.0.1:2375' | sudo tee /etc/systemd/system/docker.service.d/override.conf

# 2. Recharger
sudo systemctl daemon-reload
sudo systemctl restart docker

# 3. Tester
curl http://localhost:2375/version
```

### Vérification

```bash
# Vérifier le port
sudo ss -tlnp | grep 2375

# Tester l'API
curl http://localhost:2375/containers/json

# Avec Docker CLI
export DOCKER_HOST=tcp://localhost:2375
docker ps
```

## Références

- [Docker daemon socket options](https://docs.docker.com/engine/reference/commandline/dockerd/#daemon-socket-option)
- [Protect the Docker daemon socket](https://docs.docker.com/engine/security/protect-access/)
- [Use TLS](https://docs.docker.com/engine/security/https/)
