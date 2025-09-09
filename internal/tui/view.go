package tui

import "github.com/charmbracelet/lipgloss"

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

func (m *model) View() string {
	if m.selectedFeed != nil {
		return docStyle.Render(m.li.View())
	}
	t := renderedTabs(m)
	return docStyle.Render(t, m.lf.View())
}
