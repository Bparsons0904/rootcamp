package tui

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

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
	stateGuidedCourseOverview = iota
	stateGuidedLessonDetail
	stateGuidedCodeInput
	stateGuidedSuccess
)

type guidedShellFinishedMsg struct{}

type GuidedLearningModel struct {
	database         *sql.DB
	isOpen           bool
	width            int
	height           int
	state            int
	form             *huh.Form
	selectedLessonID string
	courseLessons    []types.CourseLessonItem
	allLessons       []types.Lesson
	renderedAbout    map[string]string
	progressMap      map[string]*types.UserProgress
	viewport         viewport.Model
	sandboxPath      string
	codeInput        textinput.Model
	feedback         string
	currentLesson    *types.Lesson
	settings         *types.Settings
	generatedSecret  string
}

func NewGuidedLearningModel(database *sql.DB) GuidedLearningModel {
	data, err := lessons.LoadLessons()
	allLessons := []types.Lesson{}
	renderedAbout := make(map[string]string)
	progressMap := make(map[string]*types.UserProgress)
	courseLessons := []types.CourseLessonItem{}

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

		courseLessons, _ = lessons.GetCourseLessons(progressMap, allLessons)
	}

	ti := textinput.New()
	ti.Placeholder = "Enter your answer here..."
	ti.CharLimit = 200
	ti.Width = 60

	return GuidedLearningModel{
		database:      database,
		isOpen:        false,
		state:         stateGuidedCourseOverview,
		allLessons:    allLessons,
		renderedAbout: renderedAbout,
		progressMap:   progressMap,
		courseLessons: courseLessons,
		codeInput:     ti,
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
	case guidedShellFinishedMsg:
		m.state = stateGuidedCodeInput
		m.codeInput.Focus()
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case stateGuidedSuccess:
			if msg.String() == "enter" || msg.String() == " " {
				m.state = stateGuidedCourseOverview
				m.selectedLessonID = ""
				m.currentLesson = nil
				m.feedback = ""

				if m.database != nil {
					progressMap, _ := db.GetAllProgress(m.database)
					m.progressMap = progressMap
				}
				courseLessons, _ := lessons.GetCourseLessons(m.progressMap, m.allLessons)
				m.courseLessons = courseLessons

				m.createForm()
				if m.form != nil {
					return m, m.form.Init()
				}
				return m, nil
			}

		case stateGuidedCodeInput:
			switch msg.String() {
			case "esc", "q":
				m.state = stateGuidedLessonDetail
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

		case stateGuidedLessonDetail:
			switch msg.String() {
			case "esc", "q":
				m.state = stateGuidedCourseOverview
				m.selectedLessonID = ""
				m.currentLesson = nil
				m.createForm()
				return m, m.form.Init()
			case "s":
				if m.currentLesson != nil && m.currentLesson.SkipSandbox {
					m.state = stateGuidedCodeInput
					m.codeInput.Focus()
					return m, nil
				}
				return m, m.startLab()
			case "c":
				m.state = stateGuidedCodeInput
				m.codeInput.Focus()
				return m, nil
			default:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}

		case stateGuidedCourseOverview:
			if msg.String() == "esc" || msg.String() == "q" {
				m.isOpen = false
				m.state = stateGuidedCourseOverview
				m.selectedLessonID = ""
				return m, nil
			}
		}
	}

	if m.state == stateGuidedCourseOverview && m.form != nil {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		if m.form.State == huh.StateCompleted {
			if m.selectedLessonID != "" {
				for i := range m.courseLessons {
					if m.courseLessons[i].Lesson.ID == m.selectedLessonID {
						if m.courseLessons[i].Status != types.LessonLocked {
							m.currentLesson = &m.courseLessons[i].Lesson
							break
						}
					}
				}
				if m.currentLesson != nil {
					m.state = stateGuidedLessonDetail
					m.setupDetailView()
				} else {
					m.createForm()
					return m, m.form.Init()
				}
			}
			return m, cmd
		}

		return m, cmd
	}

	return m, nil
}

