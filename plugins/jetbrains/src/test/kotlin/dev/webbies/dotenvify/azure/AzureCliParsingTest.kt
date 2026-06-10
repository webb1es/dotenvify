package dev.webbies.dotenvify.azure

import com.google.gson.Gson
import com.google.gson.JsonObject
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import java.time.Instant

class AzureCliParsingTest {

    private val gson = Gson()
    private fun obj(json: String) = gson.fromJson(json, JsonObject::class.java)

    // --- get-access-token / expiry parsing ---

    @Test
    fun `access token JSON exposes accessToken`() {
        val json = obj("""{"accessToken":"abc.def.ghi","expiresOn":"2099-01-01 00:00:00.000000"}""")
        assertEquals("abc.def.ghi", json.str("accessToken"))
    }

    @Test
    fun `parseExpiresOn prefers epoch expires_on`() {
        val json = obj("""{"accessToken":"t","expires_on":4102444800,"expiresOn":"ignored"}""")
        assertEquals(4102444800L, AzureCliAuthProvider.parseExpiresOn(json))
    }

    @Test
    fun `parseExpiresOn handles ISO offset expiresOn`() {
        val json = obj("""{"expiresOn":"2099-01-01T00:00:00+00:00"}""")
        assertEquals(4070908800L, AzureCliAuthProvider.parseExpiresOn(json))
    }

    @Test
    fun `parseExpiresOn handles classic local expiresOn`() {
        val json = obj("""{"expiresOn":"2099-01-01 12:00:00.000000"}""")
        // Local-zone dependent; just assert it parses to a far-future instant.
        assertTrue(AzureCliAuthProvider.parseExpiresOn(json) > Instant.now().epochSecond)
    }

    @Test
    fun `parseExpiresOn falls back to future default when missing`() {
        val json = obj("""{"accessToken":"t"}""")
        assertTrue(AzureCliAuthProvider.parseExpiresOn(json) > Instant.now().epochSecond)
    }

    // --- error classification ---

    @Test
    fun `classifyError detects not logged in`() {
        val r = AzureCliAuthProvider.classifyError("ERROR: Please run 'az login' to setup account.")
        assertEquals(AzureCliAuthProvider.TokenResult.NotLoggedIn, r)
    }

    @Test
    fun `classifyError detects not installed`() {
        val r = AzureCliAuthProvider.classifyError("zsh: command not found: az")
        assertEquals(AzureCliAuthProvider.TokenResult.NotInstalled, r)
    }

    @Test
    fun `classifyError falls back to Failed`() {
        val r = AzureCliAuthProvider.classifyError("some unexpected failure")
        assertTrue(r is AzureCliAuthProvider.TokenResult.Failed)
    }

    // --- account show parsing ---

    @Test
    fun `account show JSON exposes user and subscription`() {
        val json = obj(
            """{"name":"My Subscription","user":{"name":"jane@contoso.com","type":"user"}}"""
        )
        assertEquals("My Subscription", json.str("name"))
        assertEquals("jane@contoso.com", json.getAsJsonObject("user")?.str("name"))
    }
}
