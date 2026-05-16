package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Swif7ify/Obelisk-CLI/internal/config"
	"github.com/Swif7ify/Obelisk-CLI/internal/detector"
	"github.com/Swif7ify/Obelisk-CLI/internal/engine"
	"github.com/Swif7ify/Obelisk-CLI/internal/report"
	"github.com/Swif7ify/Obelisk-CLI/ui"
)

var flagCheckPath string
var flagSkipAI bool
var flagOutputFile string

var checkCmd = &cobra.Command{
	Use:   "check [path]",
	Short: "Run a full health check on your project",
	Long:  "Scans the project for security issues, architectural problems, and code quality, then produces an AI-graded health report.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := flagCheckPath
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
		
		// Strip spaces, quotes, and control characters
		projectPath = strings.TrimFunc(projectPath, func(r rune) bool {
			return r < 32 || r == '"' || r == '\'' || r == ' '
		})

		// Validate project path exists
		if _, err := os.Stat(projectPath); err != nil {
			return fmt.Errorf("invalid project path: %w", err)
		}

		// Load config to get format preference
		cfg, _ := config.Load()
		format := cfg.GetReportFormat()

		if flagNoColor || cfg.NoColor {
			os.Setenv("NO_COLOR", "1")
		}

		// Detect project type
		detection := detector.Detect(projectPath)
		
		// For now, only support JavaScript/TypeScript projects
		supportedTypes := map[detector.ProjectType]bool{
			detector.TypeJavaScript: true,
			detector.TypeTypeScript: true,
			detector.TypeReact:      true,
			detector.TypeNextJS:     true,
		}

		if !supportedTypes[detection.Type] {
			return fmt.Errorf(
				"unsupported project type: %s\n"+
					"Currently, Obelisk only supports JavaScript/TypeScript projects.\n"+
					"Detected: %s",
				detection.Type,
				detection.Framework,
			)
		}

		fmt.Printf("✓ Detected project type: %s\n\n", detection.Framework)

		// Create the dashboard
		dashboard := ui.NewDashboard()
		p := tea.NewProgram(dashboard)

		// Channel to capture the scan result
		resultChan := make(chan *engine.Result, 1)
		errChan := make(chan error, 1)

		// Run engine in background
		go func() {
			cfg := engine.Config{
				ProjectPath: projectPath,
				APIKey:      flagAPIKey,
				Model:       flagModel,
				SkipAI:      flagSkipAI || cfg.SkipAI || cfg.GetAPIKey() == "",
				Verbose:     flagVerbose,
			}

			result, err := engine.Run(cfg, func(phase string) {
				// Update spinner phase
				p.Send(ui.PhaseUpdateMsg{Phase: phase})
				if flagVerbose {
					fmt.Fprintf(os.Stderr, "  → %s\n", phase)
				}
			})

			if err != nil {
				errChan <- err
				p.Send(ui.ScanCompleteMsg{Err: err})
				return
			}

			resultChan <- result
			p.Send(ui.ScanCompleteMsg{
				Result: result.ScanResult,
				Report: result.Report,
			})
		}()

		_, err := p.Run()
		if err != nil {
			return err
		}

		// After the TUI exits, save the report to file
		// Check if we have a result from the scan
		select {
		case result := <-resultChan:
			outputPath := flagOutputFile
			if outputPath == "" {
				outputPath = report.GetDefaultOutputPath(projectPath, format)
			}

			// Write the report to file
			if err := report.WriteToFile(result.ScanResult, result.Report, outputPath, format); err != nil {
				fmt.Fprintf(os.Stderr, "\nWarning: Could not save report to file: %v\n", err)
			} else {
				fmt.Printf("\n✓ Report saved to: %s\n", outputPath)
			}
		case err := <-errChan:
			fmt.Fprintf(os.Stderr, "\nWarning: Scan failed, report not saved: %v\n", err)
		default:
			// No result available (shouldn't happen, but handle gracefully)
			fmt.Fprintf(os.Stderr, "\nWarning: No scan result available to save\n")
		}

		return nil
	},
}

func init() {
	checkCmd.Flags().StringVarP(&flagCheckPath, "path", "p", "", "Path to project (default: current directory)")
	checkCmd.Flags().BoolVar(&flagSkipAI, "skip-ai", false, "Skip AI analysis (offline mode)")
	checkCmd.Flags().StringVarP(&flagOutputFile, "output", "o", "", "Output file path (default: obelisk-report-<timestamp>.md in project directory)")
	rootCmd.AddCommand(checkCmd)
}
