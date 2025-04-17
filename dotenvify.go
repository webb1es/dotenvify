package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"dotenvify/plugins/azure"
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
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Process variables from Azure DevOps
func processAzureDevOpsVariables(org, project, groupName, outputFile string) {
	// Process variables using the Azure plugin
	_, err := azure.ProcessVariables(org, project, groupName, outputFile, msg)
	if err != nil {
		// Check if the error is related to Azure CLI login
		if strings.Contains(err.Error(), "not logged in to Azure CLI") {
			msg("error", "🔑 Authentication required: You need to log in to Azure CLI first")
			msg("info", "Run 'az login' in your terminal and then try again")
			os.Exit(1)
		} else if strings.Contains(err.Error(), "Azure CLI not found") {
			msg("error", "🔧 Azure CLI not installed: You need to install the Azure CLI")
			msg("info", "Visit https://docs.microsoft.com/en-us/cli/azure/install-azure-cli to install it")
			os.Exit(1)
		} else if strings.Contains(err.Error(), "failed to get access token") {
			msg("error", "🔑 Authentication failed: Could not get access token from Azure CLI")
			msg("info", "Run 'az login' in your terminal and then try again")
			os.Exit(1)
		} else {
			exitOnError(err, "Failed to process Azure DevOps variables")
		}
	}
}

func main() {
	// Define command line flags
	azureMode := flag.Bool("azure", false, "Enable Azure DevOps mode")
	projectURL := flag.String("url", azure.GetDefaultURL(), "Azure DevOps project URL (default: $AZURE_DEVOPS_URL)")
	organization := flag.String("org", azure.GetDefaultOrganization(), "Azure DevOps organization name (default: $AZURE_DEVOPS_ORG)")
	project := flag.String("project", azure.GetDefaultProject(), "Azure DevOps project name (default: $AZURE_DEVOPS_PROJECT)")
	varGroup := flag.String("group", "", "Azure DevOps variable group name (required)")
	outputFilePath := flag.String("output", ".env", "Output file path (default: .env)")

	// Set custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "🧙‍♂️ DotEnvify - Convert key-value pairs to environment variables\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  Local file mode:  dotenvify source_file [output_file]\n")
		fmt.Fprintf(os.Stderr, "  Azure DevOps:     dotenvify -azure -group \"variable-group-name\" [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nFor more information, visit: https://github.com/webb1es/dotenvify\n")
	}

	// Parse command line flags
	flag.Parse()

	// Check if we're in Azure DevOps mode
	if *azureMode {
		org := *organization
		proj := *project

		// If URL is provided, parse it to get org and project
		if *projectURL != "" {
			var err error
			org, proj, err = azure.ParseURL(*projectURL)
			if err != nil {
				exitOnError(err, "Failed to parse Azure DevOps URL")
			}
			msg("info", fmt.Sprintf("🔗 Using Azure DevOps organization: %s, project: %s", org, proj))
		} else if org == "" || proj == "" {
			// If org or project is still empty, prompt for URL
			msg("info", "No Azure DevOps URL or organization/project provided")
			url := readUserInput("Enter your Azure DevOps project URL (e.g., https://dev.azure.com/org/project): ")
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
		processAzureDevOpsVariables(org, proj, *varGroup, *outputFilePath)
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

	// Process the lines into environment variable exports
	var outputLines []string
	var errors []string
	for i := 0; i < len(lines); i += 2 {
		if i+1 < len(lines) {
			key := lines[i]
			value := lines[i+1]
			outputLines = append(outputLines, fmt.Sprintf("export %s=\"%s\"", key, value))
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

	// Write the output file
	msg("info", fmt.Sprintf("💾 Writing %d variables to '%s'...", len(outputLines), outputFile))
	err = os.WriteFile(outputFile, []byte(strings.Join(outputLines, "\n")+"\n"), 0644)
	exitOnError(err, "Failed to write output file")

	// Show success message
	if !hasErrors {
		msg("success", fmt.Sprintf("✨ Variables successfully formatted and saved to '%s'", outputFile))
	}
}
