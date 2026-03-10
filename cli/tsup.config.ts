import {defineConfig} from "tsup";

export default defineConfig({
    entry: ["src/index.ts"],
    format: "esm",
    target: "node18",
    outDir: "dist",
    clean: true,
    // Bundle @dotenvify/core (private workspace pkg) into the output
    // Keep commander as external (it's a real npm dependency)
    noExternal: ["@dotenvify/core"],
});
