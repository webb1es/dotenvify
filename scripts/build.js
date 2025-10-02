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
if (process.env.CI) {
  console.log('CI environment detected - binaries should be provided by GoReleaser');
  
  // List available binaries in dist directory
  const distDir = path.join(__dirname, '..', 'dist');
  if (fs.existsSync(distDir)) {
    const files = fs.readdirSync(distDir);
    console.log('Available files in dist:', files);
    
    // Copy binaries from dist to bin directory
    const binaryPattern = /^dotenvify_/;
    files.forEach(file => {
      if (binaryPattern.test(file)) {
        const srcPath = path.join(distDir, file);
        const destPath = path.join(binDir, file);
        
        if (fs.statSync(srcPath).isFile()) {
          fs.copyFileSync(srcPath, destPath);
          fs.chmodSync(destPath, '755');
          console.log(`Copied binary: ${file}`);
        }
      }
    });
  } else {
    console.error('dist directory not found - GoReleaser should have created it');
    process.exit(1);
  }
} else {
  console.log('Local environment - skipping binary copy (use GoReleaser for local builds)');
}

console.log('Build complete');