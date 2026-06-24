package dev.webbies.dotenvify.actions

import com.intellij.openapi.actionSystem.ActionUpdateThread
import com.intellij.openapi.actionSystem.AnAction
import com.intellij.openapi.actionSystem.AnActionEvent
import com.intellij.openapi.actionSystem.CommonDataKeys
import com.intellij.openapi.project.Project
import com.intellij.openapi.ui.DialogWrapper
import com.intellij.ui.components.JBScrollPane
import com.intellij.ui.components.JBTextArea
import dev.webbies.dotenvify.core.FormatOptions
import dev.webbies.dotenvify.ui.FormatOptionsPanel
import dev.webbies.dotenvify.ui.MONO_FONT
import java.awt.BorderLayout
import java.awt.Dimension
import javax.swing.JComponent
import javax.swing.JPanel

class ConvertSelectionAction : AnAction() {

    override fun actionPerformed(e: AnActionEvent) {
        val editor = e.getData(CommonDataKeys.EDITOR) ?: return
        val project = e.project ?: return
        ConvertCore.convertSelection(project, editor)
    }

    override fun getActionUpdateThread(): ActionUpdateThread = ActionUpdateThread.BGT

    override fun update(e: AnActionEvent) {
        val editor = e.getData(CommonDataKeys.EDITOR)
        e.presentation.isEnabledAndVisible = editor != null && editor.selectionModel.hasSelection()
    }
}

/** Preview dialog with format option toggles and live preview. */
class PreviewDialog(
    project: Project,
    entryCount: Int,
    private val okButtonText: String,
    private val formatter: (FormatOptions) -> String,
) : DialogWrapper(project) {

    private val previewArea = JBTextArea().apply { isEditable = false; font = MONO_FONT }
    private val optionsPanel = FormatOptionsPanel(project)

    init {
        title = "Preview ($entryCount entries)"
        setOKButtonText(okButtonText)
        init()
        updatePreview()
    }

    override fun createCenterPanel(): JComponent {
        optionsPanel.onChange { updatePreview() }

        return JPanel(BorderLayout(0, 8)).apply {
            add(optionsPanel, BorderLayout.NORTH)
            add(JBScrollPane(previewArea).apply { preferredSize = Dimension(600, 400) }, BorderLayout.CENTER)
        }
    }

    fun getFormattedOutput(): String = formatter(optionsPanel.options())

    private fun updatePreview() {
        previewArea.text = getFormattedOutput()
    }
}
