#!/bin/bash
set -e

echo "🔧 Formatting all Go files with gofmt..."
echo ""

# Find all .go files and format them
find . -name "*.go" -type f ! -path "./vendor/*" ! -path "./.git/*" | while read -r file; do
    echo "Formatting: $file"
    gofmt -w "$file"
done

echo ""
echo "✅ All Go files have been formatted!"
echo ""
echo "📊 Checking for any remaining formatting issues..."
unformatted=$(gofmt -l . 2>/dev/null | grep -v vendor || true)

if [ -z "$unformatted" ]; then
    echo "✅ All files are properly formatted!"
    exit 0
else
    echo "⚠️  The following files still have formatting issues:"
    echo "$unformatted"
    exit 1
fi
