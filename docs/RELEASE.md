# Release Process

## One-Time Setup

### npm

Configure Trusted Publisher on `@webbies.dev/dotenvify`:

- **Package Settings** > **Trusted Publisher**
- Publisher: **GitHub Actions** | Org: `webb1es` | Repo: `dotenvify` | Workflow: `release.yml`

### GitHub

- [Settings > Actions > General](https://github.com/webb1es/dotenvify/settings/actions) > **Workflow permissions** > *
  *Read and write**
- No secrets needed — publishing uses OIDC

---

## Release

```bash
# 1. Test
npm run build && npm run test

# 2. Tag and push
git tag -a v<version> -m "Release v<version>"
git push origin main
git push origin v<version>
```

The workflow builds, tests, publishes to npm, and creates a GitHub Release automatically.

Verify: `npm info @webbies.dev/dotenvify`

---

## Version Guidelines

- **Patch** (v1.0.x): Bug fixes
- **Minor** (v1.x.0): New features, backward compatible
- **Major** (vx.0.0): Breaking changes

## Rollback

```bash
git tag -d v<version>
git push origin --delete v<version>
npm unpublish @webbies.dev/dotenvify@<version>   # within 72 hours
```

## Manual Publish

```bash
npm login --scope=@webbies.dev
npm run build
cd cli && npm publish --access public
```
