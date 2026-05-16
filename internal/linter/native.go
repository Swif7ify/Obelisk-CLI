package linter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
	"github.com/evanw/esbuild/pkg/api"
)

// RunNativeSyntaxCheck runs a fast, native syntax check using esbuild.
func RunNativeSyntaxCheck(projectPath string) ([]scanner.Finding, error) {
	var findings []scanner.Finding

	// Inform the user we are using the native fallback
	findings = append(findings, scanner.Finding{
		Category:    scanner.CategoryQuality,
		Severity:    scanner.SeverityWarning,
		Title:       "ESLint not found — using Obelisk Native Scanner",
		Description: "Obelisk fell back to its native esbuild-powered syntax checker.",
		Suggestion:  "For deeper code quality checks, install ESLint: 'npm install --save-dev eslint'",
	})

	matcher, err := scanner.NewGitignoreMatcher(projectPath)
	if err != nil {
		return findings, err
	}

	err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if matcher.ShouldIgnore(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		var loader api.Loader
		switch ext {
		case ".js", ".mjs", ".cjs":
			loader = api.LoaderJS
		case ".ts":
			loader = api.LoaderTS
		case ".jsx":
			loader = api.LoaderJSX
		case ".tsx":
			loader = api.LoaderTSX
		default:
			return nil // Skip non JS/TS files
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		result := api.Transform(string(data), api.TransformOptions{
			Loader: loader,
		})

		for _, msg := range result.Errors {
			relPath, _ := filepath.Rel(projectPath, path)
			if relPath == "" {
				relPath = path
			}

			findings = append(findings, scanner.Finding{
				Category:    scanner.CategoryQuality,
				Severity:    scanner.SeverityError,
				Title:       msg.Text,
				Description: fmt.Sprintf("Syntax error: %s", msg.Text),
				File:        relPath,
				Line:        msg.Location.Line,
				Suggestion:  "Fix the syntax error in the file.",
			})
		}

		return nil
	})

	return findings, err
}
