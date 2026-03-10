const formats = [
    {label: "KEY=VALUE", example: 'API_KEY=secret123'},
    {label: "quoted", example: 'DB_URL="postgres://localhost"'},
    {label: "export", example: "export NODE_ENV=production"},
    {label: "space sep", example: "REDIS_HOST localhost"},
    {label: "line pairs", example: "SECRET_KEY\\n4f9a2b..."},
];

const FormatsCell = () => (
    <div className="bento-cell p-4 h-full">
        <p className="font-mono text-[10px] text-muted-foreground mb-3">
            {"// "}supported input formats
        </p>
        <div className="space-y-2">
            {formats.map((f) => (
                <div key={f.label} className="flex items-center gap-3">
          <span className="flex-shrink-0 w-16 font-mono text-[10px] text-primary font-medium">
            {f.label}
          </span>
                    <code className="font-mono text-[11px] text-muted-foreground truncate">
                        {f.example}
                    </code>
                </div>
            ))}
        </div>
    </div>
);

export default FormatsCell;
