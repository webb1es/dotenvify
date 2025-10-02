#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');

function getBinaryName() {
  const platform = os.platform();
  const arch = os.arch();
  
  let osName;
  switch (platform) {
    case 'win32':
      osName = 'windows';
      break;
    case 'darwin':
      osName = 'darwin';
      break;
    case 'linux':
      osName = 'linux';
      break;
    default:
      throw new Error(`Unsupported platform: ${platform}`);
  }
  
  let archName;
  switch (arch) {
    case 'x64':
      archName = 'amd64';
      break;
    case 'arm64':
      archName = 'arm64';
      break;
    default:
      throw new Error(`Unsupported architecture: ${arch}`);
  }
  
  const extension = platform === 'win32' ? '.exe' : '';
  return `dotenvify_${osName}_${archName}${extension}`;
}

function main() {
  try {
    const binaryName = getBinaryName();
    const binaryPath = path.join(__dirname, binaryName);
    
    const child = spawn(binaryPath, process.argv.slice(2), {
      stdio: 'inherit',
      windowsHide: false,
    });
    
    child.on('error', (error) => {
      if (error.code === 'ENOENT') {
        console.error(`Binary not found: ${binaryPath}`);
        console.error(`Platform: ${os.platform()}, Architecture: ${os.arch()}`);
        process.exit(1);
      } else {
        console.error('Failed to start dotenvify:', error.message);
        process.exit(1);
      }
    });
    
    child.on('exit', (code, signal) => {
      if (signal) {
        process.kill(process.pid, signal);
      } else {
        process.exit(code);
      }
    });
  } catch (error) {
    console.error('Error:', error.message);
    process.exit(1);
  }
}

if (require.main === module) {
  main();
}