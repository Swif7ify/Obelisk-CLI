package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Swif7ify/Obelisk-CLI/internal/config"
	"github.com/Swif7ify/Obelisk-CLI/internal/engine"
)

func (m InteractiveModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}
	switch m.CurrentView {
	case ViewMainMenu:
		return m.handleMainMenu(msg.String())
	case ViewScanInput:
		return m.handleTextInput(msg, ViewMainMenu, m.onScanSubmit)
	case ViewScanning:
		if msg.String() == "q" || msg.String() == "esc" {
			return m, tea.Quit
		}
	case ViewScanResults:
		if msg.String() == "esc" || msg.String() == "backspace" || msg.String() == "q" || msg.String() == "b" {
			m.CurrentView = ViewMainMenu
			m.ScanResult = nil
			m.ScanReport = nil
			m.StatusMsg = ""
			return m, nil
		}
		var cmd tea.Cmd
		m.Viewport, cmd = m.Viewport.Update(msg)
		return m, cmd
	case ViewAPIKey:
		return m.handleAPIKeyMenu(msg.String())
	case ViewAPIKeyInput:
		return m.handleTextInput(msg, ViewAPIKey, m.onAPIKeySubmit)
	case ViewSettings:
		return m.handleSettingsMenu(msg.String())
	case ViewSettingsModelInput:
		return m.handleTextInput(msg, ViewSettings, m.onModelSubmit)
	case ViewSettingsPathInput:
		return m.handleTextInput(msg, ViewSettings, m.onPathSubmit)
	case ViewHelp:
		if msg.String() == "esc" || msg.String() == "backspace" || msg.String() == "q" {
			m.CurrentView = ViewMainMenu
		}
	case ViewProtect:
		return m.handleProtectMenu(msg.String())
	}
	return m, nil
}

// --- Main Menu ---

func (m InteractiveModel) handleMainMenu(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		m.Menu.MoveUp()
	case "down", "j":
		m.Menu.MoveDown()
	case "enter":
		switch m.Menu.Selected().Key {
		case "scan":
			path := m.Config.DefaultPath
			if path == "" {
				path, _ = os.Getwd()
			}
			m.Input = NewInput("Project path:", path, false)
			m.Input.SetValue(path)
			m.CurrentView = ViewScanInput
		case "protect":
			m.SubCursor = 0
			m.CurrentView = ViewProtect
		case "apikey":
			m.SubCursor = 0
			m.CurrentView = ViewAPIKey
		case "settings":
			m.SubCursor = 0
			m.CurrentView = ViewSettings
		case "help":
			m.CurrentView = ViewHelp
		case "quit":
			return m, tea.Quit
		}
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

// --- Generic text input handler ---

type submitFunc func(m InteractiveModel, value string) (InteractiveModel, tea.Cmd)

func (m InteractiveModel) handleTextInput(msg tea.KeyMsg, cancelView InteractiveView, onSubmit submitFunc) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		return onSubmit(m, strings.TrimSpace(m.Input.GetValue()))
	case tea.KeyEsc:
		m.CurrentView = cancelView
		return m, nil
	case tea.KeyBackspace:
		m.Input.DeleteChar()
	case tea.KeyDelete:
		m.Input.DeleteForward()
	case tea.KeyLeft:
		m.Input.MoveLeft()
	case tea.KeyRight:
		m.Input.MoveRight()
	case tea.KeyHome, tea.KeyCtrlA:
		m.Input.MoveToStart()
	case tea.KeyEnd, tea.KeyCtrlE:
		m.Input.MoveToEnd()
	case tea.KeyCtrlU:
		m.Input.Clear()
	case tea.KeySpace:
		m.Input.InsertChar(' ')
	case tea.KeyRunes:
		for _, r := range msg.Runes {
			if r >= 32 && r <= 126 && r != 0xFFFD {
				m.Input.InsertChar(r)
			}
		}
	}
	return m, nil
}

// --- Submit callbacks ---

func (m InteractiveModel) onScanSubmit(im InteractiveModel, value string) (InteractiveModel, tea.Cmd) {
	// Strip spaces, quotes, and control characters (like \r, \n)
	path := strings.TrimFunc(value, func(r rune) bool {
		return r < 32 || r == '"' || r == '\'' || r == ' '
	})
	
	if path == "" {
		path, _ = os.Getwd()
	}
	if _, err := os.Stat(path); err != nil {
		im.StatusMsg = "Invalid path: " + err.Error()
		im.StatusIsError = true
		return im, clearStatusAfter(4 * time.Second)
	}
	im.CurrentView = ViewScanning
	im.Spinner = NewSpinner()
	return im, tea.Batch(tickCmd(), startScanCmd(path, im.Config))
}

