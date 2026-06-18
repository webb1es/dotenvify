package dev.webbies.dotenvify.azure

import com.google.gson.Gson
import com.google.gson.JsonObject
import java.net.URI
import java.net.URLEncoder
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse
import java.nio.charset.StandardCharsets
import java.time.Duration

/** REST API client for Azure DevOps variable groups. */
class AzureDevOpsClient(private val organization: String, private val project: String) {

    private val gson = Gson()
    private val httpClient: HttpClient = HttpClient.newBuilder()
        .connectTimeout(Duration.ofSeconds(30))
        .build()

    // Encode each path segment; Azure DevOps project names may contain spaces and other chars.
    private fun encodeSegment(value: String) =
        URLEncoder.encode(value, StandardCharsets.UTF_8).replace("+", "%20")

    private val baseUrl get() = "https://dev.azure.com/${encodeSegment(organization)}/${encodeSegment(project)}/_apis"

    /** Fetches all variable groups from the Azure DevOps project. */
    fun getVariableGroups(accessToken: String): List<VariableGroup> {
        val request = HttpRequest.newBuilder()
            .uri(URI.create("$baseUrl/distributedtask/variablegroups?api-version=7.1-preview.2"))
            .header("Authorization", "Bearer $accessToken")
            .header("Accept", "application/json")
            .GET()
            .build()

        val response = httpClient.send(request, HttpResponse.BodyHandlers.ofString())
        if (response.statusCode() != 200) {
            throw RuntimeException("Azure DevOps API error (${response.statusCode()}): ${response.body()}")
        }

        val json = gson.fromJson(response.body(), JsonObject::class.java)
        val groups = json.getAsJsonArray("value")
            ?: throw RuntimeException("Unexpected API response: ${response.body().take(500)}")

        return groups.map { element ->
            val obj = element.asJsonObject
            val variables = mutableMapOf<String, Variable>()

            obj.getAsJsonObject("variables")?.entrySet()?.forEach { (key, value) ->
                val varObj = value.asJsonObject
                variables[key] = Variable(
                    value = varObj.str("value") ?: "",
                    isSecret = varObj.bool("isSecret") ?: false,
                )
            }

            VariableGroup(
                id = obj.int("id") ?: 0,
                name = obj.str("name") ?: "",
                variables = variables,
                description = obj.str("description") ?: "",
            )
        }
    }

    /** Finds a variable group by exact name match. Throws if not found. */
    fun getVariableGroupByName(name: String, accessToken: String): VariableGroup {
        val groups = getVariableGroups(accessToken)
        return groups.find { it.name == name }
            ?: throw RuntimeException("Variable group '$name' not found. Available: ${groups.joinToString { it.name }}")
    }
}
