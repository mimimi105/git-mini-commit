#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const platforms = [
    { os: 'linux', arch: 'amd64' },
    { os: 'linux', arch: 'arm64' },
    { os: 'darwin', arch: 'amd64' },
    { os: 'darwin', arch: 'arm64' },
    { os: 'windows', arch: 'amd64' },
    { os: 'windows', arch: 'arm64' }
];

const binDir = path.join(__dirname, '..', 'bin');
if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
}

console.log('Building binaries for npm package...');

platforms.forEach(({ os, arch }) => {
    const ext = os === 'windows' ? '.exe' : '';
    const outputName = `git-mini-commit-${os}-${arch}${ext}`;
    const outputPath = path.join(binDir, outputName);

    console.log(`Building for ${os}-${arch}...`);

    try {
        execSync(`GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build -ldflags="-s -w" -o ${outputPath} .`, {
            cwd: path.join(__dirname, '..'),
            stdio: 'inherit'
        });
        console.log(`✅ Built ${outputName}`);
    } catch (error) {
        console.error(`❌ Failed to build ${outputName}:`, error.message);
    }
});

console.log('Binary build completed!');
