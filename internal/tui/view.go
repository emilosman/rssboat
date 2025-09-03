package tui

import "github.com/charmbracelet/lipgloss"

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

func (m *model) View() string {
	if m.selectedFeed != nil {
		return docStyle.Render(m.li.View())
	}
	return docStyle.Render(m.lf.View())
}
