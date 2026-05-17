package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// ConfirmationDialog represents a confirmation dialog
type ConfirmationDialog struct {
	Title    string
	Message  string
	Visible  bool
	Selected int // 0 = Yes, 1 = No
}

// NewConfirmationDialog creates a new confirmation dialog
func NewConfirmationDialog(title, message string) ConfirmationDialog {
	return ConfirmationDialog{
		Title:    title,
		Message:  message,
		Visible:  false,
		Selected: 1, // Default to "No" for safety
	}
}

// Show displays the confirmation dialog
func (c *ConfirmationDialog) Show() {
	c.Visible = true
	c.Selected = 1 // Reset to "No"
}

// Hide hides the confirmation dialog
func (c *ConfirmationDialog) Hide() {
	c.Visible = false
}

// Toggle toggles the selected option
func (c *ConfirmationDialog) Toggle() {
	c.Selected = 1 - c.Selected
}

// MoveLeft moves selection to the left
func (c *ConfirmationDialog) MoveLeft() {
	if c.Selected > 0 {
		c.Selected--
	}
}

// MoveRight moves selection to the right
func (c *ConfirmationDialog) MoveRight() {
	if c.Selected < 1 {
		c.Selected++
	}
}

// IsYes returns true if "Yes" is selected
func (c ConfirmationDialog) IsYes() bool {
	return c.Selected == 0
}

// View renders the confirmation dialog
func (c ConfirmationDialog) View(width, height int) string {
	if !c.Visible {
		return ""
	}

	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(ColorWarning).
		Padding(1, 2).
		Width(50).
		Background(ColorCardBg)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorWarning).
		Align(lipgloss.Center).
		Width(46)

	messageStyle := lipgloss.NewStyle().
		Foreground(ColorText).
		Align(lipgloss.Center).
		Width(46).
		MarginTop(1).
		MarginBottom(1)

	buttonActiveStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorTextBold).
		Background(ColorPrimary).
		Padding(0, 3)

	buttonInactiveStyle := lipgloss.NewStyle().
		Foreground(ColorText).
		Border(lipgloss.NormalBorder()).
		BorderForeground(ColorMuted).
		Padding(0, 3)

	// Render buttons
	yesButton := buttonInactiveStyle.Render("Yes")
	noButton := buttonInactiveStyle.Render("No")

	if c.Selected == 0 {
		yesButton = buttonActiveStyle.Render("Yes")
	} else {
		noButton = buttonActiveStyle.Render("No")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, yesButton, "  ", noButton)
	buttonsRow := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(46).
		Render(buttons)

	hint := MutedStyle.Render("←/→ or Tab to switch  Enter to confirm  Esc to cancel")
	hintRow := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(46).
		MarginTop(1).
		Render(hint)

	content := titleStyle.Render("⚠️  " + c.Title) + "\n" +
		messageStyle.Render(c.Message) + "\n" +
		buttonsRow + "\n" +
		hintRow

	dialog := dialogStyle.Render(content)

	// Center the dialog
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

// Made with Bob
