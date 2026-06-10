import {Github, Play} from "lucide-react";

const REPO = "webb1es/dotenvify";

const ClosingCTA = () => (
    <div className="bento-cell-static p-6 lg:p-10 flex flex-col items-center text-center gap-4">
        <h2
            className="font-display font-bold gradient-text leading-tight tracking-tight"
            style={{fontSize: "clamp(20px, 2.5vw, 32px)"}}
        >
            stop fighting your env files
        </h2>
        <p className="text-sm text-muted-foreground leading-relaxed max-w-md">
            Install the CLI, grab the JetBrains plugin, or format right here in the browser.
            No accounts, no servers — it all runs on your machine.
        </p>
        <div className="flex flex-wrap gap-2 justify-center">
            <a
                href="#demo"
                className="inline-flex items-center gap-1.5 h-9 px-5 rounded-lg bg-primary text-primary-foreground font-mono text-xs font-semibold hover:opacity-90 transition-opacity focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
                <Play className="w-3.5 h-3.5"/>
                try it live
            </a>
            <a
                href={`https://github.com/${REPO}`}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1.5 h-9 px-5 rounded-lg border border-border text-foreground font-mono text-xs font-semibold hover:bg-accent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
                <Github className="w-3.5 h-3.5"/>
                view on github
            </a>
        </div>
    </div>
);

export default ClosingCTA;
