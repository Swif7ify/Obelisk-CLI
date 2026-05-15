package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	Use:   "obelisk",
	Short: "🏛️ Obelisk — The AI-Powered Automated Tech Lead",
	Long: `Obelisk is a high-performance CLI tool that acts as a final gatekeeper 
for your project. It evaluates project integrity, security, and architectural 
health using static analysis and LLM intelligence.

Run 'obelisk check' to scan the current directory.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "Gemini API key (overrides GOOGLE_API_KEY env var)")
	rootCmd.PersistentFlags().StringVar(&flagModel, "model", "gemini-2.5-flash", "AI model to use")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "Disable colored output")
}
