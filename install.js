#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const os = require('os');
const https = require('https');
const { execSync } = require('child_process');

const platform = os.platform();
const arch = os.arch();

// バイナリ名の決定
let binaryName = 'git-mini-commit';
if (platform === 'win32') {
    binaryName = arch === 'arm64' ? 'git-mini-commit-arm64.exe' : 'git-mini-commit.exe';
} else if (platform === 'darwin') {
    binaryName = arch === 'arm64' ? 'git-mini-commit-macos-arm64' : 'git-mini-commit-macos';
} else if (platform === 'linux') {
    binaryName = arch === 'arm64' ? 'git-mini-commit-linux-arm64' : 'git-mini-commit';
}

const sourcePath = path.join(__dirname, 'bin', binaryName);
const targetPath = path.join(__dirname, 'bin', 'git-mini-commit' + (platform === 'win32' ? '.exe' : ''));

console.log(`Installing git-mini-commit for ${platform}-${arch}...`);

// バイナリが存在する場合はコピー
if (fs.existsSync(sourcePath)) {
    try {
        fs.copyFileSync(sourcePath, targetPath);
        fs.chmodSync(targetPath, '755');
        console.log('✅ git-mini-commit installed successfully!');
    } catch (error) {
        console.error('❌ Failed to install git-mini-commit:', error.message);
        process.exit(1);
    }
} else {
    console.error(`❌ Binary not found: ${sourcePath}`);
    console.log('Available binaries:');
    try {
        const binDir = path.join(__dirname, 'bin');
        const files = fs.readdirSync(binDir);
        files.forEach(file => console.log(`  - ${file}`));
    } catch (error) {
        console.log('  No binaries found in bin directory');
    }
    process.exit(1);
}

// バージョン確認
try {
    const version = execSync(`${targetPath} --version`, { encoding: 'utf8' }).trim();
    console.log(`Version: ${version}`);
} catch (error) {
    console.log('Note: Version check failed, but binary is installed');
}
