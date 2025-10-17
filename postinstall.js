#!/usr/bin/env node
const https = require('https');
const fs = require('fs');
const path = require('path');

const version = process.env.npm_package_version || 'latest';
const platform = process.platform; // 'darwin', 'linux', 'win32'
const arch = process.arch;         // 'x64', 'arm64', etc.

function resolveAsset() {
  if (platform === 'darwin') {
    if (arch === 'arm64') return 'i18n-translator-darwin-arm64';
    if (arch === 'x64') return 'i18n-translator-darwin-x64';
  } else if (platform === 'linux') {
    if (arch === 'x64') return 'i18n-translator-linux-x64';
  } else if (platform === 'win32') {
    if (arch === 'x64') return 'i18n-translator-windows-x64.exe';
  }
  return null;
}

const asset = resolveAsset();
if (!asset) {
  console.error(`No prebuilt binary for platform=${platform} arch=${arch}. Falling back to source.`);
  console.error('Please install Go to build locally or use a supported platform.');
  process.exit(0); // do not fail install; user can still build manually
}

const tag = `v${version}`;
const url = `https://github.com/wowblvck/i18n-translator/releases/download/${tag}/${asset}`;
const binDir = path.join(__dirname, '..', 'bin');
const dest = path.join(binDir, asset.endsWith('.exe') ? 'i18n-translator.exe' : 'i18n-translator');

if (!fs.existsSync(binDir)) fs.mkdirSync(binDir, { recursive: true });

console.log(`Downloading binary: ${url}`);
https.get(url, (res) => {
  if (res.statusCode !== 200) {
    console.error(`Failed to download ${url}: HTTP ${res.statusCode}`);
    res.resume();
    process.exit(1);
  }
  const file = fs.createWriteStream(dest, { mode: 0o755 });
  res.pipe(file);
  file.on('finish', () => {
    file.close(() => {
      if (process.platform !== 'win32') {
        try { fs.chmodSync(dest, 0o755); } catch (_) { }
      }
      console.log(`Installed binary to ${dest}`);
    });
  });
}).on('error', (err) => {
  console.error(`Error downloading ${url}:`, err);
  process.exit(1);
});