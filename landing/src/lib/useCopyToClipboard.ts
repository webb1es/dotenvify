import {useState} from "react";

/** Copy text to the clipboard with a transient `copied` flag. No-op for empty text. */
export function useCopyToClipboard(timeoutMs = 2000) {
    const [copied, setCopied] = useState(false);
    const copy = async (text?: string) => {
        if (!text) return;
        await navigator.clipboard.writeText(text);
        setCopied(true);
        setTimeout(() => setCopied(false), timeoutMs);
    };
    return {copied, copy};
}
