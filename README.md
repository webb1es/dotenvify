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

### npm (Recommended)
```bash
npm install -g @webbies.dev/dotenvify
```

### Homebrew (macOS/Linux)
```bash
brew install webb1es/tap/dotenvify
```

### Scoop (Windows)
```bash
scoop bucket add webb1es https://github.com/webb1es/scoop-bucket.git
scoop install dotenvify
```

<details>
<summary>📦 Other Methods</summary>

#### Direct Download
Download from [GitHub Releases](https://github.com/webb1es/dotenvify/releases)

#### Build from Source
```bash
git clone https://github.com/webb1es/dotenvify.git
cd dotenvify && go build
```
</details>

## 🔮 Usage

### Basic File Mode

```bash
# Process a file and save to .env (default)
dotenvify your-vars.txt

# Save to a specific file
dotenvify your-vars.txt custom-output.env

# Ignore variables with lowercase keys
dotenvify -nl your-vars.txt
```

### Azure DevOps Mode

Fetch variables directly from Azure DevOps variable groups:

```bash
# Login to Azure CLI first
az login

# Set default organization (optional)
export DOTENVIFY_DEFAULT_ORG_URL="https://dev.azure.com/your-org/your-project"
dotenvify -azure -group "your-variable-group"

# Or specify URL directly
dotenvify -azure -url "https://dev.azure.com/your-org/your-project" -group "your-variable-group"
```

**Options:**
- `-out file.env` - Custom output file
- `-nl` - Ignore lowercase variables  
- `-export` - Add 'export' prefix
- `-urls` - Only URL values

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
- 📦 **Easy Install**: npm, Homebrew, Scoop

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
