package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Swif7ify/Obelisk-CLI/internal/ai"
	"github.com/Swif7ify/Obelisk-CLI/internal/config"
	"github.com/Swif7ify/Obelisk-CLI/internal/report"
	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// InteractiveView represents the current view in the interactive TUI.
type InteractiveView int

const (
	ViewMainMenu InteractiveView = iota
	ViewScanInput
	ViewScanning
	ViewScanResults
	ViewAPIKey
	ViewAPIKeyInput
	ViewSettings
	ViewSettingsModelInput
	ViewSettingsPathInput
	ViewHelp
	ViewProtect
)

// InteractiveModel is the root Bubble Tea model for the interactive TUI.
type InteractiveModel struct {
	CurrentView   InteractiveView
	Menu          MenuModel
	Config        *config.Config
	Input         InputModel
	Spinner       SpinnerModel
	Viewport      viewport.Model
	ScanResult    *scanner.ScanResult
	ScanReport    *ai.HealthReport
	StatusMsg     string
	StatusIsError bool
	SubCursor     int
	Width         int
	Height        int
	Version       string
}

// scanCompleteMsg is sent when a scan finishes.
type scanCompleteMsg struct {
	result *scanner.ScanResult
	report *ai.HealthReport
	err    error
}

// statusClearMsg clears the status message.
type statusClearMsg struct{}

// hookResultMsg is sent when hook install completes.
type hookResultMsg struct {
	msg string
	err error
}

// NewInteractive creates a new interactive TUI model.
func NewInteractive(cfg *config.Config, version string) InteractiveModel {
	menuItems := []MenuItem{
		{Icon: "", Title: "Scan Project", Description: "Run a full health check", Key: "scan"},
		{Icon: "", Title: "Protect", Description: "Git pre-push hook", Key: "protect"},
		{Icon: "", Title: "API Key", Description: "Manage your Gemini key", Key: "apikey"},
		{Icon: "", Title: "Settings", Description: "Configure Obelisk", Key: "settings"},
		{Icon: "", Title: "Help", Description: "Commands & usage", Key: "help"},
		{Icon: "", Title: "Quit", Description: "Exit Obelisk", Key: "quit"},
	}

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().Padding(0, 2)

	return InteractiveModel{
		CurrentView: ViewMainMenu,
		Menu:        NewMenu("Main Menu", menuItems),
		Config:      cfg,
		Spinner:     NewSpinner(),
		Viewport:    vp,
		Width:       80,
		Height:      40,
		Version:     version,
	}
}

func (m InteractiveModel) Init() tea.Cmd {
	return nil
}

func (m InteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Viewport.Width = msg.Width
		// Adjust height to leave room for banner and footer
		m.Viewport.Height = msg.Height - 12
		return m, nil
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case SpinnerTickMsg:
		if m.CurrentView == ViewScanning {
			m.Spinner, _ = m.Spinner.Update(msg)
			return m, tickCmd()
		}
		return m, nil
	case PhaseUpdateMsg:
		if m.CurrentView == ViewScanning {
			m.Spinner.Phase = MapPhaseStringToSpinnerPhase(msg.Phase)
			return m, ListenForPhaseUpdates(enginePhaseChan)
		}
		return m, nil
	case scanCompleteMsg:
		m.CurrentView = ViewScanResults
		if msg.err != nil {
			m.StatusMsg = msg.err.Error()
			m.StatusIsError = true
		} else {
			m.ScanResult = msg.result
			m.ScanReport = msg.report
			
			// Generate report path first
			format := m.Config.GetReportFormat()
			if format == "" {
				format = "txt"
			}
			outputPath := report.GetDefaultOutputPath(m.ScanResult.ProjectPath, format)
			
			var content strings.Builder
			content.WriteString(RenderStats(m.ScanResult) + "\n\n")
			content.WriteString(RenderScoreCard(m.ScanReport) + "\n")
			content.WriteString(RenderFindings(m.ScanResult.Findings, outputPath) + "\n")
			content.WriteString(RenderSummary(m.ScanReport) + "\n")
			
			m.Viewport.SetContent(content.String())
			m.Viewport.GotoTop()
			
			// Auto-save report to file (matching CLI behavior)
			if m.ScanResult != nil && m.ScanReport != nil && m.Config.AutoSave && len(m.ScanResult.Findings) > 0 {
				if err := report.WriteToFile(m.ScanResult, m.ScanReport, outputPath, format); err != nil {
					// Don't show error in TUI, just silently fail
					// User can still see results in the viewport
				} else {
					// Show success message briefly
					m.StatusMsg = "Report saved to: " + outputPath
					m.StatusIsError = false
					return m, clearStatusAfter(5 * time.Second)
				}
			}
		}
		return m, nil
	case statusClearMsg:
		m.StatusMsg = ""
		return m, nil
	case hookResultMsg:
		if msg.err != nil {
			m.StatusMsg = msg.err.Error()
			m.StatusIsError = true
		} else {
			m.StatusMsg = msg.msg
			m.StatusIsError = false
		}
		return m, clearStatusAfter(3 * time.Second)
	}
	return m, nil
}

