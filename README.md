<img src="landing/public/logo.gif" alt="DotEnvify" width="120" align="left" />

# DotEnvify

Convert messy key-value pairs into clean, standardized `.env` files — with zero hassle.

<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT" /></a>&nbsp;
<a href="https://www.npmjs.com/package/@webbies.dev/dotenvify"><img src="https://img.shields.io/npm/v/@webbies.dev/dotenvify.svg" alt="npm version" /></a>&nbsp;
<a href="https://www.typescriptlang.org/"><img src="https://img.shields.io/badge/TypeScript-5.8-3178C6.svg?logo=typescript&logoColor=white" alt="TypeScript" /></a>&nbsp;
<a href="https://nodejs.org/"><img src="https://img.shields.io/badge/Node.js-18+-339933.svg?logo=node.js&logoColor=white" alt="Node.js" /></a>&nbsp;
<a href="https://github.com/webb1es/dotenvify/pulls"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" alt="PRs Welcome" /></a>

<br clear="left" />

After doing this manually one too many times, this tool was rage-coded into existence. You're welcome.

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

## Installation

```bash
npm install -g @webbies.dev/dotenvify
```

## Quick Start

```bash
# Convert a file to .env
dotenvify vars.txt

# Custom output path
dotenvify vars.txt -o production.env

# Add export prefix
dotenvify vars.txt --export
```

## Features

| | Feature | Description |
|---|---|---|
| **Auto-Detect** | Smart Parsing | Handles `KEY=VALUE`, `KEY VALUE`, key-on-separate-lines, quoted values, `export` prefixes — even mixed together |
| **Backup** | Safe by Default | Automatic backups with incremental counters before any overwrite |
| **Lock** | Preserve Mode | Keep existing values for specific variables when regenerating `.env` files |
| **Filter** | Flexible Filtering | Skip lowercase keys, filter to URLs only, sort alphabetically or keep original order |
| **Quote** | Smart Quoting | Automatically quotes values containing spaces or URLs |

## Supported Input Formats

DotEnvify auto-detects and parses all of these — even when mixed in the same file:

```bash
# KEY=VALUE
API_KEY=a1b2c3d4e5f6g7h8i9j0

# Quoted values
SECRET="my secret value"

# export prefix (stripped automatically)
export NODE_ENV=production

# Space-separated
REDIS_HOST localhost

# Key on one line, value on the next
DATABASE_URL
postgres://user:password@localhost:5432/db
```

> Lines starting with `#` are treated as comments and ignored.

## CLI Reference

```
dotenvify <source> [options]
```

| Option | Alias | Description |
|---|---|---|
| `--output <file>` | `-o` | Output file path (default: `.env`) |
| `--export` | `-e` | Add `export` prefix to all variables |
| `--overwrite` | `-f` | Overwrite output without creating a backup |
| `--preserve <vars>` | `-k` | Comma-separated variables to keep existing values for |
| `--skip-sort` | | Maintain original order (default: sorted alphabetically) |
| `--skip-lower` | | Skip variables with lowercase keys |
| `--url-only` | | Include only variables with HTTP/HTTPS URL values |

### Examples

```bash
# Overwrite without backup
dotenvify vars.txt -f

# Preserve DB creds when regenerating
dotenvify vars.txt --preserve "DATABASE_URL,API_SECRET"

# Only extract URLs, skip lowercase keys
dotenvify vars.txt --url-only --skip-lower

# Full pipeline: export-prefixed, custom output, no sorting
dotenvify vars.txt -o .env.local --export --skip-sort
```

## IDE Plugins

Use DotEnvify directly in your editor — with features the CLI can't offer, like Azure DevOps integration, paste-and-format, and real-time diagnostics.

| Plugin | Highlights | Status |
|---|---|---|
| <a href="https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify"><img src="https://img.shields.io/badge/JetBrains-IntelliJ_%2F_WebStorm_%2F_GoLand_%2F_PyCharm_%2F_Rider-000000?style=flat-square&logo=jetbrains&logoColor=white" alt="JetBrains" /></a> | Azure DevOps variable groups, paste & format, `.env` diagnostics | <img src="https://img.shields.io/badge/available-brightgreen?style=flat-square" alt="Available" /> |
| <a href="./plugins/vscode"><img src="https://img.shields.io/badge/VS_Code-Extension-007ACC?style=flat-square&logo=visualstudiocode&logoColor=white" alt="VS Code" /></a> | Parser, formatter, diagnostics — built on @dotenvify/core | <img src="https://img.shields.io/badge/coming_soon-lightgrey?style=flat-square" alt="Coming soon" /> |

> **Migrating from v0.x?** Azure DevOps support has moved from the CLI to the [JetBrains plugin](https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify) — a better fit with richer IDE integration. See the [CLI README](./cli#upgrading-from-v0x-go-version) for the full migration guide.

## Ecosystem

| Package | Description |
|---|---|
| <a href="./cli"><img src="https://img.shields.io/badge/CLI-Command--line_tool-4A154B?style=flat-square&logo=windowsterminal&logoColor=white" alt="CLI" /></a> | File conversion, scripting, CI/CD pipelines |
| <a href="./packages/core"><img src="https://img.shields.io/badge/@dotenvify/core-Shared_library-3178C6?style=flat-square&logo=typescript&logoColor=white" alt="Core" /></a> | Parser, formatter, IO — powers CLI and plugins |
| <a href="./landing"><img src="https://img.shields.io/badge/Landing_Page-Product_site-FF6F61?style=flat-square&logo=vercel&logoColor=white" alt="Landing" /></a> | Live demo and docs |

## Development

```bash
npm install          # Install dependencies
npm run build        # Build all packages
npm run test         # Run tests
npm run dev:landing  # Dev mode (landing page)
```

<details>
<summary><strong>Project Structure</strong></summary>

```
dotenvify/
├── packages/core/       # @dotenvify/core — shared TS library
├── cli/                 # CLI tool (Commander.js)
├── plugins/
│   ├── jetbrains/       # Kotlin — Gradle build
│   └── vscode/          # VS Code extension
├── landing/             # Product landing page
└── docs/                # Shared docs & assets
```

</details>

## Contributing

Found a bug? Have a feature idea? PRs welcome!

Check out the <a href="https://github.com/webb1es/dotenvify/issues"><img src="https://img.shields.io/badge/issues-GitHub-red?style=flat-square&logo=github&logoColor=white" alt="Issues" /></a> or submit a <a href="https://github.com/webb1es/dotenvify/pulls"><img src="https://img.shields.io/badge/pull_requests-GitHub-blue?style=flat-square&logo=github&logoColor=white" alt="Pull Requests" /></a>

## License

<a href="./LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="MIT License" /></a>

Go wild, make millions, just don't blame us when it formats your grocery list.

---

<p align="center"><i>"Life's too short for manual formatting."</i></p>
