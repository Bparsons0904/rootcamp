package tui

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/bobparsons/rootcamp/internal/db"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type provisionTickMsg time.Time

const (
	phaseBootSequence = iota
	phaseProvisioning
	phaseComplete
)

type WelcomeModel struct {
	phase               int
	progress            int
	bootScreen          BootScreenModel
	fileTree            FileTreeModel
	mainMenu            *MainMenuModel
	architectLog        ArchitectLogModel
	width               int
	height              int
	database            *sql.DB
	settingsModel       *SettingsModel
	guidedLearningModel *GuidedLearningModel
	learnCommandModel   *LearnCommandModel
	viewProgressModel   *ViewProgressModel
	funFactsModel       *FunFactsModel
	aboutModel          *AboutModel
	skippedAnimations   bool
}

func NewWelcomeModel(database *sql.DB) WelcomeModel {
	// Check if animations should be skipped
	skipAnimations := false
	if database != nil {
		settings, err := db.GetAllSettings(database)
		if err == nil && settings.SkipIntroAnimation {
			skipAnimations = true
		}
	}

	// Calculate panel widths
	totalWidth := 120
	leftWidth := 50
	rightWidth := 50
	middleWidth := totalWidth - leftWidth - rightWidth - 10

	// Set initial phase based on animation preference
	phase := phaseBootSequence
	progress := 0

	if skipAnimations {
		phase = phaseComplete
		progress = 100
	}

	// Initialize sub-models
	settingsModel := NewSettingsModel(database)
	guidedLearningModel := NewGuidedLearningModel(database)
	learnCommandModel := NewLearnCommandModel(database)
	viewProgressModel := NewViewProgressModel(database)
	funFactsModel := NewFunFactsModel(database)
	aboutModel := NewAboutModel(database)
	mainMenuModel := NewMainMenuModel(middleWidth)

	return WelcomeModel{
		phase:               phase,
		progress:            progress,
		bootScreen:          NewBootScreenModel(),
		fileTree:            NewFileTreeModel(leftWidth, skipAnimations),
		mainMenu:            &mainMenuModel,
		architectLog:        NewArchitectLogModel(rightWidth),
		width:               totalWidth,
		height:              40,
		database:            database,
		settingsModel:       &settingsModel,
		guidedLearningModel: &guidedLearningModel,
		learnCommandModel:   &learnCommandModel,
		viewProgressModel:   &viewProgressModel,
		funFactsModel:       &funFactsModel,
		aboutModel:          &aboutModel,
		skippedAnimations:   skipAnimations,
	}
}

func (m *WelcomeModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	// Always initialize menu
	cmds = append(cmds, m.mainMenu.Init())

	// Start animations only if not skipped
	if !m.skippedAnimations {
		cmds = append(cmds, m.bootScreen.Init(), tickForProvisionAnimation())
	}

	return tea.Batch(cmds...)
}

func tickForProvisionAnimation() tea.Cmd {
	return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
		return provisionTickMsg(t)
	})
}

