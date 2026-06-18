import {useState} from "react";
import {Check, Copy, Download, Star} from "lucide-react";

const REPO = "webb1es/dotenvify";
const INSTALL_CMD = "npx @webbies.dev/dotenvify input.txt -o .env";

const ClosingCTA = () => {
    const [copied, setCopied] = useState(false);

    const handleCopy = async () => {
        await navigator.clipboard.writeText(INSTALL_CMD);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <div className="bento-cell-static p-6 lg:p-10 flex flex-col items-center text-center gap-4">
            <h2
                className="font-display font-bold gradient-text leading-tight tracking-tight"
                style={{fontSize: "clamp(20px, 2.5vw, 32px)"}}
            >
                ready when you are
            </h2>
            <p className="text-sm text-muted-foreground leading-relaxed max-w-md">
                One command for the CLI, one click for the IDE. Runs on your machine — no accounts, no servers.
            </p>

            <button
                onClick={handleCopy}
                className="w-full max-w-md flex items-center gap-2 px-2.5 py-1.5 rounded-md bg-de-surface border border-de-surface-border font-mono text-[11px] text-left hover:border-primary/30 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring group"
                aria-label={copied ? "Copied" : `Copy: ${INSTALL_CMD}`}
            >
                <span className="text-muted-foreground select-none" aria-hidden="true">$</span>
                <code className="flex-1 text-foreground truncate">{INSTALL_CMD}</code>
                {copied ? (
                    <Check className="w-3 h-3 text-de-green flex-shrink-0"/>
                ) : (
                    <Copy className="w-3 h-3 text-muted-foreground group-hover:text-foreground flex-shrink-0"/>
                )}
            </button>

            <div className="flex flex-wrap gap-2 justify-center">
                <a
                    href="https://plugins.jetbrains.com/plugin/32351-dotenvify"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1.5 h-9 px-5 rounded-lg bg-primary text-primary-foreground font-mono text-xs font-semibold hover:opacity-90 transition-opacity focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                >
                    <Download className="w-3.5 h-3.5"/>
                    get plugin
                </a>
                <a
                    href={`https://github.com/${REPO}/stargazers`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1.5 h-9 px-5 rounded-lg border border-border text-foreground font-mono text-xs font-semibold hover:bg-accent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                >
                    <Star className="w-3.5 h-3.5"/>
                    star on github
                </a>
            </div>
        </div>
    );
};

export default ClosingCTA;
