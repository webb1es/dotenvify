import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { execFileSync } from "node:child_process";
import { writeFileSync, readFileSync, readdirSync, unlinkSync, mkdtempSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";

const BIN = join(import.meta.dirname, "..", "bin", "dotenvify.js");

let dir: string;

beforeEach(() => {
  dir = mkdtempSync(join(tmpdir(), "dotenvify-cli-test-"));
});

afterEach(() => {
  try {
    for (const f of readdirSync(dir)) unlinkSync(join(dir, f));
  } catch {}
});

function run(args: string[]): string {
  return execFileSync("node", [BIN, ...args], { encoding: "utf-8" });
}

describe("dotenvify CLI", () => {
  it("converts a file to .env", () => {
    const input = join(dir, "input.txt");
    const output = join(dir, ".env");
    writeFileSync(input, "API_KEY=abc123\nDB_HOST=localhost\n");

    const stdout = run([input, "-o", output, "-f"]);
    expect(stdout).toContain("Written 2 variable(s)");

    const result = readFileSync(output, "utf-8");
    expect(result).toContain("API_KEY=abc123");
    expect(result).toContain("DB_HOST=localhost");
  });

  it("adds export prefix with -e", () => {
    const input = join(dir, "input.txt");
    const output = join(dir, ".env");
    writeFileSync(input, "KEY=val\n");

    run([input, "-o", output, "-f", "-e"]);
    expect(readFileSync(output, "utf-8")).toBe("export KEY=val\n");
  });

  it("exits with error for missing source", () => {
    expect(() => run([join(dir, "missing.txt")])).toThrow();
  });

  it("handles mixed input formats", () => {
    const input = join(dir, "input.txt");
    const output = join(dir, ".env");
    writeFileSync(
      input,
      `API_KEY=abc
SECRET "my-secret"
DATABASE_URL
postgres://localhost/db
export NODE_ENV=production`,
    );

    run([input, "-o", output, "-f"]);
    const result = readFileSync(output, "utf-8");
    const lines = result.trim().split("\n");
    expect(lines).toHaveLength(4);
  });
});
