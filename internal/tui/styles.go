package tui

import "github.com/charmbracelet/lipgloss"

var (
	viewStyle = lipgloss.NewStyle().Margin(4, 2, 1, 2)

	listStyle = lipgloss.NewStyle().Margin(1, 2)

	unreadStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00cf42"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e53636ff"))

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ff4fff")).
			Padding(0, 1)

	unreadTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00cf42")).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#767676")).
				Padding(0, 1)
)
