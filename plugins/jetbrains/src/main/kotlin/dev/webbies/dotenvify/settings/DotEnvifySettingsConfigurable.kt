package dev.webbies.dotenvify.settings

import com.intellij.openapi.options.Configurable
import com.intellij.util.ui.FormBuilder
import javax.swing.*

class DotEnvifySettingsConfigurable : Configurable {

    private val exportCheckbox = JCheckBox(FormatOptionLabels.EXPORT)
    private val sortCheckbox = JCheckBox(FormatOptionLabels.SORT, true)
    private val noLowerCheckbox = JCheckBox(FormatOptionLabels.SKIP_LOWERCASE)
    private val urlOnlyCheckbox = JCheckBox(FormatOptionLabels.URL_ONLY)
    private val outputPathField = JTextField(".env", 20)
    private val azureOrgUrlField = JTextField("", 30).apply {
        toolTipText = "e.g. https://dev.azure.com/myorg/myproject (shared across all projects)"
    }
    private val azCliPathField = JTextField("", 30).apply {
        toolTipText = "Leave blank to auto-detect. Set if 'az' is installed in a non-standard location."
    }

    override fun getDisplayName(): String = "DotEnvify"

    override fun createComponent(): JComponent {
        reset()
        return FormBuilder.createFormBuilder()
            .addSeparator()
            .addComponent(JLabel("Format Options"))
            .addComponent(exportCheckbox)
            .addComponent(sortCheckbox)
            .addComponent(noLowerCheckbox)
            .addComponent(urlOnlyCheckbox)
            .addSeparator()
            .addLabeledComponent("Default output path:", outputPathField)
            .addSeparator()
            .addComponent(JLabel("Azure DevOps"))
            .addLabeledComponent("Organization URL:", azureOrgUrlField)
            .addLabeledComponent("Azure CLI path (optional):", azCliPathField)
            .addComponentFillVertically(JPanel(), 0)
            .panel
    }

    override fun isModified(): Boolean {
        val s = DotEnvifySettings.getInstance().state
        return exportCheckbox.isSelected != s.exportPrefix ||
                sortCheckbox.isSelected != s.sort ||
                noLowerCheckbox.isSelected != s.ignoreLowercase ||
                urlOnlyCheckbox.isSelected != s.urlOnly ||
                outputPathField.text != s.defaultOutputPath ||
                azureOrgUrlField.text != s.azureOrgUrl ||
                azCliPathField.text != s.azCliPath
    }

    override fun apply() {
        val s = DotEnvifySettings.getInstance().state
        s.exportPrefix = exportCheckbox.isSelected
        s.sort = sortCheckbox.isSelected
        s.ignoreLowercase = noLowerCheckbox.isSelected
        s.urlOnly = urlOnlyCheckbox.isSelected
        s.defaultOutputPath = outputPathField.text
        s.azureOrgUrl = azureOrgUrlField.text.trim()
        s.azCliPath = azCliPathField.text.trim()
    }

    override fun reset() {
        val s = DotEnvifySettings.getInstance().state
        exportCheckbox.isSelected = s.exportPrefix
        sortCheckbox.isSelected = s.sort
        noLowerCheckbox.isSelected = s.ignoreLowercase
        urlOnlyCheckbox.isSelected = s.urlOnly
        outputPathField.text = s.defaultOutputPath
        azureOrgUrlField.text = s.azureOrgUrl
        azCliPathField.text = s.azCliPath
    }
}
