# Release Process

Simple guide for releasing DotEnvify CLI to npm.

## One-Time Setup

These steps only need to be done once, before your first release.

### npm — Trusted Publisher (OIDC)

Publishing uses [npm Trusted Publishing](https://docs.npmjs.com/generating-provenance-statements#publishing-packages-with-provenance-via-github-actions) — no tokens or secrets needed. The trust is established via OIDC between GitHub Actions and npm.

1. **Log in to [npmjs.com](https://www.npmjs.com)**

2. **Configure Trusted Publisher** on the `@webbies.dev/dotenvify` package:
   - Go to **Package Settings** > **Trusted Publisher**
   - Publisher: **GitHub Actions**
   - Organization or user: `webb1es`
   - Repository: `dotenvify`
   - Workflow filename: `release.yml`
   - Environment name: _(leave blank)_
   - Save

### GitHub

1. **Verify Actions permissions**
   - Go to [Settings > Actions > General](https://github.com/webb1es/dotenvify/settings/actions)
   - Under **Workflow permissions**, ensure **Read and write permissions** is selected
   - This allows the workflow to create GitHub Releases and request OIDC tokens

2. **No secrets needed** — trusted publishing uses OIDC via `id-token: write` permission

---

## Releasing a New Version

1. **Test locally**
   ```bash
   npm run build
   npm run test

   # Smoke test
   echo "API_KEY\nabc123" > /tmp/test.txt
   node cli/bin/dotenvify.js /tmp/test.txt -o /tmp/test.env
   cat /tmp/test.env
   ```

2. **Create and push a version tag**
   ```bash
   git tag -a v<version> -m "Release v<version>"
   git push origin main
   git push origin v<version>
   ```

   **Version guidelines:**
   - **Patch** (v2.0.x): Bug fixes, minor tweaks
   - **Minor** (v2.x.0): New features, backward compatible
   - **Major** (vx.0.0): Breaking changes

3. **Monitor the release**
   - Go to [GitHub Actions](https://github.com/webb1es/dotenvify/actions)
   - Watch the **"Release dotenvify"** workflow complete

4. **Verify on npm**
   ```bash
   npm info @webbies.dev/dotenvify
   ```

---

## What Happens Automatically

When you push a version tag, the workflow:

1. Installs dependencies and builds all packages (core → CLI)
2. Runs tests
3. Syncs the git tag version into `cli/package.json`
4. Bundles `@dotenvify/core` into the CLI via `tsup`
5. Publishes `@webbies.dev/dotenvify` to npm with provenance (via OIDC)
6. Creates a GitHub Release with auto-generated changelog

---

## Installation (for users)

```bash
npm install -g @webbies.dev/dotenvify
```

Or run without installing:
```bash
npx @webbies.dev/dotenvify input.txt -o .env
```

---

## Manual Publish (Emergency)

If CI is down, publish from your machine:

```bash
npm login --scope=@webbies.dev
npm run build
cd cli
npm version <version>
npm publish --access public
```

> Manual publish won't include provenance attestation.

---

## Troubleshooting

**Release failed?**
- Check GitHub Actions logs for specific errors
- Verify tag format is `vX.Y.Z` (e.g., `v2.0.1`)

**npm 403 Forbidden?**
- Verify Trusted Publisher is configured on npmjs.com for `@webbies.dev/dotenvify`
- Ensure the workflow filename matches exactly: `release.yml`
- Ensure the repository matches: `webb1es/dotenvify`
- Check that `id-token: write` permission is set in the workflow

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

To unpublish from npm (within 72 hours):
```bash
npm unpublish @webbies.dev/dotenvify@<version>
```
