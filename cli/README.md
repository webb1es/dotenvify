<img src="https://raw.githubusercontent.com/webb1es/dotenvify/main/landing/public/logo.gif" alt="DotEnvify" width="120" align="left" />

# DotEnvify

Convert messy key-value pairs into clean, standardized `.env` files.

<a href="https://www.npmjs.com/package/@webbies.dev/dotenvify"><img src="https://img.shields.io/npm/v/@webbies.dev/dotenvify.svg" alt="npm" /></a>
&nbsp;
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="MIT" /></a>
&nbsp;
<a href="https://nodejs.org/"><img src="https://img.shields.io/badge/Node.js-18+-339933.svg?logo=node.js&logoColor=white" alt="Node" /></a>

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

**Windows PowerShell:** If you get an execution policy error, use Command Prompt (cmd.exe) or run
`Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser` in PowerShell as Administrator.

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

DotEnvify auto-detects and handles all of these, even mixed in the same file:

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

## IDE Plugins

Prefer working in your editor? DotEnvify has IDE plugins with features beyond the CLI: Azure DevOps integration,
paste-and-format, and real-time `.env` diagnostics.

|                                                                                                                                                                                                              | Plugin                                     | Highlights                                                |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------|-----------------------------------------------------------|
| <a href="https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify"><img src="https://img.shields.io/badge/JetBrains-Plugin-000000?style=flat-square&logo=jetbrains&logoColor=white" alt="JetBrains" /></a> | IntelliJ, WebStorm, GoLand, PyCharm, Rider | Azure DevOps variable groups, paste & format, diagnostics |
| <img src="https://img.shields.io/badge/VS_Code-Coming_soon-007ACC?style=flat-square&logo=visualstudiocode&logoColor=white" alt="VS Code" />                                                                  | Visual Studio Code                         | Parser, formatter, diagnostics                            |

## Upgrading from v0.x (Go version)

v1.0 is a full rewrite in TypeScript. **If you're upgrading, read this.**

### What's different

|                  | v0.x (Go)                        | v1.x (TypeScript)                                                                                   |
|------------------|----------------------------------|-----------------------------------------------------------------------------------------------------|
| **Runtime**      | Pre-built Go binary              | Requires Node.js 18+                                                                                |
| **Install**      | Binary via postinstall script    | Standard `npm install`                                                                              |
| **Command**      | `dotenvify -azure -group "Vars"` | `dotenvify vars.txt -o .env`                                                                        |
| **Azure DevOps** | Built into CLI (`-azure` flag)   | **Removed.** Use the [JetBrains plugin](https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify) |
| **Self-update**  | Built-in (`-update`)             | `npm update -g @webbies.dev/dotenvify`                                                              |
| **Parsing**      | Same formats supported           | Same formats supported                                                                              |

### Breaking changes

- **Azure DevOps is no longer in the CLI.** Variable group fetching has moved to
  the [JetBrains plugin](https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify), which provides a richer experience
  with IDE integration. If you relied on `-azure`, `-org`, or `-group` flags, install the plugin instead.
- **Node.js 18+ is now required.** The Go binary is no longer distributed.
- **Flag names changed.** `-nl` is now `--skip-lower`, `-ns` is now `--skip-sort`, `-out` is now `--output` / `-o`.
- **No subcommand.** Usage is `dotenvify <source>` (not `dotenvify convert <source>`).

## Links

<a href="https://github.com/webb1es/dotenvify"><img src="https://img.shields.io/badge/GitHub-Source-181717?style=flat-square&logo=github&logoColor=white" alt="GitHub" /></a>
&nbsp;
<a href="https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify"><img src="https://img.shields.io/badge/JetBrains-Plugin-000000?style=flat-square&logo=jetbrains&logoColor=white" alt="JetBrains" /></a>
&nbsp;
<a href="https://github.com/webb1es/dotenvify/issues"><img src="https://img.shields.io/badge/Issues-Report-red?style=flat-square&logo=github&logoColor=white" alt="Issues" /></a>

## License

MIT. Go wild, make millions, just don't blame us when it formats your grocery list.
