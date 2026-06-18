<img src="landing/public/logo.gif" alt="DotEnvify" width="120" align="left" />

# DotEnvify

Turn messy key-value text into a clean `.env` file — from the command line or your JetBrains IDE.

<a href="https://www.npmjs.com/package/@webbies.dev/dotenvify"><img src="https://img.shields.io/npm/v/@webbies.dev/dotenvify.svg" alt="npm" /></a>
&nbsp;
<a href="https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify"><img src="https://img.shields.io/badge/JetBrains-Plugin-000000?style=flat-square&logo=jetbrains&logoColor=white" alt="JetBrains" /></a>
&nbsp;
<a href="./LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="MIT" /></a>

<br clear="left" />

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
| `--skip-lower`      |       | Drop keys that contain lowercase letters             |
| `--url-only`        |       | Keep only values that are URLs                       |

### Input it understands

You can mix any of these in one file — DotEnvify figures out each line. Lines
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

Do the same thing inside IntelliJ, WebStorm, PyCharm, GoLand, Rider, and other
JetBrains IDEs — plus a couple of extras. Install it from the
[Marketplace](https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify).

- **Convert** — paste or convert a file/selection to `.env` with a live preview.
- **Azure DevOps** — sign in with the Azure CLI (`az login`), pick a variable
  group, and pull its variables into a `.env` file. Secret values stay in Azure.
  Needs the [Azure CLI](https://aka.ms/azcli).
- **Diagnostics** — find keys your code uses but the `.env` is missing, and keys
  in the `.env` that nothing uses.

## Contributing

Building from source, the repo layout, and how releases work are in
[CONTRIBUTING.md](./CONTRIBUTING.md).

## License

MIT — see [LICENSE](./LICENSE).
