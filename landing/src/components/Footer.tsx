const Footer = () => (
  <footer className="flex items-center justify-center px-4 h-8 text-[11px] text-muted-foreground gap-1 flex-wrap">
    <span>MIT License</span>
    <span aria-hidden="true">·</span>
    <a href="https://webbies.dev" target="_blank" rel="noopener noreferrer" className="hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded">
      Built by Webbies
    </a>
    <span aria-hidden="true">·</span>
    <a href="https://github.com/webb1es/dotenvify" target="_blank" rel="noopener noreferrer" className="hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded">
      GitHub
    </a>
    <span aria-hidden="true">·</span>
    <a href="mailto:hello@webbies.dev" className="hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded">
      hello@webbies.dev
    </a>
    <span aria-hidden="true">·</span>
    <span>© 2026</span>
  </footer>
);

export default Footer;
