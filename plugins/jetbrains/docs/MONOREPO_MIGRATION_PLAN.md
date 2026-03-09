# DotEnvify Monorepo Migration Plan

> Consolidate 3 existing repos + 2 planned projects into a single monorepo.
> Repository: https://github.com/webb1es/dotenvify

---

## Current State

| Repo | Stack | Lines (src) | Status |
|---|---|---|---|
| `dotenvify` (now `dotenvify-legacy`) | Go | ~680 | Published (npm + GitHub releases) |
| `dotenvify-plugin-jetbrains` | Kotlin | ~2,286 | Functional, unpublished |
| `dotenvify-plugin-landing` | React/Vite/TS | ~380 custom | Built, undeployed |
| *(planned)* VS Code extension | TypeScript | — | Not started |
| *(planned)* CLI landing page | TypeScript | — | Not started (now merged into unified site) |

## Target State

**2 stacks** (TypeScript + Kotlin), **1 repo**, **1 landing page**, **shared core library**.

---

## Monorepo Structure

```
dotenvify/
├── packages/
│   └── core/                        # @dotenvify/core — shared TS library
│       ├── src/
│       │   ├── parser.ts            # Multi-format parser (port from Go)
│       │   ├── formatter.ts         # Output formatter + filters
│       │   ├── io.ts                # File read/write/backup/preserve
│       │   ├── models.ts            # Types: EnvEntry, FormatOptions, ParseResult
│       │   └── index.ts             # Public API exports
│       ├── tests/
│       │   ├── parser.test.ts
│       │   ├── formatter.test.ts
│       │   └── io.test.ts
│       ├── package.json             # name: "@dotenvify/core"
│       └── tsconfig.json
│
├── cli/                             # dotenvify CLI (replaces Go version)
│   ├── src/
│   │   ├── index.ts                 # Entry point, arg parsing
│   │   ├── commands/
│   │   │   ├── convert.ts           # Default command — parse + format + write
│   │   │   └── azure.ts             # Azure DevOps fetch command
│   │   └── utils/
│   │       └── azure-auth.ts        # Device code flow for CLI
│   ├── bin/
│   │   └── dotenvify.js             # Shebang entry (#!/usr/bin/env node)
│   ├── package.json                 # name: "dotenvify", bin entry
│   ├── tsconfig.json
│   └── README.md
│
├── plugins/
│   ├── jetbrains/                   # Kotlin — moved as-is
│   │   ├── src/
│   │   │   ├── main/kotlin/dev/webbies/dotenvify/
│   │   │   │   ├── core/            # Kotlin parser/formatter (own implementation)
│   │   │   │   ├── actions/
│   │   │   │   ├── ui/
│   │   │   │   ├── settings/
│   │   │   │   ├── azure/
│   │   │   │   └── diagnostics/
│   │   │   └── test/
│   │   ├── build.gradle.kts
│   │   ├── gradle.properties
│   │   ├── gradle/
│   │   ├── .run/                    # IntelliJ run configs
│   │   └── README.md
│   │
│   └── vscode/                      # VS Code extension (new)
│       ├── src/
│       │   ├── extension.ts         # Activation, command registration
│       │   ├── commands/
│       │   │   ├── convertSelection.ts
│       │   │   ├── convertFile.ts
│       │   │   └── fetchAzure.ts
│       │   ├── providers/
│       │   │   ├── diagnosticsProvider.ts
│       │   │   └── codeActionProvider.ts
│       │   ├── views/
│       │   │   └── sidebarPanel.ts  # Webview panel (paste & format)
│       │   └── utils/
│       │       └── azure-auth.ts    # Device code flow for VS Code
│       ├── package.json             # VS Code extension manifest
│       ├── tsconfig.json
│       └── README.md
│
├── site/                            # Unified landing page
│   ├── src/
│   │   ├── main.tsx
│   │   ├── App.tsx
│   │   ├── pages/
│   │   │   └── Index.tsx
│   │   ├── components/
│   │   │   ├── Header.tsx
│   │   │   ├── HeroSection.tsx
│   │   │   ├── FeatureCard.tsx
│   │   │   ├── CodeTransform.tsx
│   │   │   ├── LiveDemo.tsx         # Interactive demo using @dotenvify/core
│   │   │   ├── ProductTabs.tsx      # CLI / JetBrains / VS Code tabs
│   │   │   ├── IdeLogos.tsx
│   │   │   ├── ThemeToggle.tsx
│   │   │   ├── Footer.tsx
│   │   │   └── ui/                  # shadcn components
│   │   ├── hooks/
│   │   └── lib/
│   ├── public/
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   └── tsconfig.json
│
├── docs/                            # Shared documentation & assets
│   ├── ROADMAP.md                   # Unified roadmap (all products)
│   ├── TODO.md                      # Unified task tracker
│   ├── MIGRATION.md                 # This document (remove after migration)
│   ├── LANDING_PAGE_PROMPT.md       # Landing page design spec
│   └── assets/
│       ├── pluginIcon.svg           # Shared brand icon
│       └── screenshots/             # Plugin/CLI screenshots
│
├── .github/
│   └── workflows/
│       ├── core.yml                 # Test core on packages/core/** changes
│       ├── cli.yml                  # Build + test + publish CLI to npm
│       ├── jetbrains.yml            # Gradle build + test on plugins/jetbrains/**
│       ├── vscode.yml               # Build + test + publish VS Code extension
│       └── site.yml                 # Build + deploy landing page
│
├── package.json                     # Workspace root
├── turbo.json                       # Turborepo task runner (optional)
├── .gitignore
├── LICENSE                          # MIT
└── README.md                        # Project overview, links to sub-READMEs
```

