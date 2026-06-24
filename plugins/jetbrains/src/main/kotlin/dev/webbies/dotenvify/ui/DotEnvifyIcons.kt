package dev.webbies.dotenvify.ui

import com.intellij.openapi.util.IconLoader

/** Plugin-bundled icons, loaded from /resources/icons. */
object DotEnvifyIcons {
    @JvmField
    val ConvertSelection = IconLoader.getIcon("/icons/convertSelection.svg", DotEnvifyIcons::class.java)

    @JvmField
    val ConvertFile = IconLoader.getIcon("/icons/convertFile.svg", DotEnvifyIcons::class.java)
}
