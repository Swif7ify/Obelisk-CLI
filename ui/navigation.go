package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// NavigationHints represents keyboard shortcuts for a view
type NavigationHints struct {
	Primary   []KeyHint
	Secondary []KeyHint
}

// KeyHint represents a single keyboard shortcut
type KeyHint struct {
	Key         string
	Description string
}

// RenderNavigationFooter renders standardized navigation hints for each view
func RenderNavigationFooter(view InteractiveView) string {
	hints := getNavigationHints(view)
	
	var parts []string
	
	// Primary hints (always shown)
	for _, hint := range hints.Primary {
		parts = append(parts, renderKeyHint(hint))
	}
	
	// Secondary hints (optional)
	for _, hint := range hints.Secondary {
		parts = append(parts, renderKeyHint(hint))
	}
	
	footer := lipgloss.JoinHorizontal(lipgloss.Left, parts...)
	return "\n" + MutedStyle.Render("  "+footer) + "\n"
}

// renderKeyHint renders a single key hint
func renderKeyHint(hint KeyHint) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)
	
	return keyStyle.Render(hint.Key) + " " + hint.Description + "  "
}

// getNavigationHints returns the navigation hints for a specific view
func getNavigationHints(view InteractiveView) NavigationHints {
	hintsMap := map[InteractiveView]NavigationHints{
		ViewMainMenu: {
			Primary: []KeyHint{
				{"↑/↓", "Navigate"},
				{"Enter", "Select"},
				{"q", "Quit"},
			},
		},
		ViewScanInput: {
			Primary: []KeyHint{
				{"Enter", "Start Scan"},
				{"Esc", "Cancel"},
			},
			Secondary: []KeyHint{
				{"Ctrl+U", "Clear"},
				{"h", "Home"},
			},
		},
		ViewScanning: {
			Primary: []KeyHint{
				{"q", "Abort"},
				{"Esc", "Abort"},
			},
		},
		ViewScanResults: {
			Primary: []KeyHint{
				{"↑/↓", "Scroll"},
				{"Esc", "Back"},
			},
			Secondary: []KeyHint{
				{"s", "Save"},
				{"h", "Home"},
				{"q", "Quit"},
			},
		},
		ViewAPIKey: {
			Primary: []KeyHint{
				{"↑/↓", "Navigate"},
				{"Enter", "Select"},
				{"Esc", "Back"},
			},
			Secondary: []KeyHint{
				{"h", "Home"},
			},
		},
		ViewAPIKeyInput: {
			Primary: []KeyHint{
				{"Enter", "Save"},
				{"Esc", "Cancel"},
			},
			Secondary: []KeyHint{
				{"Ctrl+U", "Clear"},
			},
		},
		ViewSettings: {
			Primary: []KeyHint{
				{"↑/↓", "Navigate"},
				{"Enter", "Select"},
				{"Esc", "Back"},
			},
			Secondary: []KeyHint{
				{"h", "Home"},
			},
		},
		ViewSettingsModelInput: {
			Primary: []KeyHint{
				{"Enter", "Save"},
				{"Esc", "Cancel"},
			},
		},
		ViewSettingsPathInput: {
			Primary: []KeyHint{
				{"Enter", "Save"},
				{"Esc", "Cancel"},
			},
		},
		ViewHelp: {
			Primary: []KeyHint{
				{"Esc", "Back"},
				{"q", "Quit"},
			},
			Secondary: []KeyHint{
				{"h", "Home"},
			},
		},
		ViewProtect: {
			Primary: []KeyHint{
				{"↑/↓", "Navigate"},
				{"Enter", "Select"},
				{"Esc", "Back"},
			},
			Secondary: []KeyHint{
				{"h", "Home"},
			},
		},
	}
	
	hints, ok := hintsMap[view]
	if !ok {
		return NavigationHints{
			Primary: []KeyHint{
				{"Esc", "Back"},
				{"q", "Quit"},
			},
		}
	}
	return hints
}

// Made with Bob