func (m *WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check for open sub-models first
	if m.settingsModel.IsOpen() {
		var cmd tea.Cmd
		m.settingsModel, cmd = m.settingsModel.Update(msg)
		return m, cmd
	}

	if m.guidedLearningModel.IsOpen() {
		var cmd tea.Cmd
		m.guidedLearningModel, cmd = m.guidedLearningModel.Update(msg)
		return m, cmd
	}

	if m.learnCommandModel.IsOpen() {
		var cmd tea.Cmd
		m.learnCommandModel, cmd = m.learnCommandModel.Update(msg)
		return m, cmd
	}

	if m.viewProgressModel.IsOpen() {
		var cmd tea.Cmd
		m.viewProgressModel, cmd = m.viewProgressModel.Update(msg)
		return m, cmd
	}

	if m.funFactsModel.IsOpen() {
		var cmd tea.Cmd
		m.funFactsModel, cmd = m.funFactsModel.Update(msg)
		return m, cmd
	}

	if m.aboutModel.IsOpen() {
		var cmd tea.Cmd
		m.aboutModel, cmd = m.aboutModel.Update(msg)
		return m, cmd
	}

	// Handle global keys
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case bootCheckMsg:
		if m.phase == phaseBootSequence {
			var cmd tea.Cmd
			m.bootScreen, cmd = m.bootScreen.Update(msg)
			if m.bootScreen.IsComplete() {
				// Transition to provisioning
				m.phase = phaseProvisioning
				return m, tea.Batch(cmd, m.fileTree.Init(), tickForProvisionAnimation())
			}
			return m, cmd
		}

	case provisionTickMsg:
		if m.phase == phaseProvisioning {
			// Update progress
			if m.progress < 100 {
				m.progress += 1
				if m.progress >= 100 {
					m.progress = 100
					m.phase = phaseComplete
				}
			}

			return m, tickForProvisionAnimation()
		}

	case fileRevealMsg:
		if m.phase == phaseProvisioning {
			var cmd tea.Cmd
			m.fileTree, cmd = m.fileTree.Update(msg)
			return m, cmd
		}
	}

	// Update components based on phase
	if m.phase == phaseBootSequence {
		var cmd tea.Cmd
		m.bootScreen, cmd = m.bootScreen.Update(msg)
		return m, cmd
	}

	// Update all components in provisioning/complete phase
	var cmds []tea.Cmd

	fileTreeCmd := m.updateFileTree(msg)
	if fileTreeCmd != nil {
		cmds = append(cmds, fileTreeCmd)
	}

	menuCmd := m.updateMainMenu(msg)
	if menuCmd != nil {
		cmds = append(cmds, menuCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *WelcomeModel) updateFileTree(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.fileTree, cmd = m.fileTree.Update(msg)
	return cmd
}

func (m *WelcomeModel) updateMainMenu(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var updatedMenu *MainMenuModel
	updatedMenu, cmd = m.mainMenu.Update(msg)
	m.mainMenu = updatedMenu

	// Check if menu selection was made
	if m.mainMenu.IsCompleted() {
		selection := m.mainMenu.GetSelection()

		// Reset the form first
		resetCmd := m.mainMenu.Reset()

		// Only process if we have a valid selection
		if selection != "" {
			switch selection {
			case "guided_learning":
				return tea.Batch(resetCmd, m.guidedLearningModel.Open(m.width, m.height))
			case "learn_command":
				return tea.Batch(resetCmd, m.learnCommandModel.Open(m.width, m.height))
			case "view_progress":
				return tea.Batch(resetCmd, m.viewProgressModel.Open(m.width, m.height))
			case "fun_facts":
				return tea.Batch(resetCmd, m.funFactsModel.Open(m.width, m.height))
			case "about":
				return tea.Batch(resetCmd, m.aboutModel.Open(m.width, m.height))
			case "settings":
				return tea.Batch(resetCmd, m.settingsModel.Open(m.width, m.height))
			case "exit":
				return tea.Quit
			}
		}

		return resetCmd
	}

	return cmd
}

func (m WelcomeModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Show sub-model views if open
	if m.settingsModel.IsOpen() {
		return m.settingsModel.View()
	}
	if m.guidedLearningModel.IsOpen() {
		return m.guidedLearningModel.View()
	}
	if m.learnCommandModel.IsOpen() {
		return m.learnCommandModel.View()
	}
	if m.viewProgressModel.IsOpen() {
		return m.viewProgressModel.View()
	}
	if m.funFactsModel.IsOpen() {
		return m.funFactsModel.View()
	}
	if m.aboutModel.IsOpen() {
		return m.aboutModel.View()
	}

	// Render main view based on phase
	if m.phase == phaseBootSequence {
		return m.bootScreen.View()
	}

	return m.renderProvisioningView()
}

func (m WelcomeModel) renderProvisioningView() string {
	leftWidth := 50
	rightWidth := 50
	middleWidth := m.width - leftWidth - rightWidth - 10
	panelHeight := m.height - 10

	// Render each component
	left := PanelStyle(leftWidth, panelHeight, ColorBlue).Render(m.fileTree.View())
	middle := PanelStyle(middleWidth, panelHeight, ColorOrange).Render(m.mainMenu.View())
	right := PanelStyle(rightWidth, panelHeight, ColorPurple).Render(m.architectLog.View())

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, left, middle, right)

	header := HeaderStyle(m.width).Render("ROOT CAMP v0.1")
	footer := FooterStyle().Render("Use arrow keys to navigate, Enter to select")

	centeredFooter := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(footer)

	var content string
	if m.skippedAnimations {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			header,
			"",
			mainContent,
			"",
			centeredFooter,
		)
	} else {
		progressBar := m.renderProgressBar()
		centeredProgressBar := lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(progressBar)

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			header,
			"",
			mainContent,
			"",
			centeredProgressBar,
			"",
			centeredFooter,
		)
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m WelcomeModel) renderProgressBar() string {
	barWidth := 50
	filled := int(float64(barWidth) * float64(m.progress) / 100.0)
	empty := barWidth - filled

	bar := ProgressBarFilledStyle().Render(strings.Repeat("█", filled)) +
		ProgressBarEmptyStyle().Render(strings.Repeat("░", empty))

	label := ProgressLabelStyle().Render(fmt.Sprintf("Status: [Provisioning Sandbox...] %d%%", m.progress))

	return lipgloss.JoinVertical(
		lipgloss.Center,
		label,
		bar,
	)
}
