package dev.webbies.dotenvify.ui

import com.intellij.notification.NotificationGroupManager
import com.intellij.notification.NotificationType
import com.intellij.openapi.application.WriteAction
import com.intellij.openapi.fileEditor.FileDocumentManager
import com.intellij.openapi.project.Project
import com.intellij.openapi.ui.Messages
import com.intellij.openapi.vfs.LocalFileSystem
import dev.webbies.dotenvify.core.DotEnvFormatter
import dev.webbies.dotenvify.core.DotEnvIO
import dev.webbies.dotenvify.core.EnvEntry
import dev.webbies.dotenvify.core.FormatOptions
import dev.webbies.dotenvify.settings.DotEnvifyProjectSettings
import dev.webbies.dotenvify.settings.DotEnvifySettings
import java.nio.file.Path

/**
 * Drives the apply-to-.env workflow:
 * read existing, show merge dialog if needed, write with backup, refresh VFS, notify.
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

        val dialog = EnvDiffDialog(project, existingEntries, incoming, sourceName)
        if (!dialog.showAndGet()) return
        val merged = dialog.mergedEntries

        writeAndRefresh(targetPath, DotEnvFormatter.format(merged, options))
        notifySaved(project, merged.size, targetPath)
    }

    /** Confirms a successful save with a modal so it can't be missed. */
    fun notifySaved(project: Project, count: Int, targetPath: Path) {
        Messages.showInfoMessage(
            project,
            "Saved $count variables to ${targetPath.fileName}.",
            "Saved to .env",
        )
    }

    /**
     * Writes [output] to [targetPath] (with backup), refreshes the VFS entry, and reloads
     * an open editor for the file so it reflects the new content even if it had unsaved edits.
     */
    fun writeAndRefresh(targetPath: Path, output: String) {
        DotEnvIO.writeEnvFile(targetPath, output, backup = true)
        // A synchronous VFS refresh on the EDT requires the write lock.
        WriteAction.runAndWait<Throwable> {
            val vf = LocalFileSystem.getInstance().refreshAndFindFileByNioFile(targetPath)
            val doc = vf?.let { FileDocumentManager.getInstance().getCachedDocument(it) }
            if (doc != null) FileDocumentManager.getInstance().reloadFromDisk(doc)
        }
    }

    /**
     * The default .env file name for save dialogs: the project output path when overrides
     * are active, otherwise the global default. Falls back to `.env`.
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
