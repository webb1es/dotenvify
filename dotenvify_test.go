package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadExistingVariables(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]string
	}{
		{
			name: "basic key=value format",
			content: `API_KEY=test123
DATABASE_URL=postgres://localhost`,
			expected: map[string]string{
				"API_KEY":      "test123",
				"DATABASE_URL": "postgres://localhost",
			},
		},
		{
			name: "quoted values",
			content: `API_KEY="test with spaces"
DATABASE_URL='single quotes'`,
			expected: map[string]string{
				"API_KEY":      "test with spaces",
				"DATABASE_URL": "single quotes",
			},
		},
		{
			name: "export prefix",
			content: `export API_KEY=test123
export DATABASE_URL=postgres://localhost`,
			expected: map[string]string{
				"API_KEY":      "test123",
				"DATABASE_URL": "postgres://localhost",
			},
		},
		{
			name: "comments and empty lines",
			content: `# This is a comment
API_KEY=test123

# Another comment
DATABASE_URL=postgres://localhost`,
			expected: map[string]string{
				"API_KEY":      "test123",
				"DATABASE_URL": "postgres://localhost",
			},
		},
		{
			name:     "non-existent file",
			content:  "",
			expected: map[string]string{},
		},
		{
			name:     "empty file",
			content:  "",
			expected: map[string]string{},
		},
		{
			name:     "only comments",
			content:  "# Comment 1\n# Comment 2\n",
			expected: map[string]string{},
		},
		{
			name: "mixed export and non-export",
			content: `API_KEY=test1
export DATABASE_URL=test2`,
			expected: map[string]string{
				"API_KEY":      "test1",
				"DATABASE_URL": "test2",
			},
		},
		{
			name:    "value with equals sign",
			content: `CONNECTION_STRING=user=admin;password=secret`,
			expected: map[string]string{
				"CONNECTION_STRING": "user=admin;password=secret",
			},
		},
		{
			name:    "empty value",
			content: `EMPTY_VAR=`,
			expected: map[string]string{
				"EMPTY_VAR": "",
			},
		},
		{
			name:    "whitespace around key and value",
			content: `  API_KEY  =  test123  `,
			expected: map[string]string{
				"API_KEY": "test123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tmpFile string
			if tt.content != "" {
				// Create temp file
				tmp, err := os.CreateTemp("", "test-*.env")
				if err != nil {
					t.Fatal(err)
				}
				tmpFile = tmp.Name()
				defer os.Remove(tmpFile)

				if _, err := tmp.WriteString(tt.content); err != nil {
					t.Fatal(err)
				}
				tmp.Close()
			} else {
				tmpFile = "/nonexistent/file.env"
			}

			result, err := readExistingVariables(tmpFile)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d variables, got %d", len(tt.expected), len(result))
			}

			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("key %s: expected %q, got %q", key, expectedValue, result[key])
				}
			}
		})
	}
}

