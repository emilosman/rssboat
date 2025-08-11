package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mmcdole/gofeed/rss"
)

var (
	baseStyle = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("240"))
)

type dataFetchedMsg struct {
}

func fetchDataCmd(address, startBlock, apiKey string) tea.Cmd {
	return nil
}

type model struct {
	SelectedView string
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	return model{}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "1" {
			m.SelectedView = "all"
		}
		if k == "2" {
			m.SelectedView = "feeds"
		}
		if k == "3" {
			m.SelectedView = "item"
		}
		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return updateTable(msg, m)
}

func updateTable(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m model) View() string {
	var b strings.Builder
	switch m.SelectedView {
	case "all":
		b.WriteString("all")
	case "feeds":
		b.WriteString("feeds")
	case "item":
		b.WriteString("all")
	}
	return baseStyle.Render(b.String())
}

func BuildColumns(columnNames []string) ([]table.Column, error) {
	return nil, nil
}

func BuildRows(feeds []rss.Feed) ([]table.Row, error) {
	return nil, nil
}

func BuildApp() {
}
