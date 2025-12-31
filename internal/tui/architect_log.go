package tui

import (
	"strings"

	"github.com/bobparsons/rootcamp/internal/lessons"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type ArchitectLogModel struct {
	selectedFact    string
	glamourRenderer *glamour.TermRenderer
	width           int
}

func NewArchitectLogModel(width int) ArchitectLogModel {
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(40),
	)

	fact, err := lessons.GetRandomFact()
	selectedFact := ""
	if err == nil && fact != nil {
		selectedFact = fact.Short
	}

	return ArchitectLogModel{
		selectedFact:    selectedFact,
		glamourRenderer: renderer,
		width:           width,
	}
}

func (m ArchitectLogModel) Init() tea.Cmd {
	return nil
}

func (m ArchitectLogModel) Update(msg tea.Msg) (ArchitectLogModel, tea.Cmd) {
	return m, nil
}

func (m ArchitectLogModel) View() string {
	title := PanelTitleStyle(ColorPurple).Render("THE ARCHITECT'S LOG")
	centeredTitle := lipgloss.Place(
		m.width-4,
		1,
		lipgloss.Center,
		lipgloss.Center,
		title,
	)

	rendered, err := m.glamourRenderer.Render(m.selectedFact)
	if err != nil {
		return centeredTitle + "\n\n" + m.selectedFact
	}

	return centeredTitle + "\n\n" + strings.TrimSpace(rendered)
}
