package rss

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"sort"
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

func (i *RssItem) MarkRead() {
	i.Read = true
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

		f.RssItems = append(f.RssItems, RssItem{
			Item: item,
			Read: false,
		})
		existing[key] = struct{}{}
	}
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

	yamlData, _ := io.ReadAll(file)

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

func (i *RssItem) GetField(field string) string {
	switch field {
	case "Title":
		if i.Read == false {
			return fmt.Sprintf("🟢 %s", i.Item.Title)
		}
		return i.Item.Title
	default:
		return ""
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

func (fl *FeedList) ToJson() ([]byte, error) {
	return json.Marshal(fl)
}

/*
Save to file

f, _ := os.Create("data.json")
defer f.Close()
feedList.Save(f)
*/
func (fl *FeedList) Save(w io.Writer) error {
	data, err := fl.ToJson()
	_, err = w.Write(data)
	return err
}

/*
Restore from file

f, _ := os.Open("data.json")
defer f.Close()
feedList, err := Restore(f)

	if err != nil {
			log.Fatalf("failed to restore feeds: %v", err)
	}
*/
func Restore(r io.Reader) (FeedList, error) {
	var fl FeedList
	dec := json.NewDecoder(r)
	if err := dec.Decode(&fl); err != nil {
		return fl, err
	}
	return fl, nil
}
