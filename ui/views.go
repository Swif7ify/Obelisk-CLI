package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Swif7ify/Obelisk-CLI/internal/config"
)

// RenderAPIKeyView renders the API key management view.
func RenderAPIKeyView(cfg *config.Config, subCursor int) string {
	var sb strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary).
		MarginBottom(1)

	sb.WriteString(headerStyle.Render("API Key Management") + "\n\n")

	// Show current status
	statusLabel := lipgloss.NewStyle().Bold(true).Foreground(ColorText).Render("  Status: ")
	if cfg.GetAPIKey() != "" {
		sb.WriteString(statusLabel + SuccessStyle.Render("✓ Configured") + "\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(ColorMuted).Render("  Key:    ") +
			lipgloss.NewStyle().Foreground(ColorText).Render(cfg.MaskedAPIKey()) + "\n")

		source := "config file"
		if cfg.APIKey == "" {
			if os.Getenv("GOOGLE_API_KEY") != "" {
				source = "GOOGLE_API_KEY env var"
			} else {
				source = "GEMINI_API_KEY env var"
			}
		}
		sb.WriteString(lipgloss.NewStyle().Foreground(ColorMuted).Render("  Source: ") +
			lipgloss.NewStyle().Foreground(ColorMuted).Italic(true).Render(source) + "\n")
	} else {
		sb.WriteString(statusLabel + ErrorStyle.Render("✗ Not configured") + "\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(ColorMuted).Italic(true).Render(
			"  Set a key to enable AI-powered analysis") + "\n")
	}

	sb.WriteString("\n")

	// Sub-menu options
	options := []struct {
		icon string
		text string
	}{
		{"", "Set API Key"},
		{"", "Remove API Key"},
		{"", "Back to Menu"},
	}

	for i, opt := range options {
		cursor := "  "
		style := lipgloss.NewStyle().Foreground(ColorText)
		if i == subCursor {
			cursor = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render("❯ ")
			style = lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary)
		}
		sb.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, opt.icon, style.Render(opt.text)))
	}

	return CardStyle.Width(60).Render(sb.String())
}

// RenderSettingsView renders the settings configuration view.
func RenderSettingsView(cfg *config.Config, subCursor int) string {
	var sb strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary).
		MarginBottom(1)

	sb.WriteString(headerStyle.Render("Settings") + "\n\n")

	// Current settings
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorText).Width(18)
	valueStyle := lipgloss.NewStyle().Foreground(ColorSecondary)

	sb.WriteString(fmt.Sprintf("  %s %s\n",
		labelStyle.Render("Model:"),
		valueStyle.Render(cfg.GetModel())))

	defaultPath := cfg.DefaultPath
	if defaultPath == "" {
		defaultPath = "(current directory)"
	}
	sb.WriteString(fmt.Sprintf("  %s %s\n",
		labelStyle.Render("Default Path:"),
		valueStyle.Render(defaultPath)))

	noColorStr := "Off"
	if cfg.NoColor {
		noColorStr = "On"
	}

	skipAIStr := "Disabled (AI is ON)"
	if cfg.SkipAI {
		skipAIStr = "Enabled (AI is OFF)"
	}

	sb.WriteString(fmt.Sprintf("  %s %s\n",
		labelStyle.Render("Skip AI:"),
		valueStyle.Render(skipAIStr)))

	sb.WriteString(fmt.Sprintf("  %s %s\n",
		labelStyle.Render("No Color:"),
		valueStyle.Render(noColorStr)))

	sb.WriteString(fmt.Sprintf("  %s %s\n",
		labelStyle.Render("Report Format:"),
		valueStyle.Render(cfg.GetReportFormat())))

	autoSaveStr := "Disabled"
	if cfg.AutoSave {
		autoSaveStr = "Enabled"
	}
	sb.WriteString(fmt.Sprintf("  %s %s\n",
		labelStyle.Render("Auto-Save Report:"),
		valueStyle.Render(autoSaveStr)))

	sb.WriteString("\n")

	// Sub-menu
	options := []struct {
		icon string
		text string
	}{
		{"", "Change Model"},
		{"", "Set Default Path"},
		{"", "Toggle AI (Offline Mode)"},
		{"", "Toggle No Color"},
		{"", "Toggle Report Format"},
		{"", "Toggle Auto-Save Report"},
		{"", "Reset All Settings"},
		{"", "Back to Menu"},
	}

	for i, opt := range options {
		cursor := "  "
		style := lipgloss.NewStyle().Foreground(ColorText)
		if i == subCursor {
			cursor = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render("❯ ")
			style = lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary)
		}
		sb.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, opt.icon, style.Render(opt.text)))
	}

	return CardStyle.Width(60).Render(sb.String())
}

