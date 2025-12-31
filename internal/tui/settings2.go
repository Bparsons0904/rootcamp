package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// The Architect's preferred palette: Tokyo Night inspired
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	accent    = lipgloss.AdaptiveColor{Light: "#00BBFF", Dark: "#7aa2f7"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#bb9af7"}

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(accent).
			Padding(1, 2).
			Background(lipgloss.Color("#1a1b26")) // Deep midnight background

	titleStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true).
			MarginBottom(1)
)

type Settings2Model struct {
	form     *huh.Form
	width    int
	height   int
	quitting bool
}

func NewSettings2Model() Settings2Model {
	return Settings2Model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("System Configuration").
					Description("Adjust the RootCamp environment variables. Precision is key."),

				huh.NewConfirm().
					Key("skip_intro").
					Title("Skip Intro Animation?").
					Description("For those who prefer the raw speed of the kernel over visual flair.").
					Affirmative("Skip It").
					Negative("Keep the Flair"),

				huh.NewSelect[string]().
					Key("theme").
					Options(huh.NewOptions("Tokyo Night", "Dracula", "Nord", "Cellar")...).
					Title("Interface Theme").
					Description("Choose your visual workspace."),

				huh.NewInput().
					Key("username").
					Title("Operator Handle").
					Placeholder("Architect...").
					Validate(func(str string) error {
						if len(str) < 3 {
							return fmt.Errorf("Too short. Give me a real name.")
						}
						return nil
					}),
			),
		).WithTheme(huh.ThemeBase()), // Using base theme
	}
}

func (m Settings2Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Settings2Model) Update(msg tea.Msg) (Settings2Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		m.quitting = true
	}

	return m, cmd
}

func (m Settings2Model) View() string {
	if m.quitting {
		return ""
	}

	// 1. Render the form
	formView := m.form.View()

	// 2. Wrap it in our Architect-approved modal frame
	modal := modalStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render(" 🖥️  STATION SETTINGS "),
			formView,
		),
	)

	// 3. Center the modal on the screen using Lip Gloss's Place function
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars("░"), // Background pattern for depth
		lipgloss.WithWhitespaceForeground(subtle),
	)
}
