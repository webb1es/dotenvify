/** True when the value is wrapped in matching single or double quotes. */
export function isQuoted(value: string): boolean {
    return (
        (value.startsWith('"') && value.endsWith('"')) ||
        (value.startsWith("'") && value.endsWith("'"))
    );
}

/** Strip one layer of matching surrounding quotes, if present. */
export function unquote(value: string): string {
    return isQuoted(value) ? value.slice(1, -1) : value;
}
