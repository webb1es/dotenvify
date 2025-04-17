package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// Azure DevOps API structures
type VariableGroup struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Variables   map[string]Variable `json:"variables"`
	Type        string              `json:"type"`
	IsShared    bool                `json:"isShared"`
	Description string              `json:"description,omitempty"`
}

type Variable struct {
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret"`
	Enabled  bool   `json:"enabled"`
}

// Azure DevOps API client
type AzureDevOpsClient struct {
	Organization string
	Project      string
}

// Create a new Azure DevOps client
func NewAzureDevOpsClient(org, project string) *AzureDevOpsClient {
	return &AzureDevOpsClient{
		Organization: org,
		Project:      project,
	}
}

// Get access token from Azure CLI
func getAzureDevOpsToken() (string, error) {
	// Check if Azure CLI is installed
	_, err := exec.LookPath("az")
	if err != nil {
		return "", fmt.Errorf("Azure CLI not found. Please install it from https://docs.microsoft.com/en-us/cli/azure/install-azure-cli")
	}

	// Check if user is logged in
	cmd := exec.Command("az", "account", "show")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("not logged in to Azure CLI. Please run 'az login' first")
	}

	// Get access token
	cmd = exec.Command("az", "account", "get-access-token", "--resource", "499b84ac-1321-427f-aa17-267ca6975798")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get access token: %v", err)
	}

	var response struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		return "", fmt.Errorf("failed to parse access token: %v", err)
	}

	return response.AccessToken, nil
}

// Get variable groups from Azure DevOps
func (c *AzureDevOpsClient) GetVariableGroups() ([]VariableGroup, error) {
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/distributedtask/variablegroups?api-version=6.0-preview.2",
		c.Organization, c.Project)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Get token from Azure CLI
	token, err := getAzureDevOpsToken()
	if err != nil {
		return nil, err
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Count int             `json:"count"`
		Value []VariableGroup `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Value, nil
}

// Get a specific variable group by name
func (c *AzureDevOpsClient) GetVariableGroupByName(name string) (*VariableGroup, error) {
	groups, err := c.GetVariableGroups()
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.Name == name {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("variable group '%s' not found", name)
}

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

// Get default organization name
func getDefaultOrganization() string {
	// Check environment variable
	if org := os.Getenv("AZURE_DEVOPS_ORG"); org != "" {
		return org
	}
	return ""
}

// Get default project name
func getDefaultProject() string {
	// Check environment variable
	if proj := os.Getenv("AZURE_DEVOPS_PROJECT"); proj != "" {
		return proj
	}
	return ""
}

// Get default variable group name (based on current directory)
func getDefaultVariableGroup() string {
	// Use current directory name as default
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Base(dir)
}

// Parse Azure DevOps URL to extract organization and project
func parseAzureDevOpsURL(url string) (string, string, error) {
	// URL format: https://dev.azure.com/{organization}/{project}
	// or: https://{organization}.visualstudio.com/{project}

	// Try dev.azure.com format first
	re := regexp.MustCompile(`https?://dev\.azure\.com/([^/]+)/([^/]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) == 3 {
		return matches[1], matches[2], nil
	}

	// Try visualstudio.com format
	re = regexp.MustCompile(`https?://([^.]+)\.visualstudio\.com/([^/]+)`)
	matches = re.FindStringSubmatch(url)
	if len(matches) == 3 {
		return matches[1], matches[2], nil
	}

	return "", "", fmt.Errorf("invalid Azure DevOps URL format. Expected: https://dev.azure.com/{organization}/{project} or https://{organization}.visualstudio.com/{project}")
}

// Get default project URL from environment variable
func getDefaultProjectURL() string {
	return os.Getenv("AZURE_DEVOPS_URL")
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
	// Create Azure DevOps client
	client := NewAzureDevOpsClient(org, project)

	// Get variable group
	msg("info", fmt.Sprintf("🔍 Fetching variable group '%s' from Azure DevOps...", groupName))
	group, err := client.GetVariableGroupByName(groupName)
	if err != nil {
		exitOnError(err, "Failed to fetch variable group")
	}

	// Process variables
	var outputLines []string
	variableCount := 0

	msg("info", fmt.Sprintf("🧩 Processing variables from group '%s'...", groupName))

	for key, variable := range group.Variables {
		if !variable.Enabled {
			continue
		}

		outputLines = append(outputLines, fmt.Sprintf("export %s=\"%s\"", key, variable.Value))
		variableCount++
	}

	// Write the output file
	msg("info", fmt.Sprintf("💾 Writing %d variables to '%s'...", variableCount, outputFile))
	err = os.WriteFile(outputFile, []byte(strings.Join(outputLines, "\n")+"\n"), 0644)
	exitOnError(err, "Failed to write output file")

	// Show success message
	msg("success", fmt.Sprintf("✨ Variables successfully fetched and saved to '%s'", outputFile))
}

func main() {
	// Define command line flags
	azureMode := flag.Bool("azure", false, "Enable Azure DevOps mode")
	projectURL := flag.String("url", getDefaultProjectURL(), "Azure DevOps project URL (default: $AZURE_DEVOPS_URL)")
	organization := flag.String("org", getDefaultOrganization(), "Azure DevOps organization name (default: $AZURE_DEVOPS_ORG)")
	project := flag.String("project", getDefaultProject(), "Azure DevOps project name (default: $AZURE_DEVOPS_PROJECT)")
	varGroup := flag.String("group", getDefaultVariableGroup(), "Azure DevOps variable group name (default: current directory name)")
	outputFilePath := flag.String("output", ".env", "Output file path (default: .env)")

	// Hidden flags for backward compatibility, not used anymore
	_ = flag.String("pat", "", "Deprecated: Azure DevOps Personal Access Token (not used, authentication is now handled by Azure CLI)")
	_ = flag.String("env", "", "Deprecated: Environment filtering is no longer supported")

	// Parse command line flags
	flag.Parse()

	// Check if we're in Azure DevOps mode
	if *azureMode {
		org := *organization
		proj := *project

		// If URL is provided, parse it to get org and project
		if *projectURL != "" {
			var err error
			org, proj, err = parseAzureDevOpsURL(*projectURL)
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
			org, proj, err = parseAzureDevOpsURL(url)
			if err != nil {
				exitOnError(err, "Failed to parse Azure DevOps URL")
			}
			msg("info", fmt.Sprintf("🔗 Using Azure DevOps organization: %s, project: %s", org, proj))
		}

		// Validate variable group
		if *varGroup == "" {
			exitOnError(fmt.Errorf("no variable group provided"), "Failed to get variable group")
		}

		// Process Azure DevOps variables
		processAzureDevOpsVariables(org, proj, *varGroup, *outputFilePath)
		return
	}

	// Traditional file mode
	args := flag.Args()
	if len(args) < 1 {
		msg("error", "No source file provided. Usage: dotenvify source_file [output_file]")
		flag.Usage()
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
