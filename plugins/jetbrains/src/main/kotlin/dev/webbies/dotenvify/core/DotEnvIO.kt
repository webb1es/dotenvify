package dev.webbies.dotenvify.core

import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.attribute.PosixFilePermissions

/** Reads and writes .env files with backup and preserve support. */
object DotEnvIO {

    /** Reads and parses a .env file. Returns an empty list if the file does not exist. */
    fun readEnvFile(path: Path): List<EnvEntry> {
        if (!Files.exists(path)) return emptyList()
        return DotEnvParser.parse(Files.readString(path)).entries
    }

    /** Writes content to a .env file with optional backup. Sets POSIX 0600 permissions on Unix. */
    fun writeEnvFile(path: Path, content: String, backup: Boolean = true) {
        if (backup) backupFile(path)
        path.parent?.let { Files.createDirectories(it) }
        Files.writeString(path, content)
        try {
            Files.setPosixFilePermissions(path, PosixFilePermissions.fromString("rw-------"))
        } catch (_: UnsupportedOperationException) {
            // Windows: no POSIX permissions
        }
    }

    /** Merges entries, keeping existing values for keys listed in [preserveKeys]. */
    fun applyPreserve(
        newEntries: List<EnvEntry>,
        existingEntries: List<EnvEntry>,
        preserveKeys: Set<String>,
    ): List<EnvEntry> {
        if (preserveKeys.isEmpty()) return newEntries
        val existingMap = existingEntries.associate { it.key to it.value }
        return newEntries.map { entry ->
            if (entry.key in preserveKeys && entry.key in existingMap) {
                entry.copy(value = existingMap[entry.key]!!)
            } else {
                entry
            }
        }
    }

    /** Creates an incremental backup (.backup.1, .backup.2, etc.) of the file at [path]. */
    fun backupFile(path: Path) {
        if (!Files.exists(path)) return
        var counter = 1
        while (true) {
            val backupPath = path.resolveSibling("${path.fileName}.backup.$counter")
            if (!Files.exists(backupPath)) {
                Files.copy(path, backupPath)
                return
            }
            counter++
        }
    }
}
