# Release Process Documentation

This document outlines the streamlined release process for DotEnvify.

## 🎯 Overview

DotEnvify uses an automated release process that publishes to multiple package managers:
- **GitHub Releases** (binaries for manual download)
- **Homebrew** (macOS/Linux package manager)
- **Scoop** (Windows package manager)  
- **npm** (Node.js package manager)

## 🚀 Creating a Release

### Prerequisites

Ensure these secrets are configured in GitHub repository settings:
- `HOMEBREW_TAP_TOKEN` - Personal access token with repo scope for Homebrew tap
- `SCOOP_BUCKET_TOKEN` - Personal access token with repo scope for Scoop bucket
- `NPM_TOKEN` - Automation token from npmjs.com account

### Release Steps

1. **Ensure code is ready**
   ```bash
   # Test locally
   go build && ./dotenvify --help
   ```

2. **Create and push release tag**
   ```bash
   # Create tag (use semantic versioning: vMAJOR.MINOR.PATCH)
   git tag v1.0.0
   
   # Push tag to trigger release
   git push origin v1.0.0
   ```

3. **Monitor release progress**
   - Go to [GitHub Actions](https://github.com/webb1es/dotenvify/actions)
   - Watch the "Release dotenvify" workflow
   - All steps should complete successfully

## 🔧 How It Works

### 1. GoReleaser Phase
When a version tag is pushed, GitHub Actions:
- Builds binaries for all platforms (Linux, macOS, Windows × amd64/arm64)
- Creates archives and checksums
- Publishes GitHub Release with binaries
- Updates Homebrew tap formula
- Updates Scoop bucket manifest

### 2. npm Publishing Phase
After GoReleaser completes:
- Extracts version from git tag
- Updates `package.json` version
- Copies platform binaries from GoReleaser output
- Publishes scoped package `@webbies.dev/dotenvify` to npm

## 📦 Package Outputs

Each release creates:

| Platform | Package Manager | Install Command |
|----------|----------------|----------------|
| Cross-platform | npm | `npm install -g @webbies.dev/dotenvify` |
| macOS/Linux | Homebrew | `brew install webb1es/tap/dotenvify` |
| Windows | Scoop | `scoop install dotenvify` |
| Manual | GitHub Releases | Download binary from releases page |

## 🐛 Troubleshooting

### Release Failed
1. Check GitHub Actions logs for specific error
2. Common issues:
   - **Invalid tag format**: Use `vX.Y.Z` format only
   - **Missing secrets**: Verify all tokens are configured
   - **npm conflict**: Version already exists (usually harmless)

### Partial Release Success
If GoReleaser succeeds but npm fails:
1. npm publishing is non-critical - other packages still work
2. Can manually publish npm package later if needed
3. Future releases will work normally

### Rolling Back
If a release has critical issues:
1. Delete the Git tag: `git tag -d v1.0.0 && git push origin --delete v1.0.0`
2. Delete GitHub Release from web interface
3. Package managers will retain old working versions

## 🎯 Design Principles

This release process follows these principles:

1. **Simple**: Minimal scripts, leverage existing tools
2. **Reliable**: Each step can fail independently without breaking others
3. **Fast**: Parallel publishing to multiple platforms
4. **Maintainable**: Clear separation of concerns, easy to debug
5. **Automated**: Zero manual intervention required for normal releases

## 📚 Related Documentation

- [GoReleaser Configuration](.goreleaser.yml)
- [GitHub Actions Workflow](.github/workflows/release.yml)
- [Package Managers Setup](README.md#installation)