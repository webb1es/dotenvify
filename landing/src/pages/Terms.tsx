import Header from "@/components/Header";
import Footer from "@/components/Footer";

const Terms = () => (
  <div className="flex flex-col min-h-[100dvh]">
    <Header />

    <main className="flex-1 pt-12">
      <div className="max-w-2xl mx-auto px-4 lg:px-6 py-12">
        <h1 className="text-3xl font-bold mb-2">Terms of Service</h1>
        <p className="text-sm text-muted-foreground mb-10">Last updated: March 9, 2026</p>

        <div className="space-y-8">
          <section>
            <h2 className="text-lg font-semibold mb-2">1. Acceptance of Terms</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              By using DotEnvify (the "Service"), including the CLI tool, IDE plugins, and website,
              you agree to these Terms of Service. If you do not agree, do not use the Service.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">2. Description of Service</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              DotEnvify is an open-source tool that converts environment variables into standardized
              .env files. It is available as a CLI tool, JetBrains plugin, and web application.
              The Service may integrate with third-party platforms such as Azure DevOps.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">3. Use of the Service</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              You may use the Service for any lawful purpose. You are responsible for ensuring that
              your use complies with all applicable laws and regulations. You must not use the Service
              to process or store data in violation of any applicable data protection laws.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">4. Azure DevOps Integration</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              The Service may access Azure DevOps on your behalf to retrieve variable groups.
              You authorize this access when you authenticate through the OAuth flow. You are
              responsible for the security of your credentials and access tokens. The Service
              does not store your Azure DevOps credentials.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">5. Intellectual Property</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              DotEnvify is released under the MIT License. You may use, copy, modify, and distribute
              the software in accordance with the license terms.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">6. Disclaimer of Warranties</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              The Service is provided "as is" without warranty of any kind, express or implied.
              We do not guarantee that the Service will be error-free, secure, or available at all times.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">7. Limitation of Liability</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              To the fullest extent permitted by law, Webbies shall not be liable for any indirect,
              incidental, special, or consequential damages arising from your use of the Service,
              including but not limited to loss of data or unauthorized access to your environment variables.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">8. Changes to Terms</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              We may update these Terms from time to time. Continued use of the Service after
              changes constitutes acceptance of the updated Terms.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold mb-2">9. Contact</h2>
            <p className="text-sm text-muted-foreground leading-relaxed">
              For questions about these Terms, contact us at{" "}
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

export default Terms;
