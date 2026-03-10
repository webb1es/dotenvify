import {describe, expect, it} from "vitest";
import {formatDotEnv, isHTTPURL, isURL} from "../src/formatter.js";
import type {EnvEntry} from "../src/models.js";

describe("formatDotEnv", () => {
    const entries: EnvEntry[] = [
        {key: "ZULU", value: "last"},
        {key: "ALPHA", value: "first"},
        {key: "BRAVO", value: "second"},
    ];

    it("sorts keys alphabetically by default", () => {
        expect(formatDotEnv(entries)).toBe("ALPHA=first\nBRAVO=second\nZULU=last\n");
    });

    it("preserves insertion order when sort is false", () => {
        const result = formatDotEnv(entries, {sort: false});
        expect(result).toBe("ZULU=last\nALPHA=first\nBRAVO=second\n");
    });

    it("adds export prefix", () => {
        const result = formatDotEnv([{key: "API_KEY", value: "test123"}], {export: true});
        expect(result).toBe("export API_KEY=test123\n");
    });

    it("filters lowercase keys with noLower", () => {
        const mixed: EnvEntry[] = [
            {key: "API_KEY", value: "test123"},
            {key: "lowercase", value: "should-be-skipped"},
            {key: "UPPERCASE", value: "included"},
            {key: "MixedCase", value: "included"},
        ];
        const result = formatDotEnv(mixed, {noLower: true});
        expect(result).toContain("API_KEY=");
        expect(result).toContain("UPPERCASE=");
        expect(result).toContain("MixedCase=");
        expect(result).not.toContain("lowercase=");
    });

    it("filters non-HTTP URLs with urlOnly", () => {
        const urls: EnvEntry[] = [
            {key: "DATABASE_URL", value: "https://example.com/db"},
            {key: "API_URL", value: "http://api.example.com"},
            {key: "REGULAR_VAR", value: "not-a-url"},
            {key: "POSTGRES", value: "postgres://localhost:5432/db"},
        ];
        const result = formatDotEnv(urls, {urlOnly: true});
        expect(result).toContain("DATABASE_URL=");
        expect(result).toContain("API_URL=");
        expect(result).not.toContain("REGULAR_VAR=");
        expect(result).not.toContain("POSTGRES=");
    });

    it("smart-quotes URLs", () => {
        const result = formatDotEnv([{key: "URL", value: "https://example.com"}]);
        expect(result).toBe('URL="https://example.com"\n');
    });

    it("smart-quotes values with spaces", () => {
        const result = formatDotEnv([{key: "VAR", value: "value with spaces"}]);
        expect(result).toBe('VAR="value with spaces"\n');
    });

    it("does not double-quote already quoted values", () => {
        const result = formatDotEnv([{key: "VAR", value: '"already-quoted"'}]);
        expect(result).toBe('VAR="already-quoted"\n');
    });

    it("does not quote simple values", () => {
        const result = formatDotEnv([{key: "VAR", value: "simplevalue"}]);
        expect(result).toBe("VAR=simplevalue\n");
    });

    it("returns empty string for empty entries", () => {
        expect(formatDotEnv([])).toBe("");
    });
});

describe("isURL", () => {
    it.each([
        ["http://example.com", true],
        ["https://example.com", true],
        ["postgres://localhost:5432/db", true],
        ["mysql://localhost", true],
        ["mongodb://localhost", true],
        ["redis://localhost", true],
        ["ftp://example.com", true],
        ["sftp://example.com", true],
        ["ssh://example.com", true],
        ["git://github.com", true],
        ["mailto:test@example.com", true],
        ["just-a-string", false],
        ["", false],
    ])("isURL(%s) = %s", (input, expected) => {
        expect(isURL(input)).toBe(expected);
    });
});

describe("isHTTPURL", () => {
    it.each([
        ["http://example.com", true],
        ["https://example.com", true],
        ["postgres://localhost", false],
        ["ftp://example.com", false],
        ["not-a-url", false],
        ["", false],
    ])("isHTTPURL(%s) = %s", (input, expected) => {
        expect(isHTTPURL(input)).toBe(expected);
    });
});
