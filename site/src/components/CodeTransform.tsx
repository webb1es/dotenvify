import { ArrowRight } from "lucide-react";

const inputLines = [
  { text: 'API_KEY abc123', dim: false },
  { text: 'DATABASE_URL="postgres://localhost:5432/db"', dim: false },
  { text: 'secret_token mytoken', dim: true },
  { text: 'REDIS_HOST', dim: false },
  { text: 'redis.example.com', dim: false },
  { text: 'export NODE_ENV=production', dim: true },
];

const outputLines = [
  { key: 'API_KEY', value: 'abc123' },
  { key: 'DATABASE_URL', value: '"postgres://localhost:5432/db"' },
  { key: 'NODE_ENV', value: 'production' },
  { key: 'REDIS_HOST', value: 'redis.example.com' },
];

const CodeTransform = () => (
  <div
    className="flex flex-col h-full"
    role="region"
    aria-label="Code example: before and after formatting"
  >
    <div className="flex-1 grid grid-cols-1 md:grid-cols-[1fr_auto_1fr] gap-2 items-stretch min-h-0">
      {/* Before */}
      <div className="rounded-lg border border-jb-surface-border bg-card overflow-hidden flex flex-col">
        <div className="flex items-center gap-1.5 px-3 py-1.5 border-b border-jb-surface-border">
          <div className="w-2 h-2 rounded-full bg-destructive/60" aria-hidden="true" />
          <div className="w-2 h-2 rounded-full bg-yellow-500/60" aria-hidden="true" />
          <div className="w-2 h-2 rounded-full bg-jb-green/60" aria-hidden="true" />
          <span className="text-[10px] text-muted-foreground ml-1 font-mono">raw-input.txt</span>
        </div>
        <pre className="flex-1 p-3 font-mono text-[11px] leading-relaxed overflow-auto">
          {inputLines.map((line, i) => (
            <div key={i} className={line.dim ? "text-muted-foreground" : "text-foreground"}>
              {line.text}
            </div>
          ))}
        </pre>
      </div>

      {/* Arrow */}
      <div className="hidden md:flex items-center justify-center">
        <ArrowRight className="w-5 h-5 text-jb-blue" aria-hidden="true" />
      </div>

      {/* After */}
      <div className="rounded-lg border border-jb-surface-border bg-card overflow-hidden flex flex-col">
        <div className="flex items-center gap-1.5 px-3 py-1.5 border-b border-jb-surface-border">
          <div className="w-2 h-2 rounded-full bg-destructive/60" aria-hidden="true" />
          <div className="w-2 h-2 rounded-full bg-yellow-500/60" aria-hidden="true" />
          <div className="w-2 h-2 rounded-full bg-jb-green/60" aria-hidden="true" />
          <span className="text-[10px] text-muted-foreground ml-1 font-mono">.env</span>
        </div>
        <pre className="flex-1 p-3 font-mono text-[11px] leading-relaxed overflow-auto">
          {outputLines.map((line, i) => (
            <div key={i}>
              <span className="text-jb-green">{line.key}</span>
              <span className="text-muted-foreground">=</span>
              <span className="text-foreground">{line.value}</span>
            </div>
          ))}
        </pre>
      </div>
    </div>
    <p className="text-[10px] text-muted-foreground text-center mt-1.5">
      Sorted · Lowercase filtered · Auto-quoted · Export stripped
    </p>
  </div>
);

export default CodeTransform;
