const ides = [
  { name: "IntelliJ IDEA", short: "IntelliJ" },
  { name: "GoLand", short: "GoLand" },
  { name: "WebStorm", short: "WebStorm" },
  { name: "PyCharm", short: "PyCharm" },
  { name: "Rider", short: "Rider" },
  { name: "CLion", short: "CLion" },
  { name: "RubyMine", short: "RubyMine" },
];

const IdeLogos = () => (
  <div className="flex flex-col items-center lg:items-start gap-2">
    <p className="text-xs font-semibold text-foreground">
      Works across all JetBrains IDEs · 2024.1+
    </p>
    <div className="flex flex-wrap gap-1.5 justify-center lg:justify-start">
      {ides.map((ide) => (
        <span
          key={ide.short}
          className="inline-flex items-center h-7 px-2.5 rounded-md border border-jb-surface-border bg-card text-[11px] text-muted-foreground font-medium"
        >
          {ide.short}
        </span>
      ))}
    </div>
  </div>
);

export default IdeLogos;
