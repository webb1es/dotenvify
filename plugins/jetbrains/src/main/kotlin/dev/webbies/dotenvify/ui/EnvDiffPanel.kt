package dev.webbies.dotenvify.ui

import com.intellij.openapi.project.Project
import com.intellij.openapi.ui.DialogWrapper
import com.intellij.ui.Gray
import com.intellij.ui.JBColor
import com.intellij.ui.components.JBScrollPane
import dev.webbies.dotenvify.core.DotEnvParser
import dev.webbies.dotenvify.core.EnvEntry
import java.awt.BorderLayout
import java.awt.Color
import java.awt.Component
import java.awt.Dimension
import java.awt.FlowLayout
import java.awt.Font
import javax.swing.*
import javax.swing.table.AbstractTableModel
import javax.swing.table.DefaultTableCellRenderer
import javax.swing.table.TableCellRenderer

/**
 * Merge preview dialog comparing existing .env with incoming variables.
 * For conflicting keys the user picks, per row, which value to keep.
 */
class EnvDiffDialog(
    project: Project,
    private val existingEntries: List<EnvEntry>,
    private val incomingEntries: List<EnvEntry>,
    private val sourceName: String = "Incoming",
) : DialogWrapper(project) {

    private enum class Status(val label: String, val color: JBColor) {
        ADDED("added", JBColor(Color(0, 128, 0), Color(80, 200, 80))),
        CHANGED("changed", JBColor(Color(190, 140, 0), Color(220, 180, 50))),
        REMOVED("kept", JBColor(Gray._130, Gray._150)),
        UNCHANGED("same", JBColor(Gray._130, Gray._150)),
    }

    private data class MergeRow(
        val key: String,
        val existing: String?,
        val incoming: String?,
        val status: Status,
        var useIncoming: Boolean = true,
    )

    private val mergeRows = buildDiff()
    private lateinit var tableModel: AbstractTableModel

    val mergedEntries: List<EnvEntry>
        get() = mergeRows.map { row ->
            val value = when {
                row.status == Status.CHANGED -> if (row.useIncoming) row.incoming!! else row.existing!!
                row.incoming != null -> row.incoming
                else -> row.existing!!
            }
            EnvEntry(row.key, value)
        }

    init {
        title = "Merge Preview"
        init()
    }

    override fun createCenterPanel(): JComponent {
        val counts = Status.entries.associateWith { s -> mergeRows.count { it.status == s } }

        val summaryLabel = JLabel(
            "Existing: ${existingEntries.size} keys   •   $sourceName: ${incomingEntries.size} keys   •   " +
                "+${counts[Status.ADDED]} new   ~${counts[Status.CHANGED]} changed   " +
                "=${counts[Status.UNCHANGED]} same   ·${counts[Status.REMOVED]} kept"
        )

        tableModel = object : AbstractTableModel() {
            private val columns = arrayOf("Key", "Status", "Existing Value", "$sourceName Value", "Use Incoming")

            override fun getRowCount() = mergeRows.size
            override fun getColumnCount() = columns.size
            override fun getColumnName(col: Int) = columns[col]

            override fun getValueAt(row: Int, col: Int): Any {
                val r = mergeRows[row]
                return when (col) {
                    0 -> r.key
                    1 -> r.status.label
                    2 -> r.existing ?: ""
                    3 -> r.incoming ?: ""
                    4 -> r.useIncoming
                    else -> ""
                }
            }

            // Only CHANGED rows have a real choice to make.
            override fun isCellEditable(row: Int, col: Int) = col == 4 && mergeRows[row].status == Status.CHANGED

            override fun setValueAt(value: Any?, row: Int, col: Int) {
                if (col == 4) {
                    mergeRows[row].useIncoming = value as Boolean
                    fireTableRowsUpdated(row, row)
                }
            }

            override fun getColumnClass(col: Int) =
                if (col == 4) java.lang.Boolean::class.java else String::class.java
        }

        val table = JTable(tableModel).apply {
            setShowGrid(false)
            rowHeight = 22
            autoResizeMode = JTable.AUTO_RESIZE_SUBSEQUENT_COLUMNS
            setDefaultRenderer(String::class.java, StatusAwareRenderer())
            columnModel.getColumn(4).cellRenderer = UseIncomingRenderer()
            columnModel.getColumn(0).preferredWidth = 150
            columnModel.getColumn(1).preferredWidth = 80
            columnModel.getColumn(2).preferredWidth = 200
            columnModel.getColumn(3).preferredWidth = 200
            columnModel.getColumn(4).preferredWidth = 96
        }

        val hasConflicts = mergeRows.any { it.status == Status.CHANGED }
        val bulkRow = JPanel(FlowLayout(FlowLayout.LEFT, 4, 0)).apply {
            add(JButton("Use all incoming").apply {
                isEnabled = hasConflicts
                addActionListener { setAllConflicts(true, table) }
            })
            add(JButton("Keep all existing").apply {
                isEnabled = hasConflicts
                addActionListener { setAllConflicts(false, table) }
            })
        }

        val footer = JPanel(BorderLayout()).apply {
            add(bulkRow, BorderLayout.WEST)
            add(buildLegend(), BorderLayout.EAST)
        }

        return JPanel(BorderLayout(0, 8)).apply {
            add(summaryLabel, BorderLayout.NORTH)
            add(JBScrollPane(table).apply { preferredSize = Dimension(760, 420) }, BorderLayout.CENTER)
            add(footer, BorderLayout.SOUTH)
        }
    }

    private fun setAllConflicts(useIncoming: Boolean, table: JTable) {
        mergeRows.forEach { if (it.status == Status.CHANGED) it.useIncoming = useIncoming }
        tableModel.fireTableDataChanged()
        table.repaint()
    }

    private fun buildLegend(): JComponent = JPanel(FlowLayout(FlowLayout.RIGHT, 10, 0)).apply {
        Status.entries.distinctBy { it.color.rgb to it.label }.forEach { s ->
            add(JLabel("● ${s.label}").apply { foreground = s.color })
        }
    }

    private fun buildDiff(): MutableList<MergeRow> {
        val existingMap = existingEntries.associate { it.key to it.value }
        val incomingMap = incomingEntries.associate { it.key to it.value }
        return (existingMap.keys + incomingMap.keys).sorted().map { key ->
            val e = existingMap[key]
            val i = incomingMap[key]
            MergeRow(
                key, e, i, when {
                    e == null -> Status.ADDED
                    i == null -> Status.REMOVED
                    DotEnvParser.unquote(e) != DotEnvParser.unquote(i) -> Status.CHANGED
                    else -> Status.UNCHANGED
                }
            )
        }.toMutableList()
    }

    /** Colours the Status cell, and dims the value that will NOT be used. */
    private inner class StatusAwareRenderer : DefaultTableCellRenderer() {
        override fun getTableCellRendererComponent(
            table: JTable, value: Any?, isSelected: Boolean,
            hasFocus: Boolean, row: Int, col: Int,
        ): Component {
            super.getTableCellRendererComponent(table, value, isSelected, hasFocus, row, col)
            val r = mergeRows[table.convertRowIndexToModel(row)]
            if (isSelected) return this

            foreground = when (col) {
                1 -> r.status.color // Status badge
                2 -> if (winner(r) == Winner.EXISTING) table.foreground else JBColor.GRAY // existing value
                3 -> if (winner(r) == Winner.INCOMING) r.status.color else JBColor.GRAY    // incoming value
                else -> table.foreground
            }
            font = if (col == 1) font.deriveFont(Font.BOLD) else font.deriveFont(Font.PLAIN)
            return this
        }
    }

    private enum class Winner { EXISTING, INCOMING, NONE }

    private fun winner(r: MergeRow): Winner = when (r.status) {
        Status.ADDED -> Winner.INCOMING
        Status.REMOVED -> Winner.EXISTING
        Status.UNCHANGED -> Winner.INCOMING
        Status.CHANGED -> if (r.useIncoming) Winner.INCOMING else Winner.EXISTING
    }

    /** Shows a checkbox only where a choice exists (CHANGED); a greyed "—" elsewhere. */
    private inner class UseIncomingRenderer : TableCellRenderer {
        private val checkBox = JCheckBox().apply { horizontalAlignment = SwingConstants.CENTER }
        private val placeholder = JLabel("—", SwingConstants.CENTER).apply { foreground = JBColor.GRAY }

        override fun getTableCellRendererComponent(
            table: JTable, value: Any?, isSelected: Boolean,
            hasFocus: Boolean, row: Int, col: Int,
        ): Component {
            val r = mergeRows[table.convertRowIndexToModel(row)]
            val bg = if (isSelected) table.selectionBackground else table.background
            return if (r.status == Status.CHANGED) {
                checkBox.apply {
                    this.isSelected = r.useIncoming
                    background = bg
                    toolTipText = "On: use the $sourceName value · Off: keep the existing value"
                }
            } else {
                placeholder.apply {
                    background = bg
                    isOpaque = true
                    toolTipText = "No choice needed for ${r.status.label} keys"
                }
            }
        }
    }
}
