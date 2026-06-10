package dev.webbies.dotenvify.azure

import com.intellij.icons.AllIcons
import com.intellij.ide.BrowserUtil
import com.intellij.notification.NotificationType
import com.intellij.openapi.application.ApplicationManager
import com.intellij.openapi.progress.ProgressIndicator
import com.intellij.openapi.progress.ProgressManager
import com.intellij.openapi.progress.Task
import com.intellij.openapi.project.Project
import com.intellij.openapi.ui.ComboBox
import com.intellij.openapi.ui.Messages
import com.intellij.ui.JBColor
import com.intellij.ui.JBSplitter
import com.intellij.ui.components.JBScrollPane
import com.intellij.ui.components.JBTextArea
import com.intellij.ui.table.JBTable
import com.intellij.util.ui.JBUI
import dev.webbies.dotenvify.core.DotEnvFormatter
import dev.webbies.dotenvify.core.DotEnvIO
import dev.webbies.dotenvify.core.DotEnvParser
import dev.webbies.dotenvify.core.EnvEntry
import dev.webbies.dotenvify.settings.DotEnvifyProjectSettings
import dev.webbies.dotenvify.settings.DotEnvifySettings
import dev.webbies.dotenvify.ui.EnvFileApplicator
import dev.webbies.dotenvify.ui.FormatOptionsPanel
import dev.webbies.dotenvify.ui.MONO_FONT
import java.awt.*
import java.awt.datatransfer.StringSelection
import java.net.ConnectException
import java.net.http.HttpTimeoutException
import java.nio.file.Path
import javax.swing.*
import javax.swing.table.AbstractTableModel
import javax.swing.table.DefaultTableCellRenderer

class AzureVariableGroupPanel(private val project: Project) : JPanel(BorderLayout()) {

    // --- Connection fields ---
    private val orgUrlField = JTextField().apply {
        toolTipText = "e.g. https://dev.azure.com/myorg/myproject"
    }
    private val groupCombo = ComboBox<String>().apply {
        isEditable = true
        toolTipText = "Pick a variable group, or type its name. Use ↻ to load the list."
    }
    private val loadGroupsButton = JButton(AllIcons.Actions.Refresh).apply {
        toolTipText = "Load variable groups for this project"
    }
    private val targetCombo = ComboBox(
        arrayOf(".env", ".env.local", ".env.development", ".env.staging", ".env.production")
    ).apply {
        isEditable = true
        toolTipText = "File to write (relative to project root, or an absolute path)"
    }

    // --- Auth controls (Azure CLI / `az login`) ---
    private val connectButton = JButton("Connect").apply {
        icon = AllIcons.Actions.Refresh
        toolTipText = "Check the local Azure CLI (az) session"
    }
    private val loginButton = JButton("Run az login").apply {
        icon = AllIcons.Actions.Execute
        isVisible = false
        toolTipText = "Sign in to Azure in your browser via the Azure CLI"
    }
    private val switchAccountButton = JButton("Switch account").apply {
        isVisible = false
        toolTipText = "Sign in as a different account (re-runs az login)"
    }
    private val signOutButton = JButton("Sign out").apply {
        isVisible = false
        toolTipText = "Run 'az logout' (ends your Azure CLI session everywhere)"
    }
    private val authStatusIcon = JLabel(AllIcons.General.InspectionsError)
    private val authStatusLabel = JLabel("Not connected")

    // --- Action controls ---
    private val fetchButton = JButton("Fetch Variables").apply { isEnabled = false; icon = AllIcons.Actions.Download }
    private val applyButton = JButton("Apply to .env").apply { isEnabled = false; icon = AllIcons.Actions.MenuSaveall }
    private val copyButton = JButton("Copy").apply { isEnabled = false; icon = AllIcons.Actions.Copy }
    private val openInBrowserButton = JButton(AllIcons.Ide.External_link_arrow).apply {
        isVisible = false
        toolTipText = "Open this variable group in Azure DevOps"
    }