func (m *GuidedLearningModel) validateAnswer() tea.Cmd {
	if m.currentLesson == nil {
		return nil
	}

	userInput := strings.TrimSpace(m.codeInput.Value())

	lesson := *m.currentLesson
	if lesson.SkipSandbox && m.generatedSecret != "" {
		for i := range lesson.Requirements {
			lesson.Requirements[i].Expected = strings.ReplaceAll(
				lesson.Requirements[i].Expected,
				"{SECRET_CODE}",
				m.generatedSecret,
			)
		}
	}

	valid, errorMsg := lab.ValidateLesson(lesson, userInput, m.sandboxPath)

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

		m.state = stateGuidedSuccess
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

func (m *GuidedLearningModel) startLab() tea.Cmd {
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

	useBasicBash := false
	if m.settings != nil {
		useBasicBash = m.settings.UseBasicBash
	}

	var instructions string
	if useBasicBash {
		instructions = fmt.Sprintf(`clear
printf '\033[1;96m╔══════════════════════════════════════════════════════════════════════════════╗\n'
printf '║\033[0m\033[1;92m                          ROOT CAMP - LAB SESSION                             \033[1;96m║\n'
printf '╚══════════════════════════════════════════════════════════════════════════════╝\033[0m\n\n'
printf '\033[1;33mLesson:\033[0m \033[1;97m%s\033[0m\n\n'
cat << 'EOF'
%s
EOF
printf '\n\033[1;36mYour sandbox is located at:\033[0m\n'
printf '  \033[36m%s\033[0m\n\n'
printf '\033[35mWhen you'"'"'re done, type \033[1;91mexit\033[0m\033[35m to return to Root Camp and enter your answer.\033[0m\n\n'
printf '\033[1;32mGood luck!\033[0m\n\n'
cd '%s'
exec bash`, m.currentLesson.Title, m.currentLesson.Instructions, startPath, startPath)
	} else {
		userShell := os.Getenv("SHELL")
		if userShell == "" {
			userShell = "/bin/bash"
		}

		instructions = fmt.Sprintf(`clear
printf '\033[1;96m╔══════════════════════════════════════════════════════════════════════════════╗\n'
printf '║\033[0m\033[1;92m                          ROOT CAMP - LAB SESSION                             \033[1;96m║\n'
printf '╚══════════════════════════════════════════════════════════════════════════════╝\033[0m\n\n'
printf '\033[1;33mLesson:\033[0m \033[1;97m%s\033[0m\n\n'
cat << 'EOF'
%s
EOF
printf '\n\033[1;36mYour sandbox is located at:\033[0m\n'
printf '  \033[36m%s\033[0m\n\n'
printf '\033[35mWhen you'"'"'re done, type \033[1;91mexit\033[0m\033[35m to return to Root Camp and enter your answer.\033[0m\n\n'
printf '\033[1;32mGood luck!\033[0m\n\n'
cd '%s'
exec %s`, m.currentLesson.Title, m.currentLesson.Instructions, startPath, startPath, userShell)
	}

	c := exec.Command("bash", "-c", instructions)
	c.Dir = startPath
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = append(os.Environ(), fmt.Sprintf("PS1=rootcamp:%s$ ", m.currentLesson.Code))

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return guidedShellFinishedMsg{}
	})
}

func (m GuidedLearningModel) View() string {
	if !m.isOpen {
		return ""
	}

	switch m.state {
	case stateGuidedSuccess:
		return m.renderSuccessView()
	case stateGuidedCodeInput:
		return m.renderCodeInputView()
	case stateGuidedLessonDetail:
		return m.renderDetailView()
	default:
		return m.renderOverviewView()
	}
}

func (m *GuidedLearningModel) renderOverviewView() string {
	if m.form == nil {
		errorMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Align(lipgloss.Center).
			Render("No course lessons available.\n\nPlease check that lessons are properly configured.")

		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			errorMsg,
		)
	}

	courseData, err := lessons.LoadCourse()
	if err != nil {
		errorMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("Failed to load course: %v", err))

		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			errorMsg,
		)
	}

	completed, total := lessons.GetCourseProgress(m.courseLessons)
	progressPercent := 0
	if total > 0 {
		progressPercent = (completed * 100) / total
	}

	progressBarWidth := 40
	filledBlocks := (progressPercent * progressBarWidth) / 100
	emptyBlocks := progressBarWidth - filledBlocks

	progressBar := strings.Repeat("█", filledBlocks) + strings.Repeat("░", emptyBlocks)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentBlue).
		Padding(1, 0).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("🎯 Guided Learning - %s", courseData.Course.Title))

	description := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Render(courseData.Course.Description)

	progressInfo := lipgloss.NewStyle().
		Foreground(ColorGreen).
		Bold(true).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("Progress: %d/%d lessons completed", completed, total))

	progressBarView := lipgloss.NewStyle().
		Foreground(ColorGreen).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("[%s] %d%%", progressBar, progressPercent))

	formView := m.form.View()

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Render("↑/↓: Navigate | Enter: Start Lesson | ESC/Q: Return to Menu")

	legend := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Render("[✓] Completed  [→] Next  [L] Locked")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		"",
		description,
		"",
		progressInfo,
		progressBarView,
		"",
		formView,
		"",
		legend,
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