func TestWriteVariablesToFileWithPreserve(t *testing.T) {
	tests := []struct {
		name         string
		existing     string
		variables    map[string]string
		preserveVars string
		expected     map[string]string
	}{
		{
			name: "preserve single variable",
			existing: `DATABASE_URL="keep-this"
API_KEY=old-key`,
			variables: map[string]string{
				"DATABASE_URL": "new-url",
				"API_KEY":      "new-key",
			},
			preserveVars: "DATABASE_URL",
			expected: map[string]string{
				"DATABASE_URL": "keep-this",
				"API_KEY":      "new-key",
			},
		},
		{
			name: "preserve multiple variables",
			existing: `DATABASE_URL="keep-this"
API_KEY=keep-this-too
SECRET=old-secret`,
			variables: map[string]string{
				"DATABASE_URL": "new-url",
				"API_KEY":      "new-key",
				"SECRET":       "new-secret",
			},
			preserveVars: "DATABASE_URL,API_KEY",
			expected: map[string]string{
				"DATABASE_URL": "keep-this",
				"API_KEY":      "keep-this-too",
				"SECRET":       "new-secret",
			},
		},
		{
			name:     "preserve non-existent variable",
			existing: `API_KEY=old-key`,
			variables: map[string]string{
				"DATABASE_URL": "new-url",
				"API_KEY":      "new-key",
			},
			preserveVars: "DATABASE_URL,API_KEY",
			expected: map[string]string{
				"DATABASE_URL": "new-url",
				"API_KEY":      "old-key",
			},
		},
		{
			name:     "no existing file",
			existing: "",
			variables: map[string]string{
				"DATABASE_URL": "new-url",
				"API_KEY":      "new-key",
			},
			preserveVars: "DATABASE_URL",
			expected: map[string]string{
				"DATABASE_URL": "new-url",
				"API_KEY":      "new-key",
			},
		},
		{
			name:     "preserve with whitespace in list",
			existing: `API_KEY=old-value`,
			variables: map[string]string{
				"API_KEY": "new-value",
			},
			preserveVars: " API_KEY , DATABASE_URL ",
			expected: map[string]string{
				"API_KEY": "old-value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "output.env")

			// Create existing file if content provided
			if tt.existing != "" {
				if err := os.WriteFile(outputFile, []byte(tt.existing), 0600); err != nil {
					t.Fatal(err)
				}
			}

			// Write variables with preserve
			err := writeVariablesToFile(tt.variables, outputFile, false, true, false, false, true, tt.preserveVars)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Read result
			result, err := readExistingVariables(outputFile)
			if err != nil {
				t.Fatalf("failed to read result: %v", err)
			}

			// Verify
			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("key %s: expected %q, got %q", key, expectedValue, result[key])
				}
			}
		})
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"postgres://localhost:5432/db", true},
		{"mysql://localhost", true},
		{"mongodb://localhost", true},
		{"redis://localhost", true},
		{"ftp://example.com", true},
		{"sftp://example.com", true},
		{"ssh://example.com", true},
		{"git://github.com", true},
		{"mailto:test@example.com", true},
		{"just-a-string", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isURL(tt.input)
			if result != tt.expected {
				t.Errorf("isURL(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsHTTPURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"postgres://localhost", false},
		{"ftp://example.com", false},
		{"not-a-url", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isHTTPURL(tt.input)
			if result != tt.expected {
				t.Errorf("isHTTPURL(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsQuoted(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"quoted"`, true},
		{`'quoted'`, true},
		{`not quoted`, false},
		{`"mismatched'`, false},
		{`'mismatched"`, false},
		{``, false},
		{`""`, true},
		{`''`, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isQuoted(tt.input)
			if result != tt.expected {
				t.Errorf("isQuoted(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBackupFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.env")

	// Create test file
	content := []byte("TEST=value")
	if err := os.WriteFile(testFile, content, 0600); err != nil {
		t.Fatal(err)
	}

	// First backup
	if err := backupFile(testFile); err != nil {
		t.Fatalf("first backup failed: %v", err)
	}

	backup1 := testFile + ".backup.1"
	if _, err := os.Stat(backup1); os.IsNotExist(err) {
		t.Error("backup.1 was not created")
	}

	// Second backup (should create .backup.2)
	if err := backupFile(testFile); err != nil {
		t.Fatalf("second backup failed: %v", err)
	}

	backup2 := testFile + ".backup.2"
	if _, err := os.Stat(backup2); os.IsNotExist(err) {
		t.Error("backup.2 was not created")
	}

	// Verify content
	backupContent, err := os.ReadFile(backup1)
	if err != nil {
		t.Fatal(err)
	}
	if string(backupContent) != string(content) {
		t.Error("backup content doesn't match original")
	}
}

func TestBackupFileNonExistent(t *testing.T) {
	// Backing up a non-existent file should not error
	err := backupFile("/nonexistent/file.env")
	if err != nil {
		t.Errorf("backup of non-existent file should not error: %v", err)
	}
}

func TestWriteVariablesToFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.env")

	variables := map[string]string{
		"API_KEY":      "test123",
		"DATABASE_URL": "postgres://localhost",
		"SECRET":       "value with spaces",
	}

	// Test basic write
	err := writeVariablesToFile(variables, outputFile, false, false, false, false, true, "")
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Read back
	result, err := readExistingVariables(outputFile)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	for key, expected := range variables {
		if result[key] != expected {
			t.Errorf("key %s: expected %q, got %q", key, expected, result[key])
		}
	}
}

func TestWriteVariablesToFileWithExport(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.env")

	variables := map[string]string{
		"API_KEY": "test123",
	}

	// Test with export prefix
	err := writeVariablesToFile(variables, outputFile, false, false, true, false, true, "")
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Read file content directly to check export prefix
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	expected := "export API_KEY=test123\n"
	if string(content) != expected {
		t.Errorf("expected %q, got %q", expected, string(content))
	}
}

func TestWriteVariablesToFileWithNoLower(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.env")

	variables := map[string]string{
		"API_KEY":   "test123",
		"lowercase": "should-be-skipped",
		"UPPERCASE": "included",
		"MixedCase": "included",
	}

	// Test with noLower flag
	err := writeVariablesToFile(variables, outputFile, true, false, false, false, true, "")
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Read result
	result, err := readExistingVariables(outputFile)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	// Should have 3 variables (lowercase should be skipped)
	if len(result) != 3 {
		t.Errorf("expected 3 variables, got %d", len(result))
	}

	if _, exists := result["lowercase"]; exists {
		t.Error("lowercase variable should have been skipped")
	}
}

func TestWriteVariablesToFileWithURLOnly(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.env")

	variables := map[string]string{
		"DATABASE_URL": "https://example.com/db",
		"API_URL":      "http://api.example.com",
		"REGULAR_VAR":  "not-a-url",
		"POSTGRES":     "postgres://localhost:5432/db",
	}

	// Test with urlOnly flag
	err := writeVariablesToFile(variables, outputFile, false, false, false, true, true, "")
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Read result
	result, err := readExistingVariables(outputFile)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	// Should only have HTTP/HTTPS URLs
	if len(result) != 2 {
		t.Errorf("expected 2 variables, got %d", len(result))
	}

	if _, exists := result["REGULAR_VAR"]; exists {
		t.Error("non-URL variable should have been skipped")
	}

	if _, exists := result["POSTGRES"]; exists {
		t.Error("non-HTTP URL should have been skipped")
	}
}

func TestWriteVariablesToFileWithSorting(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.env")

	variables := map[string]string{
		"ZULU":  "last",
		"ALPHA": "first",
		"BRAVO": "second",
	}

	// Test with sorting (default)
	err := writeVariablesToFile(variables, outputFile, false, false, false, false, true, "")
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Read file content directly to check order
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	expected := "ALPHA=first\nBRAVO=second\nZULU=last\n"
	if string(content) != expected {
		t.Errorf("expected %q, got %q", expected, string(content))
	}

	// Test without sorting
	outputFile2 := filepath.Join(tmpDir, "output2.env")
	err = writeVariablesToFile(variables, outputFile2, false, true, false, false, true, "")
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Just verify all variables are present
	result, err := readExistingVariables(outputFile2)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 variables, got %d", len(result))
	}
}

func TestWriteVariablesToFileQuoting(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.env")

	variables := map[string]string{
		"URL_VAR":     "https://example.com",
		"SPACE_VAR":   "value with spaces",
		"REGULAR_VAR": "simplevalue",
		"QUOTED_VAR":  "\"already-quoted\"",
	}

	err := writeVariablesToFile(variables, outputFile, false, false, false, false, true, "")
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Read file content directly
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	contentStr := string(content)

	// URL should be quoted
	if !strings.Contains(contentStr, "URL_VAR=\"https://example.com\"") {
		t.Error("URL should be quoted")
	}

	// Space value should be quoted
	if !strings.Contains(contentStr, "SPACE_VAR=\"value with spaces\"") {
		t.Error("value with spaces should be quoted")
	}

	// Regular value should not be quoted
	if !strings.Contains(contentStr, "REGULAR_VAR=simplevalue") {
		t.Error("regular value should not be quoted")
	}
}

func TestWriteVariablesToFileWithOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.env")

	// Create initial file
	initialContent := "OLD_VAR=old_value\n"
	if err := os.WriteFile(outputFile, []byte(initialContent), 0600); err != nil {
		t.Fatal(err)
	}

	variables := map[string]string{
		"NEW_VAR": "new_value",
	}

	// Write with overwrite=true (no backup should be created)
	err := writeVariablesToFile(variables, outputFile, false, false, false, false, true, "")
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Check that backup was NOT created (because overwrite=true)
	backupFile := outputFile + ".backup.1"
	if _, err := os.Stat(backupFile); !os.IsNotExist(err) {
		t.Error("backup file should not be created when overwrite=true")
	}

	// Write with overwrite=false (backup should be created)
	err = writeVariablesToFile(variables, outputFile, false, false, false, false, false, "")
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Check that backup WAS created
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		t.Error("backup file should be created when overwrite=false")
	}
}

func TestProcessVariables(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.env")

	variables := map[string]string{
		"API_KEY": "test123",
		"SECRET":  "value",
	}

	// Should not panic
	ProcessVariables(variables, outputFile, false, false, false, false, true, "")

	// Verify file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("output file was not created")
	}
}

func TestMultipleFormatsInSingleFile(t *testing.T) {
	content := `# Comment
export KEY1=value1
KEY2="value2"
KEY3='value3'
KEY4=value4`

	tmp, err := os.CreateTemp("", "test-*.env")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile := tmp.Name()
	defer os.Remove(tmpFile)

	if _, err := tmp.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	result, err := readExistingVariables(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
		"KEY3": "value3",
		"KEY4": "value4",
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d variables, got %d", len(expected), len(result))
	}

	for key, expectedValue := range expected {
		if result[key] != expectedValue {
			t.Errorf("key %s: expected %q, got %q", key, expectedValue, result[key])
		}
	}
}

func TestReadUserInput(t *testing.T) {
	// This function reads from stdin, so we can only test that it exists
	// Full testing would require mocking stdin
	_ = readUserInput
}
