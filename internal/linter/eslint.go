package linter

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// ESLintMessage represents a single ESLint message.
type ESLintMessage struct {
	RuleID   string `json:"ruleId"`
	Severity int    `json:"severity"` // 1 = warning, 2 = error
	Message  string `json:"message"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
}

// ESLintResult represents ESLint output for a single file.
type ESLintResult struct {
	FilePath string          `json:"filePath"`
	Messages []ESLintMessage `json:"messages"`
}

// RunESLint executes ESLint on the project and parses results into Findings.
func RunESLint(projectPath string) ([]scanner.Finding, error) {
	var findings []scanner.Finding

	// Check if eslint is available
	eslintBin := findESLint(projectPath)
	if eslintBin == "" {
		// ESLint not found — fallback to our native esbuild syntax checker
		return RunNativeSyntaxCheck(projectPath)
	}

	// Check if eslint config exists
	if !hasESLintConfig(projectPath) {
		findings = append(findings, scanner.Finding{
			Category:    scanner.CategoryQuality,
			Severity:    scanner.SeverityWarning,
			Title:       "Missing ESLint configuration",
			Description: "No ESLint config file found (.eslintrc, .eslintrc.js, eslint.config.js, etc.)",
			Suggestion:  "Run 'npx eslint --init' to create a configuration",
		})
		return findings, nil
	}

	// Run ESLint with JSON output format
	cmd := exec.Command(eslintBin, ".", "--format", "json", "--no-error-on-unmatched-pattern")
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		// ESLint returns non-zero exit code when there are linting errors
		// We still want to parse the output
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means lint errors found — that's expected
			if exitErr.ExitCode() == 1 {
				output = exitErr.Stderr
				// Try stdout from combined output
				if len(output) == 0 {
					output, _ = cmd.Output()
				}
			}
			// Try to use whatever output we have
			if len(output) == 0 {
				output = []byte("[]")
			}
		} else {
			return findings, nil // Can't run ESLint
		}
	}

	// Parse JSON output
	var results []ESLintResult
	if parseErr := json.Unmarshal(output, &results); parseErr != nil {
		// If we can't parse, just note that ESLint ran but output was unparseable
		findings = append(findings, scanner.Finding{
			Category:    scanner.CategoryQuality,
			Severity:    scanner.SeverityInfo,
			Title:       "ESLint output could not be parsed",
			Description: "ESLint ran but produced unparseable output",
			Suggestion:  "Run 'npx eslint . --format json' manually to debug",
		})
		return findings, nil
	}

	// Convert ESLint results to unified findings
	for _, result := range results {
		relPath, _ := filepath.Rel(projectPath, result.FilePath)
		if relPath == "" {
			relPath = result.FilePath
		}

		for _, msg := range result.Messages {
			severity := scanner.SeverityWarning
			if msg.Severity == 2 {
				severity = scanner.SeverityError
			}

			title := msg.Message
			if msg.RuleID != "" {
				title = fmt.Sprintf("[%s] %s", msg.RuleID, msg.Message)
			}

			findings = append(findings, scanner.Finding{
				Category:    scanner.CategoryQuality,
				Severity:    severity,
				Title:       title,
				Description: "ESLint: " + msg.Message,
				File:        relPath,
				Line:        msg.Line,
				Suggestion:  suggestFix(msg.RuleID),
			})
		}
	}

	// Cap at 50 findings to avoid flooding
	if len(findings) > 50 {
		total := len(findings)
		findings = findings[:50]
		findings = append(findings, scanner.Finding{
			Category:    scanner.CategoryQuality,
			Severity:    scanner.SeverityInfo,
			Title:       fmt.Sprintf("... and %d more ESLint issues", total-50),
			Description: "Too many lint issues to display. Run ESLint directly for full output",
			Suggestion:  "Run 'npx eslint .' for the complete list",
		})
	}

	return findings, nil
}

// findESLint looks for ESLint binary — local node_modules first, then global.
func findESLint(projectPath string) string {
	// Check local node_modules/.bin/eslint
	localBin := filepath.Join(projectPath, "node_modules", ".bin", "eslint")
	if _, err := os.Stat(localBin); err == nil {
		return localBin
	}
	// Windows variant
	localBinExe := localBin + ".cmd"
	if _, err := os.Stat(localBinExe); err == nil {
		return localBinExe
	}

	// Check global
	globalBin, err := exec.LookPath("eslint")
	if err == nil {
		return globalBin
	}

	// Try npx
	npxBin, err := exec.LookPath("npx")
	if err == nil {
		// Verify eslint is available via npx
		cmd := exec.Command(npxBin, "--no-install", "eslint", "--version")
		cmd.Dir = projectPath
		if err := cmd.Run(); err == nil {
			return npxBin + " eslint"
		}
	}

	return ""
}

// hasESLintConfig checks if an ESLint configuration file exists.
func hasESLintConfig(projectPath string) bool {
	configs := []string{
		".eslintrc",
		".eslintrc.js",
		".eslintrc.cjs",
		".eslintrc.json",
		".eslintrc.yml",
		".eslintrc.yaml",
		"eslint.config.js",
		"eslint.config.mjs",
		"eslint.config.cjs",
		"eslint.config.ts",
	}
	for _, cfg := range configs {
		if _, err := os.Stat(filepath.Join(projectPath, cfg)); err == nil {
			return true
		}
	}

	// Also check package.json for eslintConfig field
	pkgPath := filepath.Join(projectPath, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err == nil {
		if strings.Contains(string(data), `"eslintConfig"`) {
			return true
		}
	}

	return false
}

// suggestFix provides fix suggestions for common ESLint rules.
func suggestFix(ruleID string) string {
	suggestions := map[string]string{
		"no-unused-vars":          "Remove unused variables or prefix with underscore",
		"no-console":              "Remove console.log statements or use a proper logger",
		"no-debugger":             "Remove debugger statements before committing",
		"semi":                    "Add or remove semicolons consistently",
		"quotes":                  "Use consistent quote style (single or double)",
		"indent":                  "Fix indentation to match project settings",
		"react/prop-types":        "Add PropTypes or use TypeScript for type safety",
		"react-hooks/rules-of-hooks": "Ensure hooks are called at the top level",
		"react-hooks/exhaustive-deps": "Add missing dependencies to the dependency array",
		"@typescript-eslint/no-unused-vars": "Remove unused variables or prefix with underscore",
		"@typescript-eslint/no-explicit-any": "Replace 'any' with a specific type",
		"import/order":            "Organize imports according to convention",
		"no-undef":                "Define the variable or add a global declaration",
		"eqeqeq":                  "Use === instead of == for strict equality",
	}

	if suggestion, ok := suggestions[ruleID]; ok {
		return suggestion
	}
	return "Fix the lint issue or disable the rule if intentional"
}
