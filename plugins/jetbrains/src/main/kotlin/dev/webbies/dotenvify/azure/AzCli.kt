package dev.webbies.dotenvify.azure

import com.intellij.execution.configurations.GeneralCommandLine
import com.intellij.execution.process.CapturingProcessHandler
import com.intellij.execution.process.ProcessOutput
import com.intellij.openapi.diagnostic.Logger
import com.intellij.openapi.progress.ProgressIndicator
import com.intellij.openapi.util.SystemInfo
import dev.webbies.dotenvify.azure.AzCli.resolve
import java.io.File

/**
 * Locates and invokes the Azure CLI (`az`).
 *
 * All methods block on a subprocess — call them off the EDT.
 *
 * PATH note: when the IDE is launched from the macOS Dock/Finder (or a Windows
 * shortcut), the process environment does NOT include the user's shell PATH, so a
 * bare `az` is often not found even though it works in their terminal. [resolve]
 * therefore searches a configured path, then well-known install locations, and
 * finally asks a login shell where `az` lives.
 */
object AzCli {

    private val LOG = Logger.getInstance("[DotEnvify]")
    private const val EXEC_TIMEOUT_MS = 30_000

    @Volatile private var cachedExePath: String? = null

    sealed class Resolution {
        data class Found(val path: String) : Resolution()
        object NotFound : Resolution()
    }

    /** Resolves the absolute path to `az`, caching the first successful result. */
    fun resolve(configuredPath: String?): Resolution {
        cachedExePath?.let { return Resolution.Found(it) }

        // (a) user-configured path
        configuredPath?.takeIf { it.isNotBlank() }?.let {
            val f = File(it)
            if (f.isFile && f.canExecute()) return found(it)
        }
        // (b) common install locations
        for (cand in commonLocations()) {
            val f = File(cand)
            if (f.isFile && f.canExecute()) return found(cand)
        }
        // (c) login-shell PATH resolution (inherits the user's real PATH)
        resolveViaShell()?.let { return found(it) }

        return Resolution.NotFound
    }

    private fun found(path: String): Resolution.Found {
        cachedExePath = path
        return Resolution.Found(path)
    }

    private fun commonLocations(): List<String> = when {
        SystemInfo.isWindows -> listOf(
            "C:\\Program Files\\Microsoft SDKs\\Azure\\CLI2\\wbin\\az.cmd",
            "C:\\Program Files (x86)\\Microsoft SDKs\\Azure\\CLI2\\wbin\\az.cmd",
        )
        SystemInfo.isMac -> listOf(
            "/opt/homebrew/bin/az", // Apple Silicon Homebrew
            "/usr/local/bin/az",    // Intel Homebrew / installer
        )
        else -> listOf("/usr/bin/az", "/usr/local/bin/az", "/snap/bin/az")
    }

    /** Last resort: ask the user's login shell where `az` is, to inherit their PATH. */
    private fun resolveViaShell(): String? {
        val cmd = if (SystemInfo.isWindows) {
            GeneralCommandLine("where", "az")
        } else {
            val shell = System.getenv("SHELL")?.takeIf { it.isNotBlank() } ?: "/bin/zsh"
            // -l = login shell so ~/.zprofile / ~/.bash_profile PATH is loaded.
            GeneralCommandLine(shell, "-l", "-c", "command -v az")
        }
        return try {
            val out = CapturingProcessHandler(cmd.withCharset(Charsets.UTF_8)).runProcess(EXEC_TIMEOUT_MS)
            out.stdout.lineSequence()
                .map { it.trim() }
                .firstOrNull { it.isNotEmpty() && File(it).let { f -> f.isFile && f.canExecute() } }
        } catch (e: Exception) {
            LOG.info("az shell resolution failed: ${e.message}")
            null
        }
    }

    /** Runs a short `az` command with a fixed timeout. Capture-only; never blocks the EDT. */
    fun exec(exePath: String, vararg args: String): ProcessOutput =
        CapturingProcessHandler(buildCommandLine(exePath, args)).runProcess(EXEC_TIMEOUT_MS)

    /**
     * Runs a long-running, cancelable `az` command (e.g. `az login`, which blocks
     * until the browser flow completes). Cancelling [indicator] kills the process.
     */
    fun execInteractive(exePath: String, indicator: ProgressIndicator, vararg args: String): ProcessOutput =
        CapturingProcessHandler(buildCommandLine(exePath, args)).runProcessWithProgressIndicator(indicator)

    /**
     * Builds the command line for `az`. On Windows the Azure CLI is a batch file (`az.cmd`),
     * which the JVM cannot launch directly (`CreateProcess error=193`) — it must run through
     * `cmd.exe /c`. On macOS/Linux `az` is a normal executable and runs directly.
     */
    private fun buildCommandLine(exePath: String, args: Array<out String>): GeneralCommandLine {
        val isBatch = SystemInfo.isWindows &&
                (exePath.endsWith(".cmd", ignoreCase = true) || exePath.endsWith(".bat", ignoreCase = true))
        val cmd = if (isBatch) {
            GeneralCommandLine("cmd.exe", "/c", exePath, *args)
        } else {
            GeneralCommandLine(exePath, *args)
        }
        return cmd.withCharset(Charsets.UTF_8)
    }
}