    // --- Format options ---
    private val optionsPanel = FormatOptionsPanel(project)

    // --- Variable table ---
    private val tableModel = VarTableModel()
    private val varTable = JBTable(tableModel).apply {
        setShowGrid(false)
        rowHeight = JBUI.scale(22)
        autoResizeMode = JTable.AUTO_RESIZE_LAST_COLUMN
    }
    private val groupInfoLabel = JLabel(" ").apply {
        font = font.deriveFont(Font.ITALIC)
        foreground = JBColor.GRAY
    }
    private val previewArea = JBTextArea().apply { isEditable = false; font = MONO_FONT; lineWrap = false }
    private val statusLabel = JLabel(" ")

    /** Last-seen format options, to detect when a *filter* toggle (vs sort/export) changed. */
    private var lastOptions = optionsPanel.options()

    private var currentGroupId: Int = 0

    /** Current contents of the target .env file (unquoted), for live diff status. */
    private var targetMap: Map<String, String> = emptyMap()

    init {
        border = JBUI.Borders.empty(8)

        val connectionPanel = buildConnectionPanel().apply { alignmentX = LEFT_ALIGNMENT }
        val optionsRow = JPanel(BorderLayout()).apply {
            border = JBUI.Borders.emptyTop(4)
            add(optionsPanel, BorderLayout.WEST)
            alignmentX = LEFT_ALIGNMENT
        }
        val topPanel = JPanel().apply {
            layout = BoxLayout(this, BoxLayout.Y_AXIS)
            add(connectionPanel)
            add(optionsRow)
        }

        // === CENTER: table (top) + live preview (bottom), resizable ===
        val tableHeader = JPanel(BorderLayout()).apply {
            add(groupInfoLabel, BorderLayout.CENTER)
            add(openInBrowserButton, BorderLayout.EAST)
        }
        val tablePanel = JPanel(BorderLayout(0, 4)).apply {
            border = BorderFactory.createTitledBorder("Variables")
            add(tableHeader, BorderLayout.NORTH)
            add(JBScrollPane(varTable), BorderLayout.CENTER)
        }
        val previewPanel = JPanel(BorderLayout()).apply {
            border = BorderFactory.createTitledBorder("Preview (.env output)")
            add(JBScrollPane(previewArea), BorderLayout.CENTER)
        }
        val splitter = JBSplitter(true, 0.62f).apply {
            firstComponent = tablePanel
            secondComponent = previewPanel
        }

        // === BOTTOM: Action buttons + status ===
        val bottomRow = JPanel(BorderLayout()).apply {
            border = JBUI.Borders.emptyTop(4)
            val actions = JPanel(FlowLayout(FlowLayout.LEFT, 4, 0)).apply {
                add(applyButton)
                add(copyButton)
            }
            add(actions, BorderLayout.WEST)
            add(statusLabel, BorderLayout.EAST)
        }

        add(topPanel, BorderLayout.NORTH)
        add(splitter, BorderLayout.CENTER)
        add(bottomRow, BorderLayout.SOUTH)

        configureTable()

        // --- Wire events ---
        connectButton.addActionListener { updateAuthState() }
        loginButton.addActionListener { runAzLogin() }
        switchAccountButton.addActionListener { runAzLogin() }
        signOutButton.addActionListener { signOut() }
        loadGroupsButton.addActionListener { loadGroups() }
        fetchButton.addActionListener { fetchVariables() }
        applyButton.addActionListener { applyToFile() }
        copyButton.addActionListener { copyToClipboard() }
        openInBrowserButton.addActionListener { openInBrowser() }

        groupCombo.addActionListener { persistGroup() }
        targetCombo.addActionListener {
            persistTarget()
            recomputeTargetDiff() // diff status is relative to the chosen file
        }
        optionsPanel.onChange { onOptionsChanged() }
        tableModel.addTableModelListener { refreshDerived() }

        loadPersistedFields()
        updateAuthState()
    }

