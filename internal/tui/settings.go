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

	// Selected settings (enabled = true)
	selectedSettings []string
	// All available settings
	allSettings []settingOption
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
		selectedSettings: []string{},
		allSettings:      []settingOption{},
	}
}

func (m *SettingsModel) createForm(options []settingOption, selected []string) {
	// Build huh options
	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt.title, opt.key)
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
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

func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
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
		// Store settings and create form
		m.allSettings = msg.options
		m.selectedSettings = msg.selected
		m.createForm(msg.options, msg.selected)
		return m, m.form.Init()
	}

	// Pass message to form (if it exists)
	if m.form != nil {
		var cmd tea.Cmd
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		// Check if form completed
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

	// Render form
	formView := m.form.View()

	// Add custom border and styling
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

	// Load current settings from database
	return func() tea.Msg {
		settings, err := db.GetAllSettings(m.database)
		if err != nil {
			return settingsLoadedMsg{options: []settingOption{}, selected: []string{}}
		}

		// Define all available settings
		options := []settingOption{
			{
				key:         "skip_intro_animation",
				title:       "Skip Intro Animation",
				description: "Skip the boot sequence and provisioning animation on startup",
			},
			// Add more settings here as needed
			// {
			//     key:         "another_setting",
			//     title:       "Another Setting",
			//     description: "Description",
			// },
		}

		// Build list of currently selected (enabled) settings
		selected := []string{}
		if settings.SkipIntroAnimation {
			selected = append(selected, "skip_intro_animation")
		}
		// Add more settings checks here
		// if settings.AnotherSetting {
		//     selected = append(selected, "another_setting")
		// }

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
		// Create a map of selected settings for quick lookup
		selectedMap := make(map[string]bool)
		for _, key := range m.selectedSettings {
			selectedMap[key] = true
		}

		// Save all settings based on selection
		for _, option := range m.allSettings {
			isEnabled := selectedMap[option.key]
			err := db.SetSetting(m.database, option.key, db.BoolToString(isEnabled))
			if err != nil {
				// Handle error - for now just log it
				_ = err
			}
		}
		return nil
	}
}
