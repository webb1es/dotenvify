import { describe, it, expect, beforeEach } from "vitest";
import { mkdtempSync, writeFileSync, readFileSync, existsSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { readEnvFile, writeEnvFile, backupFile, applyPreserve } from "../src/io.js";
import type { EnvEntry } from "../src/models.js";

function makeTmpDir(): string {
  return mkdtempSync(join(tmpdir(), "dotenvify-test-"));
}

describe("readEnvFile", () => {
  it("returns empty map for non-existent file", () => {
    const result = readEnvFile("/nonexistent/file.env");
    expect(result.size).toBe(0);
  });

  it("reads key-value pairs from file", () => {
    const dir = makeTmpDir();
    const file = join(dir, ".env");
    writeFileSync(file, "API_KEY=test123\nDATABASE_URL=postgres://localhost\n");

    const result = readEnvFile(file);
    expect(result.get("API_KEY")).toBe("test123");
    expect(result.get("DATABASE_URL")).toBe("postgres://localhost");
  });

  it("handles export prefix and quotes", () => {
    const dir = makeTmpDir();
    const file = join(dir, ".env");
    writeFileSync(file, 'export API_KEY="test with spaces"\n');

    const result = readEnvFile(file);
    expect(result.get("API_KEY")).toBe("test with spaces");
  });
});

describe("writeEnvFile", () => {
  it("writes content to file", () => {
    const dir = makeTmpDir();
    const file = join(dir, ".env");
    writeEnvFile(file, "API_KEY=test123\n");

    expect(readFileSync(file, "utf-8")).toBe("API_KEY=test123\n");
  });
});

describe("backupFile", () => {
  it("returns null for non-existent file", () => {
    expect(backupFile("/nonexistent/file.env")).toBeNull();
  });

  it("creates incremental backups", () => {
    const dir = makeTmpDir();
    const file = join(dir, "test.env");
    writeFileSync(file, "TEST=value");

    const backup1 = backupFile(file);
    expect(backup1).toBe(`${file}.backup.1`);
    expect(existsSync(backup1!)).toBe(true);
    expect(readFileSync(backup1!, "utf-8")).toBe("TEST=value");

    const backup2 = backupFile(file);
    expect(backup2).toBe(`${file}.backup.2`);
    expect(existsSync(backup2!)).toBe(true);
  });
});

describe("applyPreserve", () => {
  it("preserves existing values for specified keys", () => {
    const dir = makeTmpDir();
    const file = join(dir, ".env");
    writeFileSync(file, 'DATABASE_URL="keep-this"\nAPI_KEY=old-key\n');

    const entries: EnvEntry[] = [
      { key: "DATABASE_URL", value: "new-url" },
      { key: "API_KEY", value: "new-key" },
    ];

    const count = applyPreserve(entries, ["DATABASE_URL"], file);
    expect(count).toBe(1);
    expect(entries[0].value).toBe("keep-this");
    expect(entries[1].value).toBe("new-key");
  });

  it("preserves multiple variables", () => {
    const dir = makeTmpDir();
    const file = join(dir, ".env");
    writeFileSync(file, "DATABASE_URL=keep1\nAPI_KEY=keep2\nSECRET=old\n");

    const entries: EnvEntry[] = [
      { key: "DATABASE_URL", value: "new1" },
      { key: "API_KEY", value: "new2" },
      { key: "SECRET", value: "new3" },
    ];

    const count = applyPreserve(entries, ["DATABASE_URL", "API_KEY"], file);
    expect(count).toBe(2);
    expect(entries[0].value).toBe("keep1");
    expect(entries[1].value).toBe("keep2");
    expect(entries[2].value).toBe("new3");
  });

  it("keeps new value when key not in existing file", () => {
    const dir = makeTmpDir();
    const file = join(dir, ".env");
    writeFileSync(file, "API_KEY=old-key\n");

    const entries: EnvEntry[] = [
      { key: "DATABASE_URL", value: "new-url" },
      { key: "API_KEY", value: "new-key" },
    ];

    const count = applyPreserve(entries, ["DATABASE_URL", "API_KEY"], file);
    expect(count).toBe(1); // Only API_KEY existed
    expect(entries[0].value).toBe("new-url"); // Not in file, kept new
    expect(entries[1].value).toBe("old-key"); // In file, preserved
  });

  it("returns 0 when no keys to preserve", () => {
    const count = applyPreserve([], [], "/nonexistent");
    expect(count).toBe(0);
  });
});