    private fun buildConnectionPanel(): JComponent {
        val panel = JPanel(GridBagLayout())
        val gbc = GridBagConstraints().apply {
            anchor = GridBagConstraints.WEST
            insets = JBUI.insets(0, 0, 4, 8)
        }

        // Row 0: auth controls (full width)
        gbc.gridx = 0; gbc.gridy = 0; gbc.gridwidth = 3
        gbc.fill = GridBagConstraints.HORIZONTAL; gbc.weightx = 1.0
        panel.add(
            JPanel(FlowLayout(FlowLayout.LEFT, 4, 0)).apply {
                add(authStatusIcon)
                add(authStatusLabel)
                add(connectButton)
                add(loginButton)
                add(switchAccountButton)
                add(signOutButton)
            },
            gbc,
        )

        // Row 1: Azure DevOps URL
        gbc.gridy = 1; gbc.gridx = 0; gbc.gridwidth = 1
        gbc.fill = GridBagConstraints.NONE; gbc.weightx = 0.0
        panel.add(JLabel("Azure DevOps URL:"), gbc)
        gbc.gridx = 1; gbc.gridwidth = 2
        gbc.fill = GridBagConstraints.HORIZONTAL; gbc.weightx = 1.0
        panel.add(orgUrlField, gbc)

        // Row 2: Variable group picker + load + fetch
        gbc.gridy = 2; gbc.gridx = 0; gbc.gridwidth = 1
        gbc.fill = GridBagConstraints.NONE; gbc.weightx = 0.0
        panel.add(JLabel("Variable Group:"), gbc)
        gbc.gridx = 1; gbc.gridwidth = 1
        gbc.fill = GridBagConstraints.HORIZONTAL; gbc.weightx = 1.0
        panel.add(groupCombo, gbc)
        gbc.gridx = 2; gbc.gridwidth = 1
        gbc.fill = GridBagConstraints.NONE; gbc.weightx = 0.0
        panel.add(
            JPanel(FlowLayout(FlowLayout.LEFT, 4, 0)).apply {
                add(loadGroupsButton)
                add(fetchButton)
            },
            gbc,
        )

        // Row 3: Output file
        gbc.gridy = 3; gbc.gridx = 0; gbc.gridwidth = 1
        gbc.fill = GridBagConstraints.NONE; gbc.weightx = 0.0
        panel.add(JLabel("Output file:"), gbc)
        gbc.gridx = 1; gbc.gridwidth = 2
        gbc.fill = GridBagConstraints.HORIZONTAL; gbc.weightx = 1.0
        panel.add(targetCombo, gbc)

        return panel
    }

    // --- Persistence ---

    private fun loadPersistedFields() {
        val globalState = DotEnvifySettings.getInstance().state
        if (globalState.azureOrgUrl.isNotEmpty()) orgUrlField.text = globalState.azureOrgUrl

        val projectState = DotEnvifyProjectSettings.getInstance(project).state
        if (projectState.azureGroupName.isNotEmpty()) groupCombo.editor.item = projectState.azureGroupName
        targetCombo.editor.item = projectState.azureTargetFile.ifBlank { ".env" }

        orgUrlField.document.addDocumentListener(SimpleDocListener {
            DotEnvifySettings.getInstance().state.azureOrgUrl = orgUrlField.text.trim()
        })
    }

    private fun persistGroup() {
        DotEnvifyProjectSettings.getInstance(project).state.azureGroupName = currentGroupName()
    }

    private fun persistTarget() {
        DotEnvifyProjectSettings.getInstance(project).state.azureTargetFile = currentTargetFile()
    }

    private fun currentGroupName(): String =
        (groupCombo.editor.item as? String ?: groupCombo.selectedItem as? String ?: "").trim()

    private fun currentTargetFile(): String =
        (targetCombo.editor.item as? String ?: targetCombo.selectedItem as? String ?: "").trim()

