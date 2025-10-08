# 🛠️ Scripts Utilitaires - Docker Proxy

Ce document liste tous les scripts disponibles pour faciliter le développement et la maintenance du projet.

## 📋 Liste des Scripts

### 1. `gofmt-all.sh` - Formatage du code Go

Formate automatiquement tous les fichiers `.go` du projet avec `gofmt`.

**Usage:**
```bash
chmod +x gofmt-all.sh
./gofmt-all.sh
```

**Options:**
```bash
# Mode silencieux (affiche seulement le résumé)
VERBOSE=0 ./gofmt-all.sh

# Mode verbeux (par défaut)
VERBOSE=1 ./gofmt-all.sh
```

**Ce qu'il fait:**
- ✅ Trouve tous les fichiers `.go` (exclut `vendor/`, `.git/`, `bin/`, `dist/`)
- ✅ Applique `gofmt -w` sur chaque fichier
- ✅ Affiche un rapport détaillé
- ✅ Vérifie qu'il ne reste aucun problème de formatage
- ✅ Code de sortie: 0 si OK, 1 si erreur

**Exemple de sortie:**
```
🔧 Formatting all Go files with gofmt...

Searching for Go files...
📄 Processing: ./cmd/dockershield/main.go
   ✅ Formatted
📄 Processing: ./config/config.go
   ✓ Already formatted

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Summary:
   Total files scanned:  17
   Files formatted:      3
   Already formatted:    14
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔍 Final verification...
✅ All Go files are properly formatted!
```

---

### 2. `check-lint-config.sh` - Vérification config golangci-lint

Vérifie la configuration `.golangci.yml` pour détecter les conflits et erreurs.

**Usage:**
```bash
chmod +x check-lint-config.sh
./check-lint-config.sh
```

**Ce qu'il fait:**
- ✅ Détecte les conflits entre `enable` et `disable`
- ✅ Valide la syntaxe YAML (si `yamllint` installé)
- ✅ Teste la configuration avec `golangci-lint`
- ✅ Affiche les linters actifs

**Exemple de sortie:**
```
🔍 Vérification de la configuration golangci-lint...

1️⃣ Vérification des conflits enable/disable...
✅ Pas de conflits entre enable et disable

2️⃣ Validation de la syntaxe YAML...
✅ Syntaxe YAML valide

3️⃣ Test de la configuration avec golangci-lint...
✅ golangci-lint est configuré

✅ Vérification terminée!
```

---

### 3. `test-config-conflicts.sh` - Tests de conflits de configuration

Démontre et teste les différents scénarios de conflits de configuration dans dockershield.

**Usage:**
```bash
chmod +x test-config-conflicts.sh
./test-config-conflicts.sh
```

**Ce qu'il fait:**
- 🧪 Teste 7 scénarios de conflits différents
- 📝 Crée des fichiers JSON temporaires pour les tests
- 🎨 Affiche les résultats avec couleurs
- 🧹 Nettoie automatiquement les fichiers temporaires

**Scénarios testés:**
1. Configuration par défaut (sans configuration)
2. Configuration JSON uniquement
3. Conflit ENV vs JSON
4. Conflit Allowed vs Denied
5. Désactivation des défauts (DISABLE_DEFAULTS)
6. ENV partiel avec défauts actifs
7. JSON complet qui remplace les défauts

**Exemple de sortie:**
```
🧪 Test des Conflits de Configuration - dockershield
======================================================

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Test: Conflit ENV vs JSON
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
JSON:
{
  "volumes": {
    "allowed_paths": ["/app", "/data"]
  }
}

ENV:
  DKRPRX__VOLUMES__DENIED_PATHS=/var/run/docker.sock,/tmp

✅ Résultat: ENV GAGNE: denied_paths=[/var/run/docker.sock, /tmp]
⚠️  Attention: allowed_paths du JSON est PERDU !
```

---

### 4. `format-all.sh` - Formatage simple (deprecated)

**Note:** Ce script est déprécié. Utilisez `gofmt-all.sh` à la place.

