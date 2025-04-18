package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

// Read user input and ask if they want to store it in an environment variable
func readUserInputWithEnvOption(prompt, envVar string) string {
	value := readUserInput(prompt)
	if value != "" {
		saveToEnv := readUserInput(fmt.Sprintf("Would you like to save this to %s environment variable? [Y/n]: ", envVar))
		if saveToEnv == "" || strings.ToLower(saveToEnv) == "y" || strings.ToLower(saveToEnv) == "yes" {
			// Set for current process
			err := os.Setenv(envVar, value)
			if err != nil {
				msg("warning", fmt.Sprintf("Failed to set environment variable %s: %v", envVar, err))
			} else {
				msg("info", fmt.Sprintf("🔒 Saved to %s environment variable for this session", envVar))
				// Provide instructions for permanent setting
				shellType := os.Getenv("SHELL")
				if strings.Contains(shellType, "bash") {
					msg("info", fmt.Sprintf("To save permanently, add this to your ~/.bashrc: export %s=\"%s\"", envVar, value))
				} else if strings.Contains(shellType, "zsh") {
					msg("info", fmt.Sprintf("To save permanently, add this to your ~/.zshrc: export %s=\"%s\"", envVar, value))
				} else {
					msg("info", fmt.Sprintf("To save permanently, add this to your shell profile: export %s=\"%s\"", envVar, value))
				}
			}
		}
	}
	return value
}

// WriteVariablesToFile writes variables to a file in export format
func writeVariablesToFile(variables map[string]string, outputFile string, noLower bool) error {
	var outputLines []string
	variableCount := 0
	skippedCount := 0

	for key, value := range variables {
		// Skip lowercase keys if noLower is true
		if noLower && key == strings.ToLower(key) && key != strings.ToUpper(key) {
			skippedCount++
			continue
		}
		outputLines = append(outputLines, fmt.Sprintf("export %s=\"%s\"", key, value))
		variableCount++
	}

	if skippedCount > 0 {
		msg("info", fmt.Sprintf("🔍 Skipped %d variables with lowercase keys", skippedCount))
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
func ProcessVariables(variables map[string]string, outputFile string, noLower bool) {
	// Write variables to file
	err := writeVariablesToFile(variables, outputFile, noLower)
	if err != nil {
		exitOnError(err, "Failed to write variables to file")
	}
}

// Process variables from Azure DevOps
func processAzureDevOpsVariables(org, project, groupName, outputFile string, noLower bool) {
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
	ProcessVariables(variables, outputFile, noLower)
}

func main() {
	// Define command line flags
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.BoolVar(versionFlag, "v", false, "Show version information (shorthand)")

	azureMode := flag.Bool("azure", false, "Enable Azure DevOps mode")
	flag.BoolVar(azureMode, "az", false, "Enable Azure DevOps mode (shorthand)")

	projectURL := flag.String("url", azure.GetDefaultURL(), "Azure DevOps project URL (default: $AZURE_DEVOPS_URL)")
	flag.StringVar(projectURL, "u", azure.GetDefaultURL(), "Azure DevOps project URL (shorthand)")

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

	// Set custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "🧙‍♂️ DotEnvify - Convert key-value pairs to environment variables\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  Local file mode:  dotenvify source_file [output_file]\n")
		fmt.Fprintf(os.Stderr, "  Azure DevOps:     dotenvify -azure -group \"variable-group-name\" [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")

		// Custom flag printing to combine short and long forms
		fmt.Fprintf(os.Stderr, "  -v (version)\tShow version information\n")
		fmt.Fprintf(os.Stderr, "  -az (azure)\tEnable Azure DevOps mode\n")
		fmt.Fprintf(os.Stderr, "  -u (url)\tAzure DevOps project URL (default: $AZURE_DEVOPS_URL)\n")
		fmt.Fprintf(os.Stderr, "  -o (org)\tAzure DevOps organization name (inferred from URL if not provided)\n")
		fmt.Fprintf(os.Stderr, "  -p (project)\tAzure DevOps project name (inferred from URL if not provided)\n")
		fmt.Fprintf(os.Stderr, "  -g (group)\tAzure DevOps variable group name (required)\n")
		fmt.Fprintf(os.Stderr, "  -out (output)\tOutput file path (default: .env)\n")
		fmt.Fprintf(os.Stderr, "  -nl (no-lower)\tIgnore variables with lowercase keys\n")
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

		// Check for URL from environment variable
		envURL := azure.GetDefaultURL()
		if url == "" && envURL != "" {
			url = envURL
			msg("info", fmt.Sprintf("🔗 Using Azure DevOps URL from AZURE_DEVOPS_URL environment variable: %s", url))
		}

		// If URL is provided (from args or env), parse it to get org and project
		if url != "" {
			var err error
			org, proj, err = azure.ParseURL(url)
			if err != nil {
				exitOnError(err, "Failed to parse Azure DevOps URL")
			}
			msg("info", fmt.Sprintf("🔗 Using Azure DevOps organization: %s, project: %s", org, proj))
		} else if org == "" || proj == "" {
			// If org or project is still empty, prompt for URL
			msg("info", "No Azure DevOps URL provided")
			url = readUserInputWithEnvOption("Enter your Azure DevOps project URL (e.g., https://dev.azure.com/org/project): ", "AZURE_DEVOPS_URL")
			if url == "" {
				exitOnError(fmt.Errorf("no URL provided"), "Failed to get Azure DevOps URL")
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
		processAzureDevOpsVariables(org, proj, *varGroup, *outputFilePath, *noLower)
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

	msg("info", fmt.Sprintf("🔄 Converting %d lines to environment variables...", len(lines)))

	// Process the lines into a map of variables
	variables := make(map[string]string)
	var errors []string
	for i := 0; i < len(lines); i += 2 {
		if i+1 < len(lines) {
			key := lines[i]
			value := lines[i+1]
			variables[key] = value
		} else {
			errors = append(errors, fmt.Sprintf("Line %d: Key '%s' has no value", i+1, lines[i]))
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
	ProcessVariables(variables, outputFile, *noLower)
}
