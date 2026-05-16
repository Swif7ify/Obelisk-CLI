package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Swif7ify/Obelisk-CLI/internal/config"
	"github.com/Swif7ify/Obelisk-CLI/ui"
)

// Version info — injected via ldflags at build time.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// Global flags
var (
	flagVerbose bool
	flagAPIKey  string
	flagModel   string
	flagNoColor bool
)

var rootCmd = &cobra.Command{
	Use:     "obelisk",
	Short:   "Obelisk — The AI-Powered Automated Tech Lead",
	Version: Version,
	Long: `Obelisk is a high-performance CLI tool that acts as a final gatekeeper
for your project. It evaluates project integrity, security, and architectural
health using static analysis and LLM intelligence.

Run with no subcommand to launch the interactive TUI.
Run 'obelisk scan' for headless CI/CD mode.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Launch interactive TUI when no subcommand is given
		cfg, _ := config.Load()

		// Override config with flags if provided
		if flagAPIKey != "" {
			cfg.SetAPIKey(flagAPIKey)
		}
		if flagModel != "" && flagModel != "gemini-2.5-flash" {
			cfg.Model = flagModel
		}
		if flagNoColor || cfg.NoColor {
			cfg.NoColor = true
			os.Setenv("NO_COLOR", "1")
		}

		model := ui.NewInteractive(cfg, Version)
		p := tea.NewProgram(model, tea.WithAltScreen())

		_, err := p.Run()
		return err
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Set custom version template
	rootCmd.SetVersionTemplate(`Obelisk CLI
  Version:    {{.Version}}
  Commit:     ` + Commit + `
  Built:      ` + BuildDate + `
`)

	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "Gemini API key (overrides GOOGLE_API_KEY env var)")
	rootCmd.PersistentFlags().StringVar(&flagModel, "model", "gemini-2.5-flash", "AI model to use")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "Disable colored output")
}