func (m InteractiveModel) View() string {
	var sb strings.Builder

	sb.WriteString(RenderBanner())
	sb.WriteString(MutedStyle.Render(fmt.Sprintf("  v%s", m.Version)) + "\n\n")

	switch m.CurrentView {
	case ViewMainMenu:
		sb.WriteString(m.Menu.View())
		sb.WriteString("\n")
		sb.WriteString(MutedStyle.Render("  ↑/↓ Navigate  Enter Select  q Quit") + "\n")
	case ViewScanInput:
		sb.WriteString(SubtitleStyle.Render("  Scan Project") + "\n\n")
		sb.WriteString(m.Input.View() + "\n\n")
		sb.WriteString(MutedStyle.Render("  Enter to scan  Esc to cancel") + "\n")
	case ViewScanning:
		sb.WriteString(m.Spinner.View())
		sb.WriteString("\n" + MutedStyle.Render("  Scanning project... press q to abort") + "\n")
	case ViewScanResults:
		if m.StatusIsError && m.StatusMsg != "" {
			sb.WriteString(ErrorStyle.Render("  ✗ Error: "+m.StatusMsg) + "\n\n")
		} else if m.ScanReport != nil && m.ScanResult != nil {
			sb.WriteString(m.Viewport.View() + "\n")
		}
		sb.WriteString("\n" + MutedStyle.Render("  Press Esc or Backspace to go back, ↑/↓ to scroll") + "\n")
	case ViewAPIKey:
		sb.WriteString(RenderAPIKeyView(m.Config, m.SubCursor))
	case ViewAPIKeyInput:
		sb.WriteString(SubtitleStyle.Render("  API Key") + "\n\n")
		sb.WriteString(m.Input.View() + "\n\n")
		sb.WriteString(MutedStyle.Render("  Enter to save  Esc to cancel") + "\n")
	case ViewSettings:
		sb.WriteString(RenderSettingsView(m.Config, m.SubCursor))
	case ViewSettingsModelInput:
		sb.WriteString(SubtitleStyle.Render("  Set AI Model") + "\n\n")
		sb.WriteString(m.Input.View() + "\n\n")
		sb.WriteString(MutedStyle.Render("  Models: gemini-2.5-flash, gemini-2.5-pro") + "\n")
		sb.WriteString(MutedStyle.Render("  Enter to save  Esc to cancel") + "\n")
	case ViewSettingsPathInput:
		sb.WriteString(SubtitleStyle.Render("  Set Default Path") + "\n\n")
		sb.WriteString(m.Input.View() + "\n\n")
		sb.WriteString(MutedStyle.Render("  Leave empty for current directory") + "\n")
		sb.WriteString(MutedStyle.Render("  Enter to save  Esc to cancel") + "\n")
	case ViewHelp:
		sb.WriteString(RenderHelpView())
	case ViewProtect:
		sb.WriteString(RenderProtectView(m.SubCursor))
	}

	if m.StatusMsg != "" {
		sb.WriteString(RenderStatusMessage(m.StatusMsg, m.StatusIsError))
	}

	return sb.String()
}

func clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return statusClearMsg{}
	})
}
