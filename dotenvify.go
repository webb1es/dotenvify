package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"dotenvify/plugins/azure"
)

// Version information (populated by ldflags during build)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// Print a styled message
func msg(msgType, message string) {
	prefix := message
	switch msgType {
	case "success":
		prefix = fmt.Sprintf("%s✅ %s%s", colorGreen, message, colorReset)
	case "error":
		prefix = fmt.Sprintf("%s❌ %s%s", colorRed, message, colorReset)
	case "warning":
		prefix = fmt.Sprintf("%s⚠️ %s%s", colorYellow, message, colorReset)
	case "info":
		prefix = fmt.Sprintf("%sℹ️ %s%s", colorYellow, message, colorReset)
	case "debug":
		prefix = fmt.Sprintf("%s🔍 %s%s", colorYellow, message, colorReset)
	}
	fmt.Println(prefix)
}

// Handle errors with exit
func exitOnError(err error, message string) {
	if err != nil {
		msg("error", message+": "+err.Error())
		os.Exit(1)
	}
}

// Read user input from stdin
func readUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		msg("warning", fmt.Sprintf("Error reading input: %v", err))
		return ""
	}
	return strings.TrimSpace(input)
}

// isURL checks if a string is likely a URL
func isURL(s string) bool {
	// Check for common URL prefixes
	return strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "ftp://") ||
		strings.HasPrefix(s, "sftp://") ||
		strings.HasPrefix(s, "ssh://") ||
		strings.HasPrefix(s, "git://") ||
		strings.HasPrefix(s, "file://") ||
		strings.HasPrefix(s, "mailto:") ||
		strings.HasPrefix(s, "postgres://") ||
		strings.HasPrefix(s, "mysql://") ||
		strings.HasPrefix(s, "mongodb://") ||
		strings.HasPrefix(s, "redis://")
}

