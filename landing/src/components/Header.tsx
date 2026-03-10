import {Github} from "lucide-react";
import ThemeToggle from "./ThemeToggle";

const Header = () => (
    <header className="flex items-center justify-between px-4 lg:px-6 py-3">
        <a href="/" className="flex items-center gap-2" aria-label="DotEnvify home">
            <img src="/logo.gif" alt="" aria-hidden="true" className="w-7 h-7"/>
            <span className="font-display font-semibold text-foreground text-sm tracking-tight">
        dotenvify
      </span>
        </a>
        <div className="flex items-center gap-1.5">
            <ThemeToggle/>
            <a
                href="https://github.com/webb1es/dotenvify"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1.5 h-8 px-3 rounded-lg text-sm font-mono text-muted-foreground hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                aria-label="View on GitHub"
            >
                <Github className="w-4 h-4"/>
                <span className="hidden sm:inline text-xs">GitHub</span>
            </a>
        </div>
    </header>
);

export default Header;
