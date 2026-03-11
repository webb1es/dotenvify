package dev.webbies.dotenvify.core

/**
 * Parses raw key-value text in multiple formats into a list of [EnvEntry].
 *
 * Supported: KEY=VALUE, KEY="VALUE", KEY VALUE (space-separated),
 * line-pair (KEY then VALUE on next line), multiple pairs per line.
 * Comments (#) and blank lines are skipped.
 */
object DotEnvParser {

    private val WHITESPACE = "\\s+".toRegex()

    /** Parses raw text into a [ParseResult] containing entries, warnings, and format detection. */
    fun parse(input: String): ParseResult {
        if (input.isBlank()) return ParseResult(emptyList())

        val lines = input.lines()
            .map { it.trim() }
            .filter { it.isNotEmpty() && !it.startsWith("#") }

        val entries = mutableListOf<EnvEntry>()
        val warnings = mutableListOf<String>()

        var i = 0
        while (i < lines.size) {
            val line = lines[i]
            val stripped = line.removePrefix("export ").trimStart()

            // Format: KEY=VALUE or export KEY=VALUE
            if (stripped.contains('=')) {
                val eqIdx = stripped.indexOf('=')
                val key = stripped.substring(0, eqIdx).trim()
                if (key.isNotEmpty()) {
                    entries.add(EnvEntry(key, unquote(stripped.substring(eqIdx + 1).trim())))
                }
                i++
                continue
            }

            val parts = line.split(WHITESPACE)

            // Format: KEY VALUE (space-separated pair)
            if (parts.size == 2) {
                entries.add(EnvEntry(parts[0], unquote(parts[1])))
                i++
                continue
            }

            // Format: KEY1 VAL1 KEY2 VAL2 ... (multiple space-separated pairs on one line)
            if (parts.size >= 4 && parts.size % 2 == 0) {
                for (j in parts.indices step 2) {
                    if (parts[j].isNotEmpty()) entries.add(EnvEntry(parts[j], unquote(parts[j + 1])))
                }
                i++
                continue
            }

            // Format: line pair (KEY on this line, VALUE on the next)
            if (i + 1 < lines.size) {
                entries.add(EnvEntry(line, unquote(lines[i + 1])))
                i += 2
            } else {
                warnings.add("Key '$line' has no value (end of input)")
                i++
            }
        }

        val alreadyFormatted = lines.isNotEmpty() && lines.all { line ->
            val s = line.removePrefix("export ").trimStart()
            s.contains('=')
        }

        return ParseResult(entries, warnings, alreadyFormatted)
    }

    /** Strips matching single or double quotes from a value, if present. */
    fun unquote(value: String): String {
        if (value.length >= 2 &&
            ((value.startsWith('"') && value.endsWith('"')) ||
                    (value.startsWith('\'') && value.endsWith('\'')))
        ) {
            return value.substring(1, value.length - 1)
        }
        return value
    }
}