    private fun resolvedTargetPath(): Path? {
        val target = currentTargetFile().ifBlank { ".env" }
        val p = Path.of(target)
        if (p.isAbsolute) return p
        val base = project.basePath ?: return null
        return Path.of(base, target)
    }

    // --- Auth ---

    private fun updateAuthState() {
        connectButton.isEnabled = false
        // Shells out to `az` — never run on the EDT.
        ApplicationManager.getApplication().executeOnPooledThread {
            val status = AzureCliAuthProvider.status()
            ApplicationManager.getApplication().invokeLater { applyAuthState(status) }
        }
    }

    private fun applyAuthState(status: AzureCliAuthProvider.Status) {
        connectButton.isEnabled = true
        val signedIn = status is AzureCliAuthProvider.Status.SignedIn
        fetchButton.isEnabled = signedIn
        loadGroupsButton.isEnabled = signedIn
        loginButton.isVisible = status is AzureCliAuthProvider.Status.NotLoggedIn
        switchAccountButton.isVisible = signedIn
        signOutButton.isVisible = signedIn

        when (status) {
            is AzureCliAuthProvider.Status.SignedIn -> {
                authStatusIcon.icon = AllIcons.General.InspectionsOK
                authStatusLabel.text = "Connected: ${status.user}"
                authStatusLabel.toolTipText = status.subscription.ifEmpty { null }
                authStatusLabel.foreground = JBColor(Color(0, 128, 0), Color(80, 200, 80))
            }
            AzureCliAuthProvider.Status.NotLoggedIn -> {
                authStatusIcon.icon = AllIcons.General.Warning
                authStatusLabel.text = "Not signed in"
                authStatusLabel.toolTipText = "Run 'az login' (or click 'Run az login')."
                authStatusLabel.foreground = JBColor.GRAY
            }
            AzureCliAuthProvider.Status.NotInstalled -> {
                authStatusIcon.icon = AllIcons.General.InspectionsError
                authStatusLabel.text = "Azure CLI not found"
                authStatusLabel.toolTipText =
                    "Install the Azure CLI (https://aka.ms/azcli) or set its path in Settings → DotEnvify."
                authStatusLabel.foreground = JBColor.GRAY
            }
            is AzureCliAuthProvider.Status.Failed -> {
                authStatusIcon.icon = AllIcons.General.InspectionsError
                authStatusLabel.text = "Azure CLI error"
                authStatusLabel.toolTipText = status.message
                authStatusLabel.foreground = JBColor.GRAY
            }
        }
    }

    private fun runAzLogin() {
        loginButton.isEnabled = false
        switchAccountButton.isEnabled = false
        ProgressManager.getInstance()
            .run(object : Task.Backgroundable(project, "Signing in to Azure (browser)...", true) {
                override fun run(indicator: ProgressIndicator) {
                    indicator.text = "Complete sign-in in your browser..."
                    val result = AzureCliAuthProvider.login(indicator)
                    ApplicationManager.getApplication().invokeLater {
                        loginButton.isEnabled = true
                        switchAccountButton.isEnabled = true
                        when (result) {
                            is AzureCliAuthProvider.TokenResult.Success -> {
                                EnvFileApplicator.notify(project, "Signed in to Azure")
                                updateAuthState()
                            }
                            AzureCliAuthProvider.TokenResult.NotInstalled ->
                                showError("Azure CLI ('az') not found. Install it or set its path in Settings → DotEnvify.")
                            AzureCliAuthProvider.TokenResult.NotLoggedIn ->
                                showError("Sign-in did not complete. Try again.")
                            is AzureCliAuthProvider.TokenResult.Failed ->
                                showError("Azure sign-in failed: ${result.message}")
                        }
                    }
                }
            })
    }

