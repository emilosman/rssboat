package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mmcdole/gofeed/rss"
)

type dataFetchedMsg struct {
}

func fetchDataCmd(address, startBlock, apiKey string) tea.Cmd {
	return nil
}

type model struct {
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	return model{}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func updateTable(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m model) View() string {
	return ""
}

func BuildColumns(columnNames []string) ([]table.Column, error) {
	return nil, nil
}

func BuildRows(feeds []rss.Feed) ([]table.Row, error) {
	return nil, nil
}

func BuildApp() {
}
