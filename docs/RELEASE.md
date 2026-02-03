# Release Process

Simple guide for releasing DotEnvify to npm.

## Prerequisites

Ensure `NPM_TOKEN` is configured in GitHub repository secrets.

## Release Steps

1. **Test locally**
   ```bash
   go build && ./dotenvify --help
   ```

2. **Create and push tag**
   ```bash
   # Create a semantic version tag (vMAJOR.MINOR.PATCH)
   git tag v<version>

   # Push changes and tag
   git push origin main
   git push origin v<version>
   ```

   **Version Guidelines:**
   - **Patch** (v0.3.x): Bug fixes, minor tweaks
   - **Minor** (v0.x.0): New features, backward compatible
   - **Major** (vx.0.0): Breaking changes

3. **Monitor release**
   - Go to [GitHub Actions](https://github.com/webb1es/dotenvify/actions)
   - Watch the "Release dotenvify" workflow complete

## What Happens

When you push a version tag:
1. GoReleaser builds binaries for all platforms (Linux, macOS, Windows)
2. GitHub Release is created with binaries and changelog
3. Package is published to npm as `@webbies.dev/dotenvify`

## Installation

Users install via:
```bash
npm install -g @webbies.dev/dotenvify
```

## Troubleshooting

**Release failed?**
- Check GitHub Actions logs for specific errors
- Verify tag format is `vX.Y.Z` (e.g., `v1.2.3`)
- Ensure `NPM_TOKEN` secret is set in repository settings

**Version already exists on npm?**
- The workflow handles this gracefully and continues
- Choose a new version number and create a new tag

**Need to rollback?**
```bash
# Delete tag locally and remotely
git tag -d v<version>
git push origin --delete v<version>
```
Then delete the GitHub Release from the web interface.
