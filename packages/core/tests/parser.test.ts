import {describe, expect, it} from "vitest";
import {parseDotEnv} from "../src/parser.js";

describe("parseDotEnv", () => {
    it("parses KEY=VALUE format", () => {
        const {entries} = parseDotEnv("API_KEY=test123\nDATABASE_URL=postgres://localhost");
        expect(entries).toEqual([
            {key: "API_KEY", value: "test123"},
            {key: "DATABASE_URL", value: "postgres://localhost"},
        ]);
    });

    it("parses double-quoted values", () => {
        const {entries} = parseDotEnv('API_KEY="test with spaces"');
        expect(entries).toEqual([{key: "API_KEY", value: "test with spaces"}]);
    });

    it("parses single-quoted values", () => {
        const {entries} = parseDotEnv("API_KEY='single quotes'");
        expect(entries).toEqual([{key: "API_KEY", value: "single quotes"}]);
    });

    it("parses export prefix", () => {
        const {entries} = parseDotEnv("export API_KEY=test123\nexport DB=localhost");
        expect(entries).toEqual([
            {key: "API_KEY", value: "test123"},
            {key: "DB", value: "localhost"},
        ]);
    });

    it("skips comments and blank lines", () => {
        const input = `# This is a comment
API_KEY=test123

# Another comment
DATABASE_URL=postgres://localhost`;
        const {entries} = parseDotEnv(input);
        expect(entries).toEqual([
            {key: "API_KEY", value: "test123"},
            {key: "DATABASE_URL", value: "postgres://localhost"},
        ]);
    });

    it("handles value with equals sign", () => {
        const {entries} = parseDotEnv("CONNECTION_STRING=user=admin;password=secret");
        expect(entries).toEqual([
            {key: "CONNECTION_STRING", value: "user=admin;password=secret"},
        ]);
    });

    it("handles empty value", () => {
        const {entries} = parseDotEnv("EMPTY_VAR=");
        expect(entries).toEqual([{key: "EMPTY_VAR", value: ""}]);
    });

    it("trims whitespace around key and value", () => {
        const {entries} = parseDotEnv("  API_KEY  =  test123  ");
        expect(entries).toEqual([{key: "API_KEY", value: "test123"}]);
    });

    it("parses KEY VALUE space-separated format", () => {
        const {entries} = parseDotEnv("API_KEY test123");
        expect(entries).toEqual([{key: "API_KEY", value: "test123"}]);
    });

    it("parses multiple KEY VALUE pairs on one line", () => {
        const {entries} = parseDotEnv("KEY1 val1 KEY2 val2");
        expect(entries).toEqual([
            {key: "KEY1", value: "val1"},
            {key: "KEY2", value: "val2"},
        ]);
    });

    it("parses key on one line, value on next", () => {
        const {entries} = parseDotEnv("API_KEY\ntest123");
        expect(entries).toEqual([{key: "API_KEY", value: "test123"}]);
    });

    it("warns on orphan key with no value", () => {
        const {entries, warnings} = parseDotEnv("ORPHAN_KEY");
        expect(entries).toEqual([]);
        expect(warnings).toHaveLength(1);
        expect(warnings[0]).toContain("ORPHAN_KEY");
    });

    it("handles mixed formats", () => {
        const input = `# Comment
export KEY1=value1
KEY2="value2"
KEY3='value3'
KEY4=value4`;
        const {entries} = parseDotEnv(input);
        expect(entries).toEqual([
            {key: "KEY1", value: "value1"},
            {key: "KEY2", value: "value2"},
            {key: "KEY3", value: "value3"},
            {key: "KEY4", value: "value4"},
        ]);
    });

    it("last duplicate wins", () => {
        const {entries} = parseDotEnv("KEY=first\nKEY=second");
        expect(entries).toEqual([{key: "KEY", value: "second"}]);
    });

    it("handles empty input", () => {
        const {entries, warnings} = parseDotEnv("");
        expect(entries).toEqual([]);
        expect(warnings).toEqual([]);
    });

    it("mixed export and non-export", () => {
        const {entries} = parseDotEnv("API_KEY=test1\nexport DATABASE_URL=test2");
        expect(entries).toEqual([
            {key: "API_KEY", value: "test1"},
            {key: "DATABASE_URL", value: "test2"},
        ]);
    });
});
