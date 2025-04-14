package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ANSI color codes for output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// showMessage displays a styled message based on message type
func showMessage(msgType string, message string) {
	switch msgType {
	case "success":
		fmt.Printf("%s✓ %s%s\n", colorGreen, message, colorReset)
	case "error":
		fmt.Printf("%s✗ %s%s\n", colorRed, message, colorReset)
	case "warning":
		fmt.Printf("%s⚠ %s%s\n", colorYellow, message, colorReset)
	default:
		fmt.Println(message)
	}
}

func main() {
	// Check if source file is provided
	if len(os.Args) < 2 {
		showMessage("error", "No source file provided. Usage: dotenvify source_file [output_file]")
		os.Exit(1)
	}

	sourceFile := os.Args[1]
	outputFile := ""

	// Determine output file
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	} else {
		outputFile = sourceFile
	}

	// Check if source file exists
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		showMessage("error", fmt.Sprintf("Source file '%s' does not exist.", sourceFile))
		os.Exit(1)
	}

	// Create temporary files
	tempFile, err := ioutil.TempFile("", "dotenvify-")
	if err != nil {
		showMessage("error", "Failed to create temporary file: "+err.Error())
		os.Exit(1)
	}
	defer os.Remove(tempFile.Name())

	// Read source file
	file, err := os.Open(sourceFile)
	if err != nil {
		showMessage("error", "Failed to open source file: "+err.Error())
		os.Exit(1)
	}
	defer file.Close()

	// Process the file
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		showMessage("error", "Error reading source file: "+err.Error())
		os.Exit(1)
	}

	// Process lines in pairs and write to temp file
	hasErrors := false
	var errors []string

	for i := 0; i < len(lines); i += 2 {
		// Check if we have both key and value
		if i+1 < len(lines) {
			key := strings.TrimSpace(lines[i])
			value := strings.TrimSpace(lines[i+1])

			fmt.Fprintf(tempFile, "export %s=\"%s\"\n", key, value)
		} else {
			// Handle odd number of lines (missing value for last key)
			hasErrors = true
			errors = append(errors, fmt.Sprintf("Line %d: Key '%s' has no value", i+1, lines[i]))
		}
	}

	// Close the temp file before moving it
	tempFile.Close()

	// Handle error output and final file placement
	if hasErrors && outputFile == sourceFile {
		// If errors occurred and would overwrite source, save to .out file instead
		outputFile = sourceFile + ".out"
		showMessage("warning", fmt.Sprintf("Some errors occurred. Output saved to '%s'", outputFile))

		for _, err := range errors {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	// Move temp file to final destination
	err = os.Rename(tempFile.Name(), outputFile)
	if err != nil {
		// If rename fails (possible on different filesystems), copy and delete
		data, err := ioutil.ReadFile(tempFile.Name())
		if err != nil {
			showMessage("error", "Failed to read temporary file: "+err.Error())
			os.Exit(1)
		}

		err = ioutil.WriteFile(outputFile, data, 0644)
		if err != nil {
			showMessage("error", "Failed to write output file: "+err.Error())
			os.Exit(1)
		}
	}

	// Show success message
	if !hasErrors {
		showMessage("success", fmt.Sprintf("Variables successfully formatted and saved to '%s'", outputFile))
	}
}
