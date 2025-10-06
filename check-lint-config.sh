#!/bin/bash
set -e

echo "🔍 Vérification de la configuration golangci-lint..."
echo ""

# Vérifier les conflits enable/disable
echo "1️⃣ Vérification des conflits enable/disable..."
if [ -f .golangci.yml ]; then
    enabled=$(grep -A 20 "enable:" .golangci.yml | grep "^    -" | awk '{print $2}' | sort)
    disabled=$(grep -A 10 "disable:" .golangci.yml | grep "^    -" | awk '{print $2}' | sort)

    conflicts=$(comm -12 <(echo "$enabled") <(echo "$disabled"))

    if [ -n "$conflicts" ]; then
        echo "❌ Conflits détectés entre enable et disable:"
        echo "$conflicts"
    else
        echo "✅ Pas de conflits entre enable et disable"
    fi
else
    echo "⚠️  Fichier .golangci.yml non trouvé"
fi

echo ""
echo "2️⃣ Validation de la syntaxe YAML..."
if command -v yamllint &> /dev/null; then
    yamllint .golangci.yml && echo "✅ Syntaxe YAML valide"
else
    echo "ℹ️  yamllint n'est pas installé (optionnel)"
fi

echo ""
echo "3️⃣ Test de la configuration avec golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    golangci-lint linters | head -20
    echo ""
    echo "✅ golangci-lint est configuré"
else
    echo "⚠️  golangci-lint n'est pas installé"
fi

echo ""
echo "✅ Vérification terminée!"