    private fun signOut() {
        val confirm = Messages.showYesNoDialog(
            project,
            "This runs 'az logout' and ends your Azure CLI session for ALL tools on this machine, " +
                "not just DotEnvify. Continue?",
            "Sign out of Azure CLI",
            "Sign out", "Cancel", Messages.getWarningIcon(),
        )
        if (confirm != Messages.YES) return

        ProgressManager.getInstance().run(object : Task.Backgroundable(project, "Signing out of Azure...", false) {
            override fun run(indicator: ProgressIndicator) {
                val result = AzureCliAuthProvider.logout()
                ApplicationManager.getApplication().invokeLater {
                    when (result) {
                        is AzureCliAuthProvider.TokenResult.Success -> {
                            clearResults()
                            EnvFileApplicator.notify(project, "Signed out of Azure CLI")
                        }
                        AzureCliAuthProvider.TokenResult.NotInstalled ->
                            showError("Azure CLI ('az') not found.")
                        else -> showError("Sign out failed.")
                    }
                    updateAuthState()
                }
            }
        })
    }

    // --- Load groups (populates the picker) ---

    private fun loadGroups() {
        val orgUrl = orgUrlField.text.trim()
        if (orgUrl.isEmpty()) {
            EnvFileApplicator.notify(project, "Please enter your Azure DevOps URL.", NotificationType.WARNING)
            return
        }
        loadGroupsButton.isEnabled = false
        ProgressManager.getInstance()
            .run(object : Task.Backgroundable(project, "Loading variable groups...", true) {
                override fun run(indicator: ProgressIndicator) {
                    val token = resolveToken() ?: run { enableLoadGroups(); return }
                    try {
                        val (org, proj) = AzureConnection.parseUrl(orgUrl)
                        if (proj == null) throw IllegalArgumentException("URL must include the project, e.g. https://dev.azure.com/myorg/myproject")
                        val groups = AzureDevOpsClient(org, proj).getVariableGroups(token).sortedBy { it.name }
                        ApplicationManager.getApplication().invokeLater {
                            val keep = currentGroupName()
                            groupCombo.model = DefaultComboBoxModel(groups.map { it.name }.toTypedArray())
                            if (keep.isNotEmpty()) groupCombo.editor.item = keep
                            enableLoadGroups()
                            EnvFileApplicator.notify(project, "Loaded ${groups.size} variable groups")
                        }
                    } catch (e: Exception) {
                        showError("Could not load groups: ${rootMessage(e)}")
                        enableLoadGroups()
                    }
                }
            })
    }

    // --- Fetch variables into the table (always re-pulls fresh) ---

    private fun fetchVariables() {
        val orgUrl = orgUrlField.text.trim()
        val groupName = currentGroupName()
        if (orgUrl.isEmpty()) {
            EnvFileApplicator.notify(project, "Please enter your Azure DevOps URL.", NotificationType.WARNING)
            return
        }
        if (groupName.isEmpty()) {
            EnvFileApplicator.notify(project, "Please enter or pick a variable group.", NotificationType.WARNING)
            return
        }
        persistGroup()

        fetchButton.isEnabled = false
        ProgressManager.getInstance()
            .run(object : Task.Backgroundable(project, "Fetching '$groupName'...", true) {
                override fun run(indicator: ProgressIndicator) {
                    val token = resolveToken() ?: run { enableFetch(); return }
                    try {
                        val (org, proj) = AzureConnection.parseUrl(orgUrl)
                        if (proj == null) throw IllegalArgumentException("URL must include the project, e.g. https://dev.azure.com/myorg/myproject")
                        val group = AzureDevOpsClient(org, proj).getVariableGroupByName(groupName, token)
                        ApplicationManager.getApplication().invokeLater {
                            populateFromGroup(group)
                            fetchButton.isEnabled = true
                            EnvFileApplicator.notify(project, "Fetched ${group.variables.size} variables from '${group.name}'")
                        }
                    } catch (e: IllegalArgumentException) {
                        showError("Invalid input: ${e.message}"); enableFetch()
                    } catch (e: ConnectException) {
                        showError("Cannot reach Azure DevOps. Check your network and URL."); enableFetch()
                    } catch (e: HttpTimeoutException) {
                        showError("Request timed out. Try again."); enableFetch()
                    } catch (e: Exception) {
                        showError("Fetch failed: ${rootMessage(e)}"); enableFetch()
                    }
                }
            })
    }

