import { readFileSync, writeFileSync, existsSync, statSync, copyFileSync } from "node:fs";
import { parseDotEnv } from "./parser.js";
import type { EnvEntry } from "./models.js";

/**
 * Read an existing .env file and return entries as a key-value map.
 * Returns an empty map if the file doesn't exist.
 */
export function readEnvFile(filePath: string): Map<string, string> {
  if (!existsSync(filePath)) return new Map();

  const content = readFileSync(filePath, "utf-8");
  const { entries } = parseDotEnv(content);
  return new Map(entries.map((e) => [e.key, e.value]));
}

/**
 * Write formatted content to a file with secure permissions (0o600).
 */
export function writeEnvFile(filePath: string, content: string): void {
  writeFileSync(filePath, content, { mode: 0o600 });
}

/**
 * Create an incremental backup of a file (file.backup.1, file.backup.2, ...).
 * No-op if the file doesn't exist. Returns the backup path or null.
 */
export function backupFile(filePath: string): string | null {
  if (!existsSync(filePath)) return null;

  let counter = 1;
  let backupPath: string;
  do {
    backupPath = `${filePath}.backup.${counter}`;
    counter++;
  } while (existsSync(backupPath));

  copyFileSync(filePath, backupPath);
  return backupPath;
}

/**
 * Apply preserve logic: for the given keys, replace their values in `entries`
 * with existing values from the output file (if they exist there).
 *
 * Returns the number of preserved entries.
 */
export function applyPreserve(
  entries: EnvEntry[],
  preserveKeys: string[],
  outputFilePath: string,
): number {
  if (preserveKeys.length === 0) return 0;

  const existing = readEnvFile(outputFilePath);
  let count = 0;

  for (const entry of entries) {
    if (preserveKeys.includes(entry.key) && existing.has(entry.key)) {
      entry.value = existing.get(entry.key)!;
      count++;
    }
  }

  return count;
}
