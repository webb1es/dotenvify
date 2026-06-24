package dev.webbies.dotenvify.actions

import com.intellij.openapi.actionSystem.ActionUpdateThread
import com.intellij.openapi.actionSystem.AnAction
import com.intellij.openapi.actionSystem.AnActionEvent
import com.intellij.openapi.fileEditor.FileEditorManager
import com.intellij.openapi.project.Project
import dev.webbies.dotenvify.ui.DotEnvifyIcons

/**
 * Toolbar buttons that act on the active editor/file, so they work while focus is in the
 * tool window (unlike the data-context-bound menu actions).
 */
class ConvertActiveSelectionAction(private val project: Project) : AnAction(
    "Convert Selection to .env",
    "Convert the selected text in the active editor to .env format",
    DotEnvifyIcons.ConvertSelection,
) {
    override fun actionPerformed(e: AnActionEvent) {
        val editor = FileEditorManager.getInstance(project).selectedTextEditor ?: return
        ConvertCore.convertSelection(project, editor)
    }

    override fun update(e: AnActionEvent) {
        val editor = FileEditorManager.getInstance(project).selectedTextEditor
        e.presentation.isEnabled = editor?.selectionModel?.hasSelection() == true
    }

    override fun getActionUpdateThread(): ActionUpdateThread = ActionUpdateThread.EDT
}

class ConvertActiveFileAction(private val project: Project) : AnAction(
    "Convert File to .env",
    "Convert the active file to .env format",
    DotEnvifyIcons.ConvertFile,
) {
    override fun actionPerformed(e: AnActionEvent) {
        val file = FileEditorManager.getInstance(project).selectedFiles.firstOrNull() ?: return
        ConvertCore.convertFile(project, file)
    }

    override fun update(e: AnActionEvent) {
        e.presentation.isEnabled = FileEditorManager.getInstance(project).selectedFiles.isNotEmpty()
    }

    override fun getActionUpdateThread(): ActionUpdateThread = ActionUpdateThread.EDT
}
