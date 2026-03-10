import { Github } from "lucide-react";

const CURRENT_YEAR = new Date().getFullYear();
const BUILD_DATE = "2026-03-10";

const Footer = () => (
  <footer className="flex items-center justify-between px-4 lg:px-6 py-3 text-[11px] font-mono text-muted-foreground">
    <div className="flex items-center gap-2">
      <span>&copy; {CURRENT_YEAR} webbies.dev</span>
      <span aria-hidden="true">&middot;</span>
      <span>MIT</span>
    </div>

    <div className="flex items-center gap-2">
      <span className="hidden sm:inline">updated {BUILD_DATE}</span>
      <span className="hidden sm:inline" aria-hidden="true">&middot;</span>
      <a
        href="https://github.com/webb1es/dotenvify"
        target="_blank"
        rel="noopener noreferrer"
        className="inline-flex items-center gap-1 hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded"
        aria-label="View on GitHub"
      >
        <Github className="w-3 h-3" />
        <span className="hidden sm:inline">source</span>
      </a>
    </div>
  </footer>
);

export default Footer;
