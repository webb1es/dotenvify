import type { EnvEntry, ParseResult } from "./models.js";

/**
 * Parse messy environment variable input into structured entries.
 *
 * Supported formats (can be mixed):
 * - KEY=VALUE
 * - KEY="VALUE" or KEY='VALUE'
 * - export KEY=VALUE
 * - KEY VALUE (space-separated, single pair)
 * - Multiple KEY VALUE pairs on one line (even token count >= 4)
 * - Key on one line, value on the next
 */
export function parseDotEnv(input: string): ParseResult {
  const entries: EnvEntry[] = [];
  const seen = new Map<string, number>();
  const warnings: string[] = [];

  const rawLines = input.split("\n");
  // Trim and filter blank lines, but keep track of original indices for warnings
  const lines: { text: string; lineNum: number }[] = [];
  for (let i = 0; i < rawLines.length; i++) {
    const trimmed = rawLines[i].trim();
    if (trimmed !== "") {
      lines.push({ text: trimmed, lineNum: i + 1 });
    }
  }

  function addEntry(key: string, value: string) {
    if (!key) return;
    const idx = seen.get(key);
    if (idx !== undefined) {
      // Overwrite duplicate — last wins
      entries[idx] = { key, value };
    } else {
      seen.set(key, entries.length);
      entries.push({ key, value });
    }
  }

  let i = 0;
  while (i < lines.length) {
    const { text: line, lineNum } = lines[i];

    // Skip comments
    if (line.startsWith("#")) {
      i++;
      continue;
    }

    // Strip "export " prefix
    const stripped = line.startsWith("export ") ? line.slice(7) : line;

    // KEY=VALUE format
    if (stripped.includes("=")) {
      const eqIndex = stripped.indexOf("=");
      const key = stripped.slice(0, eqIndex).trim();
      let value = stripped.slice(eqIndex + 1).trim();
      value = unquote(value);
      addEntry(key, value);
      i++;
      continue;
    }

    // Space-separated: KEY VALUE on same line
    const tokens = stripped.split(/\s+/);
    if (tokens.length === 2) {
      addEntry(tokens[0], tokens[1]);
      i++;
      continue;
    }

    // Multiple KEY VALUE pairs on one line (even count >= 4)
    if (tokens.length >= 4 && tokens.length % 2 === 0) {
      for (let j = 0; j < tokens.length; j += 2) {
        addEntry(tokens[j], tokens[j + 1]);
      }
      i++;
      continue;
    }

    // Key on this line, value on next line
    if (i + 1 < lines.length) {
      addEntry(line, lines[i + 1].text);
      i += 2;
      continue;
    }

    // Orphan key — no value
    warnings.push(`Line ${lineNum}: Key '${line}' has no value`);
    i++;
  }

  return { entries, warnings };
}

function unquote(value: string): string {
  if (
    (value.startsWith('"') && value.endsWith('"')) ||
    (value.startsWith("'") && value.endsWith("'"))
  ) {
    return value.slice(1, -1);
  }
  return value;
}
