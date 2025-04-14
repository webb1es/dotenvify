# 🧙‍♂️ DotEnvify

> ✨ Transform plain text variables into magical .env exports with a single command! ✨

## 🤔 The Problem

Tired of manually converting plain text environment variables into proper `.env` format? Sick of copying and pasting KEY-VALUE pairs one by one? Annoyed by the constant reformatting required to get those variables into your development environment?

```
API_KEY
a1b2c3d4e5f6g7h8i9j0
DATABASE_URL
postgres://user:password@localhost:5432/db
DEBUG_MODE
true
```

Manually converting to:

```
export API_KEY="a1b2c3d4e5f6g7h8i9j0"
export DATABASE_URL="postgres://user:password@localhost:5432/db"
export DEBUG_MODE="true"
```

## 🚀 The Solution

**DotEnvify** is your nerdy little helper that automagically transforms plain text key-value pairs into properly formatted environment variables!

```bash
dotenvify plain-vars.txt
# ✓ Variables successfully formatted and saved to 'plain-vars.txt'
```

## ⚡ Quick Install

```bash
curl -s https://raw.githubusercontent.com/webb1es/dotenvify/main/remote-install.sh | bash
```

## 🛠️ Usage

```bash
dotenvify source_file [output_file]
```

- If no `output_file` specified, it overwrites the source (if all goes well)
- Handles errors gracefully (because we all make mistakes 🤷‍♂️)
- Shows colorful feedback (because terminal life should be vibrant 🌈)

## 🧪 Example

Turn this boring plain text:

```
SERVER_PORT
8080
LOG_LEVEL
debug
ENABLE_METRICS
true
```

Into this glorious `.env` format:

```
export SERVER_PORT="8080"
export LOG_LEVEL="debug"
export ENABLE_METRICS="true"
```

With just one command:

```bash
dotenvify config-vars.txt
```

## 💻 Features

- ⚙️ Simple, single-command operation
- 🛡️ Skips empty lines (because whitespace matters... or doesn't)
- 🎯 Smart error handling (creates separate output file when things go 💥)
- 🌈 Colorful terminal output (because life's too short for monochrome)
- 🖥️ Cross-platform (works on macOS, Linux, Windows)

## 🔧 Requirements

- Go (version 1.13 or later)
- A burning desire to optimize your workflow

## 🚀 Installation Options

### Direct Install (for the impatient)

```bash
curl -s https://raw.githubusercontent.com/webb1es/dotenvify/main/remote-install.sh | bash
```

### From Source (for the curious)

```bash
git clone https://github.com/webb1es/dotenvify.git
cd dotenvify
go build -o dotenvify dotenvify.go
sudo mv dotenvify /usr/local/bin/
```

## ❌ Uninstalling

When you've automated all possible environment variables in the universe:

```bash
sudo rm /usr/local/bin/dotenvify
```

## 🤓 For the Nerds

DotEnvify is written in Go for blazing fast performance because life's too short to wait for your scripts to run. The entire tool is less than 150 lines of code—proving once again that the best tools are often the simplest!

---

Made with ❤️ and probably too much ☕ by a developer who got tired of repetitive tasks.