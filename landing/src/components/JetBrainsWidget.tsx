import {useEffect, useState} from "react";

const SCRIPT_SRC = "https://plugins.jetbrains.com/assets/scripts/mp-widget.js";
const PLUGIN_ID = 32351;

declare global {
    interface Window {
        MarketplaceWidget?: {
            setupMarketplaceWidget: (type: string, id: number, selector: string) => void;
        };
    }
}

let scriptPromise: Promise<void> | null = null;

function loadWidgetScript(): Promise<void> {
    if (window.MarketplaceWidget) return Promise.resolve();
    if (scriptPromise) return scriptPromise;
    scriptPromise = new Promise((resolve, reject) => {
        const script = document.createElement("script");
        script.src = SCRIPT_SRC;
        script.async = true;
        script.onload = () => resolve();
        script.onerror = () => {
            scriptPromise = null;
            reject(new Error("JetBrains Marketplace widget failed to load"));
        };
        document.head.appendChild(script);
    });
    return scriptPromise;
}

// Selector-safe unique ids; the widget targets its container via a CSS selector, so useId
// (which yields ids containing ":") is unusable here.
let instanceCount = 0;

/**
 * Renders an official JetBrains Marketplace widget for the DotEnvify plugin.
 * `install` is a one-click install button; `card` shows live version/downloads/rating.
 * Stays empty when the widget can't load (e.g. before the plugin is public).
 */
const JetBrainsWidget = ({type = "install", className}: { type?: "install" | "card"; className?: string }) => {
    const [elementId] = useState(() => `jb-mp-widget-${type}-${++instanceCount}`);

    useEffect(() => {
        let cancelled = false;
        loadWidgetScript()
            .then(() => {
                const target = document.getElementById(elementId);
                if (cancelled || !window.MarketplaceWidget || !target) return;
                target.innerHTML = "";
                window.MarketplaceWidget.setupMarketplaceWidget(type, PLUGIN_ID, `#${elementId}`);
            })
            .catch(() => {
                // Unavailable (offline or plugin not yet public); leave empty.
            });
        return () => {
            cancelled = true;
        };
    }, [type, elementId]);

    return <div id={elementId} className={className}/>;
};

export default JetBrainsWidget;
