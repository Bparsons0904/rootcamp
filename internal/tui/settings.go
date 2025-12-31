package tui

import (
	"database/sql"

	"github.com/bobparsons/rootcamp/internal/db"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type SettingsModel struct {
	database *sql.DB
	form     *huh.Form
	isOpen   bool
	width    int
	height   int

	allSettings []settingOption
	// Current selected settings (bound to form)
	selectedSettings []string
}

type settingOption struct {
	key         string
	title       string
	description string
}

type settingsLoadedMsg struct {
	options  []settingOption
	selected []string
}

func NewSettingsModel(database *sql.DB) SettingsModel {
	return SettingsModel{
		database:         database,
		isOpen:           false,
		allSettings:      []settingOption{},
		selectedSettings: []string{},
	}
}

func (m *SettingsModel) createForm(options []settingOption, selected []string) {
	m.selectedSettings = make([]string, len(selected))
	copy(m.selectedSettings, selected)

	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt.title, opt.key)
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Key("settings").
				Title("Station Settings").
				Description("Select settings to enable (Space to toggle, Enter to save)").
				Options(huhOptions...).
				Value(&m.selectedSettings),
		),
	).WithWidth(70).WithTheme(huh.ThemeDracula())
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m *SettingsModel) Update(msg tea.Msg) (*SettingsModel, tea.Cmd) {
	if !m.isOpen {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle Esc to close without saving
		if msg.String() == "esc" {
			m.isOpen = false
			return m, nil
		}

	case settingsLoadedMsg:
		m.allSettings = msg.options
		m.createForm(msg.options, msg.selected)
		return m, m.form.Init()
	}

	if m.form != nil {
		var cmd tea.Cmd
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		if m.form.State == huh.StateCompleted {
			m.isOpen = false
			return m, m.saveSettings()
		}

		return m, cmd
	}

	return m, nil
}

func (m SettingsModel) View() string {
	if !m.isOpen || m.form == nil {
		return ""
	}

	formView := m.form.View()

	modal := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(AccentBlue).
		Padding(1, 2).
		Render(formView)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}),
	)
}

func (m *SettingsModel) Open(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	m.isOpen = true

	return func() tea.Msg {
		settings, err := db.GetAllSettings(m.database)
		if err != nil {
			return settingsLoadedMsg{options: []settingOption{}, selected: []string{}}
		}

		options := []settingOption{
			{
				key:         "skip_intro_animation",
				title:       "Skip Intro Animation",
				description: "Skip the boot sequence and provisioning animation on startup",
			},
		}

		selected := []string{}
		if settings.SkipIntroAnimation {
			selected = append(selected, "skip_intro_animation")
		}

		return settingsLoadedMsg{options: options, selected: selected}
	}
}

func (m *SettingsModel) Close() {
	m.isOpen = false
}

func (m SettingsModel) IsOpen() bool {
	return m.isOpen
}

func (m SettingsModel) saveSettings() tea.Cmd {
	return func() tea.Msg {
		selectedMap := make(map[string]bool)
		for _, key := range m.selectedSettings {
			selectedMap[key] = true
		}

		for _, option := range m.allSettings {
			isEnabled := selectedMap[option.key]
			value := "false"
			if isEnabled {
				value = "true"
			}
			_ = db.SetSetting(m.database, option.key, value)
		}
		return nil
	}
}