    /** Resolves an access token, surfacing the right actionable error on failure. Returns null on failure. */
    private fun resolveToken(): String? = when (val r = AzureCliAuthProvider.getAccessTokenResult()) {
        is AzureCliAuthProvider.TokenResult.Success -> r.accessToken
        AzureCliAuthProvider.TokenResult.NotInstalled -> {
            showError("Azure CLI ('az') not found. Install it or set its path in Settings → DotEnvify."); null
        }
        AzureCliAuthProvider.TokenResult.NotLoggedIn -> {
            showError("Not signed in to Azure. Run 'az login' (or click 'Run az login'), then Connect."); null
        }
        is AzureCliAuthProvider.TokenResult.Failed -> {
            showError("Azure CLI error: ${r.message}"); null
        }
    }

    private fun populateFromGroup(group: VariableGroup) {
        currentGroupId = group.id
        val rows = group.variables.entries
            .sortedBy { it.key }
            .map { (key, v) -> VarRow(key, v.value, v.isSecret, include = !v.isSecret) }
        recomputeTargetDiff(fire = false) // refresh target snapshot before showing rows
        tableModel.setRows(rows)
        applyFiltersToSelection() // honor active Ignore-lowercase / URL-only on first display

        groupInfoLabel.text = if (group.description.isNotEmpty()) "${group.name}: ${group.description}" else group.name
        openInBrowserButton.isVisible = group.id > 0
    }

    // --- Live derived state: preview, counts, diff ---

    /** Recompute the target .env snapshot used for the diff Status column. */
    private fun recomputeTargetDiff(fire: Boolean = true) {
        val path = resolvedTargetPath()
        targetMap = if (path == null) {
            emptyMap()
        } else {
            try {
                DotEnvIO.readEnvFile(path).associate { it.key to DotEnvParser.unquote(it.value) }
            } catch (_: Exception) {
                emptyMap()
            }
        }
        if (fire) tableModel.fireTableDataChanged()
    }

    private fun onOptionsChanged() {
        val now = optionsPanel.options()
        // The two *filters* drive the Include checkboxes; sort/export only re-format the output.
        val filterChanged = now.ignoreLowercase != lastOptions.ignoreLowercase || now.urlOnly != lastOptions.urlOnly
        lastOptions = now
        if (filterChanged) {
            applyFiltersToSelection() // fires a table change → refreshDerived()
        } else {
            refreshDerived()
        }
    }

    /** Tick/untick rows so the selection reflects the active Ignore-lowercase / URL-only filters. */
    private fun applyFiltersToSelection() {
        tableModel.setIncludeWhere { row -> !row.isSecret && !filteredOut(row) }
    }

    private fun refreshDerived() {
        refreshPreview()
        updateSelectionStatus()
    }

    private fun refreshPreview() {
        // Filters already applied via the checkboxes — only sort/export shape the output here.
        previewArea.text = DotEnvFormatter.format(tableModel.includedEntries(), formatOptionsForOutput())
        previewArea.caretPosition = 0
    }

    /** Output formatting only (sort, export prefix); filtering is the checkboxes' job now. */
    private fun formatOptionsForOutput() =
        optionsPanel.options().copy(ignoreLowercase = false, urlOnly = false)

    private fun updateSelectionStatus() {
        val total = tableModel.nonSecretCount()
        val secrets = tableModel.secretCount()
        val written = tableModel.includedEntries().size
        statusLabel.text = when {
            total == 0 && secrets == 0 -> " "
            secrets > 0 -> "$written of $total selected · $secrets secret skipped"
            else -> "$written of $total selected"
        }
        applyButton.isEnabled = written > 0
        copyButton.isEnabled = written > 0
    }

