#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// This script copies the appropriate binary to the bin directory during npm pack
// The actual binaries are built by GoReleaser during the release process

const packageJson = require('../package.json');
const version = packageJson.version;

console.log(`Building npm package for dotenvify v${version}`);
console.log('Environment variables:', {
  CI: process.env.CI,
  GITHUB_ACTIONS: process.env.GITHUB_ACTIONS,
  NODE_ENV: process.env.NODE_ENV,
  PWD: process.env.PWD || process.cwd()
});

// Create bin directory if it doesn't exist
const binDir = path.join(__dirname, '..', 'bin');
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

// Always assume CI environment when called by GoReleaser (dist directory exists)
const distDir = path.join(__dirname, '..', 'dist');
const isCI = process.env.CI || process.env.GITHUB_ACTIONS || fs.existsSync(distDir);

if (isCI) {
  console.log('CI environment detected - binaries should be provided by GoReleaser');
  console.log('Dist directory exists:', fs.existsSync(distDir));
  
  if (fs.existsSync(distDir)) {
    const files = fs.readdirSync(distDir);
    console.log('Available files in dist:', files);
    
    let foundBinaries = false;
    
    // Copy individual binaries from subdirectories
    files.forEach(file => {
      const filePath = path.join(distDir, file);
      if (fs.existsSync(filePath) && fs.statSync(filePath).isDirectory() && file.startsWith('dotenvify_')) {
        console.log(`Checking directory: ${file}`);
        const binaryFiles = fs.readdirSync(filePath);
        console.log(`Files in ${file}:`, binaryFiles);
        
        binaryFiles.forEach(binaryFile => {
          if (binaryFile === 'dotenvify' || binaryFile === 'dotenvify.exe') {
            const srcPath = path.join(filePath, binaryFile);
            const destPath = path.join(binDir, file + (binaryFile.endsWith('.exe') ? '.exe' : ''));
            
            fs.copyFileSync(srcPath, destPath);
            fs.chmodSync(destPath, '755');
            console.log(`Copied binary: ${file} -> ${destPath}`);
            foundBinaries = true;
          }
        });
      }
    });
    
    if (!foundBinaries) {
      console.error('No binary files found in dist subdirectories');
      console.log('Current working directory:', process.cwd());
      console.log('Script directory:', __dirname);
      process.exit(1);
    }
  } else {
    console.error('Dist directory not found');
    console.log('Looked for dist at:', distDir);
    console.log('Current working directory:', process.cwd());
    console.log('Available files in cwd:', fs.readdirSync(process.cwd()));
    process.exit(1);
  }
} else {
  console.log('Local environment - skipping binary copy (use GoReleaser for local builds)');
}

console.log('Build complete');