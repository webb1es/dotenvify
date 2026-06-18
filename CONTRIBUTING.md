# Contributing to DotEnvify

This is a monorepo (managed with [Turbo](https://turbo.build/)) holding the CLI,
the JetBrains plugin, the shared core library, and the landing page.

## Repository layout

| Path                | What it is                                              |
|---------------------|---------------------------------------------------------|
| `packages/core`     | `@dotenvify/core` — the parser, formatter, and file I/O |
| `cli`               | `@webbies.dev/dotenvify` — the command-line tool        |
| `plugins/jetbrains` | The JetBrains IDE plugin (Kotlin / Gradle)              |
| `landing`           | The product website                                     |

## Building and testing

```bash
npm install     # install dependencies
npm run build   # build all packages
npm run test    # run all tests
```

The JetBrains plugin has its own Gradle build and needs **JDK 17** (pinned via
`gradle/gradle-daemon-jvm.properties`):

```bash
cd plugins/jetbrains
./gradlew test buildPlugin   # build the .zip in build/distributions/
./gradlew runIde             # launch a sandbox IDE to try it
./gradlew verifyPlugin       # check IDE compatibility before publishing
```

## Releases

Each part of the project ships on its own, using a tag with its own prefix. The
landing page deploys automatically from `main`.

### CLI → npm

Push a `cli-v*` tag. [`release-cli.yml`](.github/workflows/release-cli.yml) builds,
tests, sets the version from the tag, and publishes to npm.

```bash
git tag -a cli-v<version> -m "CLI v<version>" && git push origin cli-v<version>
```

### JetBrains plugin → Marketplace

**The first release must be done by hand.** The Marketplace doesn't allow
automated publishing for a brand-new plugin. Build the `.zip`
(`./gradlew buildPlugin`), upload it at
<https://plugins.jetbrains.com/plugin/add>, and wait for it to be approved.

After that, push a `jetbrains-v*` tag and
[`release-jetbrains.yml`](.github/workflows/release-jetbrains.yml) builds, signs,
and publishes the update.

```bash
git tag -a jetbrains-v<version> -m "JetBrains v<version>" && git push origin jetbrains-v<version>
```

It needs these repository secrets:

- `JETBRAINS_MARKETPLACE_TOKEN`
- `JETBRAINS_CERTIFICATE_CHAIN` — contents of `chain.crt`
- `JETBRAINS_PRIVATE_KEY` — contents of `private.pem`
- `JETBRAINS_PRIVATE_KEY_PASSWORD`

The signing keys (`private.pem`, `chain.crt`) are git-ignored. Keep them safe.
