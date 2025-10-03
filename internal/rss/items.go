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
		return fmt.Sprintf("+ %s", title)
	}
	return title
}

func (i *RssItem) Content() string {
	title := i.Title()
	date := i.Item.PublishedParsed
	link := i.Link()
	desc := i.Description()
	content := clean(i.Item.Content)

	if content != "" {
		desc = ""
	}

	return fmt.Sprintf(
		"%s\n%s\n%s\n\n%s\n\n%s",
		title,
		date,
		link,
		desc,
		content,
	)
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
