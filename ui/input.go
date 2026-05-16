package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// InputModel is a simple text input component.
type InputModel struct {
	Prompt      string
	Value       []rune
	Placeholder string
	Masked      bool // hide input characters (for API keys)
	Cursor      int  // index in the rune slice
	Focused     bool
	Width       int
}

// NewInput creates a new input model.
func NewInput(prompt, placeholder string, masked bool) InputModel {
	return InputModel{
		Prompt:      prompt,
		Value:       []rune{},
		Placeholder: placeholder,
		Masked:      masked,
		Focused:     true,
		Width:       50,
	}
}

// SetValue sets the input value from a string.
func (m *InputModel) SetValue(v string) {
	m.Value = []rune(v)
	m.Cursor = len(m.Value)
}

// GetValue returns the input value as a string.
func (m *InputModel) GetValue() string {
	return string(m.Value)
}

// InsertChar adds a character at the cursor position.
func (m *InputModel) InsertChar(ch rune) {
	if m.Cursor >= len(m.Value) {
		m.Value = append(m.Value, ch)
	} else {
		m.Value = append(m.Value[:m.Cursor], append([]rune{ch}, m.Value[m.Cursor:]...)...)
	}
	m.Cursor++
}

// DeleteChar removes the character before the cursor.
func (m *InputModel) DeleteChar() {
	if m.Cursor > 0 && len(m.Value) > 0 {
		m.Value = append(m.Value[:m.Cursor-1], m.Value[m.Cursor:]...)
		m.Cursor--
	}
}

// DeleteForward removes the character at the cursor.
func (m *InputModel) DeleteForward() {
	if m.Cursor < len(m.Value) {
		m.Value = append(m.Value[:m.Cursor], m.Value[m.Cursor+1:]...)
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
	m.Value = []rune{}
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
	if len(m.Value) == 0 && !m.Focused {
		display = placeholderStyle.Render(m.Placeholder)
	} else if len(m.Value) == 0 {
		display = cursorStyle.Render(" ") + placeholderStyle.Render(m.Placeholder)
	} else {
		var valRunes []rune
		if m.Masked {
			for i := 0; i < len(m.Value); i++ {
				valRunes = append(valRunes, '•')
			}
		} else {
			valRunes = m.Value
		}

		if m.Focused && m.Cursor <= len(valRunes) {
			before := string(valRunes[:m.Cursor])
			after := ""
			cursorChar := " "
			if m.Cursor < len(valRunes) {
				cursorChar = string(valRunes[m.Cursor])
				after = string(valRunes[m.Cursor+1:])
			}
			display = inputStyle.Render(before) +
				cursorStyle.Render(cursorChar) +
				inputStyle.Render(after)
		} else {
			display = inputStyle.Render(string(valRunes))
		}
	}

	return fmt.Sprintf("  %s %s", promptStyle.Render(m.Prompt), display)
}
