#!/bin/bash
# Script pour formater automatiquement tous les fichiers Go avec gofmt
set -e

# Couleurs pour un meilleur affichage
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 Formatting all Go files with gofmt...${NC}"
echo ""

# Compteurs
total=0
formatted=0
errors=0

# Mode verbeux ou silencieux
VERBOSE=${VERBOSE:-1}

# Fonction pour formatter un fichier
format_file() {
    local file="$1"
    ((total++))

    if [ $VERBOSE -eq 1 ]; then
        echo -e "${BLUE}📄 Processing:${NC} $file"
    fi

    # Vérifier si le fichier a besoin de formatage
    if gofmt -l "$file" | grep -q .; then
        if gofmt -w "$file"; then
            ((formatted++))
            if [ $VERBOSE -eq 1 ]; then
                echo -e "   ${GREEN}✅ Formatted${NC}"
            fi
        else
            ((errors++))
            echo -e "   ${RED}❌ Error formatting${NC}"
        fi
    else
        if [ $VERBOSE -eq 1 ]; then
            echo -e "   ${GREEN}✓ Already formatted${NC}"
        fi
    fi
}

# Trouver et formatter tous les fichiers .go
echo "Searching for Go files..."
file_count=0

while IFS= read -r file; do
    format_file "$file"
    ((file_count++))
done < <(find . -name "*.go" -type f ! -path "./vendor/*" ! -path "./.git/*" ! -path "./bin/*" ! -path "./dist/*")

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}📊 Summary:${NC}"
echo -e "   Total files scanned:  ${BLUE}$total${NC}"
echo -e "   Files formatted:      ${GREEN}$formatted${NC}"
echo -e "   Already formatted:    ${GREEN}$((total - formatted - errors))${NC}"
if [ $errors -gt 0 ]; then
    echo -e "   Errors:               ${RED}$errors${NC}"
fi
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Vérification finale
echo -e "${BLUE}🔍 Final verification...${NC}"
unformatted=$(gofmt -l . 2>/dev/null | grep -v vendor | grep -v ".git" || true)

if [ -z "$unformatted" ]; then
    echo -e "${GREEN}✅ All Go files are properly formatted!${NC}"
    echo ""
    exit 0
else
    echo -e "${RED}⚠️  The following files still have formatting issues:${NC}"
    echo "$unformatted"
    echo ""
    echo -e "${YELLOW}💡 Tip: Run 'gofmt -d <file>' to see the differences${NC}"
    exit 1
fi
