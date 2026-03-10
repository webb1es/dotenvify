import {Command} from "commander";
import {convert} from "./commands/convert.js";

const program = new Command();

program
    .name("dotenvify")
    .description("Convert messy environment variables into clean .env files")
    .version("2.0.0")
    .argument("<source>", "Source file to convert")
    .option("-o, --output <path>", "Output file path", ".env")
    .option("-e, --export", "Add 'export' prefix to variables", false)
    .option("--skip-sort", "Do not sort variables alphabetically", false)
    .option("--skip-lower", "Ignore variables with lowercase keys", false)
    .option("--url-only", "Include only variables with HTTP/HTTPS URL values", false)
    .option("-f, --overwrite", "Overwrite output file without backup", false)
    .option("-k, --preserve <keys>", "Comma-separated list of variables to preserve", "")
    .action((source: string, opts) => {
        convert(source, {
            output: opts.output,
            export: opts.export,
            noSort: opts.skipSort,
            noLower: opts.skipLower,
            urlOnly: opts.urlOnly,
            overwrite: opts.overwrite,
            preserve: opts.preserve,
        });
    });

program.parse();
