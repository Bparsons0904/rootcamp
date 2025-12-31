package tui

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bobparsons/rootcamp/internal/db"
	"github.com/bobparsons/rootcamp/internal/lab"
	"github.com/bobparsons/rootcamp/internal/lessons"
	"github.com/bobparsons/rootcamp/internal/types"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	stateLessonList = iota
	stateLessonDetail
	stateCodeInput
	stateSuccess
)

type shellFinishedMsg struct{}

type LearnCommandModel struct {
	database         *sql.DB
	isOpen           bool
	width            int
	height           int
	state            int
	form             *huh.Form
	selectedLessonID string
	allLessons       []types.Lesson
	renderedAbout    map[string]string
	progressMap      map[string]*types.UserProgress
	viewport         viewport.Model
	sandboxPath      string
	codeInput        textinput.Model
	feedback         string
	currentLesson    *types.Lesson
}

func NewLearnCommandModel(database *sql.DB) LearnCommandModel {
	data, err := lessons.LoadLessons()
	allLessons := []types.Lesson{}
	renderedAbout := make(map[string]string)
	progressMap := make(map[string]*types.UserProgress)

	if err == nil {
		allLessons = data.Lessons

		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(75),
		)

		if err == nil {
			for _, lesson := range allLessons {
				content := formatLessonAbout(lesson)
				rendered, err := renderer.Render(content)
				if err != nil {
					renderedAbout[lesson.ID] = content
				} else {
					renderedAbout[lesson.ID] = strings.TrimSpace(rendered)
				}
			}
		}

		if database != nil {
			progressMap, _ = db.GetAllProgress(database)
		}
	}

	ti := textinput.New()
	ti.Placeholder = "Enter your answer here..."
	ti.CharLimit = 200
	ti.Width = 60

	return LearnCommandModel{
		database:      database,
		isOpen:        false,
		state:         stateLessonList,
		allLessons:    allLessons,
		renderedAbout: renderedAbout,
		progressMap:   progressMap,
		codeInput:     ti,
	}
}

func formatLessonAbout(lesson types.Lesson) string {
	var parts []string

	parts = append(parts, "# "+lesson.Title)
	parts = append(parts, "")

	if lesson.About.What != "" {
		parts = append(parts, "## What is "+lesson.Code+"?")
		parts = append(parts, "")
		parts = append(parts, lesson.About.What)
		parts = append(parts, "")
	}

	if lesson.About.Example != "" {
		parts = append(parts, "## Example")
		parts = append(parts, "")
		parts = append(parts, lesson.About.Example)
		parts = append(parts, "")
	}

	if lesson.About.History != "" {
		parts = append(parts, "## History")
		parts = append(parts, "")
		parts = append(parts, lesson.About.History)
		parts = append(parts, "")
	}

	if len(lesson.About.CommonUses) > 0 {
		parts = append(parts, "## Common Uses")
		parts = append(parts, "")
		for _, use := range lesson.About.CommonUses {
			parts = append(parts, "- "+use)
		}
		parts = append(parts, "")
	}

	if lesson.Instructions != "" {
		parts = append(parts, lesson.Instructions)
		parts = append(parts, "")
	}

	return strings.Join(parts, "\n")
}

func (m LearnCommandModel) Init() tea.Cmd {
	return nil
}

