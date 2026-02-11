# 🧙‍♂️ DotEnvify

> Convert key-value pairs to environment variables with zero hassle

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/webb1es/dotenvify)](https://goreportcard.com/report/github.com/webb1es/dotenvify)

**Transform this:**
```
API_KEY
a1b2c3d4e5f6g7h8i9j0
DATABASE_URL
postgres://user:password@localhost:5432/db
```

**Into this:**
```
export API_KEY="a1b2c3d4e5f6g7h8i9j0"
export DATABASE_URL="postgres://user:password@localhost:5432/db"
```

After doing this manually one too many times, I rage-coded this tool. You're welcome.

## 🚀 Installation

```bash
npm install -g "@webbies.dev/dotenvify"
```

<details>
<summary>📦 Alternative Methods</summary>

#### Build from Source
```bash
git clone https://github.com/webb1es/dotenvify.git
cd dotenvify && go build
```
</details>

## 🔄 Updating

Keep dotenvify up to date with the latest features and fixes:

### Self-Update (Recommended)
```bash
# Check for available updates
dotenvify -check-update

# Update to the latest version
dotenvify -update
```

### Via npm
```bash
npm update -g "@webbies.dev/dotenvify"
```

The self-update feature automatically downloads and installs the latest release from GitHub, making it easy to stay current regardless of how you installed dotenvify.

## 🔮 Usage

### Basic File Mode

```bash
# Process a file and save to .env (default)
dotenvify your-vars.txt

# Save to a specific file
dotenvify your-vars.txt custom-output.env

# Overwrite existing .env without backup
dotenvify -f your-vars.txt

# Preserve specific variables (keep their existing values)
# Note: Flags must come before file arguments
dotenvify -preserve "DATABASE_URL,API_KEY" your-vars.txt

# Ignore variables with lowercase keys
dotenvify -nl your-vars.txt
```

### Azure DevOps Mode

Fetch variables directly from Azure DevOps variable groups:

#### Step 1: Login to Azure CLI
```bash
az login
```

#### Step 2: Fetch variables

**Option A: Using default organization (recommended)**
```bash
# Set default organization (run once)
export DOTENVIFY_DEFAULT_ORG_URL="https://dev.azure.com/your-org/your-project"
```
```bash
# Fetch variables from a group
dotenvify -azure -group "your-variable-group"
```

**Option B: Specify organization directly**
```bash
dotenvify -azure -url "https://dev.azure.com/your-org/your-project" -group "your-variable-group"
```

**Options:**
- `-out file.env` - Custom output file
- `-preserve "VAR1,VAR2"` or `-k` - Keep existing values for specified variables
- `-nl` - Ignore lowercase variables
- `-export` - Add 'export' prefix
- `-urls` - Only URL values
- `-f` - Overwrite output file (default: backup existing file)

<details>
<summary>🔒 Security & Authentication Details</summary>

DotEnvify uses your existing Azure CLI authentication:
- No credentials stored or handled by the tool
- Tokens are used only in memory
- Respects your organization's security policies (including MFA)

Just make sure you're logged in with `az login` before running the tool.
</details>

## ✨ Features

- ⚡ **Fast**: Written in Go
- 🔄 **Azure DevOps**: Direct variable group integration
- 🔒 **Secure**: Uses existing Azure CLI auth
- 🧹 **Smart**: Auto-detects input formats
- 🔤 **Flexible**: Multiple output options
- 💾 **Safe**: Auto-backup with incremental counters
- 🛡️ **Preserve**: Keep existing values for specific variables
- 📦 **Easy Install**: npm
- 🔄 **Self-Updating**: Built-in update mechanism

## 📝 Supported Formats

```bash
# Key-value on separate lines
KEY1
value1

# KEY=VALUE format  
KEY1=value1

# KEY="VALUE" with quotes
KEY1="value with spaces"

# KEY VALUE on same line
KEY1 value1
```

Lines starting with `#` are treated as comments and ignored. The tool is smart enough to try different parsing strategies if one fails, making it robust against unfamiliar formats.

<details>
<summary>🤓 Why I Made This</summary>

Because those precious minutes you spend formatting env vars could be spent on:
- Actually coding something cool
- Optimizing your coffee brewing technique
- Staring at the wall contemplating your life choices
</details>

## 🔧 Contributing

Found a bug? 🐛 Have a feature idea? 💡 PRs welcome! Let's make this tool even more awesome together! 🚀

Check out the [issues page](https://github.com/webb1es/dotenvify/issues) or submit a [pull request](https://github.com/webb1es/dotenvify/pulls).

## 📄 License

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

MIT License - Go wild, make millions, just don't blame me when it formats your grocery list. 🛒📝

---

<div align="center">

**Made with 💻, ☕, and the burning desire to never format env vars manually again**

![Works on My Machine](https://forthebadge.com/images/badges/works-on-my-machine.svg)
![Built with Swag](https://forthebadge.com/images/badges/built-with-swag.svg)

*"Life's too short for manual formatting."* 🧙‍♂️✨

</div>