Ancien script simple qui liste les commandes `gofmt -w` pour chaque fichier.

---

## 🚀 Workflows de Développement

### Avant de Commit

```bash
# 1. Formatter le code
./gofmt-all.sh

# 2. Vérifier les tests
go test ./...

# 3. Vérifier le linter
golangci-lint run

# 4. Build de vérification
go build ./cmd/dockershield
```

### Avant de Push un Tag

```bash
# 1. Formatter le code
./gofmt-all.sh

# 2. Tests complets
go test -v -race -cover ./...

# 3. Lint complet
golangci-lint run --timeout=5m

# 4. Vérifier la config
./check-lint-config.sh

# 5. Build multi-plateforme
GOOS=linux GOARCH=amd64 go build ./cmd/dockershield
GOOS=linux GOARCH=arm64 go build ./cmd/dockershield

# 6. Créer le tag
git tag v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

### Debug des Conflits de Configuration

```bash
# 1. Comprendre les conflits
./test-config-conflicts.sh

# 2. Lire la documentation
cat docs/CONFLICT_RESOLUTION.md

# 3. Tester en local avec docker
docker run --rm -e LOG_LEVEL=debug \
  -e DKRPRX__VOLUMES__DENIED_PATHS=/var/log \
  hypolas/dockershield:latest
```

---

## 📝 Ajout d'un Nouveau Script

Pour ajouter un nouveau script utilitaire :

1. **Créer le script** avec un nom descriptif (ex: `check-dependencies.sh`)

2. **Ajouter le shebang et les options de sécurité:**
   ```bash
   #!/bin/bash
   set -e  # Exit on error
   ```

3. **Documenter le script:**
   ```bash
   # Description: Vérifie les dépendances du projet
   # Usage: ./check-dependencies.sh
   ```

4. **Rendre exécutable:**
   ```bash
   chmod +x check-dependencies.sh
   ```

5. **Mettre à jour ce document** (SCRIPTS.md)

6. **Ajouter au .gitignore si nécessaire:**
   ```gitignore
   # Exclure les artefacts du script
   /tmp-script-output/
   ```

---

## 🔧 Dépendances Requises

### Requis
- `bash` (≥ 4.0)
- `go` (≥ 1.21)
- `gofmt` (inclus avec Go)
- `find`, `grep`, `awk` (outils Unix standard)

### Optionnel
- `golangci-lint` (pour check-lint-config.sh)
- `yamllint` (pour vérification YAML)
- `docker` (pour test-config-conflicts.sh)

### Installation des outils optionnels

**golangci-lint:**
```bash
# Linux/macOS
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Ou avec Homebrew (macOS)
brew install golangci-lint
```

**yamllint:**
```bash
# Ubuntu/Debian
apt-get install yamllint

# macOS
brew install yamllint

# Python pip
pip install yamllint
```

---

## 🐛 Troubleshooting

### "Permission denied" lors de l'exécution

**Solution:**
```bash
chmod +x *.sh
```

### "command not found: gofmt"

**Solution:** Vérifiez que Go est installé et dans le PATH
```bash
go version
which gofmt
```

### Les couleurs ne s'affichent pas correctement

**Solution:** Votre terminal ne supporte peut-être pas les couleurs ANSI
```bash
# Désactiver les couleurs
TERM=dumb ./gofmt-all.sh
```

### "No such file or directory" dans find

**Solution:** Assurez-vous d'exécuter le script depuis la racine du projet
```bash
cd /path/to/dockershield
./gofmt-all.sh
```

---

## 📚 Ressources

- [Documentation Go - gofmt](https://pkg.go.dev/cmd/gofmt)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Bash Scripting Guide](https://www.gnu.org/software/bash/manual/)
- [Configuration Conflicts - docs/CONFLICT_RESOLUTION.md](docs/CONFLICT_RESOLUTION.md)

---

**Version:** 1.0
**Dernière mise à jour:** 2025-10-06
**Mainteneur:** Nicolas HYPOLITE (nicolas.hypolite@gmail.com)
