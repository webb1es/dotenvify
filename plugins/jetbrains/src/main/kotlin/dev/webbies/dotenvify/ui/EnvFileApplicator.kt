package dev.webbies.dotenvify.ui

import com.intellij.notification.NotificationGroupManager
import com.intellij.notification.NotificationType
import com.intellij.openapi.project.Project
import com.intellij.openapi.vfs.LocalFileSystem
import dev.webbies.dotenvify.core.DotEnvFormatter
import dev.webbies.dotenvify.core.DotEnvIO
import dev.webbies.dotenvify.core.EnvEntry
import dev.webbies.dotenvify.core.FormatOptions
import dev.webbies.dotenvify.settings.DotEnvifyProjectSettings
import dev.webbies.dotenvify.settings.DotEnvifySettings
import java.nio.file.Path

/**
 * Shared utility for the apply-to-.env workflow:
 * read existing → show merge dialog if needed → write with backup → VFS refresh → notify.
 */
object EnvFileApplicator {

    fun apply(
        project: Project,
        entries: List<EnvEntry>,
        targetPath: Path,
        sourceName: String,
        options: FormatOptions = FormatOptions(),
    ) {
        val existingEntries = DotEnvIO.readEnvFile(targetPath)
        val incoming = DotEnvIO.applyPreserve(entries, existingEntries, preserveKeys(project))

        if (existingEntries.isNotEmpty()) {
            val dialog = EnvDiffDialog(project, existingEntries, incoming, sourceName)
            if (!dialog.showAndGet()) return
            val output = DotEnvFormatter.format(dialog.mergedEntries, options)
            DotEnvIO.writeEnvFile(targetPath, output, backup = true)
        } else {
            val output = DotEnvFormatter.format(incoming, options)
            DotEnvIO.writeEnvFile(targetPath, output, backup = true)
        }

        LocalFileSystem.getInstance().refreshAndFindFileByNioFile(targetPath)
        notify(project, "Saved ${entries.size} variables to ${targetPath.fileName}")
    }

    /**
     * The configured default .env file name for save dialogs — the project's output path
     * when project overrides are active, otherwise the global default. Falls back to `.env`.
     */
    fun defaultOutputPath(project: Project): String {
        val projectState = DotEnvifyProjectSettings.getInstance(project).state
        val configured = if (projectState.useGlobalDefaults) {
            DotEnvifySettings.getInstance().state.defaultOutputPath
        } else {
            projectState.outputPath
        }
        return configured.ifBlank { ".env" }
    }

    /** The project's preserve-keys set: keys whose existing values survive a merge. */
    private fun preserveKeys(project: Project): Set<String> =
        DotEnvifyProjectSettings.getInstance(project).state.preserveKeys
            .split(',')
            .map { it.trim() }
            .filter { it.isNotEmpty() }
            .toSet()

    fun notify(project: Project, message: String, type: NotificationType = NotificationType.INFORMATION) {
        NotificationGroupManager.getInstance()
            .getNotificationGroup("DotEnvify.Notifications")
            .createNotification(message, type)
            .notify(project)
    }
}
