package rss

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"sync"

	yaml "github.com/goccy/go-yaml"
	"github.com/mmcdole/gofeed"
)

var ErrFeedHasNoUrl = errors.New("Feed has no URL")
var ErrNoFeedsInList = errors.New("No feeds in list")
var MsgFeedNotLoaded = "Feed not loaded yet. Press shift+r"

type RssFeed struct {
	Url      string
	Category string
	Error    string

	Feed     *gofeed.Feed
	RssItems []RssItem
}

type RssItem struct {
	Item *gofeed.Item
	Read bool
}

type FeedList struct {
	All []*RssFeed
}

func (i *RssItem) ToggleRead() {
	i.Read = !i.Read
}

func (f *RssFeed) GetFeed() error {
	if f.Url == "" {
		return ErrFeedHasNoUrl
	}

	fp := gofeed.NewParser()
	parsedFeed, err := fp.ParseURL(f.Url)
	if err != nil {
		f.Error = err.Error()
		return err
	}

	f.Feed = parsedFeed

	f.RssItems = make([]RssItem, len(parsedFeed.Items))
	for i, item := range parsedFeed.Items {
		f.RssItems[i] = RssItem{
			Item: item,
			Read: false,
		}
	}

	f.Error = ""

	return nil
}

func (l *FeedList) Add(feeds ...*RssFeed) {
	l.All = append(l.All, feeds...)
}

func (l *FeedList) UpdateAll() error {
	if len(l.All) == 0 {
		return ErrNoFeedsInList
	}

	var wg sync.WaitGroup
	wg.Add(len(l.All))

	for i := range l.All {
		go func(i int) {
			defer wg.Done()
			l.All[i].GetFeed()
		}(i)
	}

	wg.Wait()

	return nil
}

func CreateFeedsFromFS(filesystem fs.FS) ([]*RssFeed, error) {
	file, err := filesystem.Open("feeds.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	yamlData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var raw map[string][]string
	if err := yaml.Unmarshal(yamlData, &raw); err != nil {
		return nil, err
	}

	var feeds []*RssFeed
	for category, urls := range raw {
		for _, u := range urls {
			feeds = append(feeds, &RssFeed{
				Url:      u,
				Category: category,
			})
		}
	}

	return feeds, nil
}

func (f *RssFeed) HasUnread() bool {
	for _, item := range f.RssItems {
		if !item.Read {
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

func (l *FeedList) MarkAllFeedsRead() {
	for _, feed := range l.All {
		feed.MarkAllItemsRead()
	}
}

func (f *RssFeed) GetField(field string) string {
	switch field {
	case "Url":
		return f.Url

	case "Category":
		return f.Category

	case "Title":
		if f.Feed == nil || f.Feed.Title == "" {
			return f.Url
		}
		if f.HasUnread() {
			return fmt.Sprintf("🟢 %s", f.Feed.Title)
		}
		return f.Feed.Title

	case "Latest":
		switch {
		case f.Error != "":
			return f.Error
		case len(f.RssItems) > 0:
			last := f.RssItems[0]
			return last.Item.Title
		case f.Feed != nil:
			return f.Feed.Description
		default:
			return MsgFeedNotLoaded
		}

	default:
		return ""
	}
}
