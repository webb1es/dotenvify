package dev.webbies.dotenvify.ui

import com.intellij.openapi.project.Project
import dev.webbies.dotenvify.core.FormatOptions
import dev.webbies.dotenvify.settings.DotEnvifyProjectSettings
import dev.webbies.dotenvify.settings.FormatOptionLabels
import dev.webbies.dotenvify.settings.DotEnvifySettings
import javax.swing.*

/**
 * Panel of format-option checkboxes.
 * Initialises from saved settings (project or global).
 */
class FormatOptionsPanel(project: Project? = null) : JPanel() {

    private val exportCheckbox = JCheckBox(FormatOptionLabels.EXPORT)
    private val sortCheckbox = JCheckBox(FormatOptionLabels.SORT, true)
    private val noLowerCheckbox = JCheckBox(FormatOptionLabels.SKIP_LOWERCASE)
    private val urlOnlyCheckbox = JCheckBox(FormatOptionLabels.URL_ONLY)
    private val checkboxes = listOf(exportCheckbox, sortCheckbox, noLowerCheckbox, urlOnlyCheckbox)

    init {
        layout = BoxLayout(this, BoxLayout.X_AXIS)
        border = BorderFactory.createEmptyBorder(4, 4, 4, 4)
        checkboxes.forEachIndexed { i, cb ->
            if (i > 0) add(Box.createHorizontalStrut(8))
            add(cb)
        }
        loadFromSettings(project)
    }

    fun options() = FormatOptions(
        exportPrefix = exportCheckbox.isSelected,
        sort = sortCheckbox.isSelected,
        ignoreLowercase = noLowerCheckbox.isSelected,
        urlOnly = urlOnlyCheckbox.isSelected,
    )

    fun onChange(listener: () -> Unit) {
        val itemListener = java.awt.event.ItemListener { listener() }
        checkboxes.forEach { it.addItemListener(itemListener) }
    }

    private fun loadFromSettings(project: Project?) {
        if (project != null) {
            val projectSettings = DotEnvifyProjectSettings.getInstance(project).state
            if (!projectSettings.useGlobalDefaults) {
                exportCheckbox.isSelected = projectSettings.exportPrefix
                sortCheckbox.isSelected = projectSettings.sort
                noLowerCheckbox.isSelected = projectSettings.ignoreLowercase
                urlOnlyCheckbox.isSelected = projectSettings.urlOnly
                return
            }
        }
        val global = DotEnvifySettings.getInstance().state
        exportCheckbox.isSelected = global.exportPrefix
        sortCheckbox.isSelected = global.sort
        noLowerCheckbox.isSelected = global.ignoreLowercase
        urlOnlyCheckbox.isSelected = global.urlOnly
    }
}
