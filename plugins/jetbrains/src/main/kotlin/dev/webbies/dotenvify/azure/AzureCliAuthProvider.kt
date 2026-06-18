package dev.webbies.dotenvify.azure

import com.google.gson.Gson
import com.google.gson.JsonObject
import com.intellij.openapi.diagnostic.Logger
import com.intellij.openapi.progress.ProgressIndicator
import dev.webbies.dotenvify.settings.DotEnvifySettings
import java.time.Instant
import java.time.LocalDateTime
import java.time.OffsetDateTime
import java.time.ZoneId
import java.time.format.DateTimeFormatter

/**
 * Authenticates to Azure DevOps using the user's local Azure CLI (`az`) session.
 *
 * Shells out to `az account get-access-token` for a Bearer token scoped to Azure DevOps.
 * The user signs in once with `az login` (in a terminal or via the in-IDE button) and the
 * plugin reuses that session.
 *
 * `az` owns credential storage; this plugin never persists tokens. Acquired tokens
 * are only cached in memory until just before they expire.
 */
object AzureCliAuthProvider {

    private val LOG = Logger.getInstance("[DotEnvify]")

    /** Fixed, Microsoft-assigned Azure DevOps resource (application) ID in Entra. */
    private const val DEVOPS_RESOURCE = "499b84ac-1321-427f-aa17-267ca6975798"

    private const val EXPIRY_SKEW_SECONDS = 300L

    private val gson = Gson()

    private data class Cached(val token: String, val expiresAtEpoch: Long)

    @Volatile private var cache: Cached? = null

    /** Result of acquiring an access token, with actionable failure states. */
    sealed class TokenResult {
        data class Success(val accessToken: String) : TokenResult()
        object NotInstalled : TokenResult()
        object NotLoggedIn : TokenResult()
        data class Failed(val message: String) : TokenResult()
    }

    /** Identity reported by `az account show`. */
    sealed class Status {
        data class SignedIn(val user: String, val subscription: String) : Status()
        object NotInstalled : Status()
        object NotLoggedIn : Status()
        data class Failed(val message: String) : Status()
    }

    private fun configuredPath(): String? =
        DotEnvifySettings.getInstance().state.azCliPath.ifBlank { null }

    // --- Rich API ---

    /** Acquires a Bearer token for Azure DevOps, reusing the in-memory cache when valid. */
    fun getAccessTokenResult(): TokenResult {
        cache?.let {
            if (Instant.now().epochSecond < it.expiresAtEpoch - EXPIRY_SKEW_SECONDS) {
                return TokenResult.Success(it.token)
            }
        }

        val exe = when (val r = AzCli.resolve(configuredPath())) {
            is AzCli.Resolution.Found -> r.path
            AzCli.Resolution.NotFound -> return TokenResult.NotInstalled
        }

        val out = AzCli.exec(
            exe, "account", "get-access-token",
            "--resource", DEVOPS_RESOURCE, "--output", "json",
        )
        if (out.isTimeout) return TokenResult.Failed("Azure CLI timed out.")
        if (out.exitCode != 0) return classifyError(out.stderr)

        return try {
            val json = gson.fromJson(out.stdout, JsonObject::class.java)
            val token = json.str("accessToken")
                ?: return TokenResult.Failed("Azure CLI returned no access token.")
            cache = Cached(token, parseExpiresOn(json))
            TokenResult.Success(token)
        } catch (e: Exception) {
            TokenResult.Failed("Could not parse Azure CLI output: ${e.message}")
        }
    }

