package tui

import (
	"errors"

	"github.com/charmbracelet/bubbles/table"
	"github.com/emilosman/rssboat/internal/rss"
)

func buildColumns(columnNames []string) ([]table.Column, error) {
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