---

## Code Sharing Map

```
@dotenvify/core (TypeScript)
    ├── cli/              → imports parser, formatter, IO directly
    ├── plugins/vscode/   → imports parser, formatter, models
    └── site/             → imports parser, formatter for LiveDemo component

plugins/jetbrains/ (Kotlin)
    └── own implementation in core/ package (cannot share with TS)
```

**What's shared:**
- Parser logic (parse multi-format input → EnvEntry[])
- Formatter logic (EnvEntry[] → formatted .env string)
- IO logic (read/write files, backups, preserve keys)
- Models/types (EnvEntry, FormatOptions, ParseResult)

**What's NOT shared:**
- JetBrains plugin (Kotlin/JVM — maintains its own core implementation)
- Azure auth (different flows for CLI vs IDE vs VS Code)
- UI code (each platform has its own UI layer)

---

## Workspace Configuration

### Root package.json

```json
{
  "name": "dotenvify",
  "private": true,
  "workspaces": [
    "packages/*",
    "cli",
    "plugins/vscode",
    "site"
  ],
  "scripts": {
    "build": "turbo run build",
    "test": "turbo run test",
    "lint": "turbo run lint",
    "dev:site": "npm run dev --workspace=site",
    "dev:cli": "npm run dev --workspace=cli"
  },
  "devDependencies": {
    "turbo": "^2.0.0",
    "typescript": "^5.8.0"
  }
}
```

### turbo.json

```json
{
  "$schema": "https://turbo.build/schema.json",
  "tasks": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": ["dist/**"]
    },
    "test": {
      "dependsOn": ["^build"]
    },
    "lint": {},
    "dev": {
      "cache": false,
      "persistent": true
    }
  }
}
```

**Note:** The JetBrains plugin (Gradle) lives outside the npm workspace. Its CI workflow runs Gradle independently. Turborepo only orchestrates the TypeScript projects.

---

## CI/CD Workflows

### Path-Filtered Triggers

Each workflow only runs when its directory changes:

```yaml
# .github/workflows/core.yml
on:
  push:
    paths: ['packages/core/**']
  pull_request:
    paths: ['packages/core/**']

# .github/workflows/cli.yml
on:
  push:
    paths: ['cli/**', 'packages/core/**']
    tags: ['cli-v*.*.*']

# .github/workflows/jetbrains.yml
on:
  push:
    paths: ['plugins/jetbrains/**']
  pull_request:
    paths: ['plugins/jetbrains/**']

# .github/workflows/vscode.yml
on:
  push:
    paths: ['plugins/vscode/**', 'packages/core/**']
    tags: ['vscode-v*.*.*']

# .github/workflows/site.yml
on:
  push:
    branches: [main]
    paths: ['site/**', 'packages/core/**']
```

### Versioning Strategy

Each publishable package has its own version, prefixed tags:
- `cli-v1.0.0` → publishes `dotenvify` to npm
- `vscode-v0.1.0` → publishes VS Code extension to marketplace
- JetBrains plugin → manual upload or Gradle publish task
- `@dotenvify/core` → private package, not published to npm (consumed internally via workspace)