    /** Reports who is signed in via `az account show`, without acquiring a token. */
    fun status(): Status {
        val exe = when (val r = AzCli.resolve(configuredPath())) {
            is AzCli.Resolution.Found -> r.path
            AzCli.Resolution.NotFound -> return Status.NotInstalled
        }

        val out = AzCli.exec(exe, "account", "show", "--output", "json")
        if (out.isTimeout) return Status.Failed("Azure CLI timed out.")
        if (out.exitCode != 0) {
            return when (val e = classifyError(out.stderr)) {
                TokenResult.NotLoggedIn -> Status.NotLoggedIn
                TokenResult.NotInstalled -> Status.NotInstalled
                is TokenResult.Failed -> Status.Failed(e.message)
                else -> Status.Failed(out.stderr.trim().take(300))
            }
        }

        return try {
            val json = gson.fromJson(out.stdout, JsonObject::class.java)
            val user = json.getAsJsonObject("user")?.str("name") ?: "unknown"
            val subscription = json.str("name") ?: ""
            Status.SignedIn(user, subscription)
        } catch (e: Exception) {
            Status.Failed("Could not parse Azure CLI output: ${e.message}")
        }
    }

    /**
     * Runs `az login` (opens the system browser; blocks until the user finishes).
     * Cancel the [indicator] to abort and kill the process. Clears the token cache.
     */
    fun login(indicator: ProgressIndicator): TokenResult {
        val exe = when (val r = AzCli.resolve(configuredPath())) {
            is AzCli.Resolution.Found -> r.path
            AzCli.Resolution.NotFound -> return TokenResult.NotInstalled
        }

        cache = null
        val out = try {
            AzCli.execInteractive(exe, indicator, "login", "--output", "none")
        } catch (e: Exception) {
            LOG.info("az login failed: ${e.message}")
            return TokenResult.Failed("Azure sign-in failed: ${e.message}")
        }
        return if (out.exitCode == 0) TokenResult.Success("") else classifyError(out.stderr)
    }

    /**
     * Runs `az logout`, ending the user's Azure CLI session for the whole machine.
     * Clears the token cache. Returns [TokenResult.Success] even if no account was
     * active (the end state — signed out — is what the caller wants).
     */
    fun logout(): TokenResult {
        val exe = when (val r = AzCli.resolve(configuredPath())) {
            is AzCli.Resolution.Found -> r.path
            AzCli.Resolution.NotFound -> return TokenResult.NotInstalled
        }
        cache = null
        val out = AzCli.exec(exe, "logout")
        // "az logout" exits non-zero when already signed out — treat that as success.
        val stderr = out.stderr.lowercase()
        return if (out.exitCode == 0 || "no active account" in stderr || "not logged in" in stderr) {
            TokenResult.Success("")
        } else {
            TokenResult.Failed(out.stderr.trim().take(300).ifEmpty { "Sign out failed." })
        }
    }

    // --- Helpers (internal for unit testing) ---

    internal fun classifyError(stderr: String): TokenResult {
        val s = stderr.lowercase()
        return when {
            "az login" in s || "not logged in" in s || "no subscription" in s ||
                "please run 'az login'" in s || "az account set" in s ->
                TokenResult.NotLoggedIn
            "command not found" in s || "is not recognized" in s || "no such file" in s ->
                TokenResult.NotInstalled
            else -> TokenResult.Failed(stderr.trim().take(300).ifEmpty { "Azure CLI failed." })
        }
    }

    /**
     * Extracts the token expiry as an epoch second. Defensive across `az` versions:
     * prefers the epoch `expires_on`, then ISO-offset `expiresOn`, then the classic
     * local-time `expiresOn`, falling back to a safe ~50-minute window.
     */
    internal fun parseExpiresOn(json: JsonObject): Long {
        json.long("expires_on")?.let { return it }

        val raw = json.str("expiresOn")
        if (raw.isNullOrBlank()) return defaultExpiry()

        return try {
            OffsetDateTime.parse(raw, DateTimeFormatter.ISO_OFFSET_DATE_TIME).toEpochSecond()
        } catch (_: Exception) {
            try {
                val fmt = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss[.SSSSSS]")
                LocalDateTime.parse(raw, fmt).atZone(ZoneId.systemDefault()).toEpochSecond()
            } catch (_: Exception) {
                defaultExpiry()
            }
        }
    }

    private fun defaultExpiry(): Long = Instant.now().epochSecond + 3000
}
