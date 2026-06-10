import type {ReactNode} from "react";
import Header from "@/components/Header";
import Footer from "@/components/Footer";

export interface LegalSection {
    title: string;
    body: ReactNode;
}

interface LegalPageProps {
    title: string;
    lastUpdated: string;
    sections: LegalSection[];
}

const LegalPage = ({title, lastUpdated, sections}: LegalPageProps) => (
    <div className="flex flex-col min-h-[100dvh]">
        <Header/>
        <main className="flex-1 px-4 lg:px-6 py-12">
            <div className="max-w-2xl mx-auto">
                <p className="font-mono text-xs text-muted-foreground mb-2">{"// "}last updated: {lastUpdated}</p>
                <h1 className="font-display text-2xl font-bold text-foreground mb-8">{title}</h1>
                <div className="space-y-6">
                    {sections.map((s) => (
                        <section key={s.title}>
                            <h2 className="font-display text-sm font-semibold text-foreground mb-1.5">{s.title}</h2>
                            <p className="text-sm text-muted-foreground leading-relaxed">{s.body}</p>
                        </section>
                    ))}
                </div>
            </div>
        </main>
        <Footer/>
    </div>
);

export default LegalPage;