---

## Migration Steps

### Phase 1: Scaffold Monorepo

1. Clone the empty `webb1es/dotenvify` repo locally
2. Initialize npm workspace root (`package.json`, `turbo.json`)
3. Create directory structure (`packages/`, `cli/`, `plugins/`, `site/`, `docs/`)
4. Add root `.gitignore` (node_modules, dist, build, .gradle, .idea)
5. Add `LICENSE` (MIT)
6. Add root `README.md`
7. Commit: "chore: scaffold monorepo structure"

### Phase 2: Move JetBrains Plugin

1. Copy `dotenvify-plugin-jetbrains/` contents → `plugins/jetbrains/`
2. Exclude: `.git/`, `.idea/` (project-level), any build artifacts
3. Include: `src/`, `build.gradle.kts`, `settings.gradle.kts`, `gradle.properties`, `gradle/`, `.run/`, `README.md`
4. Verify build: `cd plugins/jetbrains && ./gradlew buildPlugin`
5. Verify tests: `./gradlew test`
6. Commit: "feat: move JetBrains plugin into monorepo"

### Phase 3: Move Landing Page

1. Copy `dotenvify-plugin-landing/` contents → `site/`
2. Exclude: `.git/`, `node_modules/`, `dist/`, `bun.lock*`
3. Include: `src/`, `public/`, config files (`vite.config.ts`, `tailwind.config.ts`, `tsconfig*.json`, `package.json`, etc.)
4. Update `package.json` name to `@dotenvify/site` (private)
5. Remove `package-lock.json` (will use root lockfile)
6. Verify: `npm install` from root, then `npm run dev --workspace=site`
7. Commit: "feat: move landing page into monorepo"

### Phase 4: Create Core Package

1. Create `packages/core/package.json`:
   ```json
   {
     "name": "@dotenvify/core",
     "version": "0.1.0",
     "private": true,
     "type": "module",
     "main": "dist/index.js",
     "types": "dist/index.d.ts",
     "scripts": {
       "build": "tsc",
       "test": "vitest run"
     }
   }
   ```
2. Port Go parser → `packages/core/src/parser.ts`
   - Source reference: `dotenvify-legacy/dotenvify.go` (lines ~50-200)
   - Kotlin reference: `plugins/jetbrains/src/main/kotlin/.../core/DotEnvParser.kt`
   - Support: KEY=VALUE, KEY="VALUE", KEY VALUE, line pairs, mixed, comments, blanks
3. Port Go formatter → `packages/core/src/formatter.ts`
   - Export prefix, sorting, smart quoting, no-lower filter, url-only filter
4. Port Go IO → `packages/core/src/io.ts`
   - Read/write .env, incremental backups, preserve keys, secure permissions
5. Define types → `packages/core/src/models.ts`
   - `EnvEntry { key: string, value: string }`
   - `FormatOptions { export, sort, ignoreLowercase, urlOnly }`
   - `ParseResult { entries: EnvEntry[], warnings: string[] }`
6. Port tests from `dotenvify-legacy/dotenvify_test.go` (690 lines)
7. Verify: all tests pass
8. Commit: "feat: create @dotenvify/core TypeScript package"

### Phase 5: Create TypeScript CLI

1. Create `cli/package.json`:
   ```json
   {
     "name": "dotenvify",
     "version": "2.0.0",
     "type": "module",
     "bin": { "dotenvify": "./bin/dotenvify.js" },
     "dependencies": {
       "@dotenvify/core": "workspace:*",
       "commander": "^13.0.0"
     }
   }
   ```
2. Implement CLI commands:
   - Default: `dotenvify [input] [options]` → parse + format + write
   - `--export`, `--no-sort`, `--no-lower`, `--url-only`, `--preserve`, `--overwrite`
   - `--output` / `-o` for output path
   - `dotenvify azure --org <url> --project <name> --group <name>` (Phase 2)
3. Maintain exact CLI flag compatibility with Go version
4. Verify: `npx dotenvify --help` works, all flag combinations tested
5. Commit: "feat: create TypeScript CLI (replaces Go version)"

### Phase 6: Move Shared Docs & Assets

