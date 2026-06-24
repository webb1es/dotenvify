import {Github} from "lucide-react";
import ThemeToggle from "./ThemeToggle";
import JetBrainsWidget from "@/components/JetBrainsWidget";
import {GITHUB_URL} from "@/lib/constants";

const Header = () => (
    <header className="sticky top-0 z-50 flex items-center justify-between px-4 lg:px-6 py-3 border-b border-border bg-background/80 backdrop-blur-md supports-[backdrop-filter]:bg-background/60">
        <a href="/" className="flex items-center gap-2" aria-label="DotEnvify home">
            <img src="/logo.gif" alt="" aria-hidden="true" className="w-14 h-14 -my-3 shrink-0"/>
            <span className="font-display font-semibold text-foreground text-sm tracking-tight">
        dotenvify
      </span>
        </a>
        <div className="flex items-center gap-1.5">
            <JetBrainsWidget type="install" className="flex items-center"/>
            <a
                href={GITHUB_URL}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1.5 h-8 px-3 rounded-lg text-sm font-mono text-muted-foreground hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                aria-label="View on GitHub"
            >
                <Github className="w-4 h-4"/>
                <span className="hidden sm:inline text-xs">GitHub</span>
            </a>
            <ThemeToggle/>
        </div>
    </header>
);

export default Header;
