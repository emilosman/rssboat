package tui

import "github.com/charmbracelet/lipgloss"

var (
	viewStyle = lipgloss.NewStyle().Margin(4, 0)

	listStyle = lipgloss.NewStyle().Margin(3, 0)

	unreadStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00cf42"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e53636ff"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff4fff")).
			Margin(0, 0, 0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#767676")).
			Margin(0, 0, 0, 1)

	contentStyle = lipgloss.NewStyle().
			Margin(0, 0, 0, 1)

	helpStyle = lipgloss.NewStyle().
			Margin(0, 0, 0, 2)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ff4fff")).
			Margin(0, 1, 0, 1)

	unreadTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00cf42")).
			Margin(0, 1, 0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#dfdfdf")).
				Margin(0, 1, 0, 1)
)
