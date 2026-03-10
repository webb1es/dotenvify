export interface EnvEntry {
    key: string;
    value: string;
}

export interface FormatOptions {
    export?: boolean;
    sort?: boolean;
    noLower?: boolean;
    urlOnly?: boolean;
}

export interface ParseResult {
    entries: EnvEntry[];
    warnings: string[];
}
