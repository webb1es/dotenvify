# JetBrains Marketplace Publishing

## Prerequisites (done)

- [x] Plugin icon (`pluginIcon.svg`)
- [x] Plugin descriptor (`plugin.xml`) with ID, description, vendor info, change-notes
- [x] Signing and publishing config added to `build.gradle.kts`
- [x] Generate signing certificate (run from `plugins/jetbrains/`):
  ```bash
  # Generate private key
  openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:4096 -out private.pem

  # Generate certificate chain (self-signed, valid 10 years)
  openssl req -new -x509 -key private.pem -out chain.crt -days 3650 -subj "/CN=Webbies/O=Webbies/L=Harare/ST=Harare/C=ZW"
  ```
- [x] `private.pem` and `chain.crt` added to `.gitignore`

## Steps

- [ ] Create account at [plugins.jetbrains.com](https://plugins.jetbrains.com)
- [ ] Generate a Marketplace token at [plugins.jetbrains.com/author/me/tokens](https://plugins.jetbrains.com/author/me/tokens)
- [ ] Publish:
  ```bash
  cd plugins/jetbrains

  export JETBRAINS_MARKETPLACE_TOKEN="your-token-here"
  export PRIVATE_KEY_PASSWORD=""  # skip if no password on private key

  ./gradlew publishPlugin
  ```
- [ ] Wait for JetBrains review (1-2 business days)
