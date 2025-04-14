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

**The lazy way** (we don't judge, we encourage):
```bash
curl -s https://raw.githubusercontent.com/webb1es/dotenvify/main/remote-install.sh | bash
```

**The "I read every line of code before running it" way**:
```bash
git clone https://github.com/webb1es/dotenvify.git
cd dotenvify
go build -o dotenvify dotenvify.go
sudo mv dotenvify /usr/local/bin/
```

## 🔮 How To Use It

Seriously, it's two commands:

```bash
# Overwrite the same file (living dangerously)
dotenvify your-vars.txt

# Save to a new file (your therapist would approve)
dotenvify your-vars.txt .env
```

## ✨ Features That Keep Me Sane

- 🦥 Lazy-friendly: minimal typing required
- ⚡ Fast AF (written in Go because patience isn't my virtue)
- 🧹 Skips empty lines (because whitespace is only scary in Python)
- 🛡️ Won't wreck your original file if something goes wrong
- 👻 No dependencies because who has time for npm install

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

Found a bug? 🐛 Have a feature idea? 💡 PRs welcome! Just don't mess with my tabs vs. spaces setup—that debate ended my last relationship. 💔

- 🕵️‍♂️ **Bug Hunters**: If you spot something weird, open an issue faster than you can say "it works on my machine"
- 🧪 **Feature Wizards**: Got an idea? PRs are the ultimate "scratch your own itch" spell
- 📚 **Documentation Heroes**: Fixed a typo? You're saving developers from Stack Overflow shame
- 🧙‍♀️ **Code Reviewers**: Your nitpicks are actually appreciated (just don't tell anyone I said that)

Remember: every contribution puts you one commit closer to being the person your rubber duck thinks you are. 🦆✨

## 📄⚖️ License

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwebb1es%2Fdotenvify.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwebb1es%2Fdotenvify)

MIT License - Go wild, make millions, just don't blame me when it formats your grocery list. 🛒📝

This project is about as restrictive as a cat's opinion of your personal space—technically there are boundaries, but nobody seems to care. 🐱

---

<div align="center">

**Made with 💻, ☕, and the burning desire to never format env vars manually again**

![Works on My Machine](https://forthebadge.com/images/badges/works-on-my-machine.svg)
![Built with Swag](https://forthebadge.com/images/badges/built-with-swag.svg)

*"Life's too short for manual formatting."* 🧙‍♂️✨

</div>