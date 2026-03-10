import Header from "@/components/Header";
import Footer from "@/components/Footer";

const sections = [
    {
        title: "1. acceptance",
        body: 'By using DotEnvify (the "Service"), including the CLI tool, IDE plugins, and website, you agree to these Terms of Service. If you do not agree, do not use the Service.'
    },
    {
        title: "2. what this is",
        body: "DotEnvify is an open-source tool that converts environment variables into standardized .env files. It is available as a CLI tool, JetBrains plugin, and web application. The Service may integrate with third-party platforms such as Azure DevOps."
    },
    {
        title: "3. your responsibilities",
        body: "You may use the Service for any lawful purpose. You are responsible for ensuring that your use complies with all applicable laws and regulations. You must not use the Service to process or store data in violation of any applicable data protection laws."
    },
    {
        title: "4. azure devops",
        body: "The Service may access Azure DevOps on your behalf to retrieve variable groups. You authorize this access when you authenticate through the OAuth flow. You are responsible for the security of your credentials and access tokens. The Service does not store your Azure DevOps credentials."
    },
    {
        title: "5. license",
        body: "DotEnvify is released under the MIT License. You may use, copy, modify, and distribute the software in accordance with the license terms."
    },
    {
        title: "6. no warranties",
        body: 'The Service is provided "as is" without warranty of any kind, express or implied. We do not guarantee that the Service will be error-free, secure, or available at all times.'
    },
    {
        title: "7. liability",
        body: "To the fullest extent permitted by law, Webbies shall not be liable for any indirect, incidental, special, or consequential damages arising from your use of the Service, including but not limited to loss of data or unauthorized access to your environment variables."
    },
    {
        title: "8. changes",
        body: "We may update these Terms from time to time. Continued use of the Service after changes constitutes acceptance of the updated Terms."
    },
];

const Terms = () => (
    <div className="flex flex-col min-h-[100dvh]">
        <Header/>
        <main className="flex-1 px-4 lg:px-6 py-12">
            <div className="max-w-2xl mx-auto">
                <p className="font-mono text-xs text-muted-foreground mb-2">{"// "}last updated: march 9, 2026</p>
                <h1 className="font-display text-2xl font-bold text-foreground mb-8">terms of service</h1>
                <div className="space-y-6">
                    {sections.map((s) => (
                        <section key={s.title}>
                            <h2 className="font-display text-sm font-semibold text-foreground mb-1.5">{s.title}</h2>
                            <p className="text-sm text-muted-foreground leading-relaxed">{s.body}</p>
                        </section>
                    ))}
                    <section>
                        <h2 className="font-display text-sm font-semibold text-foreground mb-1.5">9. contact</h2>
                        <p className="text-sm text-muted-foreground leading-relaxed">
                            Questions? Email{" "}
                            <a href="mailto:dotenvify@webbies.dev" className="text-primary hover:underline">
                                dotenvify@webbies.dev
                            </a>
                        </p>
                    </section>
                </div>
            </div>
        </main>
        <Footer/>
    </div>
);

export default Terms;
