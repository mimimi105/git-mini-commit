#!/bin/bash

# Debian/Ubuntu 用インストールスクリプト
set -e

VERSION=${1:-"latest"}
ARCH=$(dpkg --print-architecture)

echo "Installing git-mini-commit for Debian/Ubuntu (${ARCH})..."

# 最新リリースのURLを取得
if [ "$VERSION" = "latest" ]; then
    DOWNLOAD_URL=$(curl -s https://api.github.com/repos/minoru-kinugasa-105/git-mini-commit/releases/latest | grep "browser_download_url.*linux-${ARCH}" | cut -d '"' -f 4)
else
    DOWNLOAD_URL="https://github.com/minoru-kinugasa-105/git-mini-commit/releases/download/v${VERSION}/git-mini-commit-linux-${ARCH}"
fi

if [ -z "$DOWNLOAD_URL" ]; then
    echo "❌ No binary found for architecture: ${ARCH}"
    echo "Available architectures: amd64, arm64"
    exit 1
fi

# 一時ディレクトリにダウンロード
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo "Downloading from: $DOWNLOAD_URL"
wget -O git-mini-commit "$DOWNLOAD_URL"
chmod +x git-mini-commit

# システムにインストール
sudo mv git-mini-commit /usr/local/bin/
sudo chmod +x /usr/local/bin/git-mini-commit

# バージョン確認
echo "✅ git-mini-commit installed successfully!"
/usr/local/bin/git-mini-commit --version

# クリーンアップ
cd /
rm -rf "$TEMP_DIR"

echo ""
echo "Usage:"
echo "  git mini-commit -m \"Your mini-commit message\""
echo "  git mini-commit list"
echo "  git mini-commit show <hash>"
