const HeroSection = () => (
  <section className="relative flex flex-col items-center text-center px-4 pt-2 pb-3 lg:pb-4">
    <div className="absolute inset-0 hero-glow gradient-pulse pointer-events-none" aria-hidden="true" />
    <h1
      className="relative gradient-text font-bold leading-tight"
      style={{ fontSize: "clamp(24px, 4vw, 44px)" }}
    >
      Azure DevOps variables → .env in one click.
    </h1>
    <p
      className="relative mt-2 text-muted-foreground max-w-2xl"
      style={{ fontSize: "clamp(14px, 1.5vw, 18px)" }}
    >
      A JetBrains plugin that pulls, formats, and manages your environment variables — without leaving your IDE.
    </p>
    <div className="relative flex flex-col sm:flex-row gap-2 mt-3">
      <a
        href="https://plugins.jetbrains.com/plugin/dev.webbies.dotenvify"
        target="_blank"
        rel="noopener noreferrer"
        className="inline-flex items-center justify-center h-10 px-6 rounded-lg bg-jb-blue text-white font-semibold text-[15px] hover:opacity-90 transition-opacity focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring w-full sm:w-auto"
      >
        Get Plugin
      </a>
      <a
        href="https://github.com/webb1es/dotenvify"
        target="_blank"
        rel="noopener noreferrer"
        className="inline-flex items-center justify-center h-10 px-6 rounded-lg border border-border text-foreground font-semibold text-[15px] hover:bg-accent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring w-full sm:w-auto"
      >
        View on GitHub
      </a>
    </div>
  </section>
);

export default HeroSection;