// RenderHelpView renders the help/commands view.
func RenderHelpView() string {
	var sb strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary)

	cmdStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Width(30)

	descStyle := lipgloss.NewStyle().
		Foreground(ColorText)

	sb.WriteString(headerStyle.Render("Commands & Usage") + "\n\n")

	commands := []struct {
		cmd  string
		desc string
	}{
		{"obelisk", "Launch interactive TUI"},
		{"obelisk scan", "Headless scan (CI/CD friendly)"},
		{"obelisk scan --format json", "Output scan results as JSON"},
		{"obelisk scan --strict", "Exit code 1 on critical/errors"},
		{"obelisk check [path]", "Run health check with dashboard"},
		{"obelisk report [path]", "Generate & export report"},
		{"obelisk protect --install", "Install git pre-push hook"},
		{"obelisk config set api-key <k>", "Store your Gemini API key"},
		{"obelisk config list", "Show all config values"},
		{"obelisk version", "Show version info"},
	}

	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(ColorText).Render("  Commands:") + "\n\n")
	for _, c := range commands {
		sb.WriteString(fmt.Sprintf("  %s %s\n", cmdStyle.Render(c.cmd), descStyle.Render(c.desc)))
	}

	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(ColorText).Render("  Global Flags:") + "\n\n")

	flags := []struct {
		flag string
		desc string
	}{
		{"--api-key <key>", "Gemini API key"},
		{"--model <name>", "AI model (default: gemini-2.5-flash)"},
		{"--verbose, -v", "Enable verbose output"},
		{"--no-color", "Disable colored output"},
	}

	for _, f := range flags {
		sb.WriteString(fmt.Sprintf("  %s %s\n", cmdStyle.Render(f.flag), descStyle.Render(f.desc)))
	}

	sb.WriteString("\n")
	sb.WriteString(MutedStyle.Render("  Press Esc or Backspace to go back") + "\n")

	return CardStyle.Width(70).Render(sb.String())
}

// RenderProtectView renders the protect sub-menu view.
func RenderProtectView(subCursor int) string {
	var sb strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary)

	sb.WriteString(headerStyle.Render("Git Protection") + "\n\n")

	sb.WriteString(lipgloss.NewStyle().Foreground(ColorText).Render(
		"  Install Obelisk as a Git pre-push hook to automatically\n"+
			"  block pushes containing critical security issues.") + "\n\n")

	options := []struct {
		icon string
		text string
	}{
		{"", "Install Pre-Push Hook"},
		{"", "Uninstall Pre-Push Hook"},
		{"", "Run Protection Check"},
		{"", "Back to Menu"},
	}

	for i, opt := range options {
		cursor := "  "
		style := lipgloss.NewStyle().Foreground(ColorText)
		if i == subCursor {
			cursor = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render("❯ ")
			style = lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary)
		}
		sb.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, opt.icon, style.Render(opt.text)))
	}

	return CardStyle.Width(60).Render(sb.String())
}

// RenderStatusMessage renders a temporary status message.
func RenderStatusMessage(msg string, isError bool) string {
	if isError {
		return "\n" + ErrorStyle.Render("  ✗ "+msg) + "\n"
	}
	return "\n" + SuccessStyle.Render("  ✓ "+msg) + "\n"
}
