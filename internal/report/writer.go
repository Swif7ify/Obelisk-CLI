package report

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Swif7ify/Obelisk-CLI/internal/ai"
	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// WriteToFile writes the formatted report to a file.
func WriteToFile(result *scanner.ScanResult, report *ai.HealthReport, outputPath string, format string) error {
	var content string
	if format == "txt" {
		content = FormatPlainText(result, report)
	} else {
		content = FormatMarkdown(result, report)
	}

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
func GetDefaultOutputPath(projectPath string, format string) string {
	timestamp := time.Now().Format("20060102-150405")
	ext := format
	if ext == "" {
		ext = "md"
	}
	filename := fmt.Sprintf("obelisk-report-%s.%s", timestamp, ext)
	return filepath.Join(projectPath, filename)
}

// Made with Bob
