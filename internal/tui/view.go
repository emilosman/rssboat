package tui

func (m *model) View() string {
	switch {
	case m.i != nil:
		return viewStyle.Render(m.v.View())
	case m.f != nil:
		return listStyle.Render(m.li.View())
	default:
		t := renderedTabs(m)
		return listStyle.Render(t, m.lf.View())
	}
}
