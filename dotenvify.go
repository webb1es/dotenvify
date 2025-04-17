package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		prefix = fmt.Sprintf("%s✓ %s%s", colorGreen, message, colorReset)
	case "error":
		prefix = fmt.Sprintf("%s✗ %s%s", colorRed, message, colorReset)
	case "warning":
		prefix = fmt.Sprintf("%s⚠ %s%s", colorYellow, message, colorReset)
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

func main() {
	// Check for source file argument
	if len(os.Args) < 2 {
		msg("error", "No source file provided. Usage: dotenvify source_file [output_file]")
		os.Exit(1)
	}

	sourceFile := os.Args[1]
	outputFile := sourceFile
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
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
	if hasErrors && outputFile == sourceFile {
		outputFile = sourceFile + ".out"
		msg("warning", fmt.Sprintf("Some errors occurred. Output saved to '%s'", outputFile))
		for _, err := range errors {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	// Write the output file
	err = os.WriteFile(outputFile, []byte(strings.Join(outputLines, "\n")+"\n"), 0644)
	exitOnError(err, "Failed to write output file")

	// Show success message
	if !hasErrors {
		msg("success", fmt.Sprintf("Variables successfully formatted and saved to '%s'", outputFile))
	}
}
