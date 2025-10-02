# 🧙‍♂️ DotEnvify

> Because copying and pasting environment variables shouldn't be your midnight cardio 🏃‍♂️💨

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/webb1es/dotenvify)](https://goreportcard.com/report/github.com/webb1es/dotenvify)

## 🤦‍♂️ The Problem

Your team lead drops this in the chat:

```
API_KEY
a1b2c3d4e5f6g7h8i9j0
DATABASE_URL
postgres://user:password@localhost:5432/db
```

And you need:

```
export API_KEY="a1b2c3d4e5f6g7h8i9j0"
export DATABASE_URL="postgres://user:password@localhost:5432/db"
```

After doing this manually one too many times, I rage-coded this tool. You're welcome.

## 🚀 Installation

### 🍺 Homebrew (macOS and Linux)

```bash
brew install webb1es/tap/dotenvify
```

### 🪣 Scoop (Windows)

```bash
scoop bucket add webb1es https://github.com/webb1es/scoop-bucket.git
scoop install dotenvify
```

<details>
<summary>📦 Other Installation Methods</summary>

#### Direct Download
Download the latest release from the [GitHub Releases page](https://github.com/webb1es/dotenvify/releases).

#### Build from Source
```bash
git clone https://github.com/webb1es/dotenvify.git
cd dotenvify
go build -o dotenvify dotenvify.go
sudo mv dotenvify /usr/local/bin/
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

### 🚀 Azure DevOps Mode

Fetch variables directly from your Azure DevOps variable groups:

```bash
# First, make sure you're logged in
az login

# Using environment variable for default URL
export DOTENVIFY_DEFAULT_ORG_URL="https://dev.azure.com/your-org/your-project"
dotenvify -az -g "your-variable-group"

# Using custom project URL
dotenvify -az -u "https://dev.azure.com/your-org/your-project" -g "your-variable-group"

# Save to a specific file
dotenvify -az -g "your-variable-group" -out "custom-output.env"

# Ignore variables with lowercase keys
dotenvify -az -g "your-variable-group" -nl
```

> **Note:** You can set a default Azure DevOps URL using the `DOTENVIFY_DEFAULT_ORG_URL` environment variable. You can override this with the `-u` or `-url` flag.

<details>
<summary>🔒 Security & Authentication Details</summary>

DotEnvify uses your existing Azure CLI authentication:
- No credentials stored or handled by the tool
- Tokens are used only in memory
- Respects your organization's security policies (including MFA)

Just make sure you're logged in with `az login` before running the tool.
</details>

## ✨ Key Features

- ⚡ **Fast**: Written in Go for blazing speed
- 🔄 **Azure Integration**: Fetch variables directly from Azure DevOps
- 🔒 **Secure Auth**: Uses your existing Azure CLI credentials
- 🧹 **Smart Parsing**: Skips empty lines automatically
- 🔍 **Format Detection**: Recognizes when files are already in the expected format
- 🔤 **Case Control**: Option to ignore lowercase variables with `-nl`
- 👻 **Zero Dependencies**: No package nightmares

## 📝 Input Format

DotEnvify now supports multiple input formats:

```
# Traditional format (key and value on separate lines)
KEY1
value1
KEY2
value2

# KEY=VALUE format
KEY1=value1
KEY2=value2

# KEY="VALUE" format with quotes
KEY1="value1"
KEY2="value1 with spaces"

# KEY VALUE format (on same line)
KEY1 value1
KEY2 value2

# Multiple KEY VALUE pairs on one line
KEY1 value1 KEY2 value2
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
