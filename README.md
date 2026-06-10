<img src="landing/public/logo.gif" alt="DotEnvify" width="120" align="left" />

# DotEnvify

Convert messy key-value pairs into clean, standardized `.env` files — from the CLI or your IDE.

<a href="https://www.npmjs.com/package/@webbies.dev/dotenvify"><img src="https://img.shields.io/npm/v/@webbies.dev/dotenvify.svg" alt="npm" /></a>
&nbsp;
<a href="https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify"><img src="https://img.shields.io/badge/JetBrains-Plugin-000000?style=flat-square&logo=jetbrains&logoColor=white" alt="JetBrains" /></a>
&nbsp;
<a href="./LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="MIT" /></a>

<br clear="left" />

**Transform this:**

```
API_KEY
a1b2c3d4e5f6g7h8i9j0
DATABASE_URL
postgres://user:password@localhost:5432/db
```

**Into this:**

```env
API_KEY=a1b2c3d4e5f6g7h8i9j0
DATABASE_URL="postgres://user:password@localhost:5432/db"
```

## Install

```bash
npm install -g @webbies.dev/dotenvify
# or run without installing:
npx @webbies.dev/dotenvify vars.txt -o .env
```

## CLI usage

```bash
dotenvify <source> [options]
```

| Option              | Alias | Description                                              |
|---------------------|-------|----------------------------------------------------------|
| `--output <file>`   | `-o`  | Output file path (default: `.env`)                       |
| `--export`          | `-e`  | Add `export` prefix to all variables                     |
| `--overwrite`       | `-f`  | Overwrite output without creating a backup               |
| `--preserve <vars>` | `-k`  | Comma-separated variables to keep existing values for    |
| `--skip-sort`       |       | Maintain original order (default: sorted alphabetically) |
| `--skip-lower`      |       | Skip variables with lowercase keys                       |
| `--url-only`        |       | Include only variables with HTTP/HTTPS URL values        |

Backups (`.env.backup.N`) are written before any overwrite unless `-f` is given.

### Supported formats

Auto-detected and handled, even when mixed in one file (`#` lines are ignored):

```bash
API_KEY=a1b2c3d4e5f6g7h8i9j0          # KEY=VALUE
SECRET="my secret value"              # quoted
export NODE_ENV=production            # export prefix (stripped)
REDIS_HOST localhost                  # space-separated
DATABASE_URL                          # key on one line,
postgres://user:pass@localhost/db     #   value on the next
```

## JetBrains plugin

Pull Azure DevOps variable groups into `.env` and format right inside the IDE
(IntelliJ, WebStorm, GoLand, PyCharm, Rider, …). Install from the
[JetBrains Marketplace](https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify).

- **Azure DevOps** — sign in with the local Azure CLI (`az login`); load a variable
  group, pick which variables to apply, edit values inline, and write to a chosen
  `.env*` file with a merge preview. Secret values are shown but skipped (Azure does
  not expose them via API). Requires the [Azure CLI](https://aka.ms/azcli) installed.
- **Convert** — paste or convert files/selections to `.env` with live preview.
- **Diagnostics** — detect missing/unused keys across many languages.

## Repository

| Path               | Description                              |
|--------------------|------------------------------------------|
| `packages/core`    | `@dotenvify/core` — parser, formatter, IO |
| `cli`              | `@webbies.dev/dotenvify` — CLI tool       |
| `plugins/jetbrains`| IntelliJ-platform plugin (Kotlin/Gradle)  |
| `landing`          | Product site                              |

## Development

```bash
npm install     # install dependencies
npm run build   # build all packages
npm run test    # run tests
```

JetBrains plugin (from `plugins/jetbrains`, builds on **JDK 17** — pinned via
`gradle/gradle-daemon-jvm.properties`):

```bash
./gradlew test buildPlugin   # artifact in build/distributions/
./gradlew runIde             # launch a sandbox IDE
```

## Releasing

- **CLI** — `npm run build && npm run test`, then push a tag; GitHub Actions
  publishes to npm via OIDC:
  ```bash
  git tag -a v<version> -m "Release v<version>" && git push origin main v<version>
  ```
- **JetBrains** — first upload the built `.zip` manually to the Marketplace, then:
  ```bash
  JETBRAINS_MARKETPLACE_TOKEN=<token> ./gradlew publishPlugin
  ```
  Signing keys (`private.pem`, `chain.crt`) are git-ignored — keep them safe.

## License

MIT — see [LICENSE](./LICENSE).
