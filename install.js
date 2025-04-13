#!/usr/bin/env node

const os = require('os');
const fs = require('fs');
const path = require('path');
const { promisify } = require('util');
const got = require('got');
const extract = require('extract-zip');
const { pipeline } = require('stream');
const { execSync } = require('child_process');

const streamPipeline = promisify(pipeline);

// Package version - should match package.json
const version = '0.1.0';
// GitHub repository information
const repo = {
  owner: 'VesperAkshay',
  name: 'lazynode'
};

// Determine platform and architecture
const platform = os.platform();
const arch = os.arch();

// Map Node.js arch to Go arch
const archMap = {
  'x64': 'amd64',
  'arm64': 'arm64',
  'ia32': '386'
};

// Map Node.js platform to Go platform
const platformMap = {
  'win32': 'windows',
  'darwin': 'darwin',
  'linux': 'linux'
};

// Get appropriate binary name
const goPlatform = platformMap[platform] || platform;
const goArch = archMap[arch] || arch;
const extension = platform === 'win32' ? '.zip' : '.tar.gz';

// Binary name in the archive
const binaryName = `lazynode_${version}_${goPlatform}_${goArch}${extension}`;

// Download URL
const downloadUrl = `https://github.com/${repo.owner}/${repo.name}/releases/download/v${version}/${binaryName}`;

// Binary path
const binPath = path.join(__dirname, 'bin');
const executablePath = path.join(binPath, platform === 'win32' ? 'lazynode.exe' : 'lazynode');

async function main() {
  try {
    // Create bin directory if it doesn't exist
    if (!fs.existsSync(binPath)) {
      fs.mkdirSync(binPath, { recursive: true });
    }

    console.log(`Downloading LazyNode v${version} for ${goPlatform} ${goArch}...`);
    console.log(`From: ${downloadUrl}`);

    // Download the binary
    const tempFile = path.join(os.tmpdir(), binaryName);
    await streamPipeline(
      got.stream(downloadUrl),
      fs.createWriteStream(tempFile)
    );

    console.log('Download complete. Extracting...');

    // Extract the archive
    if (extension === '.zip') {
      await extract(tempFile, { dir: binPath });
    } else {
      // For tar.gz files
      execSync(`tar -xzf "${tempFile}" -C "${binPath}"`);
    }

    // Make the binary executable on Unix platforms
    if (platform !== 'win32') {
      fs.chmodSync(executablePath, 0o755); // rwxr-xr-x
    }

    // Clean up temp file
    fs.unlinkSync(tempFile);

    console.log('LazyNode has been installed successfully!');
    console.log(`Binary location: ${executablePath}`);
    console.log('You can now run "lazynode" from your terminal.');

  } catch (error) {
    console.error('Error installing LazyNode:', error.message);
    console.error('Please report this issue at:', `https://github.com/${repo.owner}/${repo.name}/issues`);
    process.exit(1);
  }
}

main().catch(console.error); 