import type {EnvEntry, FormatOptions} from "./models.js";

const URL_PREFIXES = [
    "http://",
    "https://",
    "ftp://",
    "sftp://",
    "ssh://",
    "git://",
    "file://",
    "mailto:",
    "postgres://",
    "mysql://",
    "mongodb://",
    "redis://",
];

export function isURL(value: string): boolean {
    return URL_PREFIXES.some((prefix) => value.startsWith(prefix));
}

export function isHTTPURL(value: string): boolean {
    return value.startsWith("http://") || value.startsWith("https://");
}

function isQuoted(value: string): boolean {
    return (
        (value.startsWith('"') && value.endsWith('"')) ||
        (value.startsWith("'") && value.endsWith("'"))
    );
}

function isLowercase(key: string): boolean {
    return key === key.toLowerCase() && key !== key.toUpperCase();
}

function smartQuote(value: string): string {
    if (isQuoted(value)) return value;
    if (isURL(value) || value.includes(" ")) return `"${value}"`;
    return value;
}

/**
 * Format environment entries into a .env string.
 *
 * Options:
 * - sort: alphabetically sort keys (default: true)
 * - export: prepend "export " to each line
 * - noLower: skip keys that are entirely lowercase
 * - urlOnly: only include entries with HTTP/HTTPS URL values
 */
export function formatDotEnv(
    entries: EnvEntry[],
    options: FormatOptions = {},
): string {
    const {export: useExport = false, sort = true, noLower = false, urlOnly = false} = options;

    let filtered = entries;

    if (noLower) {
        filtered = filtered.filter((e) => !isLowercase(e.key));
    }

    if (urlOnly) {
        filtered = filtered.filter((e) => isHTTPURL(e.value));
    }

    if (sort) {
        filtered = [...filtered].sort((a, b) => a.key.localeCompare(b.key));
    }

    const lines = filtered.map((e) => {
        const value = smartQuote(e.value);
        return useExport ? `export ${e.key}=${value}` : `${e.key}=${value}`;
    });

    return lines.length > 0 ? lines.join("\n") + "\n" : "";
}
