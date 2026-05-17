package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// EnhancedSpinnerModel is an improved spinner with progress tracking
type EnhancedSpinnerModel struct {
	Phase         SpinnerPhase
	frame         int
	FilesScanned  int
	TotalFiles    int
	StartTime     time.Time
	CurrentFile   string
	RecentFiles   []string
	MaxRecentShow int
}

// NewEnhancedSpinner creates a new enhanced spinner model
func NewEnhancedSpinner() EnhancedSpinnerModel {
	return EnhancedSpinnerModel{
		Phase:         PhaseSecrets,
		frame:         0,
		StartTime:     time.Now(),
		RecentFiles:   []string{},
		MaxRecentShow: 3,
	}
}

// UpdateProgress updates the progress information
func (m *EnhancedSpinnerModel) UpdateProgress(filesScanned, totalFiles int, currentFile string) {
	m.FilesScanned = filesScanned
	m.TotalFiles = totalFiles
	m.CurrentFile = currentFile
	
	// Add to recent files
	if currentFile != "" && (len(m.RecentFiles) == 0 || m.RecentFiles[len(m.RecentFiles)-1] != currentFile) {
		m.RecentFiles = append(m.RecentFiles, currentFile)
		if len(m.RecentFiles) > m.MaxRecentShow {
			m.RecentFiles = m.RecentFiles[1:]
		}
	}
}

func (m EnhancedSpinnerModel) Update(msg tea.Msg) (EnhancedSpinnerModel, tea.Cmd) {
	switch msg.(type) {
	case SpinnerTickMsg:
		m.frame = (m.frame + 1) % len(spinnerFrames)
	}
	return m, nil
}

func (m EnhancedSpinnerModel) View() string {
	if m.Phase == PhaseDone {
		elapsed := time.Since(m.StartTime).Round(time.Second)
		return SuccessStyle.Render(fmt.Sprintf("✓ Scan complete! (%s)", elapsed))
	}

	frame := spinnerFrames[m.frame]
	label := phaseLabels[m.Phase]

	var sb strings.Builder
	
	// Phase header with spinner
	sb.WriteString("\n")
	sb.WriteString(SubtitleStyle.Render(frame) + " " + label + "\n\n")

	// Progress bar
	total := int(PhaseDone)
	current := int(m.Phase)
	barWidth := 40
	filled := current * barWidth / total
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	
	progressText := fmt.Sprintf("Phase %d/%d", current+1, total)
	sb.WriteString(MutedStyle.Render("  ["+bar+"] ") + 
		lipgloss.NewStyle().Foreground(ColorPrimary).Render(progressText) + "\n\n")

	// File progress (if available)
	if m.TotalFiles > 0 {
		filePercent := int(float64(m.FilesScanned) / float64(m.TotalFiles) * 100)
		fileBar := strings.Repeat("▓", filePercent/2) + strings.Repeat("░", 50-filePercent/2)
		
		sb.WriteString(lipgloss.NewStyle().Foreground(ColorText).Render("  Files: ") +
			lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true).
				Render(fmt.Sprintf("%d/%d", m.FilesScanned, m.TotalFiles)) + "\n")
		sb.WriteString(MutedStyle.Render("  ["+fileBar+"] ") +
			lipgloss.NewStyle().Foreground(ColorSecondary).Render(fmt.Sprintf("%d%%", filePercent)) + "\n\n")
	}

	// Time tracking
	elapsed := time.Since(m.StartTime)
	sb.WriteString(MutedStyle.Render(fmt.Sprintf("  Elapsed: %s", elapsed.Round(time.Second))) + "\n")

	// ETA calculation
	if m.FilesScanned > 0 && m.TotalFiles > 0 {
		avgPerFile := elapsed / time.Duration(m.FilesScanned)
		remaining := time.Duration(m.TotalFiles-m.FilesScanned) * avgPerFile
		sb.WriteString(MutedStyle.Render(fmt.Sprintf("  ETA: ~%s", remaining.Round(time.Second))) + "\n")
	}

	// Current file
	if m.CurrentFile != "" {
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(ColorPrimary).Render("  ⟳ Current: ") +
			lipgloss.NewStyle().Foreground(ColorText).Render(truncateFilePath(m.CurrentFile, 60)) + "\n")
	}

	// Recent files
	if len(m.RecentFiles) > 0 {
		sb.WriteString("\n" + MutedStyle.Render("  Recently scanned:") + "\n")
		for _, file := range m.RecentFiles {
			sb.WriteString(SuccessStyle.Render("    ✓ ") +
				MutedStyle.Render(truncateFilePath(file, 60)) + "\n")
		}
	}

	return sb.String()
}

// truncateFilePath truncates a file path to fit within maxLen
func truncateFilePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	// Show beginning and end
	half := (maxLen - 3) / 2
	return path[:half] + "..." + path[len(path)-half:]
}

// Made with Bob
