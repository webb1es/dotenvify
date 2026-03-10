import { Play, Github, Star, GitFork, CircleDot, Users } from "lucide-react";

const REPO = "webb1es/dotenvify";

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

    <div className="mt-6 space-y-4">
      <div className="flex flex-wrap gap-2">
        <a
          href="#demo"
          className="inline-flex items-center gap-1.5 h-9 px-5 rounded-lg bg-primary text-primary-foreground font-mono text-xs font-semibold hover:opacity-90 transition-opacity focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <Play className="w-3.5 h-3.5" />
          try it live
        </a>
        <a
          href={`https://github.com/${REPO}`}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-1.5 h-9 px-5 rounded-lg border border-border text-foreground font-mono text-xs font-semibold hover:bg-accent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <Github className="w-3.5 h-3.5" />
          git clone
        </a>
      </div>

      {/* GitHub badges */}
      <div className="flex flex-wrap gap-1.5">
        <a href={`https://github.com/${REPO}/stargazers`} target="_blank" rel="noopener noreferrer" className="inline-flex items-center gap-1 h-5 px-1.5 rounded border border-de-surface-border text-[10px] font-mono text-muted-foreground hover:text-foreground hover:border-foreground/20 transition-colors">
          <Star className="w-2.5 h-2.5" />
          star
        </a>
        <a href={`https://github.com/${REPO}/fork`} target="_blank" rel="noopener noreferrer" className="inline-flex items-center gap-1 h-5 px-1.5 rounded border border-de-surface-border text-[10px] font-mono text-muted-foreground hover:text-foreground hover:border-foreground/20 transition-colors">
          <GitFork className="w-2.5 h-2.5" />
          fork
        </a>
        <a href={`https://github.com/${REPO}/issues`} target="_blank" rel="noopener noreferrer" className="inline-flex items-center gap-1 h-5 px-1.5 rounded border border-de-surface-border text-[10px] font-mono text-muted-foreground hover:text-foreground hover:border-foreground/20 transition-colors">
          <CircleDot className="w-2.5 h-2.5" />
          issues
        </a>
        <a href={`https://github.com/${REPO}/graphs/contributors`} target="_blank" rel="noopener noreferrer" className="inline-flex items-center gap-1 h-5 px-1.5 rounded border border-de-surface-border text-[10px] font-mono text-muted-foreground hover:text-foreground hover:border-foreground/20 transition-colors">
          <Users className="w-2.5 h-2.5" />
          contribute
        </a>
      </div>
    </div>
  </div>
);

export default HeroCell;
