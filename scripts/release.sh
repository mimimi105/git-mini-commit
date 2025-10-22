#!/bin/bash

# Manual release script for git-mini-commit
set -e

VERSION=${1:-"latest"}
echo "Creating release for version: $VERSION"

# Build all platforms
echo "Building binaries for all platforms..."

# Linux AMD64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o builds/git-mini-commit-linux-amd64 .

# Linux ARM64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o builds/git-mini-commit-linux-arm64 .

# macOS Intel
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o builds/git-mini-commit-darwin-amd64 .

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o builds/git-mini-commit-darwin-arm64 .

# Windows AMD64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o builds/git-mini-commit-windows-amd64.exe .

# Windows ARM64
GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o builds/git-mini-commit-windows-arm64.exe .

echo "✅ All binaries built successfully!"

# Create release
echo "Creating GitHub release..."
gh release create "v$VERSION" \
  --title "v$VERSION" \
  --notes "Release v$VERSION

- Linux AMD64/ARM64
- macOS Intel/Apple Silicon  
- Windows AMD64/ARM64
- All binaries tested and working

## Installation

### Homebrew (macOS/Linux)
\`\`\`bash
brew install mimimi105/git-mini-commit/git-mini-commit
\`\`\`

### Direct Download
\`\`\`bash
# macOS (Apple Silicon)
curl -L https://github.com/mimimi105/mini-commit/releases/latest/download/git-mini-commit-darwin-arm64 -o git-mini-commit
chmod +x git-mini-commit
sudo mv git-mini-commit /usr/local/bin/
\`\`\`

### Scoop (Windows)
\`\`\`bash
scoop bucket add git-mini-commit https://github.com/mimimi105/scoop-git-mini-commit
scoop install git-mini-commit
\`\`\`" \
  builds/*

echo "✅ Release created successfully!"
echo "Release URL: https://github.com/mimimi105/mini-commit/releases/tag/v$VERSION"
