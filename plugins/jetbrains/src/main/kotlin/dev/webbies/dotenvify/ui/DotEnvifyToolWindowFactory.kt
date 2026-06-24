package dev.webbies.dotenvify.ui

import com.intellij.icons.AllIcons
import com.intellij.openapi.project.Project
import com.intellij.openapi.util.Disposer
import com.intellij.openapi.wm.ToolWindow
import com.intellij.openapi.wm.ToolWindowFactory
import com.intellij.ui.content.ContentFactory
import dev.webbies.dotenvify.actions.ConvertActiveFileAction
import dev.webbies.dotenvify.actions.ConvertActiveSelectionAction
import dev.webbies.dotenvify.azure.AzureVariableGroupPanel
import dev.webbies.dotenvify.diagnostics.DiagnosticsPanel

class DotEnvifyToolWindowFactory : ToolWindowFactory {

    override fun createToolWindowContent(project: Project, toolWindow: ToolWindow) {
        val contentFactory = ContentFactory.getInstance()

        // Azure first: it's the primary feature
        val azureTab = contentFactory.createContent(AzureVariableGroupPanel(project), "Azure DevOps", false).apply {
            icon = AllIcons.Providers.Azure
        }
        val convertPanel = DotEnvifyToolWindowPanel(project)
        val convertTab = contentFactory.createContent(convertPanel, "Paste & Convert", false).apply {
            icon = AllIcons.Actions.RealIntentionBulb
        }
        // Dispose panels with their content so listeners/timers release on unload.
        Disposer.register(convertTab, convertPanel)

        val diagnosticsPanel = DiagnosticsPanel(project)
        val diagnosticsTab = contentFactory.createContent(diagnosticsPanel, "Diagnostics", false).apply {
            icon = AllIcons.Actions.Find
        }
        Disposer.register(diagnosticsTab, diagnosticsPanel)

        toolWindow.contentManager.addContent(azureTab)
        toolWindow.contentManager.addContent(convertTab)
        toolWindow.contentManager.addContent(diagnosticsTab)

        // Convert actions as toolbar buttons in the tool-window header
        toolWindow.setTitleActions(
            listOf(
                ConvertActiveSelectionAction(project),
                ConvertActiveFileAction(project),
            )
        )
    }
}
