package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type MainMenuModel struct {
	form             *huh.Form
	selectedMenuItem *string
	width            int
}

func NewMainMenuModel(width int) MainMenuModel {
	selection := ""
	m := MainMenuModel{
		width:            width,
		selectedMenuItem: &selection,
	}
	m.createForm()
	return m
}

func (m *MainMenuModel) createForm() {
	selection := ""
	m.selectedMenuItem = &selection

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Main Menu").
				Description("Select an option").
				Options(
					huh.NewOption("Guided Learning", "guided_learning"),
					huh.NewOption("Learn Command", "learn_command"),
					huh.NewOption("View Progress", "view_progress"),
					huh.NewOption("Fun Facts", "fun_facts"),
					huh.NewOption("About Root Camp", "about"),
					huh.NewOption("Settings", "settings"),
					huh.NewOption("Exit", "exit"),
				).
				Value(m.selectedMenuItem),
		),
	).WithWidth(60).WithTheme(huh.ThemeDracula())
}

func (m *MainMenuModel) Init() tea.Cmd {
	if m.form != nil {
		return m.form.Init()
	}
	return nil
}

func (m *MainMenuModel) Update(msg tea.Msg) (*MainMenuModel, tea.Cmd) {
	if m.form == nil {
		return m, nil
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	return m, cmd
}

func (m *MainMenuModel) View() string {
	title := PanelTitleStyle(ColorOrange).Render("MAIN MENU")
	centeredTitle := lipgloss.Place(
		m.width-4,
		1,
		lipgloss.Center,
		lipgloss.Center,
		title,
	)

	formView := m.form.View()
	content := lipgloss.Place(
		m.width-4,
		20,
		lipgloss.Center,
		lipgloss.Top,
		formView,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		centeredTitle,
		"",
		content,
	)
}

func (m *MainMenuModel) GetSelection() string {
	if m.selectedMenuItem != nil {
		return *m.selectedMenuItem
	}
	return ""
}

func (m *MainMenuModel) IsCompleted() bool {
	return m.form != nil && m.form.State == huh.StateCompleted
}

func (m *MainMenuModel) Reset() tea.Cmd {
	m.createForm()
	if m.form != nil {
		return m.form.Init()
	}
	return nil
}