func (m *GuidedLearningModel) renderDetailView() string {
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
		Render("Arrow keys to scroll | [S] Start Lab | [C] Enter Code | ESC/Q to return")

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

func (m *GuidedLearningModel) renderCodeInputView() string {
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
		Render("Type your answer and press Enter to submit | ESC/Q to cancel")

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

func (m *GuidedLearningModel) renderSuccessView() string {
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
		Render("Press Enter or Space to return to course overview")

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

func (m *GuidedLearningModel) setupDetailView() {
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

	if m.currentLesson != nil && m.currentLesson.SkipSandbox {
		m.generatedSecret = generateGuidedSecretCode()
		rendered = replaceGuidedPlaceholders(rendered, m.generatedSecret)
	}

	m.viewport.SetContent(rendered)
	m.feedback = ""
}

func (m *GuidedLearningModel) createForm() {
	m.selectedLessonID = ""

	if len(m.courseLessons) == 0 {
		return
	}

	options := make([]huh.Option[string], len(m.courseLessons))
	nextFound := false

	for i, item := range m.courseLessons {
		var statusMark string
		var label string

		switch item.Status {
		case types.LessonComplete:
			statusMark = "✓"
			label = fmt.Sprintf("[%s] %d. %s - %s", statusMark, item.Sequence, item.Lesson.Code, item.Lesson.Title)
		case types.LessonUnlocked:
			if !nextFound {
				statusMark = "→"
				nextFound = true
			} else {
				statusMark = " "
			}
			label = fmt.Sprintf("[%s] %d. %s - %s", statusMark, item.Sequence, item.Lesson.Code, item.Lesson.Title)
		case types.LessonLocked:
			statusMark = "L"
			label = fmt.Sprintf("[%s] %d. %s - %s", statusMark, item.Sequence, item.Lesson.Code, item.Lesson.Title)
		}

		options[i] = huh.NewOption(label, item.Lesson.ID)
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a lesson:").
				Options(options...).
				Value(&m.selectedLessonID).
				Height(15),
		),
	).WithWidth(90)
}

func generateGuidedSecretCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	const length = 12
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	result := string(b)
	return result[:4] + "-" + result[4:8] + "-" + result[8:]
}

func getGuidedOSShortcuts() (copy, paste string) {
	if runtime.GOOS == "darwin" {
		return "Cmd + C", "Cmd + V"
	}
	return "Ctrl + Shift + C", "Ctrl + Shift + V"
}

func replaceGuidedPlaceholders(text, secretCode string) string {
	copyShortcut, pasteShortcut := getGuidedOSShortcuts()

	text = strings.ReplaceAll(text, "{SECRET_CODE}", secretCode)
	text = strings.ReplaceAll(text, "{COPY_SHORTCUT}", copyShortcut)
	text = strings.ReplaceAll(text, "{PASTE_SHORTCUT}", pasteShortcut)

	return text
}

func (m *GuidedLearningModel) Open(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	m.isOpen = true
	m.state = stateGuidedCourseOverview
	m.selectedLessonID = ""
	m.currentLesson = nil
	m.feedback = ""

	if m.database != nil {
		progressMap, _ := db.GetAllProgress(m.database)
		m.progressMap = progressMap

		settings, _ := db.GetAllSettings(m.database)
		m.settings = settings
	}

	courseLessons, _ := lessons.GetCourseLessons(m.progressMap, m.allLessons)
	m.courseLessons = courseLessons

	m.createForm()
	if m.form != nil {
		return m.form.Init()
	}
	return nil
}

func (m *GuidedLearningModel) Close() {
	if m.sandboxPath != "" {
		lab.Cleanup(m.sandboxPath)
		m.sandboxPath = ""
	}
	m.isOpen = false
	m.state = stateGuidedCourseOverview
	m.selectedLessonID = ""
	m.currentLesson = nil
	m.feedback = ""
}

func (m GuidedLearningModel) IsOpen() bool {
	return m.isOpen
}
