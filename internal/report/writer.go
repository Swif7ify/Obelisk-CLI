package report

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Swif7ify/Obelisk-CLI/internal/ai"
	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// WriteToFile writes the formatted report to a file.
func WriteToFile(result *scanner.ScanResult, report *ai.HealthReport, outputPath string) error {
	// Generate the plain text report
	content := FormatPlainText(result, report)

	// Ensure the directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}

	return nil
}

// GetDefaultOutputPath returns the default output path for the report.
func GetDefaultOutputPath(projectPath string) string {
	return filepath.Join(projectPath, "obelisk-report.txt")
}

// Made with Bob
