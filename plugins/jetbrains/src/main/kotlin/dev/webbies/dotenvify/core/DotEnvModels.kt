package dev.webbies.dotenvify.core

/** A single key-value pair from a parsed .env source. */
data class EnvEntry(
    val key: String,
    val value: String,
)

/** Output of [DotEnvParser.parse], containing parsed entries and any parse warnings. */
data class ParseResult(
    val entries: List<EnvEntry>,
    val warnings: List<String> = emptyList(),
    /** True if every line was already in KEY=VALUE format. */
    val alreadyFormatted: Boolean = false,
)

/** Controls how [DotEnvFormatter] transforms entries into .env output. */
data class FormatOptions(
    /** Prefix each line with `export `. */
    val exportPrefix: Boolean = false,
    /** Sort keys alphabetically. */
    val sort: Boolean = true,
    /** Exclude keys that contain lowercase letters. */
    val ignoreLowercase: Boolean = true,
    /** Only include entries whose values are URLs. */
    val urlOnly: Boolean = false,
    /** Keys whose existing values should not be overwritten on merge. */
    val preserveKeys: Set<String> = emptySet(),
)
