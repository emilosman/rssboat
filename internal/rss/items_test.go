package rss_test

import (
	"testing"

	"github.com/emilosman/rssboat/internal/rss"
	"github.com/mmcdole/gofeed"
)

func TestItemLink(t *testing.T) {
	tests := []struct {
		name string
		item rss.RssItem
		want string
	}{
		{
			name: "handles no gofeed item",
			item: rss.RssItem{},
			want: "",
		},
		{
			name: "uses item link",
			item: rss.RssItem{
				Item: &gofeed.Item{Link: "example.com/items/2"},
			},
			want: "example.com/items/2",
		},
		{
			name: "falls back to enclosure",
			item: rss.RssItem{
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
