package rss

import (
	"errors"
	"io"
	"io/fs"
	"sync"

	yaml "github.com/goccy/go-yaml"
	"github.com/mmcdole/gofeed"
)

var ErrFeedHasNoUrl = errors.New("Feed has no URL")
var ErrNoFeedsInList = errors.New("No feeds in list")

type Feed struct {
	*gofeed.Feed
	Url      string
	Category string
	Error    string
	Items    []RssItem
}

type RssItem struct {
	*gofeed.Item
	Read bool
}

type FeedList struct {
	All []*Feed
}

func (i *RssItem) ToggleRead() {
	i.Read = !i.Read
}

func (f *Feed) GetFeed() error {
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

	f.Error = ""

	return nil
}

func (l *FeedList) Add(feeds ...*Feed) {
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

func CreateFeedsFromFS(filesystem fs.FS) ([]*Feed, error) {
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

	var feeds []*Feed
	for category, urls := range raw {
		for _, u := range urls {
			feeds = append(feeds, &Feed{
				Url:      u,
				Category: category,
			})
		}
	}

	return feeds, nil
}

func (f *Feed) HasUnread() bool {
	for _, item := range f.Items {
		if !item.Read {
			return true
		}
	}
	return false
}

func (f *Feed) MarkAllItemsRead() {
	for i := range f.Items {
		f.Items[i].Read = true
	}
}

func (l *FeedList) MarkAllFeedsRead() {
	for _, feed := range l.All {
		feed.MarkAllItemsRead()
	}
}

func (f *Feed) GetFields(fields ...string) []string {
	var result []string
	for _, field := range fields {
		switch field {
		case "Url":
			result = append(result, f.Url)
		case "Category":
			result = append(result, f.Category)
		case "Title":
			if f.Feed == nil || f.Title == "" {
				result = append(result, f.Url)
			} else {
				result = append(result, f.Title)
			}
		case "Latest":
			if f.Feed != nil && f.Items != nil {
				result = append(result, f.Items[0].Title)
			}
			if f.Feed != nil {
				result = append(result, f.Feed.Description)
			} else {
				result = append(result, "")
			}
		}
	}
	return result
}
