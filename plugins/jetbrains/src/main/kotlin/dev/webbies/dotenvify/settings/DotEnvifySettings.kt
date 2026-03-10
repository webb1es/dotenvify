package dev.webbies.dotenvify.settings

import com.intellij.openapi.application.ApplicationManager
import com.intellij.openapi.components.PersistentStateComponent
import com.intellij.openapi.components.State
import com.intellij.openapi.components.Storage

/**
 * Global (application-level) settings for DotEnvify.
 */
@State(name = "DotEnvifySettings", storages = [Storage("DotEnvifySettings.xml")])
class DotEnvifySettings : PersistentStateComponent<DotEnvifySettings.State> {

    data class State(
        /** Prefix each output line with `export`. */
        var exportPrefix: Boolean = false,
        /** Sort keys alphabetically. */
        var sort: Boolean = true,
        /** Exclude keys containing lowercase letters. */
        var ignoreLowercase: Boolean = true,
        /** Only include entries whose values are URLs. */
        var urlOnly: Boolean = false,
        /** Default file path for .env output. */
        var defaultOutputPath: String = ".env",
        /** Persisted Azure DevOps organization URL for re-use across sessions. */
        var azureOrgUrl: String = "",
    )

    private var state = State()

    override fun getState(): State = state

    override fun loadState(state: State) {
        this.state = state
    }

    companion object {
        fun getInstance(): DotEnvifySettings =
            ApplicationManager.getApplication().getService(DotEnvifySettings::class.java)
    }
}
