package rss

import (
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestItemLink(t *testing.T) {
	tests := []struct {
		name string
		item RssItem
		want string
	}{
		{
			name: "handles no gofeed item",
			item: RssItem{},
			want: "",
		},
		{
			name: "uses item link",
			item: RssItem{
				Item: &gofeed.Item{Link: "example.com/items/2"},
			},
			want: "example.com/items/2",
		},
		{
			name: "falls back to enclosure",
			item: RssItem{
				Item: &gofeed.Item{
					Enclosures: []*gofeed.Enclosure{{URL: "example.com/enclosure/2"}},
				},
			},
			want: "example.com/enclosure/2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.Link()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestItems(t *testing.T) {
	t.Run("Should get read status of unread feed item", func(t *testing.T) {
		unreadRssItem, _, _, _, _, _ := newTestData()

		want := "🟢 Latest item title"
		got := unreadRssItem.Title()
		if got != want {
			t.Errorf("Did not get correct field value, want %s, got %s", want, got)
		}
	})

	t.Run("Should get title of read feed item", func(t *testing.T) {
		_, rssItem, _, _, _, _ := newTestData()

		want := "Latest item title"
		got := rssItem.Title()
		if got != want {
			t.Errorf("Did not get correct field value, want %s, got %s", want, got)
		}
	})

	t.Run("Toggle feed item read flag", func(t *testing.T) {
		var feedItem RssItem

		feedItem.ToggleRead()

		if feedItem.Read != true {
			t.Error("Error toggling feed item read flag")
		}

		feedItem.ToggleRead()

		if feedItem.Read != false {
			t.Error("Error toggling feed item read flag")
		}
	})

	t.Run("Mark item read", func(t *testing.T) {
		var feedItem RssItem

		feedItem.MarkRead()

		if !feedItem.Read {
			t.Error("Item should be read")
		}
	})
}
