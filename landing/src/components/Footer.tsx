import { Link } from "react-router-dom";

const Footer = () => (
  <footer className="flex items-center justify-center px-4 py-3 text-[11px] font-mono text-muted-foreground gap-x-2 gap-y-1 flex-wrap">
    <span>MIT</span>
    <span aria-hidden="true">·</span>
    <a href="https://webbies.dev" target="_blank" rel="noopener noreferrer" className="hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded">
      webbies.dev
    </a>
    <span aria-hidden="true">·</span>
    <a href="mailto:dotenvify@webbies.dev" className="hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded">
      dotenvify@webbies.dev
    </a>
    <span aria-hidden="true">·</span>
    <Link to="/terms" className="hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded">
      terms
    </Link>
    <span aria-hidden="true">·</span>
    <Link to="/privacy" className="hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring rounded">
      privacy
    </Link>
  </footer>
);

export default Footer;
