package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	viewStyle = lipgloss.NewStyle().Margin(10, 20)
	listStyle = lipgloss.NewStyle().Margin(1, 2)
)

func (m *model) View() string {
	switch {
	case m.i != nil:
		return viewStyle.Render(m.v.View())
	case m.selectedFeed != nil:
		return listStyle.Render(m.li.View())
	default:
		t := renderedTabs(m)
		return listStyle.Render(t, m.lf.View())
	}
}
