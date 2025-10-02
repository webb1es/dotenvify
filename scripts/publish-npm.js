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
  // Remove 'v' prefix from version if present
  const npmVersion = version.startsWith('v') ? version.slice(1) : version;
  
  // Update package.json version to match the release
  const packageJsonPath = path.join(__dirname, '..', 'package.json');
  const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
  const packageName = packageJson.name;
  packageJson.version = npmVersion;
  
  fs.writeFileSync(packageJsonPath, JSON.stringify(packageJson, null, 2) + '\n');
  console.log(`Updated package.json version to ${npmVersion}`);
  
  // Set npm authentication first (needed for npm view to work properly)
  const npmToken = process.env.NPM_TOKEN;
  if (!npmToken) {
    console.error('NPM_TOKEN environment variable is required');
    process.exit(1);
  }
  
  // Configure npm registry authentication
  execSync(`npm config set //registry.npmjs.org/:_authToken=${npmToken}`, { stdio: 'inherit' });

  // Idempotency: if the version already exists in the registry, skip publish successfully
  try {
    const viewCmd = `npm view ${packageName}@${npmVersion} version`;
    const existing = execSync(viewCmd, { stdio: ['ignore', 'pipe', 'pipe'] }).toString().trim();
    if (existing === npmVersion) {
      console.log(`✅ ${packageName}@${npmVersion} already exists on npm. Skipping publish.`);
      process.exit(0);
    }
  } catch (_) {
    // If npm view fails, we assume the version does not exist yet and proceed.
  }
  
  console.log(`🚀 Proceeding with npm publish for version ${npmVersion}`);
  
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
  try {
    const publishOutput = execSync('npm publish --access public', { stdio: ['ignore', 'pipe', 'pipe'] }).toString();
    // Echo npm publish output for visibility in CI logs
    if (publishOutput) process.stdout.write(publishOutput);
    console.log(`✅ Successfully published ${packageName}@${npmVersion} to npm`);
    console.log(`🎉 Package is now available for installation via: npm install -g ${packageName}`);
  } catch (error) {
    // Check if it's an idempotent publish error (version already exists)
    const errOut = (error && (error.stderr?.toString() || error.stdout?.toString() || error.message)) || '';
    const alreadyPublished = (
      errOut.includes('E409') ||
      errOut.includes('409') ||
      errOut.includes('Conflict') ||
      errOut.includes('E403') ||
      errOut.includes('EPUBLISHCONFLICT') ||
      /cannot publish over (the )?previously published versions/i.test(errOut)
    );
    if (alreadyPublished) {
      console.log(`✅ Version ${npmVersion} already published to npm (idempotent: treating conflict/forbidden as success)`);
      console.log('📦 npm package is up to date');
      // Exit successfully for conflicts indicating the version already exists
      process.exit(0);
    } else {
      // Re-throw other errors
      throw error;
    }
  }
  
} catch (error) {
  console.error('Failed to publish npm package:', error.message);
  process.exit(1);
}