func (m *LearnCommandModel) Update(msg tea.Msg) (*LearnCommandModel, tea.Cmd) {
	if !m.isOpen {
		return m, nil
	}

	switch msg := msg.(type) {
	case shellFinishedMsg:
		m.state = stateCodeInput
		m.codeInput.Focus()
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case stateSuccess:
			if msg.String() == "enter" || msg.String() == " " {
				m.state = stateLessonList
				m.selectedLessonID = ""
				m.currentLesson = nil
				m.feedback = ""
				m.createForm()
				return m, m.form.Init()
			}

		case stateCodeInput:
			switch msg.String() {
			case "esc":
				m.state = stateLessonDetail
				m.codeInput.SetValue("")
				m.feedback = ""
				return m, nil
			case "enter":
				return m, m.validateAnswer()
			default:
				var cmd tea.Cmd
				m.codeInput, cmd = m.codeInput.Update(msg)
				return m, cmd
			}

		case stateLessonDetail:
			switch msg.String() {
			case "esc", "q":
				m.state = stateLessonList
				m.selectedLessonID = ""
				m.currentLesson = nil
				m.createForm()
				return m, m.form.Init()
			case "s":
				return m, m.startLab()
			case "c":
				m.state = stateCodeInput
				m.codeInput.Focus()
				return m, nil
			default:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}

		case stateLessonList:
			if msg.String() == "esc" || msg.String() == "q" {
				m.isOpen = false
				m.state = stateLessonList
				m.selectedLessonID = ""
				return m, nil
			}
		}
	}

	if m.state == stateLessonList && m.form != nil {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		if m.form.State == huh.StateCompleted {
			if m.selectedLessonID != "" {
				for i := range m.allLessons {
					if m.allLessons[i].ID == m.selectedLessonID {
						m.currentLesson = &m.allLessons[i]
						break
					}
				}
				if m.currentLesson != nil {
					m.state = stateLessonDetail
					m.setupDetailView()
				}
			}
			return m, cmd
		}

		return m, cmd
	}

	return m, nil
}

func (m *LearnCommandModel) validateAnswer() tea.Cmd {
	if m.currentLesson == nil {
		return nil
	}

	userInput := strings.TrimSpace(m.codeInput.Value())
	valid, errorMsg := lab.ValidateLesson(*m.currentLesson, userInput, m.sandboxPath)

	if valid {
		db.MarkComplete(m.database, m.currentLesson.ID)

		if m.sandboxPath != "" {
			lab.Cleanup(m.sandboxPath)
			m.sandboxPath = ""
		}

		progress, _ := db.GetProgress(m.database, m.currentLesson.ID)
		if progress != nil {
			m.progressMap[m.currentLesson.ID] = progress
		}

		m.state = stateSuccess
		m.feedback = "Congratulations! You've completed this lesson! 🎉"
	} else {
		db.IncrementAttempts(m.database, m.currentLesson.ID)

		progress, _ := db.GetProgress(m.database, m.currentLesson.ID)
		if progress != nil {
			m.progressMap[m.currentLesson.ID] = progress
		}

		if errorMsg != "" {
			m.feedback = fmt.Sprintf("❌ Incorrect. Hint: %s", errorMsg)
		} else {
			m.feedback = "❌ Incorrect answer. Try again!"
		}
	}

	m.codeInput.SetValue("")
	return nil
}

func (m *LearnCommandModel) startLab() tea.Cmd {
	if m.currentLesson == nil {
		return nil
	}

	sandboxPath, err := lab.Create(*m.currentLesson)
	if err != nil {
		m.feedback = fmt.Sprintf("Failed to create sandbox: %v", err)
		return nil
	}

	m.sandboxPath = sandboxPath
	startPath := lab.GetStartPath(sandboxPath, *m.currentLesson)

	instructions := fmt.Sprintf(`clear
cat << 'EOF'
╔══════════════════════════════════════════════════════════════════════════════╗
║                          ROOT CAMP - LAB SESSION                             ║
╚══════════════════════════════════════════════════════════════════════════════╝

Lesson: %s

%s

Your sandbox is located at:
  %s

When you're done, type 'exit' to return to Root Camp and enter your answer.

Good luck!

EOF
exec bash`, m.currentLesson.Title, m.currentLesson.Instructions, startPath)

	c := exec.Command("bash", "-c", instructions)
	c.Dir = startPath
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = append(os.Environ(), fmt.Sprintf("PS1=rootcamp:%s$ ", m.currentLesson.Code))

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return shellFinishedMsg{}
	})
}

func (m LearnCommandModel) View() string {
	if !m.isOpen {
		return ""
	}

	switch m.state {
	case stateSuccess:
		return m.renderSuccessView()
	case stateCodeInput:
		return m.renderCodeInputView()
	case stateLessonDetail:
		return m.renderDetailView()
	default:
		return m.renderListView()
	}
}

