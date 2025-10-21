#!/bin/bash

set -e

# バージョン情報
VERSION=${1:-"0.1.0"}
ARCH=${2:-"amd64"}

echo "Building Debian package for git-mini-commit v${VERSION} (${ARCH})"

# 一時ディレクトリを作成
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# ソースをコピー
cp -r /Users/SchoolAccount/Documents/mini-commit/* .

# バージョンを更新
sed -i "s/0\.1\.0/${VERSION}/g" debian/changelog
sed -i "s/0\.1\.0/${VERSION}/g" debian/control

# バイナリをビルド
GOOS=linux GOARCH=${ARCH} CGO_ENABLED=0 go build -ldflags="-s -w" -o git-mini-commit .

# Debian パッケージをビルド
dpkg-buildpackage -us -uc -b

# 結果をコピー
cp ../git-mini-commit_${VERSION}-1_${ARCH}.deb /Users/SchoolAccount/Documents/mini-commit/

echo "✅ Debian package built: git-mini-commit_${VERSION}-1_${ARCH}.deb"

# クリーンアップ
cd /
rm -rf "$TEMP_DIR"
