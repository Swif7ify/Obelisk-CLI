package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Swif7ify/Obelisk-CLI/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Obelisk configuration",
	Long: `View and modify Obelisk configuration settings.
Configuration is stored at ~/.obelisk/config.json.

Examples:
  obelisk config list                    # Show all config
  obelisk config set api-key <key>       # Store API key
  obelisk config get api-key             # Show API key status
  obelisk config set model gemini-2.5-pro
  obelisk config set auto-save false     # Disable auto-saving reports
  obelisk config reset                   # Reset to defaults`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all configuration values",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		configPath, _ := config.GetConfigPath()

		fmt.Printf("Obelisk Configuration\n\n")
		fmt.Printf("  %-18s %s\n", "Config file:", configPath)
		fmt.Printf("  %-18s %s\n", "API Key:", cfg.MaskedAPIKey())
		fmt.Printf("  %-18s %s\n", "Model:", cfg.GetModel())

		defaultPath := cfg.DefaultPath
		if defaultPath == "" {
			defaultPath = "(current directory)"
		}
		fmt.Printf("  %-18s %s\n", "Default Path:", defaultPath)
		fmt.Printf("  %-18s %s\n", "Report Format:", cfg.GetReportFormat())

		autoSave := "On"
		if !cfg.AutoSave {
			autoSave = "Off"
		}
		fmt.Printf("  %-18s %s\n", "Auto Save Report:", autoSave)

		noColor := "Off"
		if cfg.NoColor {
			noColor = "On"
		}
		fmt.Printf("  %-18s %s\n", "No Color:", noColor)

		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Available keys:
  api-key       Your Gemini API key
  model         AI model (e.g. gemini-2.5-flash)
  path          Default project path
  no-color      Disable colored output (true/false)
  auto-save     Auto-save report if issues found (true/false)
  report-format Report format (md/txt)`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		key := strings.ToLower(args[0])
		value := args[1]
		val := strings.ToLower(value)

		switch key {
		case "api-key", "apikey":
			cfg.SetAPIKey(value)
			fmt.Println("✓ API key saved")
		case "model":
			cfg.Model = value
			fmt.Printf("✓ Model set to: %s\n", value)
		case "path", "default-path":
			cfg.DefaultPath = value
			fmt.Printf("✓ Default path set to: %s\n", value)
		case "no-color", "nocolor":
			cfg.NoColor = val == "true" || val == "1"
			fmt.Printf("✓ No color: %v\n", cfg.NoColor)
		case "report-format", "format":
			if val == "txt" {
				cfg.ReportFormat = "txt"
			} else {
				cfg.ReportFormat = "md"
			}
			fmt.Printf("✓ Report format set to: %s\n", cfg.ReportFormat)
		case "auto-save", "autosave":
			if val == "true" || val == "1" || val == "yes" {
				cfg.AutoSave = true
			} else {
				cfg.AutoSave = false
			}
			fmt.Printf("✓ Auto-save report set to: %v\n", cfg.AutoSave)
		default:
			return fmt.Errorf("unknown config key: %s\nAvailable: api-key, model, path, no-color, auto-save, report-format", key)
		}

		return cfg.Save()
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		key := strings.ToLower(args[0])
		switch key {
		case "api-key", "apikey":
			fmt.Println(cfg.MaskedAPIKey())
		case "model":
			fmt.Println(cfg.GetModel())
		case "path", "default-path":
			p := cfg.DefaultPath
			if p == "" {
				p = "(not set)"
			}
			fmt.Println(p)
		case "no-color", "nocolor":
			fmt.Println(cfg.NoColor)
		case "report-format", "format":
			fmt.Println(cfg.GetReportFormat())
		case "auto-save", "autosave":
			fmt.Println(cfg.AutoSave)
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}
		return nil
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all configuration to defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		cfg.Reset()
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Println("✓ Configuration reset to defaults")
		return nil
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configResetCmd)
	rootCmd.AddCommand(configCmd)
}
