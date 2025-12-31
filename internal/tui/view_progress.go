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
		Width(90).
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

	viewportHeight := height - 8
	m.viewport = viewport.New(90, viewportHeight)
	m.viewport.SetContent(content)
	m.ready = true

	return nil
}

func (m *ViewProgressModel) buildProgressView(progress stats.OverallProgress) string {
	var content string

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentPurple).
		Render("YOUR PROGRESS")
	content += header + "\n\n"

	content += m.renderOverallProgress(progress.Overall)
	content += "\n\n"

	content += m.renderLevelProgress(progress.ByLevel)
	content += "\n\n"

	content += m.renderModuleProgress(progress.ByModule)

	return content
}

func (m *ViewProgressModel) renderOverallProgress(overall stats.ProgressStats) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(TextPrimary).
		Render("Overall Progress")

	statsLine := fmt.Sprintf("%d/%d lessons completed (%.0f%%)",
		overall.Completed, overall.Total, overall.Percentage)

	bar := stats.RenderProgressBar(overall.Percentage)
	barStyled := lipgloss.NewStyle().
		Foreground(AccentGreen).
		Render(bar)

	return fmt.Sprintf("%s\n%s\n%s", title, statsLine, barStyled)
}

func (m *ViewProgressModel) renderLevelProgress(levels []stats.LevelStats) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(TextPrimary).
		Render("Progress by Level")

	content := title + "\n"

	colorMap := map[string]lipgloss.TerminalColor{
		"beginner":     AccentGreen,
		"intermediate": AccentOrange,
		"advanced":     AccentPurple,
		"expert":       AccentPurple,
	}

	for _, level := range levels {
		levelName := capitalizeFirst(level.Level)
		statsLine := fmt.Sprintf("  %-15s %2d/%-2d  ",
			levelName,
			level.Stats.Completed,
			level.Stats.Total)

		bar := stats.RenderProgressBar(level.Stats.Percentage)
		color := colorMap[level.Level]
		barStyled := lipgloss.NewStyle().Foreground(color).Render(bar)

		percentStr := fmt.Sprintf(" %.0f%%", level.Stats.Percentage)

		content += statsLine + barStyled + percentStr + "\n"
	}

	return content
}

func (m *ViewProgressModel) renderModuleProgress(modules []stats.ModuleStats) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(TextPrimary).
		Render("Progress by Module")

	content := title + "\n"

	if len(modules) == 0 {
		content += "  No modules available yet\n"
		return content
	}

	for _, module := range modules {
		moduleName := formatModuleName(module.Module)
		statsLine := fmt.Sprintf("  %-20s %2d/%-2d  ",
			moduleName,
			module.Stats.Completed,
			module.Stats.Total)

		bar := stats.RenderProgressBar(module.Stats.Percentage)
		barStyled := lipgloss.NewStyle().Foreground(AccentBlue).Render(bar)

		percentStr := fmt.Sprintf(" %.0f%%", module.Stats.Percentage)

		content += statsLine + barStyled + percentStr + "\n"
	}

	return content
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
