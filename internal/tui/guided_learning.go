package tui

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GuidedLearningModel struct {
	database *sql.DB
	isOpen   bool
	width    int
	height   int
}

func NewGuidedLearningModel(database *sql.DB) GuidedLearningModel {
	return GuidedLearningModel{
		database: database,
		isOpen:   false,
	}
}

func (m GuidedLearningModel) Init() tea.Cmd {
	return nil
}

func (m *GuidedLearningModel) Update(msg tea.Msg) (*GuidedLearningModel, tea.Cmd) {
	if !m.isOpen {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			m.isOpen = false
			return m, nil
		}
	}

	return m, nil
}

func (m GuidedLearningModel) View() string {
	if !m.isOpen {
		return ""
	}

	content := lipgloss.NewStyle().
		Width(80).
		Height(20).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(AccentBlue).
		Padding(2, 4).
		Render("Guided Learning View\n\n(Coming Soon)\n\nPress 'q' or 'esc' to return to menu")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}),
	)
}

func (m *GuidedLearningModel) Open(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	m.isOpen = true
	return nil
}

func (m *GuidedLearningModel) Close() {
	m.isOpen = false
}

func (m GuidedLearningModel) IsOpen() bool {
	return m.isOpen
}
