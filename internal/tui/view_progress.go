package tui

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bobparsons/rootcamp/internal/db"
	"github.com/bobparsons/rootcamp/internal/lessons"
	"github.com/bobparsons/rootcamp/internal/stats"
	"github.com/bobparsons/rootcamp/internal/types"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	progressViewWidth     = 90
	progressViewPadding   = 4
	progressViewportChrome = 8
	progressLabelWidth    = 15
	progressStatsWidth    = 2
)

type ViewProgressModel struct {
	database *sql.DB
	isOpen   bool
	width    int
	height   int
	viewport viewport.Model
	ready    bool
}

func NewViewProgressModel(database *sql.DB) ViewProgressModel {
	return ViewProgressModel{
		database: database,
		isOpen:   false,
	}
}

func (m ViewProgressModel) Init() tea.Cmd {
	return nil
}

func (m *ViewProgressModel) Update(msg tea.Msg) (*ViewProgressModel, tea.Cmd) {
	if !m.isOpen {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.isOpen = false
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ViewProgressModel) View() string {
	if !m.isOpen {
		return ""
	}

	if !m.ready {
		return "Loading progress..."
	}

	content := lipgloss.NewStyle().
		Width(progressViewWidth).
		PaddingTop(progressViewPadding).
		Render(m.viewport.View())

	footer := lipgloss.NewStyle().
		Foreground(TextMuted).
		Render("\nPress 'q' or 'esc' to return to menu")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Top,
		content+footer,
	)
}

func (m *ViewProgressModel) Open(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	m.isOpen = true

	lessonsData, err := lessons.LoadLessons()
	if err != nil {
		m.viewport.SetContent("Error loading lessons: " + err.Error())
		m.ready = true
		return nil
	}

	progressMap, err := db.GetAllProgress(m.database)
	if err != nil {
		progressMap = make(map[string]*types.UserProgress)
	}

	progress := stats.CalculateProgress(lessonsData.Lessons, progressMap)
	content := m.buildProgressView(progress)

	viewportHeight := height - progressViewportChrome
	m.viewport = viewport.New(progressViewWidth, viewportHeight)
	m.viewport.SetContent(content)
	m.ready = true

	return nil
}

func (m *ViewProgressModel) buildProgressView(progress stats.OverallProgress) string {
	var content strings.Builder

	header := m.renderSectionTitle("YOUR PROGRESS", AccentPurple)
	content.WriteString(header + "\n\n")

	overallTitle := m.renderSectionTitle("Overall Progress", TextPrimary)
	content.WriteString(overallTitle + "\n")
	content.WriteString(m.renderOverallProgress(progress.Overall) + "\n\n")

	content.WriteString(m.renderLevelProgress(progress.ByLevel) + "\n\n")
	content.WriteString(m.renderModuleProgress(progress.ByModule))

	return content.String()
}

func (m *ViewProgressModel) renderSectionTitle(title string, color lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(color).
		Render(title)
}

func (m *ViewProgressModel) renderProgressLine(label string, progStats stats.ProgressStats, color lipgloss.TerminalColor) string {
	statsLine := fmt.Sprintf("  %-*s %*d/%-*d  ",
		progressLabelWidth,
		label,
		progressStatsWidth,
		progStats.Completed,
		progressStatsWidth,
		progStats.Total)

	bar := lipgloss.NewStyle().
		Foreground(color).
		Render(stats.RenderProgressBar(progStats.Percentage))

	percentStr := fmt.Sprintf(" %.0f%%", progStats.Percentage)

	return statsLine + bar + percentStr
}

func (m *ViewProgressModel) renderOverallProgress(overall stats.ProgressStats) string {
	return m.renderProgressLine("Overall", overall, AccentGreen)
}

func (m *ViewProgressModel) renderLevelProgress(levels []stats.LevelStats) string {
	var content strings.Builder

	content.WriteString(m.renderSectionTitle("Progress by Level", TextPrimary) + "\n")

	colorMap := map[string]lipgloss.TerminalColor{
		"beginner":     AccentGreen,
		"intermediate": AccentOrange,
		"advanced":     AccentPurple,
		"expert":       AccentPurple,
	}

	for _, level := range levels {
		levelName := capitalizeFirst(level.Level)
		color := colorMap[level.Level]
		content.WriteString(m.renderProgressLine(levelName, level.Stats, color) + "\n")
	}

	return content.String()
}

func (m *ViewProgressModel) renderModuleProgress(modules []stats.ModuleStats) string {
	var content strings.Builder

	content.WriteString(m.renderSectionTitle("Progress by Module", TextPrimary) + "\n")

	if len(modules) == 0 {
		content.WriteString("  No modules available yet\n")
		return content.String()
	}

	for _, module := range modules {
		moduleName := formatModuleName(module.Module)
		content.WriteString(m.renderProgressLine(moduleName, module.Stats, AccentBlue) + "\n")
	}

	return content.String()
}

func (m *ViewProgressModel) Close() {
	m.isOpen = false
	m.ready = false
}

func (m ViewProgressModel) IsOpen() bool {
	return m.isOpen
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	if runes[0] >= 'a' && runes[0] <= 'z' {
		runes[0] = runes[0] - 32
	}
	return string(runes)
}

func formatModuleName(s string) string {
	result := ""
	capitalize := true
	for _, char := range s {
		if char == '-' {
			result += " "
			capitalize = true
		} else if capitalize {
			if char >= 'a' && char <= 'z' {
				result += string(char - 32)
			} else {
				result += string(char)
			}
			capitalize = false
		} else {
			result += string(char)
		}
	}
	return strings.TrimSpace(result)
}
