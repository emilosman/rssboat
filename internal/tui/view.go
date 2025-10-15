package tui

import "github.com/charmbracelet/lipgloss"

func (m *model) View() string {
	switch {
	case m.i != nil:
		content := lipgloss.JoinVertical(lipgloss.Left, m.v.View(), m.vh.View(viewKeyMap{}))
		return viewStyle.Render(content)
	case m.f != nil:
		return listStyle.Render(m.li.View())
	default:
		t := renderedTabs(m)
		return listStyle.Render(t, m.lf.View())
	}
}
