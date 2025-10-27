package tui

import "github.com/charmbracelet/lipgloss"

func (m *model) View() string {
	switch {
	case m.i != nil:
		title := renderedTitle(m)
		content := lipgloss.JoinVertical(lipgloss.Left, title, m.v.View(), m.vh.View(viewKeyMap{}))
		return viewStyle.Render(content)
	case m.f != nil:
		title := renderedTitle(m)
		status := renderedStatus(m)
		content := lipgloss.JoinVertical(lipgloss.Left, title, status, m.li.View())
		return listStyle.Render(content)
	default:
		tabs := renderedTabs(m)
		status := renderedStatus(m)
		content := lipgloss.JoinVertical(lipgloss.Left, tabs, status, m.lf.View())
		return listStyle.Render(content)
	}
}
