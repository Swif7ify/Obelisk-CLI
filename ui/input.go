package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// InputModel is a simple text input component.
type InputModel struct {
	Prompt      string
	Value       string
	Placeholder string
	Masked      bool // hide input characters (for API keys)
	Cursor      int
	Focused     bool
	Width       int
}

// NewInput creates a new input model.
func NewInput(prompt, placeholder string, masked bool) InputModel {
	return InputModel{
		Prompt:      prompt,
		Placeholder: placeholder,
		Masked:      masked,
		Focused:     true,
		Width:       50,
	}
}

// InsertChar adds a character at the cursor position.
func (m *InputModel) InsertChar(ch rune) {
	if m.Cursor >= len(m.Value) {
		m.Value += string(ch)
	} else {
		m.Value = m.Value[:m.Cursor] + string(ch) + m.Value[m.Cursor:]
	}
	m.Cursor++
}

// DeleteChar removes the character before the cursor.
func (m *InputModel) DeleteChar() {
	if m.Cursor > 0 && len(m.Value) > 0 {
		m.Value = m.Value[:m.Cursor-1] + m.Value[m.Cursor:]
		m.Cursor--
	}
}

// DeleteForward removes the character at the cursor.
func (m *InputModel) DeleteForward() {
	if m.Cursor < len(m.Value) {
		m.Value = m.Value[:m.Cursor] + m.Value[m.Cursor+1:]
	}
}

// MoveLeft moves the cursor left.
func (m *InputModel) MoveLeft() {
	if m.Cursor > 0 {
		m.Cursor--
	}
}

// MoveRight moves the cursor right.
func (m *InputModel) MoveRight() {
	if m.Cursor < len(m.Value) {
		m.Cursor++
	}
}

// MoveToStart moves cursor to start.
func (m *InputModel) MoveToStart() {
	m.Cursor = 0
}

// MoveToEnd moves cursor to end.
func (m *InputModel) MoveToEnd() {
	m.Cursor = len(m.Value)
}

// Clear resets the input value.
func (m *InputModel) Clear() {
	m.Value = ""
	m.Cursor = 0
}

// View renders the input field.
func (m InputModel) View() string {
	promptStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary)

	inputStyle := lipgloss.NewStyle().
		Foreground(ColorText)

	cursorStyle := lipgloss.NewStyle().
		Background(ColorPrimary).
		Foreground(lipgloss.Color("#FFFFFF"))

	placeholderStyle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Italic(true)

	var display string
	if m.Value == "" && !m.Focused {
		display = placeholderStyle.Render(m.Placeholder)
	} else if m.Value == "" {
		display = cursorStyle.Render(" ") + placeholderStyle.Render(m.Placeholder)
	} else {
		val := m.Value
		if m.Masked {
			val = strings.Repeat("•", len(val))
		}

		if m.Focused && m.Cursor <= len(val) {
			before := val[:m.Cursor]
			after := ""
			cursorChar := " "
			if m.Cursor < len(val) {
				cursorChar = string(val[m.Cursor])
				after = val[m.Cursor+1:]
			}
			display = inputStyle.Render(before) +
				cursorStyle.Render(cursorChar) +
				inputStyle.Render(after)
		} else {
			display = inputStyle.Render(val)
		}
	}

	return fmt.Sprintf("  %s %s", promptStyle.Render(m.Prompt), display)
}
