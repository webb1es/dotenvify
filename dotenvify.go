package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"dotenvify/plugins/azure"
	"github.com/creativeprojects/go-selfupdate"
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

// readExistingVariables reads variables from an existing .env file
func readExistingVariables(filePath string) (map[string]string, error) {
	variables := make(map[string]string)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return variables, nil // Return empty map if file doesn't exist
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Remove "export " prefix if present
		line = strings.TrimPrefix(line, "export ")

		// Parse KEY=VALUE format
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}

			if key != "" {
				variables[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return variables, nil
}

// backupFile creates a backup of an existing file with an incremental counter
func backupFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	counter := 1
	var backupPath string
	for {
		backupPath = fmt.Sprintf("%s.backup.%d", filePath, counter)
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			break
		}
		counter++
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file for backup: %v", err)
	}

	if err := os.WriteFile(backupPath, content, 0600); err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}

	msg("info", fmt.Sprintf("💾 Backed up existing file to '%s'", backupPath))
	return nil
}

// WriteVariablesToFile writes variables to a file, optionally with export prefix
func writeVariablesToFile(variables map[string]string, outputFile string, noLower bool, noSort bool, useExport bool, urlOnly bool, overwrite bool, preserveVars string) error {
	var outputLines []string
	variableCount := 0
	skippedCount := 0
	preservedCount := 0

	// Parse preserve list and read existing variables BEFORE backup
	preserveList := make(map[string]bool)
	if preserveVars != "" {
		for _, varName := range strings.Split(preserveVars, ",") {
			varName = strings.TrimSpace(varName)
			if varName != "" {
				preserveList[varName] = true
			}
		}
	}

	// Read existing variables if preserve list is provided (BEFORE backup)
	var existingVars map[string]string
	if len(preserveList) > 0 {
		var err error
		existingVars, err = readExistingVariables(outputFile)
		if err != nil {
			return fmt.Errorf("failed to read existing variables: %v", err)
		}

		// Replace new values with existing ones for preserved variables
		for varName := range preserveList {
			if existingValue, exists := existingVars[varName]; exists {
				variables[varName] = existingValue
				preservedCount++
			}
			// If variable doesn't exist in the file, it will be written with the new value
		}
	}

	// Handle existing file backup AFTER reading preserved values
	if !overwrite {
		if err := backupFile(outputFile); err != nil {
			return err
		}
	}

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

		// Check if the value needs quoting (URLs or values with spaces) and not already quoted
		if (isURL(value) || strings.Contains(value, " ")) && !isQuoted(value) {
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

	if preservedCount > 0 {
		msg("info", fmt.Sprintf("🔒 Preserved %d existing variables", preservedCount))
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
func ProcessVariables(variables map[string]string, outputFile string, noLower bool, noSort bool, useExport bool, urlOnly bool, overwrite bool, preserveVars string) {
	// Write variables to file
	err := writeVariablesToFile(variables, outputFile, noLower, noSort, useExport, urlOnly, overwrite, preserveVars)
	if err != nil {
		exitOnError(err, "Failed to write variables to file")
	}
}

// handleUpdate handles both checking and performing updates
func handleUpdate(performUpdate bool) error {
	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("webb1es/dotenvify"))
	if err != nil {
		return fmt.Errorf("failed to check for updates: %v", err)
	}
	if !found {
		return fmt.Errorf("no releases found")
	}

	if version == "dev" {
		msg("warning", "⚠️ Running development version")
		msg("info", fmt.Sprintf("Latest available version: %s", latest.Version()))
		return nil
	}

	if latest.LessOrEqual(version) {
		msg("success", fmt.Sprintf("✨ Already up to date (version %s)", version))
		return nil
	}

	msg("info", fmt.Sprintf("📦 New version available: %s → %s", version, latest.Version()))
	if !performUpdate {
		msg("info", "Run 'dotenvify -update' to update")
		return nil
	}

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	msg("success", fmt.Sprintf("✨ Successfully updated to version %s", latest.Version()))
	return nil
}

// Process variables from Azure DevOps
func processAzureDevOpsVariables(org, project, groupName, outputFile string, noLower bool, noSort bool, useExport bool, urlOnly bool, overwrite bool, preserveVars string) {
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
	ProcessVariables(variables, outputFile, noLower, noSort, useExport, urlOnly, overwrite, preserveVars)
}

func main() {
	// Define command line flags
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.BoolVar(versionFlag, "v", false, "Show version information (shorthand)")

	updateFlag := flag.Bool("update", false, "Update dotenvify to the latest version")
	flag.BoolVar(updateFlag, "up", false, "Update dotenvify to the latest version (shorthand)")

	checkUpdateFlag := flag.Bool("check-update", false, "Check if a new version is available")
	flag.BoolVar(checkUpdateFlag, "cu", false, "Check if a new version is available (shorthand)")

	azureMode := flag.Bool("azure", false, "Enable Azure DevOps mode")
	flag.BoolVar(azureMode, "az", false, "Enable Azure DevOps mode (shorthand)")

	organization := flag.String("org", "", "Azure DevOps organization URL (e.g., https://dev.azure.com/org/project)")
	flag.StringVar(organization, "o", "", "Azure DevOps organization URL (shorthand)")

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

	overwrite := flag.Bool("overwrite", false, "Overwrite output file if it exists")
	flag.BoolVar(overwrite, "f", false, "Overwrite output file if it exists (shorthand)")

	preserve := flag.String("preserve", "", "Comma-separated list of variables to preserve (keep unchanged if they exist)")
	flag.StringVar(preserve, "k", "", "Comma-separated list of variables to preserve (shorthand)")

	// Set custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "🧙‍♂️ DotEnvify - Convert key-value pairs to environment variables\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  Local file mode:  dotenvify source_file [output_file]\n")
		fmt.Fprintf(os.Stderr, "  Azure DevOps:     dotenvify -azure -group \"group1,group2,group3\" [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")

		// Custom flag printing to combine short and long forms
		fmt.Fprintf(os.Stderr, "  -v (version)\t\tShow version information\n")
		fmt.Fprintf(os.Stderr, "  -up (update)\t\tUpdate dotenvify to the latest version\n")
		fmt.Fprintf(os.Stderr, "  -cu (check-update)\tCheck if a new version is available\n")
		fmt.Fprintf(os.Stderr, "  -az (azure)\t\tEnable Azure DevOps mode\n")
		fmt.Fprintf(os.Stderr, "  -o (org)\tAzure DevOps organization URL (or set DOTENVIFY_DEFAULT_ORG_URL env var)\n")
		fmt.Fprintf(os.Stderr, "  -g (group)\tAzure DevOps variable group name(s) - comma-separated for multiple groups\n")
		fmt.Fprintf(os.Stderr, "  -out (output)\tOutput file path (default: .env)\n")
		fmt.Fprintf(os.Stderr, "  -nl (no-lower)\tIgnore variables with lowercase keys\n")
		fmt.Fprintf(os.Stderr, "  -ns (no-sort)\tDo not sort variables alphabetically\n")
		fmt.Fprintf(os.Stderr, "  -e (export)\tAdd 'export' prefix to variables\n")
		fmt.Fprintf(os.Stderr, "  -urls (url-only)\tInclude only variables with HTTP/HTTPS URL values\n")
		fmt.Fprintf(os.Stderr, "  -f (overwrite)\tOverwrite output file if it exists\n")
		fmt.Fprintf(os.Stderr, "  -k (preserve)\tComma-separated list of variables to preserve (keep unchanged)\n")
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

	// Handle update flags
	if *updateFlag || *checkUpdateFlag {
		exitOnError(handleUpdate(*updateFlag), "Update check failed")
		os.Exit(0)
	}

	// Check if we're in Azure DevOps mode
	if *azureMode {
		var org, proj string
		orgURL := *organization

		// If org URL is provided, parse it to get org and project
		if orgURL != "" {
			var err error
			org, proj, err = azure.ParseURL(orgURL)
			if err != nil {
				exitOnError(err, "Failed to parse Azure DevOps URL")
			}
			msg("info", fmt.Sprintf("🔗 Using Azure DevOps organization: %s, project: %s", org, proj))
		} else {
			// Check for default URL from environment variable
			defaultURL := os.Getenv("DOTENVIFY_DEFAULT_ORG_URL")
			if defaultURL == "" {
				msg("error", "Azure DevOps organization URL is required")
				msg("info", "Please provide URL using one of the following methods:")
				msg("info", "  1. CLI flag: -org or -o")
				msg("info", "  2. Environment variable: DOTENVIFY_DEFAULT_ORG_URL")
				msg("info", "  3. Interactive input (prompted below)")

				orgURL = readUserInput("Enter your Azure DevOps project URL (e.g., https://dev.azure.com/org/project): ")
				if orgURL == "" {
					exitOnError(fmt.Errorf("no URL provided"), "Azure DevOps URL is required")
				}
			} else {
				orgURL = defaultURL
				msg("info", fmt.Sprintf("🔗 Using default Azure DevOps URL from DOTENVIFY_DEFAULT_ORG_URL: %s", defaultURL))
			}

			var err error
			org, proj, err = azure.ParseURL(orgURL)
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
		processAzureDevOpsVariables(org, proj, *varGroup, *outputFilePath, *noLower, *noSort, *useExport, *urlOnly, *overwrite, *preserve)
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
	ProcessVariables(variables, outputFile, *noLower, *noSort, *useExport, *urlOnly, *overwrite, *preserve)
}
