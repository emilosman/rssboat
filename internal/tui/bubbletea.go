package tui

import (
	"errors"
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

var columnNames = []string{"Category", "Title"}

func newModel(feedList *rss.FeedList) (model, error) {
	var m model

	columns, err := BuildColumns(columnNames)

	rows, err := buildRows(feedList.All, columnNames)
	if err != nil {
		return m, fmt.Errorf("Error building rows %q", err)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	m.table = t
	m.feedList = feedList

	return m, nil
}

func BuildColumns(columnNames []string) ([]table.Column, error) {
	var columns []table.Column

	if len(columnNames) == 0 {
		return columns, errors.New("No column names given")
	}

	for _, name := range columnNames {
		column := table.Column{Title: name, Width: 20}
		columns = append(columns, column)
	}

	return columns, nil
}

func buildRows(feeds []*rss.Feed, columnNames []string) ([]table.Row, error) {
	var rows []table.Row

	if len(feeds) == 0 {
		return rows, errors.New("No feeds given")
	}

	for _, f := range feeds {
		fields := f.GetFields(columnNames)
		row := table.Row(fields)
		rows = append(rows, row)
	}

	return rows, nil
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
			rows, err := buildRows(m.feedList.All, columnNames)
			if err != nil {
				fmt.Println("Error building rows", err)
			}

			m.table.SetRows(rows)
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

	model, err := newModel(&feedList)
	if err != nil {
		fmt.Printf("Error creating model %v", err)
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
