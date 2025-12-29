package tui

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/Bparsons0904/rootcamp/internal/db"
	"github.com/Bparsons0904/rootcamp/internal/lab"
	"github.com/Bparsons0904/rootcamp/internal/types"
)

type View int

const (
	DashboardView View = iota
	LessonView
)

type shellExitedMsg struct{}

type Model struct {
	db              *sql.DB
	lessons         []types.Lesson
	progress        map[string]*types.UserProgress
	view            View
	selected        int
	currentLesson   *types.Lesson
	activeLab       string
	labStarted      bool
	codeInputActive bool
	input           textinput.Model
	feedback        string
	feedbackSuccess bool
	width           int
	height          int
}

func NewModel(database *sql.DB, lessons []types.Lesson, progress map[string]*types.UserProgress) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter code here"
	ti.CharLimit = 50
	ti.Width = 30

	return Model{
		db:              database,
		lessons:         lessons,
		progress:        progress,
		view:            DashboardView,
		selected:        0,
		currentLesson:   nil,
		activeLab:       "",
		labStarted:      false,
		codeInputActive: false,
		input:           ti,
		feedback:        "",
		feedbackSuccess: false,
		width:           80,
		height:          24,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case shellExitedMsg:
		m.labStarted = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.activeLab != "" {
				lab.Cleanup(m.activeLab)
			}
			return m, tea.Quit

		case "ctrl+d":
			if m.view == LessonView && !m.codeInputActive {
				m.codeInputActive = true
				m.input.Focus()
				return m, nil
			}

		case "esc":
			if m.view == LessonView {
				if m.codeInputActive {
					m.codeInputActive = false
					m.input.Blur()
					m.input.SetValue("")
					return m, nil
				}

				if m.activeLab != "" {
					lab.Cleanup(m.activeLab)
					m.activeLab = ""
				}
				m.view = DashboardView
				m.currentLesson = nil
				m.feedback = ""
				m.input.SetValue("")
				m.labStarted = false
				m.codeInputActive = false
				return m, nil
			}

		case "up", "k":
			if m.view == DashboardView && m.selected > 0 {
				m.selected--
			}

		case "down", "j":
			if m.view == DashboardView && m.selected < len(m.lessons)-1 {
				m.selected++
			}

		case "enter":
			if m.view == DashboardView {
				m.currentLesson = &m.lessons[m.selected]
				m.view = LessonView
				m.feedback = ""
				m.input.SetValue("")
				m.labStarted = false
				m.codeInputActive = false
				return m, nil

			} else if m.view == LessonView && m.codeInputActive {
				userInput := m.input.Value()
				if lab.Validate(userInput, m.currentLesson.SecretCode) {
					m.feedback = "Correct! Lesson completed!"
					m.feedbackSuccess = true

					if err := db.MarkComplete(m.db, m.currentLesson.ID); err == nil {
						progress, _ := db.GetProgress(m.db, m.currentLesson.ID)
						m.progress[m.currentLesson.ID] = progress
					}

					if m.activeLab != "" {
						lab.Cleanup(m.activeLab)
						m.activeLab = ""
					}
				} else {
					m.feedback = "Incorrect code. Try again!"
					m.feedbackSuccess = false
					db.IncrementAttempts(m.db, m.currentLesson.ID)
				}
				return m, nil

			} else if m.view == LessonView && !m.codeInputActive && m.activeLab == "" {
				sandboxPath, err := lab.Create(*m.currentLesson)
				if err != nil {
					m.feedback = "Failed to create lab environment: " + err.Error()
					m.feedbackSuccess = false
					return m, nil
				}
				m.activeLab = sandboxPath

				// Create init script with welcome message
				initScript := fmt.Sprintf(`clear
echo "════════════════════════════════════════════════════════════════"
echo "  RootCamp Lab: %s"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "  Task: %s"
echo ""
echo "  Type 'exit' when done to return to RootCamp"
echo "════════════════════════════════════════════════════════════════"
echo ""
`, m.currentLesson.ID, m.currentLesson.Lab.Task)

				initPath := filepath.Join(sandboxPath, ".init.sh")
				os.WriteFile(initPath, []byte(initScript), 0644)

				// Use clean bash with init file for welcome message
				cmd := exec.Command("/bin/bash", "--init-file", initPath)
				cmd.Dir = sandboxPath

				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					return shellExitedMsg{}
				})
			}
		}
	}

	if m.view == LessonView && m.codeInputActive {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	switch m.view {
	case DashboardView:
		return m.renderDashboard()
	case LessonView:
		return m.renderLesson()
	default:
		return "Unknown view"
	}
}
