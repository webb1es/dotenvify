# 🧙‍♂️ DotEnvify

Sick of reformatting environment variables by hand? Yeah, me too.

This little Go tool converts plain text key-value pairs into proper `.env` format. I built it because I was tired of doing this manually a million times.

## What it does

Turns this:
```
API_KEY
a1b2c3d4e5f6g7h8i9j0
DEBUG
true
```

Into this:
```
export API_KEY="a1b2c3d4e5f6g7h8i9j0"
export DEBUG="true"
```

## Install

```bash
curl -s https://raw.githubusercontent.com/webb1es/dotenvify/main/remote-install.sh | bash
```

## Use it

```bash
# Basic usage (overwrites the source file)
dotenvify plain-vars.txt

# Or specify an output file
dotenvify plain-vars.txt .env
```

## Features

- ⚡ Simple and fast
- 🤦‍♂️ Handles weird formatting and errors
- 🎨 Colorful output (because why not?)
- 🔥 Works on Mac, Linux, Windows

Check the [wiki](https://github.com/webb1es/dotenvify/wiki) for more details if you're into that sort of thing.

## Requirements

- Go 1.13+

---

Made during a late-night coding session fueled by too much coffee ☕