func (m *LearnCommandModel) renderListView() string {
	if m.form == nil {
		return ""
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentBlue).
		Padding(1, 0).
		Align(lipgloss.Center).
		Render("📚 Learn Command - Interactive Lessons")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Render("Use arrow keys to navigate, Enter to view lesson, ESC to return to menu")

	formView := m.form.View()

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		"",
		formView,
		"",
		instructions,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m *LearnCommandModel) renderDetailView() string {
	if m.currentLesson == nil {
		return ""
	}

	contentWidth := 90

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentBlue).
		Padding(1, 0).
		Width(contentWidth).
		Render(fmt.Sprintf("📖 Lesson: %s", m.currentLesson.Title))

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(contentWidth).
		Render("Arrow keys to scroll | [S] Start Lab | [C] Enter Code | ESC to return")

	feedbackStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true).
		Width(contentWidth)

	var feedbackView string
	if m.feedback != "" {
		feedbackView = feedbackStyle.Render(m.feedback)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		"",
		m.viewport.View(),
		"",
		feedbackView,
		"",
		instructions,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Top,
		content,
	)
}

func (m *LearnCommandModel) renderCodeInputView() string {
	if m.currentLesson == nil {
		return ""
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentBlue).
		Padding(1, 0).
		Width(m.width).
		Align(lipgloss.Center).
		Render("🔑 Enter Your Answer")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(m.width).
		Align(lipgloss.Center).
		Render("Type your answer and press Enter to submit | ESC to cancel")

	feedbackStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true).
		Width(m.width).
		Align(lipgloss.Center)

	var feedbackView string
	if m.feedback != "" {
		feedbackView = feedbackStyle.Render(m.feedback)
	}

	inputView := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.codeInput.View())

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		"",
		"",
		title,
		"",
		inputView,
		"",
		feedbackView,
		"",
		instructions,
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(content)
}

func (m *LearnCommandModel) renderSuccessView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorGreen).
		Padding(1, 0).
		Width(m.width).
		Align(lipgloss.Center).
		Render("🎉 Lesson Complete!")

	message := lipgloss.NewStyle().
		Foreground(ColorGreen).
		Bold(true).
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.feedback)

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(m.width).
		Align(lipgloss.Center).
		Render("Press Enter or Space to return to lesson list")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		"",
		"",
		title,
		"",
		message,
		"",
		instructions,
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(content)
}

func (m *LearnCommandModel) setupDetailView() {
	contentWidth := 90
	viewportHeight := m.height - 8
	if viewportHeight < 10 {
		viewportHeight = 10
	}

	m.viewport = viewport.New(contentWidth, viewportHeight)
	m.viewport.YPosition = 0

	rendered, ok := m.renderedAbout[m.selectedLessonID]
	if !ok {
		rendered = "Lesson content not found"
	}

	m.viewport.SetContent(rendered)
	m.feedback = ""
}

func (m *LearnCommandModel) createForm() {
	m.selectedLessonID = ""

	if len(m.allLessons) == 0 {
		return
	}

	options := make([]huh.Option[string], len(m.allLessons))
	for i, lesson := range m.allLessons {
		completionMark := " "
		if progress, ok := m.progressMap[lesson.ID]; ok && progress.Completed {
			completionMark = "✓"
		}
		label := fmt.Sprintf("[%s] %s (%s)", completionMark, lesson.Title, lesson.Level)
		options[i] = huh.NewOption(label, lesson.ID)
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a lesson to begin:").
				Options(options...).
				Value(&m.selectedLessonID).
				Height(15),
		),
	).WithWidth(90)
}

func (m *LearnCommandModel) Open(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	m.isOpen = true
	m.state = stateLessonList
	m.selectedLessonID = ""
	m.currentLesson = nil
	m.feedback = ""

	if m.database != nil {
		progressMap, _ := db.GetAllProgress(m.database)
		m.progressMap = progressMap
	}

	m.createForm()
	return m.form.Init()
}

func (m *LearnCommandModel) Close() {
	if m.sandboxPath != "" {
		lab.Cleanup(m.sandboxPath)
		m.sandboxPath = ""
	}
	m.isOpen = false
	m.state = stateLessonList
	m.selectedLessonID = ""
	m.currentLesson = nil
	m.feedback = ""
}

func (m LearnCommandModel) IsOpen() bool {
	return m.isOpen
}
