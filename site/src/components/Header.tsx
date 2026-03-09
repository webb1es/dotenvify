import { Github } from "lucide-react";
import ThemeToggle from "./ThemeToggle";

const Header = () => (
  <header className="fixed top-0 left-0 right-0 z-50 h-12 flex items-center justify-between px-4 lg:px-6 bg-background/80 backdrop-blur-xl border-b border-border">
    <a href="/" className="flex items-center gap-2" aria-label="DotEnvify home">
      <div className="w-6 h-6 rounded bg-gradient-to-br from-jb-blue to-jb-purple flex items-center justify-center">
        <span className="text-[10px] font-bold text-white leading-none">.e</span>
      </div>
      <span className="font-semibold text-foreground text-sm">DotEnvify</span>
    </a>
    <div className="flex items-center gap-2">
      <ThemeToggle />
      <a
        href="https://github.com/webb1es/dotenvify"
        target="_blank"
        rel="noopener noreferrer"
        className="inline-flex items-center gap-1.5 h-8 px-3 rounded-lg border border-border text-sm font-medium text-foreground hover:bg-accent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        <Github className="w-4 h-4" />
        <span className="hidden sm:inline">GitHub</span>
      </a>
      <a
        href="https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify"
        target="_blank"
        rel="noopener noreferrer"
        className="inline-flex items-center h-8 px-4 rounded-lg bg-jb-blue text-white text-sm font-semibold hover:opacity-90 transition-opacity focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        Get Plugin
      </a>
    </div>
  </header>
);

export default Header;
