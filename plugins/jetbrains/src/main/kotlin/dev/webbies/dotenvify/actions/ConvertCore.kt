package dev.webbies.dotenvify.actions

import com.intellij.notification.NotificationType
import com.intellij.openapi.command.WriteCommandAction
import com.intellij.openapi.editor.Editor
import com.intellij.openapi.fileChooser.FileChooserFactory
import com.intellij.openapi.fileChooser.FileSaverDescriptor
import com.intellij.openapi.project.Project
import com.intellij.openapi.vfs.LocalFileSystem
import com.intellij.openapi.vfs.VirtualFile
import dev.webbies.dotenvify.core.DotEnvFormatter
import dev.webbies.dotenvify.core.DotEnvIO
import dev.webbies.dotenvify.core.DotEnvParser
import dev.webbies.dotenvify.ui.EnvFileApplicator
import java.nio.file.Path

/**
 * Conversion flows shared by the menu actions (target from the action's data context) and the
 * toolbar buttons (target from the active editor/file).
 */
object ConvertCore {

    fun convertSelection(project: Project, editor: Editor) {
        val selectedText = editor.selectionModel.selectedText
        if (selectedText.isNullOrBlank()) {
            EnvFileApplicator.notify(project, "Please select some text to convert.", NotificationType.WARNING)
            return
        }

        val parseResult = DotEnvParser.parse(selectedText)
        if (parseResult.entries.isEmpty()) {
            EnvFileApplicator.notify(project, "No key-value pairs found in selection.", NotificationType.WARNING)
            return
        }

        val dialog = PreviewDialog(project, parseResult.entries.size, okButtonText = "Replace selection") { options ->
            DotEnvFormatter.format(parseResult.entries, options)
        }
        if (dialog.showAndGet()) {
            WriteCommandAction.runWriteCommandAction(project) {
                editor.document.replaceString(
                    editor.selectionModel.selectionStart,
                    editor.selectionModel.selectionEnd,
                    dialog.getFormattedOutput().trimEnd('\n'),
                )
            }
        }
    }

    fun convertFile(project: Project, virtualFile: VirtualFile) {
        val content = String(virtualFile.contentsToByteArray(), Charsets.UTF_8)
        val parseResult = DotEnvParser.parse(content)

        if (parseResult.entries.isEmpty()) {
            EnvFileApplicator.notify(
                project,
                "No key-value pairs found in '${virtualFile.name}'.",
                NotificationType.WARNING,
            )
            return
        }

        val dialog = PreviewDialog(project, parseResult.entries.size, okButtonText = "Save .env…") { options ->
            DotEnvFormatter.format(parseResult.entries, options)
        }

        if (dialog.showAndGet()) {
            val output = dialog.getFormattedOutput()
            val descriptor = FileSaverDescriptor("Save .env file", "Choose where to save the .env file")
            val wrapper = FileChooserFactory.getInstance()
                .createSaveFileDialog(descriptor, project)
                .save(virtualFile.parent, EnvFileApplicator.defaultOutputPath(project)) ?: return

            val targetPath = Path.of(wrapper.file.absolutePath)
            DotEnvIO.writeEnvFile(targetPath, output, backup = true)
            LocalFileSystem.getInstance().refreshAndFindFileByNioFile(targetPath)
            EnvFileApplicator.notify(project, "Saved ${parseResult.entries.size} variables to ${wrapper.file.name}")
        }
    }
}