    // --- Apply / Copy / Open ---

    private fun applyToFile() {
        val entries = tableModel.includedEntries()
        if (entries.isEmpty()) {
            EnvFileApplicator.notify(project, "Select at least one variable to apply.", NotificationType.WARNING)
            return
        }
        val path = resolvedTargetPath() ?: return
        persistTarget()
        EnvFileApplicator.apply(project, entries, path, "Azure DevOps", formatOptionsForOutput())
        recomputeTargetDiff() // file changed — refresh the diff status
    }

    private fun copyToClipboard() {
        val output = previewArea.text
        if (output.isBlank()) return
        Toolkit.getDefaultToolkit().systemClipboard.setContents(StringSelection(output), null)
        EnvFileApplicator.notify(project, "Copied ${output.count { it == '\n' }} variables to clipboard")
    }

    private fun openInBrowser() {
        val orgUrl = orgUrlField.text.trim()
        if (currentGroupId <= 0) return
        try {
            val (org, proj) = AzureConnection.parseUrl(orgUrl)
            if (proj == null) return
            BrowserUtil.browse(
                "https://dev.azure.com/$org/$proj/_library" +
                    "?itemType=VariableGroups&view=VariableGroupView&variableGroupId=$currentGroupId"
            )
        } catch (_: IllegalArgumentException) {
            // invalid URL — nothing to open
        }
    }

    // --- Helpers ---

    private fun clearResults() {
        tableModel.setRows(emptyList())
        groupInfoLabel.text = " "
        openInBrowserButton.isVisible = false
        currentGroupId = 0
        refreshDerived()
    }

    private fun enableFetch() {
        ApplicationManager.getApplication().invokeLater { fetchButton.isEnabled = true }
    }

    private fun enableLoadGroups() {
        ApplicationManager.getApplication().invokeLater { loadGroupsButton.isEnabled = true }
    }

    private fun rootMessage(e: Throwable): String {
        val root = generateSequence<Throwable>(e) { it.cause }.last()
        return if (root !== e) "${e.message} (${root.javaClass.simpleName}: ${root.message})" else (e.message ?: e.javaClass.simpleName)
    }

    private fun showError(message: String) {
        ApplicationManager.getApplication().invokeLater {
            EnvFileApplicator.notify(project, message, NotificationType.ERROR)
        }
    }

    private class SimpleDocListener(private val action: () -> Unit) : javax.swing.event.DocumentListener {
        override fun insertUpdate(e: javax.swing.event.DocumentEvent) = action()
        override fun removeUpdate(e: javax.swing.event.DocumentEvent) = action()
        override fun changedUpdate(e: javax.swing.event.DocumentEvent) = action()
    }

    // --- Diff status of a variable vs the current target file ---

    private enum class DiffStatus(val label: String, val color: JBColor) {
        SECRET("secret", JBColor.GRAY),
        NEW("new", JBColor(Color(0, 128, 0), Color(80, 200, 80))),
        CHANGED("changed", JBColor(Color(190, 140, 0), Color(220, 180, 50))),
        SAME("same", JBColor.GRAY),
    }

    private fun statusOf(row: VarRow): DiffStatus = when {
        row.isSecret -> DiffStatus.SECRET
        !targetMap.containsKey(row.key) -> DiffStatus.NEW
        DotEnvParser.unquote(row.value) != targetMap[row.key] -> DiffStatus.CHANGED
        else -> DiffStatus.SAME
    }

    /** True when the row would be dropped by the active format-option filters. */
    private fun filteredOut(row: VarRow): Boolean {
        val opts = optionsPanel.options()
        if (opts.ignoreLowercase && row.key == row.key.lowercase() && row.key != row.key.uppercase()) return true
        if (opts.urlOnly && !DotEnvFormatter.isHTTPURL(row.value)) return true
        return false
    }

    // --- Table setup & model ---

