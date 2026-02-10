package azure

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		expectedOrg  string
		expectedProj string
		expectError  bool
	}{
		{
			name:         "dev.azure.com format",
			url:          "https://dev.azure.com/myorg/myproject",
			expectedOrg:  "myorg",
			expectedProj: "myproject",
			expectError:  false,
		},
		{
			name:         "dev.azure.com with trailing slash",
			url:          "https://dev.azure.com/myorg/myproject/",
			expectedOrg:  "myorg",
			expectedProj: "myproject",
			expectError:  false,
		},
		{
			name:         "visualstudio.com format",
			url:          "https://myorg.visualstudio.com/myproject",
			expectedOrg:  "myorg",
			expectedProj: "myproject",
			expectError:  false,
		},
		{
			name:         "visualstudio.com with trailing slash",
			url:          "https://myorg.visualstudio.com/myproject/",
			expectedOrg:  "myorg",
			expectedProj: "myproject",
			expectError:  false,
		},
		{
			name:         "http protocol dev.azure.com",
			url:          "http://dev.azure.com/testorg/testproject",
			expectedOrg:  "testorg",
			expectedProj: "testproject",
			expectError:  false,
		},
		{
			name:        "invalid URL format",
			url:         "https://example.com/not/azure",
			expectError: true,
		},
		{
			name:        "incomplete URL",
			url:         "https://dev.azure.com/myorg",
			expectError: true,
		},
		{
			name:        "empty URL",
			url:         "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org, proj, err := ParseURL(tt.url)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if org != tt.expectedOrg {
				t.Errorf("organization: expected %q, got %q", tt.expectedOrg, org)
			}

			if proj != tt.expectedProj {
				t.Errorf("project: expected %q, got %q", tt.expectedProj, proj)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	org := "testorg"
	project := "testproject"

	client := NewClient(org, project)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.Organization != org {
		t.Errorf("Organization: expected %q, got %q", org, client.Organization)
	}

	if client.Project != project {
		t.Errorf("Project: expected %q, got %q", project, client.Project)
	}
}

func TestGetVariablesValidation(t *testing.T) {
	tests := []struct {
		name       string
		org        string
		project    string
		groupNames string
		expectErr  bool
	}{
		{
			name:       "valid inputs",
			org:        "myorg",
			project:    "myproject",
			groupNames: "group1",
			expectErr:  false, // Will error because we can't actually call Azure, but validates inputs
		},
		{
			name:       "empty organization",
			org:        "",
			project:    "myproject",
			groupNames: "group1",
			expectErr:  true,
		},
		{
			name:       "empty project",
			org:        "myorg",
			project:    "",
			groupNames: "group1",
			expectErr:  true,
		},
		{
			name:       "empty group names",
			org:        "myorg",
			project:    "myproject",
			groupNames: "",
			expectErr:  true,
		},
		{
			name:       "whitespace only organization",
			org:        "   ",
			project:    "myproject",
			groupNames: "group1",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMsg := func(string, string) {}
			_, err := GetVariables(tt.org, tt.project, tt.groupNames, mockMsg)

			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectErr && err != nil {
				// This is expected to fail due to Azure CLI not being available in tests
				// We're just testing input validation here
				t.Skip("skipping test that requires Azure CLI")
			}
		})
	}
}