func (m InteractiveModel) onAPIKeySubmit(im InteractiveModel, value string) (InteractiveModel, tea.Cmd) {
	// Clean the API key - remove any control characters, spaces, and invalid Unicode
	cleaned := strings.Map(func(r rune) rune {
		// Keep only valid ASCII alphanumeric, dash, and underscore
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		// Remove everything else (spaces, control chars, invalid Unicode)
		return -1
	}, value)
	
	if cleaned == "" {
		im.StatusMsg = "API key cannot be empty"
		im.StatusIsError = true
		return im, clearStatusAfter(3 * time.Second)
	}
	
	im.Config.SetAPIKey(cleaned)
	if err := im.Config.Save(); err != nil {
		im.StatusMsg = "Failed to save: " + err.Error()
		im.StatusIsError = true
	} else {
		im.StatusMsg = "API key saved successfully!"
		im.StatusIsError = false
	}
	im.CurrentView = ViewAPIKey
	return im, clearStatusAfter(3 * time.Second)
}

func (m InteractiveModel) onModelSubmit(im InteractiveModel, value string) (InteractiveModel, tea.Cmd) {
	if value != "" {
		im.Config.Model = value
	}
	if err := im.Config.Save(); err != nil {
		im.StatusMsg = "Failed to save: " + err.Error()
		im.StatusIsError = true
	} else {
		im.StatusMsg = "Model set to: " + im.Config.GetModel()
		im.StatusIsError = false
	}
	im.CurrentView = ViewSettings
	return im, clearStatusAfter(3 * time.Second)
}

func (m InteractiveModel) onPathSubmit(im InteractiveModel, value string) (InteractiveModel, tea.Cmd) {
	im.Config.DefaultPath = value
	if err := im.Config.Save(); err != nil {
		im.StatusMsg = "Failed to save: " + err.Error()
		im.StatusIsError = true
	} else {
		d := value
		if d == "" {
			d = "(current directory)"
		}
		im.StatusMsg = "Default path: " + d
		im.StatusIsError = false
	}
	im.CurrentView = ViewSettings
	return im, clearStatusAfter(3 * time.Second)
}

// --- API Key Sub-Menu ---

func (m InteractiveModel) handleAPIKeyMenu(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		if m.SubCursor > 0 {
			m.SubCursor--
		}
	case "down", "j":
		if m.SubCursor < 2 {
			m.SubCursor++
		}
	case "enter":
		switch m.SubCursor {
		case 0: // Set
			m.Input = NewInput("API Key:", "paste your Gemini API key", true)
			m.CurrentView = ViewAPIKeyInput
		case 1: // Remove
			m.Config.ClearAPIKey()
			_ = m.Config.Save()
			m.StatusMsg = "API key removed"
			m.StatusIsError = false
			return m, clearStatusAfter(3 * time.Second)
		case 2: // Back
			m.CurrentView = ViewMainMenu
		}
	case "esc", "backspace":
		m.CurrentView = ViewMainMenu
	}
	return m, nil
}

// --- Settings Sub-Menu ---

func (m InteractiveModel) handleSettingsMenu(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		if m.SubCursor > 0 {
			m.SubCursor--
		}
	case "down", "j":
		if m.SubCursor < 7 {
			m.SubCursor++
		}
	case "enter":
		switch m.SubCursor {
		case 0: // Model
			m.Input = NewInput("Model:", m.Config.GetModel(), false)
			m.Input.SetValue(m.Config.GetModel())
			m.CurrentView = ViewSettingsModelInput
		case 1: // Path
			m.Input = NewInput("Default path:", "(current directory)", false)
			m.Input.SetValue(m.Config.DefaultPath)
			m.CurrentView = ViewSettingsPathInput
		case 2: // Toggle AI
			m.Config.SkipAI = !m.Config.SkipAI
			_ = m.Config.Save()
			s := "ON"
			if m.Config.SkipAI {
				s = "OFF"
			}
			m.StatusMsg = "AI is now " + s
			m.StatusIsError = false
			return m, clearStatusAfter(3 * time.Second)
		case 3: // Toggle no-color
			m.Config.NoColor = !m.Config.NoColor
			_ = m.Config.Save()
			s := "disabled"
			if m.Config.NoColor {
				s = "enabled"
			}
			m.StatusMsg = "No color: " + s
			m.StatusIsError = false
			return m, clearStatusAfter(3 * time.Second)
		case 4: // Toggle Report Format
			if m.Config.GetReportFormat() == "md" {
				m.Config.ReportFormat = "txt"
			} else {
				m.Config.ReportFormat = "md"
			}
			_ = m.Config.Save()
			m.StatusMsg = "Report format: " + m.Config.ReportFormat
			m.StatusIsError = false
			return m, clearStatusAfter(3 * time.Second)
		case 5: // Toggle Auto-Save
			m.Config.AutoSave = !m.Config.AutoSave
			_ = m.Config.Save()
			s := "Disabled"
			if m.Config.AutoSave {
				s = "Enabled"
			}
			m.StatusMsg = "Auto-save report: " + s
			m.StatusIsError = false
			return m, clearStatusAfter(3 * time.Second)
		case 6: // Reset
			m.Config.Reset()
			_ = m.Config.Save()
			m.StatusMsg = "Settings reset to defaults"
			m.StatusIsError = false
			return m, clearStatusAfter(3 * time.Second)
		case 7: // Back
			m.CurrentView = ViewMainMenu
		}
	case "esc", "backspace":
		m.CurrentView = ViewMainMenu
	}
	return m, nil
}

