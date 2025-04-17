package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// VariableGroup represents an Azure DevOps variable group
type VariableGroup struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Variables   map[string]Variable `json:"variables"`
	Type        string              `json:"type"`
	IsShared    bool                `json:"isShared"`
	Description string              `json:"description,omitempty"`
}

// Variable represents a variable in an Azure DevOps variable group
type Variable struct {
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret"`
	Enabled  bool   `json:"enabled"`
}

// Client is the Azure DevOps API client
type Client struct {
	Organization string
	Project      string
}

// NewClient creates a new Azure DevOps client
func NewClient(org, project string) *Client {
	return &Client{
		Organization: org,
		Project:      project,
	}
}

// getToken gets an access token from Azure CLI
func getToken() (string, error) {
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

// GetVariableGroups gets all variable groups from Azure DevOps
func (c *Client) GetVariableGroups() ([]VariableGroup, error) {
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/distributedtask/variablegroups?api-version=6.0-preview.2",
		c.Organization, c.Project)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Get token from Azure CLI
	token, err := getToken()
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

// GetVariableGroupByName gets a specific variable group by name
func (c *Client) GetVariableGroupByName(name string) (*VariableGroup, error) {
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

// GetDefaultOrganization gets the default organization name from environment variables
func GetDefaultOrganization() string {
	// Check environment variable
	if org := os.Getenv("AZURE_DEVOPS_ORG"); org != "" {
		return org
	}
	return ""
}

// GetDefaultProject gets the default project name from environment variables
func GetDefaultProject() string {
	// Check environment variable
	if proj := os.Getenv("AZURE_DEVOPS_PROJECT"); proj != "" {
		return proj
	}
	return ""
}

// ParseURL parses an Azure DevOps URL to extract organization and project
func ParseURL(url string) (string, string, error) {
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

// GetDefaultURL gets the default project URL from environment variables
func GetDefaultURL() string {
	return os.Getenv("AZURE_DEVOPS_URL")
}

// ProcessVariables processes variables from Azure DevOps
func ProcessVariables(org, project, groupName, outputFile string, msgFunc func(string, string)) ([]string, error) {
	// Create Azure DevOps client
	client := NewClient(org, project)

	// Get variable group
	msgFunc("info", fmt.Sprintf("🔍 Fetching variable group '%s' from Azure DevOps...", groupName))
	group, err := client.GetVariableGroupByName(groupName)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch variable group: %v", err)
	}

	// Process variables
	var outputLines []string
	variableCount := 0

	msgFunc("info", fmt.Sprintf("🧩 Processing variables from group '%s'...", groupName))

	for key, variable := range group.Variables {
		outputLines = append(outputLines, fmt.Sprintf("export %s=\"%s\"", key, variable.Value))
		variableCount++
	}

	// Write the output file
	msgFunc("info", fmt.Sprintf("💾 Writing %d variables to '%s'...", variableCount, outputFile))
	err = os.WriteFile(outputFile, []byte(strings.Join(outputLines, "\n")+"\n"), 0644)
	if err != nil {
		return nil, fmt.Errorf("Failed to write output file: %v", err)
	}

	// Show success message
	msgFunc("success", fmt.Sprintf("✨ Variables successfully fetched and saved to '%s'", outputFile))

	return outputLines, nil
}
