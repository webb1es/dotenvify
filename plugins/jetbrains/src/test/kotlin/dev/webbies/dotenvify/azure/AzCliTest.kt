package dev.webbies.dotenvify.azure

import org.junit.Assert.assertEquals
import org.junit.Test

class AzCliTest {

    private val winAzCmd = "C:\\Program Files\\Microsoft SDKs\\Azure\\CLI2\\wbin\\az.cmd"

    @Test
    fun `windows wraps az_cmd through cmd_exe`() {
        val tokens = AzCli.commandTokens(
            isWindows = true,
            exePath = winAzCmd,
            args = listOf("account", "show", "--output", "json"),
        )
        assertEquals(
            listOf("cmd.exe", "/c", winAzCmd, "account", "show", "--output", "json"),
            tokens,
        )
    }

    @Test
    fun `windows wraps bat files case-insensitively`() {
        val tokens = AzCli.commandTokens(isWindows = true, exePath = "C:\\tools\\az.BAT", args = listOf("logout"))
        assertEquals(listOf("cmd.exe", "/c", "C:\\tools\\az.BAT", "logout"), tokens)
    }

    @Test
    fun `windows runs a real exe directly`() {
        val tokens =
            AzCli.commandTokens(isWindows = true, exePath = "C:\\tools\\az.exe", args = listOf("account", "show"))
        assertEquals(listOf("C:\\tools\\az.exe", "account", "show"), tokens)
    }

    @Test
    fun `unix runs az directly`() {
        val tokens =
            AzCli.commandTokens(isWindows = false, exePath = "/opt/homebrew/bin/az", args = listOf("account", "show"))
        assertEquals(listOf("/opt/homebrew/bin/az", "account", "show"), tokens)
    }

    @Test
    fun `unix never wraps even a cmd-suffixed path`() {
        // `.cmd` is meaningless on Unix; the wrapper must only ever apply on Windows.
        val tokens = AzCli.commandTokens(isWindows = false, exePath = "/weird/az.cmd", args = listOf("logout"))
        assertEquals(listOf("/weird/az.cmd", "logout"), tokens)
    }
}
