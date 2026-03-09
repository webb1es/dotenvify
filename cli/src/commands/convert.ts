import { readFileSync, existsSync } from "node:fs";
import {
  parseDotEnv,
  formatDotEnv,
  writeEnvFile,
  backupFile,
  applyPreserve,
} from "@dotenvify/core";
import type { FormatOptions } from "@dotenvify/core";

export interface ConvertOptions {
  output: string;
  export: boolean;
  noSort: boolean;
  noLower: boolean;
  urlOnly: boolean;
  overwrite: boolean;
  preserve: string;
}

export function convert(sourceFile: string, opts: ConvertOptions): void {
  // Validate source file
  if (!existsSync(sourceFile)) {
    console.error(`Error: Source file '${sourceFile}' does not exist`);
    process.exit(1);
  }

  // Validate paths
  if (sourceFile.includes("..") || opts.output.includes("..")) {
    console.error("Error: File paths cannot contain '..'");
    process.exit(1);
  }

  // Read and parse source
  const content = readFileSync(sourceFile, "utf-8");
  const { entries, warnings } = parseDotEnv(content);

  if (warnings.length > 0) {
    for (const w of warnings) {
      console.warn(`Warning: ${w}`);
    }
  }

  // Apply preserve logic before writing
  const preserveKeys = opts.preserve
    ? opts.preserve.split(",").map((k) => k.trim()).filter(Boolean)
    : [];

  if (preserveKeys.length > 0) {
    const count = applyPreserve(entries, preserveKeys, opts.output);
    if (count > 0) {
      console.log(`Preserved ${count} existing variable(s)`);
    }
  }

  // Backup unless overwrite
  if (!opts.overwrite) {
    const backupPath = backupFile(opts.output);
    if (backupPath) {
      console.log(`Backed up existing file to '${backupPath}'`);
    }
  }

  // Format
  const formatOpts: FormatOptions = {
    export: opts.export,
    sort: !opts.noSort,
    noLower: opts.noLower,
    urlOnly: opts.urlOnly,
  };

  const formatted = formatDotEnv(entries, formatOpts);

  // Write
  writeEnvFile(opts.output, formatted);
  const lineCount = formatted ? formatted.trimEnd().split("\n").length : 0;
  console.log(`Written ${lineCount} variable(s) to '${opts.output}'`);
}
