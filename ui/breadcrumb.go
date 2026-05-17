package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// BreadcrumbModel represents a navigation breadcrumb trail
type BreadcrumbModel struct {
	Path []string
}

// NewBreadcrumb creates a new breadcrumb model
func NewBreadcrumb(path ...string) BreadcrumbModel {
	return BreadcrumbModel{Path: path}
}

// Push adds a new item to the breadcrumb trail
func (b *BreadcrumbModel) Push(item string) {
	b.Path = append(b.Path, item)
}

// Pop removes the last item from the breadcrumb trail
func (b *BreadcrumbModel) Pop() {
	if len(b.Path) > 0 {
		b.Path = b.Path[:len(b.Path)-1]
	}
}

// Reset clears the breadcrumb trail
func (b *BreadcrumbModel) Reset() {
	b.Path = []string{}
}

// View renders the breadcrumb trail
func (b BreadcrumbModel) View() string {
	if len(b.Path) == 0 {
		return ""
	}

	separatorStyle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render(" › ")

	itemStyle := lipgloss.NewStyle().
		Foreground(ColorText)

	lastItemStyle := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	var parts []string
	for i, item := range b.Path {
		if i == len(b.Path)-1 {
			parts = append(parts, lastItemStyle.Render(item))
		} else {
			parts = append(parts, itemStyle.Render(item))
		}
	}

	breadcrumb := strings.Join(parts, separatorStyle)
	return lipgloss.NewStyle().
		Foreground(ColorMuted).
		Padding(0, 2).
		Render(breadcrumb) + "\n"
}

// GetViewBreadcrumb returns the breadcrumb path for a given view
func GetViewBreadcrumb(view InteractiveView) []string {
	breadcrumbs := map[InteractiveView][]string{
		ViewMainMenu:             {"Home"},
		ViewScanInput:            {"Home", "Scan Project"},
		ViewScanning:             {"Home", "Scan Project", "Scanning"},
		ViewScanResults:          {"Home", "Scan Project", "Results"},
		ViewAPIKey:               {"Home", "API Key"},
		ViewAPIKeyInput:          {"Home", "API Key", "Set Key"},
		ViewSettings:             {"Home", "Settings"},
		ViewSettingsModelInput:   {"Home", "Settings", "Change Model"},
		ViewSettingsPathInput:    {"Home", "Settings", "Set Path"},
		ViewHelp:                 {"Home", "Help"},
		ViewProtect:              {"Home", "Git Protection"},
	}
	return breadcrumbs[view]
}

// Made with Bob
