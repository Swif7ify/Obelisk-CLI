package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Swif7ify/Obelisk-CLI/internal/engine"
	"github.com/Swif7ify/Obelisk-CLI/ui"
)

var flagCheckPath string
var flagSkipAI bool

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

		// Create the dashboard
		dashboard := ui.NewDashboard()
		p := tea.NewProgram(dashboard, tea.WithAltScreen())

		// Run engine in background
		go func() {
			cfg := engine.Config{
				ProjectPath: projectPath,
				APIKey:      flagAPIKey,
				Model:       flagModel,
				SkipAI:      flagSkipAI,
				Verbose:     flagVerbose,
			}

			result, err := engine.Run(cfg, func(phase string) {
				// Update spinner phase
				if flagVerbose {
					fmt.Fprintf(os.Stderr, "  → %s\n", phase)
				}
			})

			if err != nil {
				p.Send(ui.ScanCompleteMsg{Err: err})
				return
			}

			p.Send(ui.ScanCompleteMsg{
				Result: result.ScanResult,
				Report: result.Report,
			})
		}()

		_, err := p.Run()
		return err
	},
}

func init() {
	checkCmd.Flags().StringVarP(&flagCheckPath, "path", "p", "", "Path to project (default: current directory)")
	checkCmd.Flags().BoolVar(&flagSkipAI, "skip-ai", false, "Skip AI analysis (offline mode)")
	rootCmd.AddCommand(checkCmd)
}
