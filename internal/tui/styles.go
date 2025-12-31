package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorBlue        = lipgloss.Color("#5FB3FF")
	ColorOrange      = lipgloss.Color("#FFB86C")
	ColorPurple      = lipgloss.Color("#9D7CFF")
	ColorCyan        = lipgloss.Color("#7DCFFF")
	ColorGreen       = lipgloss.Color("#9ECE6A")
	ColorBrightGreen = lipgloss.Color("#00FF00")
	ColorDarkBg      = lipgloss.Color("#1A1B26")
	ColorGray        = lipgloss.Color("#565F89")
	ColorDarkGray    = lipgloss.Color("#414868")
	ColorLightGray   = lipgloss.Color("#999999")
	ColorWhite       = lipgloss.Color("#FFFFFF")

	// Tokyo Night inspired palette for modals
	AccentBlue      = lipgloss.AdaptiveColor{Light: "#00BBFF", Dark: "#7aa2f7"}
	AccentGreen     = lipgloss.AdaptiveColor{Light: "#00AA00", Dark: "#9ece6a"}
	AccentOrange    = lipgloss.AdaptiveColor{Light: "#FF8800", Dark: "#ff9e64"}
	AccentPurple    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#bb9af7"}
	HighlightPurple = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#bb9af7"}
	SubtleGray      = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	DeepMidnight    = lipgloss.Color("#1a1b26")
	TextPrimary     = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#c0caf5"}
	TextMuted       = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#565f89"}
)

func PanelStyle(width, height int, borderColor lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)
}

func PanelTitleStyle(bgColor lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorDarkBg).
		Background(bgColor).
		Bold(true).
		Padding(0, 1)
}

func HeaderStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorCyan).
		Bold(true).
		Width(width).
		Align(lipgloss.Center)
}

func FooterStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorGray).
		Italic(true)
}

func MenuOptionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorCyan)
}

func DisabledOptionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorGray).
		Italic(true)
}

func FileTreeStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorGreen)
}

func BootOKStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorBrightGreen).
		Bold(true)
}

func BootMessageStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorWhite)
}

func BootStartingStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorLightGray)
}

func ProgressBarFilledStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorGreen).
		Bold(true)
}

func ProgressBarEmptyStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorDarkGray)
}

func ProgressLabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(ColorGray).
		Italic(true)
}

func ModalContainerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.ThickBorder()).
		BorderForeground(AccentBlue).
		Background(DeepMidnight)
}

