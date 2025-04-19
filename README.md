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

# Using project URL (recommended)
dotenvify -az -u "https://dev.azure.com/your-org/your-project" -g "your-variable-group"

# Save to a specific file
dotenvify -az -g "your-variable-group" -out "custom-output.env"

# Ignore variables with lowercase keys
dotenvify -az -g "your-variable-group" -nl
```

<details>
<summary>🔒 Security & Authentication Details</summary>

DotEnvify uses your existing Azure CLI authentication:
- No credentials stored or handled by the tool
- Tokens are used only in memory
- Respects your organization's security policies (including MFA)

Just make sure you're logged in with `az login` before running the tool.
</details>

> 💡 **Tip**: Set `AZURE_DEVOPS_URL` in your environment to avoid typing the URL each time!

## ✨ Key Features

- ⚡ **Fast**: Written in Go for blazing speed
- 🔄 **Azure Integration**: Fetch variables directly from Azure DevOps
- 🔒 **Secure Auth**: Uses your existing Azure CLI credentials
- 🧹 **Smart Parsing**: Skips empty lines automatically
- 🔤 **Case Control**: Option to ignore lowercase variables with `-nl`
- 👻 **Zero Dependencies**: No package nightmares

## 📝 Input Format

```
KEY1
value1
KEY2
value2
```

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