1. Merge docs:
   - `dotenvify-plugin-jetbrains/docs/ROADMAP.md` → `docs/ROADMAP.md` (expand to cover all products)
   - `dotenvify-plugin-jetbrains/docs/TODO.md` → `docs/TODO.md` (expand with CLI + VS Code tasks)
   - `dotenvify-plugin-landing/docs/LANDING_PAGE_PROMPT.md` → `docs/LANDING_PAGE_PROMPT.md`
2. Move shared assets:
   - `plugins/jetbrains/src/main/resources/META-INF/pluginIcon.svg` → copy to `docs/assets/pluginIcon.svg`
3. Commit: "docs: consolidate documentation and shared assets"

### Phase 7: Update Landing Page to Unified Site

1. Add product tabs to landing page: CLI / JetBrains / VS Code
2. Add `@dotenvify/core` as dependency to `site/package.json`
3. Build `LiveDemo.tsx` component — interactive paste-and-format using core package
4. Each tab shows:
   - Install command (npm / JetBrains Marketplace / VS Code Marketplace)
   - Key features specific to that platform
   - Screenshot placeholder
5. Shared sections: Hero, code transform example, live demo, footer
6. Commit: "feat: unified landing page with product tabs and live demo"

### Phase 8: CI/CD Setup

1. Create `.github/workflows/core.yml` — test core on changes
2. Create `.github/workflows/cli.yml` — build, test, publish to npm on `cli-v*` tags
3. Create `.github/workflows/jetbrains.yml` — Gradle build + test
4. Create `.github/workflows/vscode.yml` — build, test, publish on `vscode-v*` tags
5. Create `.github/workflows/site.yml` — deploy to Vercel/Netlify on main push
6. Commit: "ci: add path-filtered workflows for all projects"

### Phase 9: Archive Old Repos

1. Update `dotenvify-legacy` README: add deprecation notice pointing to `webb1es/dotenvify`
2. Update `dotenvify-plugin-jetbrains` README: same
3. Update `dotenvify-plugin-landing` README: same
4. Archive all three repos on GitHub (Settings → Archive)
5. npm: publish `dotenvify@2.0.0` from new monorepo CLI, deprecate old Go-based versions

### Phase 10: Build VS Code Extension (future)

1. Scaffold `plugins/vscode/` with `yo code` or manual setup
2. Import `@dotenvify/core` for parser/formatter
3. Implement:
   - Convert selection command
   - Convert file command
   - Sidebar webview panel (paste & format)
   - Azure DevOps fetch command
   - Diagnostics (missing/unused keys)
4. Publish to VS Code Marketplace

---

## npm Distribution

### Before (Go)
```
npm install -g @webbies.dev/dotenvify
# → postinstall script downloads Go binary from GitHub releases
# → binary placed in node_modules/.bin/
# → platform-specific (darwin/linux/windows × amd64/arm64)
```

### After (TypeScript)
```
npm install -g dotenvify
# → pure JavaScript, runs on Node.js directly
# → no postinstall, no binary download, no platform issues
# → works everywhere Node.js works
npx dotenvify --help
# → zero-install usage
```

**Breaking change:** Package name changes from `@webbies.dev/dotenvify` to `dotenvify` (if available on npm). If not available, keep scoped name `@dotenvify/cli` or `@webbies/dotenvify`.

---

## Risk Mitigation

| Risk | Mitigation |
|---|---|
| Go CLI users broken by npm package name change | Publish deprecation notice on old package, point to new one |
| JetBrains Gradle build conflicts with npm workspace | Gradle is fully isolated — not in `workspaces` array, own build system |
| Core package breaking change affects CLI + VS Code + site | Core has its own test suite; Turborepo `dependsOn: ["^build"]` ensures downstream rebuilds |
| Large repo clone for contributors | `git clone --filter=blob:none` (partial clone) or `--sparse-checkout` for specific directories |
| Kotlin core diverges from TypeScript core | Acceptable tradeoff — JetBrains plugin is self-contained. Document parity expectations in ROADMAP.md |

---

## Success Criteria

- [ ] All existing tests pass (JetBrains Gradle tests, Go test ports → TS)
- [ ] `npx dotenvify --help` works from monorepo CLI
- [ ] `cd plugins/jetbrains && ./gradlew buildPlugin` works
- [ ] Landing page builds and serves (`npm run dev --workspace=site`)
- [ ] CI runs correct workflow per path change
- [ ] Old repos archived with redirect notices
- [ ] npm package published from new repo
