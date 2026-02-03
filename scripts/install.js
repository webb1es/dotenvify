#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const https = require('https');
const { createWriteStream, mkdirSync, chmodSync, existsSync } = fs;

const REPO = 'webb1es/dotenvify';
const BINARY_NAME = 'dotenvify';

function getPlatformInfo() {
  const platform = process.platform;
  const arch = process.arch;

  const platformMap = {
    darwin: 'darwin',
    linux: 'linux',
    win32: 'windows',
  };

  const archMap = {
    x64: 'amd64',
    arm64: 'arm64',
  };

  const os = platformMap[platform];
  const goArch = archMap[arch];

  if (!os || !goArch) {
    console.error(`Unsupported platform: ${platform} ${arch}`);
    process.exit(1);
  }

  // Windows arm64 not supported
  if (os === 'windows' && goArch === 'arm64') {
    console.error('Windows arm64 is not supported');
    process.exit(1);
  }

  return { os, arch: goArch, isWindows: platform === 'win32' };
}

function getPackageVersion() {
  const packageJson = require('../package.json');
  return packageJson.version;
}

function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    const follow = (url) => {
      https.get(url, (response) => {
        if (response.statusCode === 302 || response.statusCode === 301) {
          follow(response.headers.location);
          return;
        }

        if (response.statusCode !== 200) {
          reject(new Error(`Download failed: ${response.statusCode}`));
          return;
        }

        const file = createWriteStream(dest);
        response.pipe(file);
        file.on('finish', () => {
          file.close();
          resolve();
        });
        file.on('error', (err) => {
          fs.unlinkSync(dest);
          reject(err);
        });
      }).on('error', reject);
    };
    follow(url);
  });
}

async function extractTarGz(archive, dest) {
  execSync(`tar -xzf "${archive}" -C "${dest}"`, { stdio: 'inherit' });
}

async function extractZip(archive, dest) {
  if (process.platform === 'win32') {
    execSync(`powershell -command "Expand-Archive -Path '${archive}' -DestinationPath '${dest}'"`, { stdio: 'inherit' });
  } else {
    execSync(`unzip -o "${archive}" -d "${dest}"`, { stdio: 'inherit' });
  }
}

async function install() {
  const { os, arch, isWindows } = getPlatformInfo();
  const version = getPackageVersion();
  const ext = isWindows ? 'zip' : 'tar.gz';
  const binaryExt = isWindows ? '.exe' : '';

  const archiveName = `${BINARY_NAME}_${version}_${os}_${arch}.${ext}`;
  const downloadUrl = `https://github.com/${REPO}/releases/download/v${version}/${archiveName}`;

  const binDir = path.join(__dirname, '..', 'bin');
  const tmpDir = path.join(__dirname, '..', '.tmp');
  const archivePath = path.join(tmpDir, archiveName);
  const binaryDest = path.join(binDir, `${BINARY_NAME}${binaryExt}`);

  // Create directories
  if (!existsSync(binDir)) mkdirSync(binDir, { recursive: true });
  if (!existsSync(tmpDir)) mkdirSync(tmpDir, { recursive: true });

  console.log(`Downloading ${BINARY_NAME} v${version} for ${os}/${arch}...`);
  console.log(`URL: ${downloadUrl}`);

  try {
    await downloadFile(downloadUrl, archivePath);
  } catch (err) {
    console.error(`Failed to download binary: ${err.message}`);
    console.error('');
    console.error('This may happen if:');
    console.error(`  - Version v${version} has not been released on GitHub yet`);
    console.error('  - Your platform is not supported');
    console.error('');
    console.error('You can manually download from:');
    console.error(`  https://github.com/${REPO}/releases`);
    process.exit(1);
  }

  console.log('Extracting...');
  if (isWindows) {
    await extractZip(archivePath, tmpDir);
  } else {
    await extractTarGz(archivePath, tmpDir);
  }

  // Find and move the binary
  const extractedBinary = path.join(tmpDir, `${BINARY_NAME}${binaryExt}`);
  if (existsSync(extractedBinary)) {
    fs.renameSync(extractedBinary, binaryDest);
    if (!isWindows) {
      chmodSync(binaryDest, 0o755);
    }
  } else {
    console.error('Binary not found in archive');
    process.exit(1);
  }

  // Cleanup
  fs.rmSync(tmpDir, { recursive: true, force: true });

  console.log(`Successfully installed ${BINARY_NAME} to ${binaryDest}`);
}

install().catch((err) => {
  console.error('Installation failed:', err);
  process.exit(1);
});
