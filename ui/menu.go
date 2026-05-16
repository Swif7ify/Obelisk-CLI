package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a single menu entry.
type MenuItem struct {
	Icon        string
	Title       string
	Description string
	Key         string // internal key for identification
}

// MenuModel is a reusable navigable menu component.
type MenuModel struct {
	Items    []MenuItem
	Cursor   int
	Title    string
	Width    int
}

// NewMenu creates a new menu model.
func NewMenu(title string, items []MenuItem) MenuModel {
	return MenuModel{
		Items:  items,
		Cursor: 0,
		Title:  title,
		Width:  60,
	}
}

// MoveUp moves the cursor up.
func (m *MenuModel) MoveUp() {
	if m.Cursor > 0 {
		m.Cursor--
	}
}

// MoveDown moves the cursor down.
func (m *MenuModel) MoveDown() {
	if m.Cursor < len(m.Items)-1 {
		m.Cursor++
	}
}

// Selected returns the currently highlighted menu item.
func (m *MenuModel) Selected() MenuItem {
	if m.Cursor >= 0 && m.Cursor < len(m.Items) {
		return m.Items[m.Cursor]
	}
	return MenuItem{}
}

// View renders the menu.
func (m MenuModel) View() string {
	var sb strings.Builder

	activeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorTextBold).
		Background(ColorPrimary).
		Padding(0, 2).
		MarginLeft(2)

	inactiveStyle := lipgloss.NewStyle().
		Foreground(ColorText).
		Padding(0, 1).
		MarginLeft(2)

	descActiveStyle := lipgloss.NewStyle().
		Foreground(ColorHighlight).
		Italic(true)

	descInactiveStyle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Italic(true)

	cursorActive := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("▶")

	cursorInactive := lipgloss.NewStyle().
		Render(" ")

	for i, item := range m.Items {
		cursor := cursorInactive
		var label string
		var desc string

		if i == m.Cursor {
			cursor = cursorActive
			label = activeStyle.Render(fmt.Sprintf("%s %s", item.Icon, item.Title))
			desc = descActiveStyle.Render(item.Description)
		} else {
			label = inactiveStyle.Render(fmt.Sprintf("%s %s", item.Icon, item.Title))
			desc = descInactiveStyle.Render(item.Description)
		}

		line := fmt.Sprintf(" %s %s  %s", cursor, label, desc)
		sb.WriteString(line + "\n")
	}

	return sb.String()
}
