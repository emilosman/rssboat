package tui

import "github.com/charmbracelet/bubbles/key"

func listShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("←/h"),
			key.WithHelp("←/h", "left"),
		),
		key.NewBinding(
			key.WithKeys("→/l"),
			key.WithHelp("→/l", "right"),
		),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view feed"),
		),
		key.NewBinding(
			key.WithKeys("shift+r"),
			key.WithHelp("shift+r", "refresh all feeds"),
		),
		key.NewBinding(
			key.WithKeys("q/esc"),
			key.WithHelp("q/esc", "quit"),
		),
	}
}

func listFullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			key.NewBinding(
				key.WithKeys("←/h"),
				key.WithHelp("←/h", "previous tab"),
			),
			key.NewBinding(
				key.WithKeys("→/l/tab"),
				key.WithHelp("→/l/tab", "next tab"),
			),
			key.NewBinding(
				key.WithKeys("o"),
				key.WithHelp("o", "open website"),
			),
			key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "refresh single feed"),
			),
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "view feed"),
			),
			key.NewBinding(
				key.WithKeys("shift+a"),
				key.WithHelp("shift+a", "mark feed as read"),
			),
			key.NewBinding(
				key.WithKeys("shift+c"),
				key.WithHelp("shift+c", "mark all items as read"),
			),
			key.NewBinding(
				key.WithKeys("shift+e"),
				key.WithHelp("shift+e", "edit urls file"),
			),
			key.NewBinding(
				key.WithKeys("shift+r"),
				key.WithHelp("shift+r", "refresh all feeds"),
			),
			key.NewBinding(
				key.WithKeys("ctrl+a"),
				key.WithHelp("ctrl+a", "mark tab as read"),
			),
			key.NewBinding(
				key.WithKeys("ctrl+r"),
				key.WithHelp("ctrl+r", "refresh tab"),
			),
			key.NewBinding(
				key.WithKeys("q/esc"),
				key.WithHelp("q/esc", "quit"),
			),
		},
	}

}

func itemsShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "toggle read"),
		),
		key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open item url"),
		),
		key.NewBinding(
			key.WithKeys("b/q/esc"),
			key.WithHelp("b/q/esc", "back"),
		),
		key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh feed"),
		),
		key.NewBinding(
			key.WithKeys("shift+a"),
			key.WithHelp("shift+a", "mark all items read"),
		),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "preview item"),
		),
	}
}

func itemsFullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "toggle read"),
			),
			key.NewBinding(
				key.WithKeys("n"),
				key.WithHelp("n", "next unread item"),
			),
			key.NewBinding(
				key.WithKeys("o"),
				key.WithHelp("o", "open item url"),
			),
			key.NewBinding(
				key.WithKeys("b/q/esc"),
				key.WithHelp("b/q/esc", "back"),
			),
			key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "refresh feed"),
			),
			key.NewBinding(
				key.WithKeys("shift+a"),
				key.WithHelp("shift+a", "mark all items read"),
			),
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "preview item"),
			),
		},
	}
}
