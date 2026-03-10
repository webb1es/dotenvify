import {useState} from "react";
import {Check, Clock, Copy, Download, ExternalLink, Github} from "lucide-react";

const IDE_ICONS: Record<string, { src: string; color: string }> = {
    IntelliJ: {src: "/icons/intellij.svg", color: "#E44332"},
    WebStorm: {src: "/icons/webstorm.svg", color: "#07C3F2"},
    GoLand: {src: "/icons/goland.svg", color: "#9B59D0"},
    PyCharm: {src: "/icons/pycharm.svg", color: "#21D789"},
    Rider: {src: "/icons/rider.svg", color: "#DD1265"},
    CLion: {src: "/icons/clion.svg", color: "#21D789"},
    "VS Code": {src: "/icons/vscode.svg", color: "#007ACC"},
};

const CTA_ICONS: Record<string, React.ReactNode> = {
    "view source": <Github className="w-3.5 h-3.5"/>,
    "get plugin": <Download className="w-3.5 h-3.5"/>,
    "get extension": <Download className="w-3.5 h-3.5"/>,
    "learn more": <ExternalLink className="w-3.5 h-3.5"/>,
};

interface ProductCardProps {
    name: string;
    icon: React.ReactNode;
    description: string;
    installCmd?: string;
    href?: string;
    ctaLabel?: string;
    comingSoon?: boolean;
    badges?: string[];
}

const ProductCard = ({
                         name,
                         icon,
                         description,
                         installCmd,
                         href,
                         ctaLabel,
                         comingSoon,
                         badges,
                     }: ProductCardProps) => {
    const [copied, setCopied] = useState(false);

    const handleCopy = async () => {
        if (!installCmd) return;
        await navigator.clipboard.writeText(installCmd);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <div className={`bento-cell p-4 flex flex-col justify-between h-full ${comingSoon ? "opacity-60" : ""}`}>
            <div>
                <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center gap-2">
                        <span className="text-primary">{icon}</span>
                        <span className="font-display text-sm font-semibold text-foreground">{name}</span>
                    </div>
                    {comingSoon && (
                        <span
                            className="inline-flex items-center gap-1 h-5 px-2 rounded-full bg-de-orange/10 text-de-orange text-[10px] font-mono font-semibold">
              <Clock className="w-2.5 h-2.5"/>
              soon
            </span>
                    )}
                </div>

                <p className="text-xs text-muted-foreground leading-relaxed">{description}</p>

                {installCmd && (
                    <button
                        onClick={handleCopy}
                        className="w-full flex items-center gap-2 mt-3 px-2.5 py-1.5 rounded-md bg-de-surface border border-de-surface-border font-mono text-[11px] text-left hover:border-primary/30 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring group"
                        aria-label={copied ? "Copied" : `Copy: ${installCmd}`}
                    >
                        <span className="text-muted-foreground select-none" aria-hidden="true">$</span>
                        <code className="flex-1 text-foreground truncate">{installCmd}</code>
                        {copied ? (
                            <Check className="w-3 h-3 text-de-green flex-shrink-0"/>
                        ) : (
                            <Copy className="w-3 h-3 text-muted-foreground group-hover:text-foreground flex-shrink-0"/>
                        )}
                    </button>
                )}

                {badges && (
                    <div className="flex flex-wrap gap-1.5 mt-3">
                        {badges.map((b) => {
                            const ide = IDE_ICONS[b];
                            return (
                                <span
                                    key={b}
                                    className="inline-flex items-center gap-1 h-5 px-1.5 rounded border border-de-surface-border text-[10px] font-mono text-muted-foreground"
                                >
                  {ide && (
                      <img
                          src={ide.src}
                          alt=""
                          aria-hidden="true"
                          className="w-3 h-3 dark:invert"
                      />
                  )}
                                    <span style={{color: ide?.color}}>{b}</span>
                </span>
                            );
                        })}
                    </div>
                )}
            </div>

            {href && ctaLabel && (
                <a
                    href={href}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center justify-center gap-1.5 h-8 px-4 mt-4 rounded-lg bg-primary text-primary-foreground font-mono text-xs font-semibold hover:opacity-90 transition-opacity focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                >
                    {CTA_ICONS[ctaLabel]}
                    {ctaLabel}
                </a>
            )}
        </div>
    );
};

export default ProductCard;
