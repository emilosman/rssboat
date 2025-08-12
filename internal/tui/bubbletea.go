package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emilosman/rssboat/internal/rss"
)

type model struct {
	table    table.Model
	feedList *rss.FeedList
}

func newModel(feedList *rss.FeedList) model {
	columns := []table.Column{
		{Title: "Category", Width: 20},
		{Title: "URL", Width: 50},
	}

	rows := buildRows(feedList.All)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	return model{table: t, feedList: feedList}
}

func buildRows(feeds []*rss.Feed) []table.Row {
	var rows []table.Row
	for _, f := range feeds {
		rows = append(rows, table.Row{
			f.Category,
			f.Url,
		})
	}
	return rows
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		// Shift+R to refresh feeds
		if msg.String() == "R" {
			go func() {
				// fetch feeds
				if err := m.feedList.UpdateAll(); err != nil {
					fmt.Println("Update error:", err)
				}
			}()
			// after network fetch, update table rows
			m.table.SetRows(buildRows(m.feedList.All))
		}

		// Quit with q or ctrl+c
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.table.View()
}

func BuildApp() {
	filesystem := os.DirFS(".")
	feeds, err := rss.CreateFeedsFromFS(filesystem)
	if err != nil {
		log.Fatal(err)
	}

	var feedList rss.FeedList
	feedList.Add(feeds...)

	p := tea.NewProgram(newModel(&feedList))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
