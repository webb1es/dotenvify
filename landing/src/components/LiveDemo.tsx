import {useCallback, useEffect, useRef, useState} from "react";
import type {FormatOptions} from "@dotenvify/core/browser";
import {formatDotEnv, parseDotEnv} from "@dotenvify/core/browser";
import {ArrowDownAZ, CaseLower, Check, Copy, FileOutput, Link} from "lucide-react";

const PLACEHOLDER = `API_KEY abc123
DATABASE_URL="postgres://localhost:5432/db"
secret_token mytoken
REDIS_HOST
redis.example.com
export NODE_ENV=production`;

const LiveDemo = () => {
    const [input, setInput] = useState("");
    const [options, setOptions] = useState<FormatOptions>({
        sort: true,
        export: false,
        noLower: false,
        urlOnly: false,
    });
    const [output, setOutput] = useState("");
    const [copied, setCopied] = useState(false);
    const timerRef = useRef<ReturnType<typeof setTimeout>>();

    const process = useCallback((text: string, opts: FormatOptions) => {
        const source = text.trim() ? text : PLACEHOLDER;
        const result = parseDotEnv(source);
        setOutput(formatDotEnv(result.entries, opts));
    }, []);

    useEffect(() => {
        clearTimeout(timerRef.current);
        timerRef.current = setTimeout(() => process(input, options), 100);
        return () => clearTimeout(timerRef.current);
    }, [input, options, process]);

    const handleCopy = async () => {
        if (!output) return;
        await navigator.clipboard.writeText(output);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    const toggle = (key: keyof FormatOptions) =>
        setOptions((prev) => ({...prev, [key]: !prev[key]}));

    return (
        <div className="bento-cell-static flex flex-col h-full" id="demo">
            {/* Toolbar */}
            <div className="flex items-center justify-between px-3 py-2 border-b border-de-surface-border">
                <div className="flex items-center gap-1.5 overflow-x-auto" role="group" aria-label="Format options">
                    {([
                        ["sort", "sort", <ArrowDownAZ key="sort-icon" className="w-3 h-3"/>],
                        ["export", "export", <FileOutput key="export-icon" className="w-3 h-3"/>],
                        ["noLower", "skip lowercase", <CaseLower key="lower-icon" className="w-3 h-3"/>],
                        ["urlOnly", "urls only", <Link key="url-icon" className="w-3 h-3"/>],
                    ] as [keyof FormatOptions, string, React.ReactNode][]).map(([key, label, icon]) => (
                        <button
                            key={key}
                            onClick={() => toggle(key)}
                            aria-pressed={!!options[key]}
                            className={`flex-shrink-0 inline-flex items-center gap-1 h-6 px-2 rounded font-mono text-[10px] font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring ${
                                options[key]
                                    ? "bg-primary/15 text-primary border border-primary/30"
                                    : "text-muted-foreground hover:text-foreground border border-transparent"
                            }`}
                        >
                            {icon}
                            {label}
                        </button>
                    ))}
                </div>
                <button
                    onClick={handleCopy}
                    aria-label={copied ? "Copied" : "Copy output"}
                    className="flex-shrink-0 inline-flex items-center gap-1 h-6 px-2 rounded font-mono text-[10px] text-muted-foreground hover:text-foreground transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                >
                    {copied ? <Check className="w-3 h-3 text-de-green"/> : <Copy className="w-3 h-3"/>}
                    <span>{copied ? "copied!" : "copy"}</span>
                </button>
            </div>

            {/* Editor */}
            <div className="flex-1 grid grid-cols-1 md:grid-cols-2 min-h-0">
                {/* Input */}
                <div className="flex flex-col border-b md:border-b-0 md:border-r border-de-surface-border">
                    <div className="px-3 py-1.5 border-b border-de-surface-border">
                        <span className="font-mono text-[10px] text-muted-foreground">input.txt</span>
                    </div>
                    <div className="flex-1 relative">
            <textarea
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder={PLACEHOLDER}
                spellCheck={false}
                aria-label="Paste environment variables"
                className="absolute inset-0 w-full h-full p-3 font-mono text-xs leading-relaxed bg-transparent text-foreground placeholder:text-muted-foreground/30 resize-none focus:outline-none"
            />
                    </div>
                </div>

                {/* Output */}
                <div className="flex flex-col">
                    <div className="px-3 py-1.5 border-b border-de-surface-border">
                        <span className="font-mono text-[10px] text-de-green">.env</span>
                    </div>
                    <pre
                        className="flex-1 p-3 font-mono text-xs leading-relaxed overflow-auto"
                        aria-live="polite"
                        aria-label="Formatted output"
                    >
            {output.split("\n").map((line, i) => {
                if (!line) return null;
                const eqIdx = line.indexOf("=");
                if (eqIdx === -1) return <div key={i}>{line}</div>;

                const prefix = line.startsWith("export ") ? "export " : "";
                const rest = prefix ? line.slice(7) : line;
                const restEq = rest.indexOf("=");
                const key = rest.slice(0, restEq);
                const value = rest.slice(restEq + 1);

                return (
                    <div key={i}>
                        {prefix && <span className="text-de-purple">{prefix}</span>}
                        <span className="text-de-green">{key}</span>
                        <span className="text-muted-foreground">=</span>
                        <span className="text-foreground">{value}</span>
                    </div>
                );
            })}
          </pre>
                </div>
            </div>
        </div>
    );
};

export default LiveDemo;
