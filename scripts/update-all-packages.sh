#!/bin/bash

# 全パッケージマネージャーのレシピを更新するスクリプト
set -e

echo "🔄 Updating all package manager recipes..."

# 最新リリース情報を取得
LATEST_RELEASE=$(curl -s https://api.github.com/repos/mimimi105/git-mini-commit/releases/latest)
VERSION=$(echo "$LATEST_RELEASE" | jq -r '.tag_name' | sed 's/v//')

echo "📦 Latest version: $VERSION"

# Homebrew Formula を更新
echo "🍺 Updating Homebrew Formula..."
ruby scripts/update-homebrew.rb

# Scoop レシピを更新
echo "🪣 Updating Scoop recipe..."
ruby scripts/update-scoop.rb

# Debian パッケージをビルド
echo "📦 Building Debian package..."
chmod +x scripts/build-debian.sh
./scripts/build-debian.sh "$VERSION" amd64

echo "✅ All package recipes updated!"
echo ""
echo "Next steps:"
echo "1. Commit and push changes"
echo "2. Create Homebrew tap repository"
echo "3. Submit Scoop recipe to main bucket"
echo "4. Upload Debian packages to GitHub Release"
