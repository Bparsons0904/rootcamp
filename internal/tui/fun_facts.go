package tui

import (
	"database/sql"
	"strings"

	"github.com/bobparsons/rootcamp/internal/lessons"
	"github.com/bobparsons/rootcamp/internal/types"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	stateList = iota
	stateDetail
)

type FunFactsModel struct {
	database      *sql.DB
	isOpen        bool
	width         int
	height        int
	state         int
	form          *huh.Form
	selectedFactID string
	allFacts      []types.FunFact
	renderedFacts map[string]string
	viewport      viewport.Model
}

func NewFunFactsModel(database *sql.DB) FunFactsModel {
	data, err := lessons.LoadFunFacts()
	allFacts := []types.FunFact{}
	renderedFacts := make(map[string]string)

	if err == nil {
		allFacts = data.Facts

		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(75),
		)

		if err == nil {
			for _, fact := range allFacts {
				rendered, err := renderer.Render(fact.Full)
				if err != nil {
					renderedFacts[fact.ID] = fact.Full
				} else {
					renderedFacts[fact.ID] = strings.TrimSpace(rendered)
				}
			}
		}
	}

	return FunFactsModel{
		database:      database,
		isOpen:        false,
		state:         stateList,
		allFacts:      allFacts,
		renderedFacts: renderedFacts,
	}
}

func (m FunFactsModel) Init() tea.Cmd {
	return nil
}

func (m *FunFactsModel) Update(msg tea.Msg) (*FunFactsModel, tea.Cmd) {
	if !m.isOpen {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == stateDetail {
			switch msg.String() {
			case "esc", "q":
				m.state = stateList
				m.createForm()
				return m, m.form.Init()
			default:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}

		if m.state == stateList {
			if msg.String() == "esc" || msg.String() == "q" {
				m.isOpen = false
				m.state = stateList
				m.selectedFactID = ""
				return m, nil
			}
		}
	}

	if m.state == stateList && m.form != nil {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		if m.form.State == huh.StateCompleted {
			if m.selectedFactID != "" {
				m.state = stateDetail
				m.setupDetailView()
			}
			return m, cmd
		}

		return m, cmd
	}

	return m, nil
}

func (m FunFactsModel) View() string {
	if !m.isOpen {
		return ""
	}

	if m.state == stateDetail {
		return m.renderDetailView()
	}

	return m.renderListView()
}

func (m *FunFactsModel) renderListView() string {
	if m.form == nil {
		return ""
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentBlue).
		Padding(1, 0).
		Align(lipgloss.Center).
		Render("📚 Fun Facts About Unix & The Terminal")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Render("Use arrow keys to navigate, Enter to view, ESC/Q to return to menu")

	formView := m.form.View()

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		formView,
		"",
		instructions,
	)

	bordered := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(AccentBlue).
		Padding(1, 2).
		Width(100).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		bordered,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}),
	)
}

func (m *FunFactsModel) renderDetailView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentBlue).
		Padding(1, 0).
		Render("📖 Fun Fact Detail")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("Use arrow keys to scroll, ESC/Q to return to list")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		m.viewport.View(),
		"",
		instructions,
	)

	bordered := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(AccentBlue).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		bordered,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}),
	)
}

func (m *FunFactsModel) setupDetailView() {
	if m.viewport.Width == 0 {
		m.viewport = viewport.New(80, 20)
		m.viewport.YPosition = 0
	}

	rendered, ok := m.renderedFacts[m.selectedFactID]
	if !ok {
		rendered = "Fact not found"
	}

	m.viewport.SetContent(rendered)
}

func (m *FunFactsModel) createForm() {
	m.selectedFactID = ""

	if len(m.allFacts) == 0 {
		return
	}

	options := make([]huh.Option[string], len(m.allFacts))
	for i, fact := range m.allFacts {
		options[i] = huh.NewOption(fact.Title, fact.ID)
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a fact to learn more:").
				Options(options...).
				Value(&m.selectedFactID).
				Height(15),
		),
	).WithWidth(90)
}

func (m *FunFactsModel) Open(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	m.isOpen = true
	m.state = stateList
	m.selectedFactID = ""

	m.createForm()
	return m.form.Init()
}

func (m *FunFactsModel) Close() {
	m.isOpen = false
	m.state = stateList
	m.selectedFactID = ""
}

func (m FunFactsModel) IsOpen() bool {
	return m.isOpen
}
