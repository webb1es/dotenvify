<p align="center">
  <img src="https://raw.githubusercontent.com/webb1es/dotenvify/main/docs/assets/logo.png" alt="DotEnvify" width="80" />
</p>

<h3 align="center">DotEnvify</h3>

<p align="center">Convert key-value pairs to environment variables with zero hassle</p>

<p align="center">
  <a href="https://www.npmjs.com/package/@webbies.dev/dotenvify"><img src="https://img.shields.io/npm/v/@webbies.dev/dotenvify.svg" alt="npm" /></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="MIT" /></a>
  <a href="https://nodejs.org/"><img src="https://img.shields.io/badge/Node.js-18+-339933.svg?logo=node.js&logoColor=white" alt="Node" /></a>
</p>

---

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

After doing this manually one too many times, this tool was rage-coded into existence. You're welcome.

## Install

```bash
npm install -g @webbies.dev/dotenvify
```

Or run directly without installing:

```bash
npx @webbies.dev/dotenvify vars.txt -o .env
```

<details>
<summary><strong>Troubleshooting</strong></summary>

**Permission errors (EACCES):**
```bash
# Option 1: Use npx (no install needed)
npx @webbies.dev/dotenvify [commands]

# Option 2: Fix npm permissions (recommended)
mkdir ~/.npm-global
npm config set prefix '~/.npm-global'
export PATH=~/.npm-global/bin:$PATH  # Add to ~/.bashrc or ~/.zshrc
npm install -g @webbies.dev/dotenvify
```

**Windows PowerShell:** If you get an execution policy error, use Command Prompt (cmd.exe) or run `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser` in PowerShell as Administrator.

</details>

## Usage

```bash
dotenvify <source> [options]
```

```bash
# Convert a file to .env (default)
dotenvify vars.txt

# Custom output path
dotenvify vars.txt -o production.env

# Add export prefix to all variables
dotenvify vars.txt --export

# Overwrite existing .env without backup
dotenvify vars.txt -f

# Preserve specific variables (keep their existing values)
dotenvify vars.txt --preserve "DATABASE_URL,API_KEY"

# Only extract URLs, skip lowercase keys
dotenvify vars.txt --url-only --skip-lower
```

## Options

| Option              | Alias | Description                                              |
|---------------------|-------|----------------------------------------------------------|
| `--output <file>`   | `-o`  | Output file path (default: `.env`)                       |
| `--export`          | `-e`  | Add `export` prefix to all variables                     |
| `--overwrite`       | `-f`  | Overwrite output without creating a backup               |
| `--preserve <vars>` | `-k`  | Comma-separated variables to keep existing values for    |
| `--skip-sort`       |       | Maintain original order (default: sorted alphabetically) |
| `--skip-lower`      |       | Skip variables with lowercase keys                       |
| `--url-only`        |       | Include only variables with HTTP/HTTPS URL values        |

## Supported Formats

DotEnvify auto-detects and handles all of these — even mixed in the same file:

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

Lines starting with `#` are treated as comments and ignored.

## Upgrading from v0.x (Go version)

v1.0.0 is a full rewrite in TypeScript. What changed:

|                        | v0.x (Go)                         | v1.x (TypeScript)                                                                            |
|------------------------|-----------------------------------|----------------------------------------------------------------------------------------------|
| **Runtime**            | Pre-built binary                  | Node.js 18+                                                                                  |
| **Install**            | Binary downloaded via postinstall | Standard npm install                                                                         |
| **Package**            | `@webbies.dev/dotenvify`          | `@webbies.dev/dotenvify` (same)                                                              |
| **Azure DevOps**       | Built into CLI                    | Available via [JetBrains plugin](https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify) |
| **Self-update**        | Built-in                          | Use `npm update -g`                                                                          |
| **Parsing/formatting** | Same                              | Same                                                                                         |

## Links

<a href="https://github.com/webb1es/dotenvify"><img src="https://img.shields.io/badge/GitHub-Source-181717?style=flat-square&logo=github&logoColor=white" alt="GitHub" /></a>&nbsp;
<a href="https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify"><img src="https://img.shields.io/badge/JetBrains-Plugin-000000?style=flat-square&logo=jetbrains&logoColor=white" alt="JetBrains" /></a>&nbsp;
<a href="https://github.com/webb1es/dotenvify/issues"><img src="https://img.shields.io/badge/Issues-Report-red?style=flat-square&logo=github&logoColor=white" alt="Issues" /></a>

## License

MIT — Go wild, make millions, just don't blame us when it formats your grocery list.
