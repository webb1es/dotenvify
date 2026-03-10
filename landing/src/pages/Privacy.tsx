import Header from "@/components/Header";
import Footer from "@/components/Footer";

const sections = [
    {
        title: "1. overview",
        body: "DotEnvify is built by Webbies. This Privacy Policy explains how we handle information when you use our CLI tool, IDE plugins, and website."
    },
    {
        title: "2. data collection",
        body: "We do not collect personal data. DotEnvify processes environment variables entirely on your local machine (CLI and IDE plugins) or in your browser (website). No data is sent to our servers."
    },
    {
        title: "3. azure devops",
        body: "When you use the Azure DevOps integration, the plugin communicates directly with Microsoft's Azure DevOps APIs using your OAuth credentials. We do not proxy, intercept, or store your credentials, access tokens, or variable group data. Authentication is handled entirely between your device and Microsoft."
    },
    {
        title: "4. local data",
        body: "The CLI and IDE plugins read and write .env files on your local filesystem. These files may contain sensitive information such as API keys, database credentials, and other secrets. You are responsible for securing these files. DotEnvify writes output files with restrictive permissions (owner read/write only)."
    },
    {
        title: "5. website",
        body: "The DotEnvify website is a static site. We do not use cookies, analytics, or tracking scripts. Any environment variable conversion done on the website happens entirely in your browser. No data leaves your device."
    },
    {
        title: "6. third parties",
        body: "The website is hosted on Vercel. Vercel may collect standard web server logs (IP addresses, browser type, access times). Refer to Vercel's Privacy Policy for details."
    },
    {
        title: "7. children",
        body: "The Service is not directed at children under 13. We do not knowingly collect information from children."
    },
    {
        title: "8. changes",
        body: "We may update this Privacy Policy from time to time. Changes will be reflected on this page with an updated date."
    },
];

const Privacy = () => (
    <div className="flex flex-col min-h-[100dvh]">
        <Header/>
        <main className="flex-1 px-4 lg:px-6 py-12">
            <div className="max-w-2xl mx-auto">
                <p className="font-mono text-xs text-muted-foreground mb-2">{"// "}last updated: march 9, 2026</p>
                <h1 className="font-display text-2xl font-bold text-foreground mb-8">privacy policy</h1>
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

export default Privacy;
