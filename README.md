# DotEnvify

Convert messy environment variables into clean, standardized `.env` files.

## Packages

| Package                                 | Description                                       | Status      |
|-----------------------------------------|---------------------------------------------------|-------------|
| [`@dotenvify/core`](./packages/core)    | Shared TypeScript library — parser, formatter, IO | In progress |
| [`dotenvify` CLI](./cli)                | Command-line tool (replaces Go version)           | In progress |
| [JetBrains Plugin](./plugins/jetbrains) | IntelliJ/WebStorm plugin                          | Functional  |
| [VS Code Extension](./plugins/vscode)   | VS Code extension                                 | Planned     |
| [Landing Page](./landing)               | Unified product site                              | In progress |

## Getting Started

```bash
# Install dependencies
npm install

# Build all packages
npm run build

# Run tests
npm run test

# Dev mode (landing page)
npm run dev:landing
```

## Structure

```
dotenvify/
├── packages/core/       # @dotenvify/core — shared TS library
├── cli/                 # CLI tool
├── plugins/
│   ├── jetbrains/       # Kotlin — Gradle build (independent)
│   └── vscode/          # VS Code extension
├── landing/             # Unified landing page
└── docs/                # Shared docs & assets
```

## License

[MIT](./LICENSE)
