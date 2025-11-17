package rss

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

type RssItem struct {
	Item     *gofeed.Item
	Bookmark bool
	Read     bool
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
	title := i.Item.Title
	if i.Bookmark {
		title = fmt.Sprintf("* %s", title)
	}
	if !i.Read {
		title = fmt.Sprintf("+ %s", title)
	}
	return title
}

func (i *RssItem) FilterContent() string {
	return fmt.Sprintf("%s %s", i.Title(), i.Description())
}

func (i *RssItem) Content() string {
	title := i.Title()
	date := i.Item.PublishedParsed
	link := i.Link()
	desc := i.Description()
	content := i.Item.Content

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
	desc := i.Item.Description
	if desc == "" {
		desc = i.Item.Content
	}
	return desc
}

func (i *RssItem) ToggleRead() {
	i.Read = !i.Read
}

func (i *RssItem) ToggleBookmark() {
	i.Bookmark = !i.Bookmark
}

func (i *RssItem) MarkRead() {
	i.Read = true
}

func sanitizeItem(item *gofeed.Item) {
	item.Title = clean(item.Title)
	item.Description = clean(item.Description)
	item.Content = clean(item.Content)
}