    private fun configureTable() {
        varTable.columnModel.getColumn(COL_INCLUDE).apply { maxWidth = JBUI.scale(34); minWidth = JBUI.scale(34) }
        varTable.columnModel.getColumn(COL_KEY).preferredWidth = JBUI.scale(200)
        varTable.columnModel.getColumn(COL_STATUS).apply { maxWidth = JBUI.scale(90); minWidth = JBUI.scale(70) }

        val textRenderer = object : DefaultTableCellRenderer() {
            override fun getTableCellRendererComponent(
                table: JTable, value: Any?, isSelected: Boolean, hasFocus: Boolean, row: Int, col: Int,
            ): Component {
                super.getTableCellRendererComponent(table, value, isSelected, hasFocus, row, col)
                val r = tableModel.rowAt(table.convertRowIndexToModel(row))
                font = if (r.isSecret) font.deriveFont(Font.ITALIC) else font.deriveFont(Font.PLAIN)
                if (!isSelected) {
                    foreground = when {
                        col == COL_STATUS -> statusOf(r).color
                        r.isSecret || !r.include -> JBColor.GRAY
                        else -> table.foreground
                    }
                }
                toolTipText = when {
                    r.isSecret -> "Secret — value not exposed by Azure; set it manually"
                    filteredOut(r) -> "Excluded by the current format options (Ignore lowercase / URL-only)"
                    else -> null
                }
                return this
            }
        }
        varTable.columnModel.getColumn(COL_KEY).cellRenderer = textRenderer
        varTable.columnModel.getColumn(COL_VALUE).cellRenderer = textRenderer
        varTable.columnModel.getColumn(COL_STATUS).cellRenderer = textRenderer
    }

    private data class VarRow(val key: String, var value: String, val isSecret: Boolean, var include: Boolean)

    private inner class VarTableModel : AbstractTableModel() {
        private val columns = arrayOf("", "Key", "Value", "Status")
        private var rows: List<VarRow> = emptyList()

        fun setRows(newRows: List<VarRow>) {
            rows = newRows
            fireTableDataChanged()
        }

        fun includedEntries(): List<EnvEntry> =
            rows.filter { it.include && !it.isSecret }.map { EnvEntry(it.key, it.value) }

        fun nonSecretCount(): Int = rows.count { !it.isSecret }
        fun secretCount(): Int = rows.count { it.isSecret }
        fun rowAt(index: Int): VarRow = rows[index]

        /** Bulk-set the Include flag from a rule (used by the format-option filters). */
        fun setIncludeWhere(rule: (VarRow) -> Boolean) {
            rows.forEach { it.include = rule(it) }
            fireTableDataChanged()
        }

        override fun getRowCount(): Int = rows.size
        override fun getColumnCount(): Int = columns.size
        override fun getColumnName(col: Int): String = columns[col]

        override fun getColumnClass(col: Int): Class<*> =
            if (col == COL_INCLUDE) java.lang.Boolean::class.java else String::class.java

        override fun isCellEditable(row: Int, col: Int): Boolean {
            val r = rows[row]
            if (r.isSecret) return false
            return col == COL_INCLUDE || col == COL_VALUE
        }

        override fun getValueAt(row: Int, col: Int): Any {
            val r = rows[row]
            return when (col) {
                COL_INCLUDE -> r.include && !r.isSecret
                COL_KEY -> r.key
                COL_VALUE -> if (r.isSecret) "(secret — set in Azure)" else r.value
                COL_STATUS -> statusOf(r).label
                else -> ""
            }
        }

        override fun setValueAt(value: Any?, row: Int, col: Int) {
            val r = rows[row]
            when (col) {
                COL_INCLUDE -> r.include = value as? Boolean ?: r.include
                COL_VALUE -> r.value = value as? String ?: r.value
            }
            fireTableRowsUpdated(row, row)
        }
    }

    private companion object {
        const val COL_INCLUDE = 0
        const val COL_KEY = 1
        const val COL_VALUE = 2
        const val COL_STATUS = 3
    }
}
