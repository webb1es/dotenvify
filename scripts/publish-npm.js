#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const version = process.env.RELEASE_VERSION;
const tag = process.env.RELEASE_TAG;

if (!version || !tag) {
  console.error('RELEASE_VERSION and RELEASE_TAG environment variables are required');
  process.exit(1);
}

console.log(`Publishing npm package version ${version} (${tag})`);

try {
  // Update package.json version to match the release
  const packageJsonPath = path.join(__dirname, '..', 'package.json');
  const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
  
  // Remove 'v' prefix from version if present
  const npmVersion = version.startsWith('v') ? version.slice(1) : version;
  packageJson.version = npmVersion;
  
  fs.writeFileSync(packageJsonPath, JSON.stringify(packageJson, null, 2) + '\n');
  console.log(`Updated package.json version to ${npmVersion}`);
  
  // Check if this version already exists on npm
  try {
    const checkResult = execSync(`npm view @webbies.dev/dotenvify@${npmVersion} version 2>/dev/null`, { encoding: 'utf8' });
    if (checkResult.trim() === npmVersion) {
      console.log(`✅ Version ${npmVersion} already exists on npm - skipping publish (this is expected behavior)`);
      console.log('📦 npm package is up to date');
      process.exit(0);
    }
  } catch (error) {
    // Version doesn't exist, proceed with publish
    console.log(`🔍 Version ${npmVersion} not found on npm - proceeding with publish`);
  }
  
  // Set npm authentication
  const npmToken = process.env.NPM_TOKEN;
  if (!npmToken) {
    console.error('NPM_TOKEN environment variable is required');
    process.exit(1);
  }
  
  // Configure npm registry authentication
  execSync(`npm config set //registry.npmjs.org/:_authToken=${npmToken}`, { stdio: 'inherit' });
  
  // Build the package (copies binaries to bin directory)
  console.log('Building npm package...');
  execSync('npm run build', { stdio: 'inherit' });
  
  // Check if binaries exist
  const binDir = path.join(__dirname, '..', 'bin');
  const binFiles = fs.readdirSync(binDir).filter(f => f.startsWith('dotenvify_'));
  
  if (binFiles.length === 0) {
    console.error('No binary files found in bin directory');
    process.exit(1);
  }
  
  console.log(`Found ${binFiles.length} binary files:`, binFiles);
  
  // Publish to npm
  console.log('📦 Publishing to npm...');
  execSync('npm publish --access public', { stdio: 'inherit' });
  
  console.log(`✅ Successfully published @webbies.dev/dotenvify@${npmVersion} to npm`);
  console.log('🎉 Package is now available for installation via: npm install -g @webbies.dev/dotenvify');
  
} catch (error) {
  console.error('Failed to publish npm package:', error.message);
  process.exit(1);
}