// --- Protect Sub-Menu ---

func (m InteractiveModel) handleProtectMenu(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		if m.SubCursor > 0 {
			m.SubCursor--
		}
	case "down", "j":
		if m.SubCursor < 3 {
			m.SubCursor++
		}
	case "enter":
		switch m.SubCursor {
		case 0: // Install hook
			return m, installHookCmd()
		case 1: // Uninstall hook
			return m, uninstallHookCmd()
		case 2: // Run check
			m.CurrentView = ViewScanning
			m.Spinner = NewSpinner()
			return m, tea.Batch(tickCmd(), runProtectCheckCmd(m.Config))
		case 3: // Back
			m.CurrentView = ViewMainMenu
		}
	case "esc", "backspace":
		m.CurrentView = ViewMainMenu
	}
	return m, nil
}

// --- Tea Commands ---

var enginePhaseChan chan string

func startScanCmd(path string, cfg *config.Config) tea.Cmd {
	enginePhaseChan = make(chan string, 100)
	
	// Create a channel for the final result
	resultChan := make(chan scanCompleteMsg)
	
	// Start the engine in a goroutine
	go func() {
		ecfg := engine.Config{
			ProjectPath: path,
			APIKey:      cfg.GetAPIKey(),
			Model:       cfg.GetModel(),
			SkipAI:      cfg.SkipAI || cfg.GetAPIKey() == "",
		}
		result, err := engine.Run(ecfg, func(phase string) {
			enginePhaseChan <- phase
		})
		
		if err != nil {
			resultChan <- scanCompleteMsg{err: err}
		} else {
			resultChan <- scanCompleteMsg{result: result.ScanResult, report: result.Report}
		}
	}()
	
	// Return a batch command that listens for BOTH the first phase update AND the final result
	return tea.Batch(
		ListenForPhaseUpdates(enginePhaseChan),
		waitForScanResult(resultChan),
	)
}

func waitForScanResult(c chan scanCompleteMsg) tea.Cmd {
	return func() tea.Msg {
		return <-c
	}
}

func installHookCmd() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("git", "rev-parse", "--git-dir")
		output, err := cmd.Output()
		if err != nil {
			return hookResultMsg{err: fmt.Errorf("not a git repository")}
		}
		dir := strings.TrimSpace(string(output))
		absDir, _ := filepath.Abs(dir)
		hookPath := filepath.Join(absDir, "hooks", "pre-push")
		hook := "#!/bin/sh\necho \"🏛️ Obelisk pre-push check...\"\nobelisk protect --strict\nif [ $? -ne 0 ]; then\n  echo \"❌ Push blocked\"\n  exit 1\nfi\necho \"✅ Passed!\"\nexit 0\n"
		_ = os.MkdirAll(filepath.Dir(hookPath), 0755)
		perm := os.FileMode(0755)
		if runtime.GOOS == "windows" {
			perm = 0644
		}
		if err := os.WriteFile(hookPath, []byte(hook), perm); err != nil {
			return hookResultMsg{err: err}
		}
		return hookResultMsg{msg: "Pre-push hook installed!"}
	}
}

func uninstallHookCmd() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("git", "rev-parse", "--git-dir")
		output, err := cmd.Output()
		if err != nil {
			return hookResultMsg{err: fmt.Errorf("not a git repository")}
		}
		dir := strings.TrimSpace(string(output))
		absDir, _ := filepath.Abs(dir)
		hookPath := filepath.Join(absDir, "hooks", "pre-push")
		
		if _, err := os.Stat(hookPath); os.IsNotExist(err) {
			return hookResultMsg{err: fmt.Errorf("no pre-push hook found to uninstall")}
		}
		
		if err := os.Remove(hookPath); err != nil {
			return hookResultMsg{err: err}
		}
		return hookResultMsg{msg: "Pre-push hook uninstalled!"}
	}
}

func runProtectCheckCmd(cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		path, _ := os.Getwd()
		ecfg := engine.Config{
			ProjectPath: path,
			APIKey:      cfg.GetAPIKey(),
			Model:       cfg.GetModel(),
			SkipAI:      true, // no AI for protect check
		}
		
		result, err := engine.Run(ecfg, func(phase string) {
			// Do nothing for phase updates since we're not listening to them here
		})
		if err != nil {
			return hookResultMsg{err: err}
		}
		
		criticals := result.ScanResult.CountBySeverity(3) // SeverityCritical
		errors := result.ScanResult.CountBySeverity(2)    // SeverityError
		
		if criticals > 0 || errors > 0 {
			return hookResultMsg{err: fmt.Errorf("protection check failed: %d critical, %d errors", criticals, errors)}
		}
		
		return hookResultMsg{msg: "Protection Check Passed! No critical issues."}
	}
}
