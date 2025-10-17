#!/usr/bin/env node
const https = require('https');
const fs = require('fs');
const path = require('path');
const { URL } = require('url');

const version = process.env.npm_package_version || 'latest';
const platform = process.platform; // 'darwin', 'linux', 'win32'
const arch = process.arch;         // 'x64', 'arm64', etc.

function resolveAsset() {
  if (platform === 'darwin') {
    if (arch === 'arm64') return 'i18n-ai-translator-darwin-arm64';
    if (arch === 'x64') return 'i18n-ai-translator-darwin-x64';
  } else if (platform === 'linux') {
    if (arch === 'x64') return 'i18n-ai-translator-linux-x64';
  } else if (platform === 'win32') {
    if (arch === 'x64') return 'i18n-ai-translator-windows-x64.exe';
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
const url = `https://github.com/wowblvck/i18n-ai-translator/releases/download/${tag}/${asset}`;
const binDir = path.join(__dirname, '..', 'bin');
const dest = path.join(binDir, asset.endsWith('.exe') ? 'i18n-ai-translator.exe' : 'i18n-ai-translator');

if (!fs.existsSync(binDir)) fs.mkdirSync(binDir, { recursive: true });

console.log(`Downloading binary: ${url}`);

/**
 * Download a file with automatic redirect handling (302, 301, etc.)
 */
function downloadWithRedirect(url, dest, depth = 0) {
  if (depth > 5) {
    console.error('Too many redirects');
    process.exit(1);
  }

  https.get(url, (res) => {
    if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
      const redirectUrl = new URL(res.headers.location, url).toString();
      console.log(`Redirected to ${redirectUrl}`);
      res.resume();
      return downloadWithRedirect(redirectUrl, dest, depth + 1);
    }

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
}

downloadWithRedirect(url, dest);
