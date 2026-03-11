package dev.webbies.dotenvify.azure

/** An Azure DevOps Variable Group containing named variables. */
data class VariableGroup(
    val id: Int,
    val name: String,
    val variables: Map<String, Variable>,
    val description: String = "",
)

/** A single variable from an Azure DevOps Variable Group. */
data class Variable(
    val value: String,
    /** Secret variables are masked in Azure and cannot be read via API. */
    val isSecret: Boolean = false,
)

/** Response from the Azure AD Device Code Flow initiation endpoint. */
data class DeviceCodeResponse(
    val deviceCode: String,
    val userCode: String,
    val verificationUri: String,
    /** Seconds until the device code expires. */
    val expiresIn: Int,
    /** Minimum seconds to wait between poll attempts. */
    val interval: Int,
    val message: String,
)

/** OAuth token response from Azure AD. */
data class TokenResponse(
    val accessToken: String,
    val refreshToken: String,
    /** Seconds until the access token expires. */
    val expiresIn: Int,
    val tokenType: String,
)

data class AzureConnection(
    val organizationUrl: String,
    val project: String,
) {
    val organization: String
        get() = parseOrg(organizationUrl)

    companion object {
        private val DEV_AZURE_REGEX = Regex("""https?://dev\.azure\.com/([^/]+)/([^/]+)""")
        private val VSTUDIO_REGEX = Regex("""https?://([^.]+)\.visualstudio\.com/([^/]+)""")

        /** Extracts (organization, project) from an Azure DevOps URL. Throws on invalid format. */
        fun parseUrl(url: String): Pair<String, String?> {
            DEV_AZURE_REGEX.find(url)?.let {
                return it.groupValues[1] to it.groupValues[2].ifEmpty { null }
            }
            VSTUDIO_REGEX.find(url)?.let {
                return it.groupValues[1] to it.groupValues[2].ifEmpty { null }
            }
            throw IllegalArgumentException(
                "Invalid Azure DevOps URL. Expected: https://dev.azure.com/{org}/{project} " +
                        "or https://{org}.visualstudio.com/{project}"
            )
        }

        private fun parseOrg(url: String): String = parseUrl(url).first
    }
}