// isHTTPURL checks if a string is specifically an HTTP/HTTPS URL
func isHTTPURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// isQuoted checks if a string is already quoted
func isQuoted(s string) bool {
	return (strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) ||
		(strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'"))
}

// WriteVariablesToFile writes variables to a file, optionally with export prefix
func writeVariablesToFile(variables map[string]string, outputFile string, noLower bool, noSort bool, useExport bool, urlOnly bool) error {
	var outputLines []string
	variableCount := 0
	skippedCount := 0

	// Collect the keys
	var keys []string
	for key := range variables {
		// Skip lowercase keys if noLower is true
		if noLower && key == strings.ToLower(key) && key != strings.ToUpper(key) {
			skippedCount++
			continue
		}
		// Skip non-URL values if urlOnly is true
		if urlOnly && !isHTTPURL(variables[key]) {
			skippedCount++
			continue
		}
		keys = append(keys, key)
		variableCount++
	}

	// Sort keys alphabetically unless noSort is true
	if !noSort {
		sort.Strings(keys)
	}

	// Create output lines using sorted keys
	for _, key := range keys {
		value := variables[key]

		// Check if the value is a URL and not already quoted
		if isURL(value) && !isQuoted(value) {
			value = fmt.Sprintf("\"%s\"", value)
		}

		if useExport {
			outputLines = append(outputLines, fmt.Sprintf("export %s=%s", key, value))
		} else {
			outputLines = append(outputLines, fmt.Sprintf("%s=%s", key, value))
		}
	}

	if skippedCount > 0 {
		if urlOnly {
			msg("info", fmt.Sprintf("🔍 Skipped %d variables (filtered)", skippedCount))
		} else {
			msg("info", fmt.Sprintf("🔍 Skipped %d variables with lowercase keys", skippedCount))
		}
	}

	// Write the output file
	msg("info", fmt.Sprintf("💾 Writing %d variables to '%s'...", variableCount, outputFile))
	// Use more restrictive file permissions (0600) for files containing sensitive data
	err := os.WriteFile(outputFile, []byte(strings.Join(outputLines, "\n")+"\n"), 0600)
	if err != nil {
		return fmt.Errorf("Failed to write output file: %v", err)
	}

	// Show a success message
	msg("success", fmt.Sprintf("✨ Variables successfully saved to '%s'", outputFile))

	return nil
}

// ProcessVariables processes variables from a source and writes them to a file
func ProcessVariables(variables map[string]string, outputFile string, noLower bool, noSort bool, useExport bool, urlOnly bool) {
	// Write variables to file
	err := writeVariablesToFile(variables, outputFile, noLower, noSort, useExport, urlOnly)
	if err != nil {
		exitOnError(err, "Failed to write variables to file")
	}
}

// Process variables from Azure DevOps
func processAzureDevOpsVariables(org, project, groupName, outputFile string, noLower bool, noSort bool, useExport bool, urlOnly bool) {
	// Process variables using the Azure plugin
	variables, err := azure.GetVariables(org, project, groupName, msg)
	if err != nil {
		// Check if the error is related to Azure CLI login
		errMsg := err.Error()
		if strings.Contains(errMsg, "not logged in to Azure CLI") {
			msg("error", "🔑 Authentication required: You need to log in to Azure CLI first")
			msg("info", "Run 'az login' in your terminal and then try again")
			os.Exit(1)
		} else if strings.Contains(errMsg, "Azure CLI not found") {
			msg("error", "🔧 Azure CLI not installed: You need to install the Azure CLI")
			msg("info", "Visit https://docs.microsoft.com/en-us/cli/azure/install-azure-cli to install it")
			os.Exit(1)
		} else if strings.Contains(errMsg, "failed to get access token") {
			msg("error", "🔑 Authentication failed: Could not get access token from Azure CLI")
			msg("info", "Run 'az login' in your terminal and then try again")
			os.Exit(1)
		} else {
			exitOnError(err, "Failed to process Azure DevOps variables")
		}
	}

	// Process variables
	ProcessVariables(variables, outputFile, noLower, noSort, useExport, urlOnly)
}

func main() {
	// Define command line flags
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.BoolVar(versionFlag, "v", false, "Show version information (shorthand)")

	azureMode := flag.Bool("azure", false, "Enable Azure DevOps mode")
	flag.BoolVar(azureMode, "az", false, "Enable Azure DevOps mode (shorthand)")

	projectURL := flag.String("url", "", "Azure DevOps project URL")
	flag.StringVar(projectURL, "u", "", "Azure DevOps project URL (shorthand)")

	organization := flag.String("org", "", "Azure DevOps organization name (inferred from URL if not provided)")
	flag.StringVar(organization, "o", "", "Azure DevOps organization name (shorthand)")

	project := flag.String("project", "", "Azure DevOps project name (inferred from URL if not provided)")
	flag.StringVar(project, "p", "", "Azure DevOps project name (shorthand)")

	varGroup := flag.String("group", "", "Azure DevOps variable group name (required)")
	flag.StringVar(varGroup, "g", "", "Azure DevOps variable group name (shorthand)")

	outputFilePath := flag.String("output", ".env", "Output file path (default: .env)")
	flag.StringVar(outputFilePath, "out", ".env", "Output file path (shorthand)")

	noLower := flag.Bool("no-lower", false, "Ignore variables with lowercase keys")
	flag.BoolVar(noLower, "nl", false, "Ignore variables with lowercase keys (shorthand)")

	noSort := flag.Bool("no-sort", false, "Do not sort variables alphabetically")
	flag.BoolVar(noSort, "ns", false, "Do not sort variables alphabetically (shorthand)")

	useExport := flag.Bool("export", false, "Add 'export' prefix to variables")
	flag.BoolVar(useExport, "e", false, "Add 'export' prefix to variables (shorthand)")

	urlOnly := flag.Bool("url-only", false, "Include only variables with HTTP/HTTPS URL values")
	flag.BoolVar(urlOnly, "urls", false, "Include only variables with HTTP/HTTPS URL values (shorthand)")

	// Set custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "🧙‍♂️ DotEnvify - Convert key-value pairs to environment variables\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  Local file mode:  dotenvify source_file [output_file]\n")
		fmt.Fprintf(os.Stderr, "  Azure DevOps:     dotenvify -azure -group \"group1,group2,group3\" [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")

		// Custom flag printing to combine short and long forms
		fmt.Fprintf(os.Stderr, "  -v (version)\tShow version information\n")
		fmt.Fprintf(os.Stderr, "  -az (azure)\tEnable Azure DevOps mode\n")
		fmt.Fprintf(os.Stderr, "  -u (url)\tAzure DevOps project URL (or set DOTENVIFY_DEFAULT_ORG_URL env var)\n")
		fmt.Fprintf(os.Stderr, "  -o (org)\tAzure DevOps organization name (inferred from URL if not provided)\n")
		fmt.Fprintf(os.Stderr, "  -p (project)\tAzure DevOps project name (inferred from URL if not provided)\n")
		fmt.Fprintf(os.Stderr, "  -g (group)\tAzure DevOps variable group name(s) - comma-separated for multiple groups\n")
		fmt.Fprintf(os.Stderr, "  -out (output)\tOutput file path (default: .env)\n")
		fmt.Fprintf(os.Stderr, "  -nl (no-lower)\tIgnore variables with lowercase keys\n")
		fmt.Fprintf(os.Stderr, "  -ns (no-sort)\tDo not sort variables alphabetically\n")
		fmt.Fprintf(os.Stderr, "  -e (export)\tAdd 'export' prefix to variables\n")
		fmt.Fprintf(os.Stderr, "  -urls (url-only)\tInclude only variables with HTTP/HTTPS URL values\n")
		fmt.Fprintf(os.Stderr, "  -h (help)\tShow this help message\n")

		fmt.Fprintf(os.Stderr, "\nFor more information, visit: https://github.com/webb1es/dotenvify\n")
	}

	// Parse command line flags
	flag.Parse()

	// Check if a version flag was provided
	if *versionFlag {
		fmt.Printf("dotenvify version: %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("build date: %s\n", date)
		os.Exit(0)
	}

	// Check if we're in Azure DevOps mode
	if *azureMode {
		org := *organization
		proj := *project
		url := *projectURL

		// If URL is provided, parse it to get org and project
		if url != "" {
			var err error
			org, proj, err = azure.ParseURL(url)
			if err != nil {
				exitOnError(err, "Failed to parse Azure DevOps URL")
			}
			msg("info", fmt.Sprintf("🔗 Using Azure DevOps organization: %s, project: %s", org, proj))
		} else if org == "" || proj == "" {
			// Check for default URL from environment variable
			defaultURL := os.Getenv("DOTENVIFY_DEFAULT_ORG_URL")
			if defaultURL == "" {
				msg("error", "Azure DevOps URL is required")
				msg("info", "Please provide URL using one of the following methods:")
				msg("info", "  1. CLI flag: -url or -u")
				msg("info", "  2. Environment variable: DOTENVIFY_DEFAULT_ORG_URL")
				msg("info", "  3. Interactive input (prompted below)")

				url = readUserInput("Enter your Azure DevOps project URL (e.g., https://dev.azure.com/org/project): ")
				if url == "" {
					exitOnError(fmt.Errorf("no URL provided"), "Azure DevOps URL is required")
				}
			} else {
				url = defaultURL
				msg("info", fmt.Sprintf("🔗 Using default Azure DevOps URL from DOTENVIFY_DEFAULT_ORG_URL: %s", defaultURL))
			}

			var err error
			org, proj, err = azure.ParseURL(url)
			if err != nil {
				exitOnError(err, "Failed to parse Azure DevOps URL")
			}
			msg("info", fmt.Sprintf("🔗 Using Azure DevOps organization: %s, project: %s", org, proj))
		}

		// Validate variable group
		if *varGroup == "" {
			msg("info", "No variable group name provided")
			groupName := readUserInput("Enter your Azure DevOps variable group name: ")
			if groupName == "" {
				exitOnError(fmt.Errorf("no variable group provided"), "Failed to get variable group")
			}
			*varGroup = groupName
		}

		// Process Azure DevOps variables
		processAzureDevOpsVariables(org, proj, *varGroup, *outputFilePath, *noLower, *noSort, *useExport, *urlOnly)
		return
	}

	// Traditional file mode
	args := flag.Args()
	if len(args) < 1 {
		msg("error", "No source file provided. Use -h or -help for usage information.")
		os.Exit(1)
	}

	sourceFile := args[0]
	outputFile := *outputFilePath
	if len(args) > 1 {
		outputFile = args[1]
	}

	msg("info", fmt.Sprintf("📄 Processing source file: %s", sourceFile))

	// Validate file paths to prevent path traversal
	if strings.Contains(sourceFile, "..") || strings.Contains(outputFile, "..") {
		exitOnError(fmt.Errorf("path contains potentially unsafe '..' sequence"), "Invalid file path")
	}

	// Check if source file exists
	_, err := os.Stat(sourceFile)
	exitOnError(err, fmt.Sprintf("Source file '%s' does not exist", sourceFile))

	// Read source file
	file, err := os.Open(sourceFile)
	exitOnError(err, "Failed to open source file")
	defer file.Close()

	// Read lines, skip empty ones
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	exitOnError(scanner.Err(), "Error reading source file")

	// Check if the file is already in the expected output format
	alreadyInCorrectFormat := true
	for _, line := range lines {
		// Skip comment lines
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Check if line follows the expected pattern based on useExport flag
		if *useExport {
			// If useExport is true, expect "export KEY=VALUE"
			if !strings.HasPrefix(line, "export ") || !strings.Contains(line, "=") {
				alreadyInCorrectFormat = false
				break
			}
		} else {
			// If useExport is false, expect "KEY=VALUE" without export prefix
			if strings.HasPrefix(line, "export ") || !strings.Contains(line, "=") {
				alreadyInCorrectFormat = false
				break
			}
		}
	}

	if alreadyInCorrectFormat && len(lines) > 0 && !*urlOnly && !*noLower {
		msg("success", fmt.Sprintf("✨ File '%s' is already in the expected format!", sourceFile))
		if sourceFile != outputFile {
			// If output file is different, copy the source file to the output file
			sourceContent, err := os.ReadFile(sourceFile)
			exitOnError(err, "Failed to read source file")
			err = os.WriteFile(outputFile, sourceContent, 0600)
			exitOnError(err, "Failed to write output file")
			msg("success", fmt.Sprintf("✨ File copied to '%s'", outputFile))
		}
		return
	}

	msg("info", fmt.Sprintf("🔄 Converting %d lines to environment variables...", len(lines)))

	// Parse the lines into a map of variables
	variables := make(map[string]string)
	var errors []string

	// Try to parse in various formats
	i := 0
	for i < len(lines) {
		line := lines[i]

		// Skip comment lines
		if strings.HasPrefix(line, "#") {
			i++
			continue
		}

		// Check for KEY=VALUE or KEY="VALUE" format
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Handle quoted values: KEY="VALUE" or KEY='VALUE'
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}

			if key != "" {
				variables[key] = value
			}
			i++
			continue
		}

		// Check for KEY VALUE format on the same line
		parts := strings.Fields(line)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			if key != "" {
				variables[key] = value
			}
			i++
			continue
		}

		// Check for multiple KEY VALUE pairs on the same line
		if len(parts) >= 4 && len(parts)%2 == 0 {
			for j := 0; j < len(parts); j += 2 {
				key := parts[j]
				value := parts[j+1]
				if key != "" {
					variables[key] = value
				}
			}
			i++
			continue
		}

		// Traditional format: KEY on one line, VALUE on the next
		if i+1 < len(lines) {
			key := line
			value := lines[i+1]
			if key != "" {
				variables[key] = value
			}
			i += 2
		} else {
			errors = append(errors, fmt.Sprintf("Line %d: Key '%s' has no value", i+1, line))
			i++
		}
	}

	// Handle errors if any
	hasErrors := len(errors) > 0
	if hasErrors {
		if outputFile == sourceFile {
			outputFile = sourceFile + ".out"
		}
		msg("warning", fmt.Sprintf("⚠️ Some errors occurred. Output saved to '%s'", outputFile))
		for _, err := range errors {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	// Process variables
	ProcessVariables(variables, outputFile, *noLower, *noSort, *useExport, *urlOnly)
}
