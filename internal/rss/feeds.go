package rss

import (
	"fmt"
	"sort"
	"sync"

	"github.com/mmcdole/gofeed"
)

type RssFeed struct {
	Url      string
	Category string
	Error    string

	Feed     *gofeed.Feed
	RssItems []*RssItem
}

func (f *RssFeed) existingKeys() map[string]struct{} {
	existing := make(map[string]struct{}, len(f.RssItems))
	for _, item := range f.RssItems {
		if item.Item.GUID != "" {
			existing[item.Item.GUID] = struct{}{}
		} else if item.Item.Link != "" {
			existing[item.Item.Link] = struct{}{}
		}
	}
	return existing
}

func (f *RssFeed) SortByDate() {
	sort.Slice(f.RssItems, func(i, j int) bool {
		ti := f.RssItems[i].Item.PublishedParsed
		tj := f.RssItems[j].Item.PublishedParsed

		switch {
		case ti == nil && tj == nil:
			return false
		case ti == nil:
			return false
		case tj == nil:
			return true
		default:
			return ti.After(*tj)
		}
	})
}

func (f *RssFeed) HasUnread() bool {
	for i := range f.RssItems {
		if !f.RssItems[i].Read {
			return true
		}
	}
	return false
}

func (f *RssFeed) MarkAllItemsRead() {
	for i := range f.RssItems {
		f.RssItems[i].Read = true
	}
}

func (f *RssFeed) Title() string {
	var title string
	if f.Feed == nil || f.Feed.Title == "" {
		title = f.Url
	} else {
		title = clean(f.Feed.Title)
	}

	if f.HasUnread() {
		return fmt.Sprintf("+ %s", title)
	}

	return title
}

func (f *RssFeed) Description() string {
	return clean(f.Feed.Description)
}

func (f *RssFeed) Latest() string {
	switch {
	case f.Error != "":
		return f.Error
	case len(f.RssItems) > 0:
		last := f.RssItems[0]
		// reverse RssItems and return first unread item
		for i := len(f.RssItems) - 1; i >= 0; i-- {
			if !f.RssItems[i].Read {
				last = f.RssItems[i]
			}
		}
		return clean(last.Item.Title)
	case f.Feed != nil:
		return f.Description()
	default:
		return MsgFeedNotLoaded
	}
}

func (f *RssFeed) GetFeed() error {
	if f.Url == "" {
		return ErrFeedHasNoUrl
	}

	parsedFeed, err := gofeed.NewParser().ParseURL(f.Url)
	if err != nil {
		f.Error = err.Error()
		return err
	}

	f.Feed = parsedFeed
	f.mergeItems(parsedFeed.Items)
	f.SortByDate()
	f.Error = ""
	return nil
}

func (f *RssFeed) NextAfter(prev *RssItem) *RssItem {
	n := len(f.RssItems)
	if n == 0 {
		return nil
	}

	for i, item := range f.RssItems {
		if item == prev {
			if i < n-1 {
				return f.RssItems[i+1]
			}
			return nil
		}
	}
	return nil
}

func (f *RssFeed) mergeItems(items []*gofeed.Item) {
	existing := f.existingKeys()

	for _, item := range items {
		key := item.GUID
		if key == "" {
			key = item.Link
		}

		if _, ok := existing[key]; ok {
			continue
		}

		f.RssItems = append(f.RssItems, &RssItem{
			Item: item,
			Read: false,
		})
		existing[key] = struct{}{}
	}
}

func UpdateFeeds(feeds ...*RssFeed) (<-chan FeedResult, error) {
	if len(feeds) == 0 {
		return nil, ErrNoFeedsInList
	}

	results := make(chan FeedResult, len(feeds))
	var wg sync.WaitGroup
	wg.Add(len(feeds))

	for _, feed := range feeds {
		go func(f *RssFeed) {
			defer wg.Done()
			err := f.GetFeed()
			results <- FeedResult{Feed: f, Err: err}
		}(feed)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results, nil
}

func MarkFeedsAsRead(feeds ...*RssFeed) {
	for i := range feeds {
		feeds[i].MarkAllItemsRead()
	}
}
