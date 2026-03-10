package dev.webbies.dotenvify.settings

import com.intellij.openapi.components.PersistentStateComponent
import com.intellij.openapi.components.State
import com.intellij.openapi.components.Storage
import com.intellij.openapi.project.Project

/**
 * Per-project settings for DotEnvify.
 */
@State(name = "DotEnvifyProjectSettings", storages = [Storage("dotenvify.xml")])
class DotEnvifyProjectSettings : PersistentStateComponent<DotEnvifyProjectSettings.State> {

    data class State(
        /** When true, ignore project-level overrides and use global settings. */
        var useGlobalDefaults: Boolean = true,
        var exportPrefix: Boolean = false,
        var sort: Boolean = true,
        var ignoreLowercase: Boolean = true,
        var urlOnly: Boolean = false,
        /** Project-specific .env output path. */
        var outputPath: String = ".env",
        /** Comma-separated list of keys whose values should be preserved on merge. */
        var preserveKeys: String = "",
        /** Persisted Azure DevOps variable group name for this project. */
        var azureGroupName: String = "",
    )

    private var state = State()

    override fun getState(): State = state

    override fun loadState(state: State) {
        this.state = state
    }

    companion object {
        fun getInstance(project: Project): DotEnvifyProjectSettings =
            project.getService(DotEnvifyProjectSettings::class.java)
    }
}
