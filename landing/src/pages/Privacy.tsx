import Header from "@/components/Header";
import Footer from "@/components/Footer";

const Privacy = () => (
  <div className="flex flex-col min-h-[100dvh]">
    <Header />

    <main className="flex-1 pt-12">
      <div className="max-w-2xl mx-auto px-4 lg:px-6 py-12">
        <h1 className="text-3xl font-bold mb-2">Privacy Policy</h1>
        <p className="text-sm text-muted-foreground mb-10">Last updated: March 9, 2026</p>

        <div className="space-y-8">
          <section>
            <h2 className="text-lg font-semibold mb-2">1. Overview</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              DotEnvify is built by Webbies. This Privacy Policy explains how we handle
              information when you use our CLI tool, IDE plugins, and website.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">2. Information We Collect</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              <strong className="text-foreground">We do not collect personal data.</strong> DotEnvify processes environment
              variables entirely on your local machine (CLI and IDE plugins) or in your browser
              (website). No data is sent to our servers.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">3. Azure DevOps Integration</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              When you use the Azure DevOps integration, the plugin communicates directly with
              Microsoft's Azure DevOps APIs using your OAuth credentials. We do not proxy,
              intercept, or store your credentials, access tokens, or variable group data.
              Authentication is handled entirely between your device and Microsoft.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">4. Local Data</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              The CLI and IDE plugins read and write .env files on your local filesystem.
              These files may contain sensitive information such as API keys, database credentials,
              and other secrets. You are responsible for securing these files. DotEnvify writes
              output files with restrictive permissions (owner read/write only).
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">5. Website</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              The DotEnvify website is a static site. We do not use cookies, analytics,
              or tracking scripts. Any environment variable conversion done on the website
              happens entirely in your browser — no data leaves your device.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">6. Third-Party Services</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              The website is hosted on Vercel. Vercel may collect standard web server logs
              (IP addresses, browser type, access times). Refer to{" "}
              <a
                href="https://vercel.com/legal/privacy-policy"
                target="_blank"
                rel="noopener noreferrer"
                className="text-foreground underline hover:opacity-80 transition-opacity"
              >
                Vercel's Privacy Policy
              </a>{" "}
              for details.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">7. Children's Privacy</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              The Service is not directed at children under 13. We do not knowingly collect
              information from children.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">8. Changes to This Policy</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              We may update this Privacy Policy from time to time. Changes will be reflected
              on this page with an updated date.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">9. Contact</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              For questions about this Privacy Policy, contact us at{" "}
              <a href="mailto:dotenvify@webbies.dev" className="text-foreground underline hover:opacity-80 transition-opacity">
                dotenvify@webbies.dev
              </a>.
            </p>
          </section>
        </div>
      </div>
    </main>

    <Footer />
  </div>
);

export default Privacy;
