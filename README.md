# 🧙‍♂️ DotEnvify

> Because copying and pasting environment variables shouldn't be your midnight cardio 🏃‍♂️💨

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/webb1es/dotenvify)](https://goreportcard.com/report/github.com/webb1es/dotenvify)

## 🤦‍♂️ The Problem

We've all been there. It's 11:47 PM, you're setting up a new project, and your team lead drops this in the chat:

```
API_KEY
a1b2c3d4e5f6g7h8i9j0
DATABASE_URL
postgres://user:password@localhost:5432/db
SECRET_TOKEN
shhhh-its-a-secret
```

And you need:

```
export API_KEY="a1b2c3d4e5f6g7h8i9j0"
export DATABASE_URL="postgres://user:password@localhost:5432/db"
export SECRET_TOKEN="shhhh-its-a-secret"
```

After doing this manually for the 17th time last month, I rage-coded this tool. You're welcome.

## 🚀 Install This Sanity-Saver

### 🍺 Homebrew (macOS and Linux)

```bash
# Install from the Homebrew tap
brew install webb1es/tap/dotenvify
```

### 🪣 Scoop (Windows)

```bash
# Add the scoop bucket
scoop bucket add webb1es https://github.com/webb1es/scoop-bucket.git

# Install dotenvify
scoop install dotenvify
```

### 📦 Direct Download (All Platforms)

Download the latest release for your platform from the [GitHub Releases page](https://github.com/webb1es/dotenvify/releases).


### 👨‍💻 The "I Read Every Line of Code" Way

```bash
git clone https://github.com/webb1es/dotenvify.git
cd dotenvify
go build -o dotenvify dotenvify.go
sudo mv dotenvify /usr/local/bin/
```

## 🔮 How To Use It

### Basic Usage

Seriously, it's simple:

```bash
# Process a file and save to .env (default)
dotenvify your-vars.txt

# Save to a specific file instead of .env
dotenvify your-vars.txt custom-output.env

# Get help with available options
dotenvify -h

# Ignore variables with lowercase keys
dotenvify -no-lower your-vars.txt
# or using the shorthand flag
dotenvify -nl your-vars.txt

# The output is always in export KEY="VALUE" format
# .env will be overwritten if it already exists
```

## 🔌 Supported Plugins

DotEnvify uses a plugin architecture that allows it to integrate with various services. Currently supported plugins:

### 🚀 Azure DevOps Integration

Fetch variables directly from your Azure DevOps variable groups:

```bash
# First, make sure you're logged in to Azure CLI
az login

# Fetch variables using your project URL (easiest method)
dotenvify -azure -url "https://dev.azure.com/your-org/your-project" -group "your-variable-group"
# or using shorthand flags
dotenvify -az -u "https://dev.azure.com/your-org/your-project" -g "your-variable-group"

# Or use the traditional way with organization and project
dotenvify -azure -org "your-org" -project "your-project" -group "your-variable-group"
# or using shorthand flags
dotenvify -az -o "your-org" -p "your-project" -g "your-variable-group"

# Set a default URL in your environment to make it even easier
export AZURE_DEVOPS_URL="https://dev.azure.com/your-org/your-project"
dotenvify -azure -group "your-variable-group"
# or using shorthand flags
dotenvify -az -g "your-variable-group"

# If you don't provide a URL or org/project, dotenvify will ask you interactively
dotenvify -azure -group "your-variable-group"

# Ignore variables with lowercase keys
dotenvify -azure -group "your-variable-group" -no-lower
# or using shorthand flags
dotenvify -az -g "your-variable-group" -nl

# Save to a specific file instead of .env
dotenvify -azure -group "your-variable-group" -output "custom-output.env"
# or using shorthand flags
dotenvify -az -g "your-variable-group" -out "custom-output.env"
```

All commands above will save the variables to `.env` by default, overwriting any existing file.

#### 🔒 Security First

DotEnvify prioritizes security:

- Uses your existing Azure CLI authentication - no credentials stored or handled
- Tokens are used only in memory and never written to disk
- Respects your organization's security policies (including MFA)
- No sensitive information is ever logged or exposed

Just make sure you're logged in with `az login` before running the tool.

#### 🔍 Finding Your Azure DevOps Project URL

To get your Azure DevOps project URL:

1. Go to [https://dev.azure.com/](https://dev.azure.com/)
2. Sign in with your Azure credentials
3. Select your organization from the list
4. Select your project
5. Copy the URL from your browser's address bar, which should look like:
   `https://dev.azure.com/your-org/your-project`

#### 🧠 Smart Defaults

DotEnvify tries to minimize the parameters you need to provide:

- Project URL: Uses `AZURE_DEVOPS_URL` environment variable if set (recommended approach)
- Organization name: Inferred from URL
- Project name: Inferred from URL
- Variable group name: Is **required** - will prompt if not provided
- Output file: Always defaults to `.env` in the current directory

**Environment Variable Precedence:**
1. Command-line arguments (highest precedence)
2. `AZURE_DEVOPS_URL` environment variable
3. Interactive prompt (if needed)

DotEnvify will offer to save your URL to the AZURE_DEVOPS_URL environment variable when entered interactively, so you won't need to type it again.

This is perfect for developers who need to run microservices locally with the same environment variables used in Azure DevOps pipelines!

## ✨ Features That Keep Me Sane

- 🦥 Lazy-friendly: minimal typing required
- ⚡ Fast AF (written in Go because patience isn't my virtue)
- 🧹 Skips empty lines (because whitespace is only scary in Python)
- 🛡️ Won't wreck your original file if something goes wrong
- 👻 No dependencies because who has time for npm install
- 🚀 Blazing fast execution for those tight deadlines
- 🔌 Plugin architecture: Core functionality is platform-agnostic with plugins for different integrations
- 🔄 Azure DevOps integration to fetch variables directly from variable groups
- 🔒 Secure: Uses your existing Azure CLI authentication - no credentials stored or handled by dotenvify
- 🧠 Smart defaults: Minimizes required parameters with environment variables and sensible defaults
- 🎯 Simple: Just give it a variable group name, and it figures out the rest
- 💬 Interactive: Will ask for input if needed rather than just failing
- 🎨 Beautiful terminal output with fun emojis because we're all nerds here
- 📝 Always defaults to `.env` output file, overwriting any existing file
- 🔮 Extensible: More integrations coming soon!
- 🔤 Case-sensitive: Option to ignore variables with lowercase keys using -no-lower flag
- 🔍 Short flags: All flags have short alternatives for even less typing

## 📝 Format It Understands

Your input should be in this format:
```
KEY1
value1
KEY2
value2
```

Not like this (that's a whole different tool):
```
KEY1 value1
KEY2 value2
```

## 🤓 Why I Made This

Because I value my remaining sanity points and figured you might too. Those precious minutes you spend formatting env vars could be spent on:
- Actually coding something cool
- Optimizing your coffee brewing technique
- Staring at the wall contemplating your life choices

One of these is clearly better than manual formatting. Probably.

## 🔧🛠️ Contributing

Found a bug? 🐛 Have a feature idea? 💡 PRs welcome! Just don't mess with my tabs vs. spaces setup—that debate ended my last relationship. 💔 Let's make this tool even more awesome together! 🚀

- 🕵️‍♂️ **Bug Hunters**: If you spot something weird, open an issue faster than you can say "it works on my machine"
- 🧪 **Feature Wizards**: Got an idea? PRs are the ultimate "scratch your own itch" spell
- 📚 **Documentation Heroes**: Fixed a typo? You're saving developers from Stack Overflow shame
- 🧙‍♀️ **Code Reviewers**: Your nitpicks are actually appreciated (just don't tell anyone I said that) ✨

Remember: every contribution puts you one commit closer to being the person your rubber duck thinks you are. 🦆✨

### 🏷️ Versioning and Releases

DotEnvify follows [Semantic Versioning](https://semver.org/) (SemVer):
- Version format: `vMAJOR.MINOR.PATCH` (e.g., `v0.1.3`)
- MAJOR: Breaking changes
- MINOR: New features, no breaking changes
- PATCH: Bug fixes, no breaking changes

⚠️ **Important**: When creating release tags, strictly follow the `vX.Y.Z` format. Tags like `v0.1.3.1` are not valid semantic versions and will cause the release process to fail.

## 📄⚖️ License

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![FOSSA](https://img.shields.io/badge/FOSSA-Approved-success)

### What the MIT License Means for You

DotEnvify is released under the MIT License, which is one of the most permissive and user-friendly open source licenses out there. Here's what it means in plain English:

#### 🆓 You Can:
- **Use it however you want**: Personal projects, commercial products, secret government operations (we won't ask)
- **Modify it**: Change the code, add features, remove the emojis (but why would you?)
- **Distribute it**: Share it with friends, colleagues, or include it in your own projects
- **Sell it**: Incorporate it into commercial products without paying royalties

#### ⚠️ The Only Requirements:
- Keep the original copyright notice and MIT license text in your copy/modification
- That's literally it

#### 🛡️ What You Get:
- **No Warranty**: The software is provided "as is" without any guarantees
- **Limited Liability**: The author isn't responsible if DotEnvify accidentally formats your hard drive (it won't, but legally speaking, we're covered)

#### 💼 In Business Terms:
The MIT License is business-friendly and perfect for both open source projects and commercial use. It doesn't require you to open-source your own code, even if you modify DotEnvify.

MIT License - Go wild, make millions, just don't blame me when it formats your grocery list. 🛒📝

This project is about as restrictive as a cat's opinion of your personal space—technically there are boundaries, but nobody seems to care. 🐱

---

<div align="center">

**Made with 💻, ☕, and the burning desire to never format env vars manually again**

![Works on My Machine](https://forthebadge.com/images/badges/works-on-my-machine.svg)
![Built with Swag](https://forthebadge.com/images/badges/built-with-swag.svg)

*"Life's too short for manual formatting."* 🧙‍♂️✨

</div>
