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

		fmt.Printf("🏛️ Obelisk Configuration\n\n")
		fmt.Printf("  %-18s %s\n", "Config file:", configPath)
		fmt.Printf("  %-18s %s\n", "API Key:", cfg.MaskedAPIKey())
		fmt.Printf("  %-18s %s\n", "Model:", cfg.GetModel())

		defaultPath := cfg.DefaultPath
		if defaultPath == "" {
			defaultPath = "(current directory)"
		}
		fmt.Printf("  %-18s %s\n", "Default Path:", defaultPath)

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
  api-key     Your Gemini API key
  model       AI model name (e.g., gemini-2.5-flash)
  path        Default project scan path
  no-color    Disable colors (true/false)`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		key := strings.ToLower(args[0])
		value := args[1]

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
			cfg.NoColor = strings.ToLower(value) == "true" || value == "1"
			fmt.Printf("✓ No color: %v\n", cfg.NoColor)
		default:
			return fmt.Errorf("unknown config key: %s\nAvailable: api-key, model, path, no-color", key)
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
