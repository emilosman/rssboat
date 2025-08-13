package tui

import (
	"testing"

	"github.com/emilosman/rssboat/internal/rss"
	"github.com/mmcdole/gofeed"
)

var (
	feed = rss.Feed{
		Url:      "example.com",
		Category: "Fun",
		Feed: &gofeed.Feed{
			Title: "Feed title",
		},
	}

	feedList = rss.FeedList{}
)

func TestTui(t *testing.T) {
	feedList.Add(&feed)

	t.Run("Build columns", func(t *testing.T) {
		columns, err := BuildColumns(columnNames)
		if err != nil {
			t.Errorf("Error building columns: %q", err)
		}

		if len(columnNames) != len(columns) {
			t.Error("Wrong number of columns returned")
		}
	})

	t.Run("Build rows", func(t *testing.T) {
		rows, err := buildRows(feedList.All, columnNames)
		if err != nil {
			t.Errorf("Error building rows: %q", err)
		}

		if len(rows) != len(feedList.All) {
			t.Errorf("Wrong number of rows returned, wanted %d, got %d", len(feedList.All), len(rows))
		}

		if len(rows[0]) != len(columnNames) {
			t.Errorf("Wrong number of data in rows, wanted %d, got %d", len(columnNames), len(rows[0]))
		}
	})
}
