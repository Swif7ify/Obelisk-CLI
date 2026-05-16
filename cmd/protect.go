package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Swif7ify/Obelisk-CLI/internal/engine"
)

var flagInstall bool
var flagStrict bool

var protectCmd = &cobra.Command{
	Use:   "protect",
	Short: "Git pre-push hook integration",
	Long:  "Install Obelisk as a Git pre-push hook or run a protection check.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			switch args[0] {
			case "install":
				return installHook()
			case "uninstall":
				return uninstallHook()
			default:
				return fmt.Errorf("unknown command %q for protect", args[0])
			}
		}
		if flagInstall {
			return installHook()
		}
		return runProtect()
	},
}

func installHook() error {
	// Find git directory
	gitDir, err := findGitDir()
	if err != nil {
		return fmt.Errorf("not a git repository: %w", err)
	}

	hookPath := filepath.Join(gitDir, "hooks", "pre-push")
	hookContent := `#!/bin/sh
# Obelisk CLI pre-push hook
# Automatically installed by 'obelisk protect --install'

echo "Obelisk pre-push check..."

obelisk protect --strict
exit_code=$?

if [ $exit_code -ne 0 ]; then
    echo ""
    echo "❌ Push blocked by Obelisk. Fix critical issues before pushing."
    echo "   Run 'obelisk check' for details."
    exit 1
fi

echo "✅ Obelisk check passed!"
exit 0
`

	// Create hooks directory if needed
	hooksDir := filepath.Dir(hookPath)
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
		return fmt.Errorf("failed to write hook: %w", err)
	}

	fmt.Println("✅ Obelisk pre-push hook installed!")
	fmt.Printf("   Hook path: %s\n", hookPath)
	fmt.Println("   Pushes will now be checked for critical security issues.")
	return nil
}

func uninstallHook() error {
	gitDir, err := findGitDir()
	if err != nil {
		return fmt.Errorf("not a git repository: %w", err)
	}

	hookPath := filepath.Join(gitDir, "hooks", "pre-push")
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		fmt.Println("ℹ️ No pre-push hook found to uninstall.")
		return nil
	}

	if err := os.Remove(hookPath); err != nil {
		return fmt.Errorf("failed to remove hook: %w", err)
	}

	fmt.Println("✅ Obelisk pre-push hook uninstalled successfully!")
	return nil
}

func runProtect() error {
	cfg := engine.Config{
		APIKey:  flagAPIKey,
		Model:   flagModel,
		SkipAI:  true, // Don't call AI for pre-push (too slow)
		Verbose: flagVerbose,
	}

	result, err := engine.Run(cfg, func(phase string) {
		if flagVerbose {
			fmt.Printf("  → %s\n", phase)
		}
	})
	if err != nil {
		return err
	}

	criticals := result.ScanResult.CountBySeverity(3) // SeverityCritical
	errors := result.ScanResult.CountBySeverity(2)    // SeverityError

	if flagStrict && (criticals > 0 || errors > 0) {
		fmt.Printf("❌ Found %d critical and %d error issues\n", criticals, errors)
		os.Exit(1)
	}

	if criticals > 0 {
		fmt.Printf("❌ Found %d critical issues — push blocked\n", criticals)
		os.Exit(1)
	}

	fmt.Println("✅ No critical issues found")
	return nil
}

func findGitDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	dir := string(output)
	// Trim whitespace/newlines
	for len(dir) > 0 && (dir[len(dir)-1] == '\n' || dir[len(dir)-1] == '\r' || dir[len(dir)-1] == ' ') {
		dir = dir[:len(dir)-1]
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return dir, nil
	}
	return absDir, nil
}

func init() {
	protectCmd.Flags().BoolVar(&flagInstall, "install", false, "Install as Git pre-push hook")
	protectCmd.Flags().BoolVar(&flagStrict, "strict", false, "Fail on errors and critical issues")
	rootCmd.AddCommand(protectCmd)
}
