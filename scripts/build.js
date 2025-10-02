#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// This script copies the appropriate binary to the bin directory during npm pack
// The actual binaries are built by GoReleaser during the release process

const packageJson = require('../package.json');
const version = packageJson.version;

console.log(`Building npm package for dotenvify v${version}`);

// Create bin directory if it doesn't exist
const binDir = path.join(__dirname, '..', 'bin');
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

// Check if we're in a CI environment (binaries should be provided by GoReleaser)
if (process.env.CI || process.env.GITHUB_ACTIONS) {
  console.log('CI environment detected - binaries should be provided by GoReleaser');
  
  // Look for binaries in multiple possible locations
  const possibleDirs = [
    path.join(__dirname, '..', 'dist'),
    path.join(process.cwd(), 'dist'),
    process.cwd()
  ];
  
  let distDir = null;
  for (const dir of possibleDirs) {
    if (fs.existsSync(dir)) {
      const files = fs.readdirSync(dir);
      if (files.some(f => f.includes('dotenvify_'))) {
        distDir = dir;
        break;
      }
    }
  }
  
  if (distDir) {
    const files = fs.readdirSync(distDir);
    console.log('Available files in dist:', files);
    
    // Copy individual binaries from subdirectories
    files.forEach(file => {
      const filePath = path.join(distDir, file);
      if (fs.statSync(filePath).isDirectory() && file.startsWith('dotenvify_')) {
        // This is a binary directory, look for the actual binary inside
        const binaryFiles = fs.readdirSync(filePath);
        binaryFiles.forEach(binaryFile => {
          if (binaryFile === 'dotenvify' || binaryFile === 'dotenvify.exe') {
            const srcPath = path.join(filePath, binaryFile);
            const destPath = path.join(binDir, file + (binaryFile.endsWith('.exe') ? '.exe' : ''));
            
            fs.copyFileSync(srcPath, destPath);
            fs.chmodSync(destPath, '755');
            console.log(`Copied binary: ${file} -> ${destPath}`);
          }
        });
      }
    });
  } else {
    console.error('No dist directory with binaries found');
    process.exit(1);
  }
} else {
  console.log('Local environment - skipping binary copy (use GoReleaser for local builds)');
}

console.log('Build complete');