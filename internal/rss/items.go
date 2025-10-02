package rss

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

type RssItem struct {
	Item *gofeed.Item
	Read bool
}

func (i *RssItem) Link() string {
	if i.Item == nil {
		return ""
	}
	if i.Item.Link != "" {
		return i.Item.Link
	}
	if len(i.Item.Enclosures) > 0 {
		return i.Item.Enclosures[0].URL
	}
	return ""
}

func (i *RssItem) Title() string {
	title := clean(i.Item.Title)
	if !i.Read {
		return fmt.Sprintf("🟢 %s", title)
	}
	return title
}

func (i *RssItem) Content() string {
	return fmt.Sprintf("%s\n%s\n\n%s\n\n%s", i.Title(), i.Link(), i.Description(), clean(i.Item.Content))
}

func (i *RssItem) Description() string {
	return clean(i.Item.Description)
}

func (i *RssItem) ToggleRead() {
	i.Read = !i.Read
}

func (i *RssItem) MarkRead() {
	i.Read = true
}
