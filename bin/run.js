#!/usr/bin/env node
const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

const binName = process.platform === 'win32' ? 'i18n-translator.exe' : 'i18n-translator';
const binPath = path.join(__dirname, binName);

if (!fs.existsSync(binPath)) {
  console.error(`Binary not found at ${binPath}. Try reinstalling the package or check postinstall logs.`);
  process.exit(1);
}

const args = process.argv.slice(2);
const child = spawn(binPath, args, { stdio: 'inherit' });

child.on('exit', (code) => process.exit(code));
child.on('error', (err) => {
  console.error('Failed to start binary:', err);
  process.exit(1);
});