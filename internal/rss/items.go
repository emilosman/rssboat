package rss

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

type RssItem struct {
	Item *gofeed.Item
	Read bool
}

func (i *RssItem) Title() string {
	title := Clean(i.Item.Title)
	if !i.Read {
		return fmt.Sprintf("🟢 %s", title)
	}
	return title
}

func (i *RssItem) Description() string {
	return Clean(i.Item.Description)
}

func (i *RssItem) ToggleRead() {
	i.Read = !i.Read
}

func (i *RssItem) MarkRead() {
	i.Read = true
}
