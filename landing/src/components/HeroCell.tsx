const HeroCell = () => (
  <div className="bento-cell-static p-6 lg:p-8 flex flex-col justify-between h-full">
    <div>
      {/* Comment-style label */}
      <p className="font-mono text-xs text-muted-foreground mb-4">
        {"//"} convert messy env vars
      </p>

      <h1
        className="font-display font-bold gradient-text leading-tight tracking-tight"
        style={{ fontSize: "clamp(24px, 3.5vw, 40px)" }}
      >
        paste anything,
        <br />
        get .env
      </h1>

      <p className="mt-3 text-sm text-muted-foreground leading-relaxed max-w-xs">
        CLI, JetBrains plugin, or right here in the browser.
        Parse any format. Zero config.
      </p>
    </div>

    <div className="flex flex-wrap gap-2 mt-6">
      <a
        href="#demo"
        className="inline-flex items-center h-9 px-5 rounded-lg bg-primary text-primary-foreground font-mono text-xs font-semibold hover:opacity-90 transition-opacity focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        try it live ↓
      </a>
      <a
        href="https://github.com/webb1es/dotenvify"
        target="_blank"
        rel="noopener noreferrer"
        className="inline-flex items-center h-9 px-5 rounded-lg border border-border text-foreground font-mono text-xs font-semibold hover:bg-accent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        git clone
      </a>
    </div>
  </div>
);

export default HeroCell;
