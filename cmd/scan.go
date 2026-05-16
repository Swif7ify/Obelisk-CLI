package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Swif7ify/Obelisk-CLI/internal/config"
	"github.com/Swif7ify/Obelisk-CLI/internal/engine"
	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

var flagScanPath string
var flagScanFormat string
var flagScanStrict bool
var flagScanSkipAI bool
var flagScanOutput string

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Run a headless health check (CI/CD friendly)",
	Long: `Runs a non-interactive health scan and prints results to stdout.
Ideal for CI/CD pipelines, scripting, and automation.

Examples:
  obelisk scan                          # Scan current directory
  obelisk scan ./my-project             # Scan specific path
  obelisk scan --format json            # Output as JSON
  obelisk scan --strict                 # Exit code 1 on critical/errors
  obelisk scan --skip-ai                # Skip AI analysis`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := flagScanPath
		if len(args) > 0 {
			projectPath = args[0]
		}
		if projectPath == "" {
			var err error
			projectPath, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		projectPath = strings.Trim(strings.TrimSpace(projectPath), `"'`)

		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			return fmt.Errorf("project path does not exist: %s", projectPath)
		}

		// Load config for API key fallback
		cfg, _ := config.Load()

		apiKey := flagAPIKey
		if apiKey == "" {
			apiKey = cfg.GetAPIKey()
		}

		model := flagModel
		if model == "" || model == "gemini-2.5-flash" {
			model = cfg.GetModel()
		}

		skipAI := flagScanSkipAI || apiKey == ""

		engineCfg := engine.Config{
			ProjectPath: projectPath,
			APIKey:      apiKey,
			Model:       model,
			SkipAI:      skipAI,
			Verbose:     flagVerbose,
		}

		if !flagNoColor {
			fmt.Fprintf(os.Stderr, "🏛️ Obelisk — Scanning %s...\n", projectPath)
		}

		result, err := engine.Run(engineCfg, func(phase string) {
			if flagVerbose {
				fmt.Fprintf(os.Stderr, "  → %s\n", phase)
			}
		})
		if err != nil {
			return err
		}

		// Output
		switch strings.ToLower(flagScanFormat) {
		case "json":
			return printScanJSON(result)
		case "markdown", "md":
			return printScanMarkdown(result)
		default:
			return printScanText(result)
		}
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		// In strict mode, exit with code 1 if critical/errors found
		if !flagScanStrict {
			return nil
		}
		// Re-read isn't needed — handled in RunE via os.Exit below
		return nil
	},
}

func printScanText(result *engine.Result) error {
	s := result.ScanResult
	r := result.Report

	criticals := s.CountBySeverity(scanner.SeverityCritical)
	errors := s.CountBySeverity(scanner.SeverityError)
	warnings := s.CountBySeverity(scanner.SeverityWarning)
	infos := s.CountBySeverity(scanner.SeverityInfo)

	fmt.Printf("\n📁 Project: %s\n", s.ProjectPath)
	fmt.Printf("📋 Type: %s (%s)\n", result.Detection.Framework, result.Detection.Type)
	fmt.Printf("📄 Files: %d | Dirs: %d\n\n", s.FileCount, s.DirCount)

	if r != nil {
		fmt.Printf("━━━ Grade: %s (%d/100) ━━━\n", r.Grade, r.OverallScore)
		fmt.Printf("  Security:     %d/100\n", r.SecurityScore)
		fmt.Printf("  Architecture: %d/100\n", r.ArchitectureScore)
		fmt.Printf("  Quality:      %d/100\n\n", r.QualityScore)
	}

	fmt.Printf("━━━ Issues: %d critical, %d errors, %d warnings, %d info ━━━\n\n",
		criticals, errors, warnings, infos)

	for _, f := range s.Findings {
		icon := "ℹ️"
		switch f.Severity {
		case scanner.SeverityCritical:
			icon = "🔴"
		case scanner.SeverityError:
			icon = "🟠"
		case scanner.SeverityWarning:
			icon = "🟡"
		}
		line := fmt.Sprintf("%s [%s] %s", icon, f.Severity, f.Title)
		if f.File != "" {
			line += fmt.Sprintf(" (%s)", f.File)
		}
		fmt.Println(line)
	}

	if r != nil && r.Summary != "" && !strings.HasPrefix(r.Summary, "Project analysis completed without AI") {
		fmt.Printf("\n━━━ Summary ━━━\n%s\n", r.Summary)
	}

	if flagScanStrict && (criticals > 0 || errors > 0) {
		fmt.Fprintf(os.Stderr, "\n❌ Strict mode: %d critical, %d errors found\n", criticals, errors)
		os.Exit(1)
	}

	return nil
}

func printScanJSON(result *engine.Result) error {
	data, err := json.MarshalIndent(map[string]interface{}{
		"scan_result": result.ScanResult,
		"report":      result.Report,
		"detection":   result.Detection,
	}, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	if flagScanStrict {
		s := result.ScanResult
		if s.CountBySeverity(scanner.SeverityCritical) > 0 || s.CountBySeverity(scanner.SeverityError) > 0 {
			os.Exit(1)
		}
	}
	return nil
}

func printScanMarkdown(result *engine.Result) error {
	r := result.Report
	s := result.ScanResult

	fmt.Printf("# 🏛️ Obelisk Scan Report\n\n")
	fmt.Printf("**Project:** %s\n", s.ProjectPath)
	fmt.Printf("**Type:** %s (%s)\n", result.Detection.Framework, result.Detection.Type)
	fmt.Printf("**Files:** %d | **Dirs:** %d\n\n", s.FileCount, s.DirCount)

	if r != nil {
		fmt.Printf("## Grade: %s (%d/100)\n\n", r.Grade, r.OverallScore)
		fmt.Printf("| Category | Score |\n|---|---|\n")
		fmt.Printf("| Security | %d/100 |\n", r.SecurityScore)
		fmt.Printf("| Architecture | %d/100 |\n", r.ArchitectureScore)
		fmt.Printf("| Quality | %d/100 |\n\n", r.QualityScore)
	}

	if len(s.Findings) > 0 {
		fmt.Printf("## Findings (%d)\n\n", len(s.Findings))
		for _, f := range s.Findings {
			fmt.Printf("- **[%s]** %s", f.Severity, f.Title)
			if f.File != "" {
				fmt.Printf(" (`%s`)", f.File)
			}
			fmt.Println()
		}
	}

	fmt.Printf("\n---\n*Generated by Obelisk CLI*\n")

	if flagScanStrict {
		if s.CountBySeverity(scanner.SeverityCritical) > 0 || s.CountBySeverity(scanner.SeverityError) > 0 {
			os.Exit(1)
		}
	}
	return nil
}

func init() {
	scanCmd.Flags().StringVarP(&flagScanPath, "path", "p", "", "Path to project (default: current directory)")
	scanCmd.Flags().StringVarP(&flagScanFormat, "format", "f", "text", "Output format: text, json, markdown")
	scanCmd.Flags().BoolVar(&flagScanStrict, "strict", false, "Exit code 1 on critical/error findings")
	scanCmd.Flags().BoolVar(&flagScanSkipAI, "skip-ai", false, "Skip AI analysis")
	scanCmd.Flags().StringVarP(&flagScanOutput, "output", "o", "", "Write output to file")
	rootCmd.AddCommand(scanCmd)
}
