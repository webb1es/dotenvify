# Dotenvify

Turn messy key-value text into a clean `.env` file — from the command line.

[![npm](https://img.shields.io/npm/v/@webbies.dev/dotenvify.svg)](https://www.npmjs.com/package/@webbies.dev/dotenvify)
[![JetBrains](https://img.shields.io/badge/JetBrains-Plugin-000000?style=flat-square&logo=jetbrains&logoColor=white)](https://plugins.jetbrains.com/plugin/32351-dotenvify)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/webb1es/dotenvify/blob/main/LICENSE)

You paste this:

```
API_KEY
a1b2c3d4e5f6g7h8i9j0
DATABASE_URL
postgres://user:password@localhost:5432/db
```

You get this:

```env
API_KEY=a1b2c3d4e5f6g7h8i9j0
DATABASE_URL="postgres://user:password@localhost:5432/db"
```

Keys are sorted, values that need quotes get quoted, and your existing `.env` is
backed up before anything is written.

## Install

```bash
npm install -g @webbies.dev/dotenvify
# or run it once, without installing:
npx @webbies.dev/dotenvify vars.txt -o .env
```

## Usage

```bash
dotenvify <source> [options]
```

`<source>` is the file you want to convert.

| Option              | Short | What it does                                         |
|---------------------|-------|------------------------------------------------------|
| `--output <file>`   | `-o`  | Where to write (default: `.env`)                     |
| `--export`          | `-e`  | Add an `export ` prefix to every line                |
| `--overwrite`       | `-f`  | Overwrite without making a backup                    |
| `--preserve <vars>` | `-k`  | Keep current values for these keys (comma-separated) |
| `--skip-sort`       |       | Keep the original order instead of sorting           |
| `--skip-lower`      |       | Drop keys that are entirely lowercase                |
| `--url-only`        |       | Keep only values that are HTTP/HTTPS URLs            |

### Input it understands

You can mix any of these in one file — Dotenvify figures out each line. Lines
starting with `#` are left alone.

```bash
API_KEY=a1b2c3d4e5f6g7h8i9j0          # already KEY=VALUE
SECRET="my secret value"              # quoted
export NODE_ENV=production            # export prefix (removed)
REDIS_HOST localhost                  # separated by a space
DATABASE_URL                          # key on one line,
postgres://user:pass@localhost/db     #   value on the next
```

## JetBrains plugin

Prefer working in your IDE? The same conversion is available inside IntelliJ,
WebStorm, PyCharm, GoLand, Rider, and other JetBrains IDEs, plus Azure DevOps
variable-group import and `.env` diagnostics. Install it from the
[Marketplace](https://plugins.jetbrains.com/plugin/32351-dotenvify).

## License

MIT — see [LICENSE](https://github.com/webb1es/dotenvify/blob/main/LICENSE).
