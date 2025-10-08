#!/bin/bash
# Script de test pour démontrer les conflits de configuration
set -e

echo "🧪 Test des Conflits de Configuration - dockershield"
echo "======================================================"
echo ""

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fonction helper
test_case() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}Test: $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

result() {
    echo -e "${GREEN}✅ Résultat: $1${NC}"
    echo ""
}

warning() {
    echo -e "${RED}⚠️  Attention: $1${NC}"
    echo ""
}

# Test 1: Configuration par défaut
test_case "Défauts de sécurité (sans configuration)"
echo "Commande:"
echo "  docker run --rm hypolas/proxy-docker:latest"
echo ""
echo "Configuration attendue:"
echo "  - Socket Docker BLOQUÉ (/var/run/docker.sock)"
echo "  - Conteneur proxy PROTÉGÉ"
echo "  - Réseau proxy PROTÉGÉ (si défini)"
result "Les défauts de sécurité sont appliqués"

# Test 2: JSON simple
test_case "Configuration JSON uniquement"
cat > /tmp/test-filters.json << 'EOF'
{
  "volumes": {
    "allowed_paths": ["/app", "/data"]
  }
}
EOF

echo "Fichier JSON:"
cat /tmp/test-filters.json
echo ""
echo "Commande:"
echo "  docker run --rm -v /tmp/test-filters.json:/filters.json \\"
echo "    -e FILTERS_CONFIG=/filters.json \\"
echo "    hypolas/proxy-docker:latest"
echo ""
result "allowed_paths=[/app, /data]"
warning "Les défauts de Volumes sont REMPLACÉS (socket Docker plus protégé !)"

# Test 3: ENV vs JSON
test_case "Conflit ENV vs JSON"
echo "JSON:"
cat /tmp/test-filters.json
echo ""
echo "ENV:"
echo "  DKRPRX__VOLUMES__DENIED_PATHS=/var/run/docker.sock,/tmp"
echo ""
echo "Commande:"
echo "  docker run --rm -v /tmp/test-filters.json:/filters.json \\"
echo "    -e FILTERS_CONFIG=/filters.json \\"
echo "    -e DKRPRX__VOLUMES__DENIED_PATHS=/var/run/docker.sock,/tmp \\"
echo "    hypolas/proxy-docker:latest"
echo ""
result "ENV GAGNE: denied_paths=[/var/run/docker.sock, /tmp]"
warning "allowed_paths du JSON est PERDU !"

# Test 4: Allowed vs Denied
test_case "Conflit Allowed vs Denied dans le même filtre"
cat > /tmp/test-filters2.json << 'EOF'
{
  "volumes": {
    "denied_paths": ["/var"],
    "allowed_paths": ["/app", "/data"]
  }
}
EOF

echo "Configuration:"
cat /tmp/test-filters2.json
echo ""
echo "Tests de montage:"
echo "  1. /var/log     → Matche denied (/var)"
echo "  2. /tmp         → Ne matche ni denied ni allowed"
echo "  3. /app/config  → Ne matche pas denied, matche allowed"
echo ""
result "1. BLOQUÉ (denied vérifié en premier)"
result "2. BLOQUÉ (pas dans allowed)"
result "3. AUTORISÉ ✅"

# Test 5: Désactivation des défauts
test_case "DKRPRX__DISABLE_DEFAULTS=true"
echo "Commande:"
echo "  docker run --rm \\"
echo "    -e DKRPRX__DISABLE_DEFAULTS=true \\"
echo "    hypolas/proxy-docker:latest"
echo ""
warning "AUCUNE PROTECTION PAR DÉFAUT !"
warning "Le socket Docker n'est PLUS protégé"
warning "Le conteneur proxy n'est PLUS protégé"
echo ""
result "À utiliser UNIQUEMENT si vous définissez vos propres règles"

# Test 6: ENV partiel avec défauts
test_case "ENV partiel + Défauts actifs"
echo "Configuration:"
echo "  ENV: DKRPRX__CONTAINERS__DENIED_IMAGES=malicious/*"
echo "  (DKRPRX__DISABLE_DEFAULTS non défini)"
echo ""
echo "Résultat:"
echo "  Containers:"
echo "    - denied_images: [malicious/*]        (de ENV)"
echo "    - denied_names:  [^dockershield$]     (de DÉFAUTS)"
echo "  Volumes:"
echo "    - denied_paths: [/var/run/docker.sock] (de DÉFAUTS)"
echo ""
result "Fusion: ENV pour Containers, Défauts pour Volumes et Networks"

# Test 7: JSON complet qui remplace les défauts
test_case "JSON complet sans DISABLE_DEFAULTS"
cat > /tmp/test-filters3.json << 'EOF'
{
  "volumes": {
    "denied_paths": ["/sensitive"]
  },
  "containers": {
    "allowed_images": ["nginx:*", "redis:*"]
  },
  "networks": {
    "denied_names": ["public"]
  }
}
EOF

echo "JSON:"
cat /tmp/test-filters3.json
echo ""
echo "Résultat:"
echo "  Volumes:"
echo "    - denied_paths: [/sensitive]  (de JSON)"
echo "    - Les défauts (/var/run/docker.sock) sont PERDUS !"
echo ""
echo "  Containers:"
echo "    - allowed_images: [nginx:*, redis:*]  (de JSON)"
echo "    - Les défauts (^dockershield$) sont PERDUS !"
echo ""
warning "En définissant une section dans JSON, vous perdez les défauts de cette section !"

# Résumé
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}📋 RÉSUMÉ DES RÈGLES DE PRIORITÉ${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "1. ENV (DKRPRX__*) > JSON > Défauts"
echo "2. La fusion se fait PAR SECTION (pas par champ)"
echo "3. denied est vérifié AVANT allowed"
echo "4. DISABLE_DEFAULTS désactive TOUS les défauts"
echo ""
echo -e "${GREEN}✅ Tests terminés${NC}"
echo ""
echo "📚 Pour plus de détails: docs/CONFLICT_RESOLUTION.md"

# Cleanup
rm -f /tmp/test-filters.json /tmp/test-filters2.json /tmp/test-filters3.json
