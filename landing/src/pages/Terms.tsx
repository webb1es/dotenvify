import LegalPage, {type LegalSection} from "@/components/LegalPage";

const sections: LegalSection[] = [
    {
        title: "1. acceptance",
        body: 'By using dotenvify (the "Service"), including the CLI tool, the JetBrains plugin, and the website, you agree to these Terms of Service. If you do not agree, do not use the Service.'
    },
    {
        title: "2. what this is",
        body: "dotenvify is an open-source tool that converts environment variables into standardized .env files. It is available as a CLI tool and a JetBrains plugin, with an interactive demo on the website. The Service may integrate with third-party platforms such as Azure DevOps."
    },
    {
        title: "3. your responsibilities",
        body: "You may use the Service for any lawful purpose. You are responsible for ensuring that your use complies with all applicable laws and regulations. You must not use the Service to process or store data in violation of any applicable data protection laws."
    },
    {
        title: "4. azure devops",
        body: "The Service may access Azure DevOps on your behalf to retrieve variable groups. Authentication uses your local Azure CLI (`az`) session: you sign in with `az login`, and the plugin reuses that session to request a short-lived access token. The Azure CLI owns credential storage; the Service does not store your Azure DevOps credentials or tokens. You are responsible for the security of your local CLI session."
    },
    {
        title: "5. license",
        body: "dotenvify is released under the MIT License. You may use, copy, modify, and distribute the software in accordance with the license terms."
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
    {
        title: "9. contact",
        body: (
            <>
                Questions? Email{" "}
                <a href="mailto:dotenvify@webbies.dev" className="text-primary hover:underline">
                    dotenvify@webbies.dev
                </a>
            </>
        )
    },
];

const Terms = () => <LegalPage title="terms of service" lastUpdated="june 19, 2026" sections={sections}/>;

export default Terms;
