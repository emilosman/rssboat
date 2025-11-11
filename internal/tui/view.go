package tui

import "github.com/charmbracelet/lipgloss"

func (m *model) View() string {
	switch {
	case m.i != nil:
		// Item view
		title := renderedTitle(m)
		status := renderedStatus(m)
		content := contentStyle.Render(m.v.View())
		help := helpStyle.Render(m.vh.View(viewKeyMap{}))
		view := lipgloss.JoinVertical(lipgloss.Left, title, status, content, help)
		return viewStyle.Render(view)
	case m.f != nil:
		// Feed view
		title := renderedTitle(m)
		status := renderedStatus(m)
		list := m.li.View()
		view := lipgloss.JoinVertical(lipgloss.Left, title, status, list)
		return listStyle.Render(view)
	default:
		// List view
		tabs := renderedTabs(m)
		status := renderedStatus(m)
		list := m.lf.View()
		view := lipgloss.JoinVertical(lipgloss.Left, tabs, status, list)
		return listStyle.Render(view)
	}
}
