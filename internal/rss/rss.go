package rss

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	yaml "github.com/goccy/go-yaml"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
)

var (
	ErrFeedHasNoUrl    = errors.New("Feed has no URL")
	ErrNoFeedsInList   = errors.New("No feeds in list")
	ErrNoCategoryGiven = errors.New("No category given")
	ErrChacheEmpty     = errors.New("Cache empty")
	MsgFeedNotLoaded   = "Feed not loaded yet. Press shift+r"
)

type RssFeed struct {
	Url      string
	Category string
	Error    string

	Feed     *gofeed.Feed
	RssItems []*RssItem
}

type FeedResult struct {
	Feed *RssFeed
	Err  error
}

type RssItem struct {
	Item *gofeed.Item
	Read bool
}

type List struct {
	Feeds     []*RssFeed
	FeedIndex map[string]*RssFeed
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

func (l *List) GetCategory(category string) ([]*RssFeed, error) {
	var feeds []*RssFeed

	if category == "" {
		return feeds, ErrNoCategoryGiven
	}

	for _, feed := range l.Feeds {
		if feed.Category == category {
			feeds = append(feeds, feed)
		}
	}

	return feeds, nil
}

func (l *List) GetAllCategories() (map[string][]*RssFeed, error) {
	categories := make(map[string][]*RssFeed)

	for _, feed := range l.Feeds {
		if feed == nil {
			continue
		}
		cat := feed.Category
		if cat == "" {
			cat = "Uncategorized"
		}
		categories[cat] = append(categories[cat], feed)
	}

	return categories, nil
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

func (l *List) Add(feeds ...*RssFeed) {
	l.Feeds = append(l.Feeds, feeds...)
}

func (l *List) UpdateAllAsync() (<-chan FeedResult, error) {
	if len(l.Feeds) == 0 {
		return nil, ErrNoFeedsInList
	}

	results := make(chan FeedResult, len(l.Feeds))
	var wg sync.WaitGroup
	wg.Add(len(l.Feeds))

	for _, feed := range l.Feeds {
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

func (l *List) CreateFeedsFromYaml(filesystem fs.FS, filename string) error {
	file, err := filesystem.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	yamlData, _ := io.ReadAll(file)

	var raw map[string][]string
	if err := yaml.Unmarshal(yamlData, &raw); err != nil {
		return err
	}

	var feeds []*RssFeed
	for category, urls := range raw {
		for _, u := range urls {
			feed := &RssFeed{
				Url:      u,
				Category: category,
			}
			l.FeedIndex[u] = feed
			feeds = append(feeds, feed)
		}
	}

	l.Feeds = feeds

	return nil
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

func (l *List) MarkAllFeedsRead() {
	for _, feed := range l.Feeds {
		feed.MarkAllItemsRead()
	}
}

func (i *RssItem) Title() string {
	title := Clean(i.Item.Title)
	if !i.Read {
		return fmt.Sprintf("🟢 %s", title)
	}
	return title
}

func (f *RssFeed) Title() string {
	if f.Feed == nil || f.Feed.Title == "" {
		return f.Url
	}

	title := Clean(f.Feed.Title)

	if f.HasUnread() {
		return fmt.Sprintf("🟢 %s", title)
	}

	return title
}

func (i *RssItem) Description() string {
	return Clean(i.Item.Description)
}

func (f *RssFeed) Description() string {
	return Clean(f.Feed.Description)
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
		return Clean(last.Item.Title)
	case f.Feed != nil:
		return f.Description()
	default:
		return MsgFeedNotLoaded
	}
}

func (l *List) ToJson() ([]byte, error) {
	return json.Marshal(l)
}

/*
Save to file

f, _ := os.Create("data.json")
defer f.Close()
l.Save(f)
*/
func (l *List) Save(w io.Writer) error {
	data, err := l.ToJson()
	_, err = w.Write(data)
	return err
}

/*
Restore from file

f, _ := os.Open("data.json")
defer f.Close()
l, err := Restore(f)

	if err != nil {
			log.Fatalf("failed to restore feeds: %v", err)
	}
*/
func (l *List) Restore(r io.Reader) error {
	var decoded List
	decoder := json.NewDecoder(r)

	err := decoder.Decode(&decoded)
	if err != nil {
		return err
	}

	if len(decoded.Feeds) == 0 {
		return ErrChacheEmpty
	}

	for _, decodedFeed := range decoded.Feeds {
		feed := l.FeedIndex[decodedFeed.Url]
		if feed != nil {
			*feed = *decodedFeed
		}
	}

	return nil
}

func LoadList(filesystem fs.FS) (*List, error) {
	l := List{
		FeedIndex: make(map[string]*RssFeed),
	}

	err := l.CreateFeedsFromYaml(filesystem, "urls.yaml")
	if err != nil {
		return &l, err
	}

	cacheFilePath, err := CacheFilePath()
	if err != nil {
		return &l, err
	}

	f, err := os.Open(cacheFilePath)
	if err != nil {
		return &l, err
	}
	defer f.Close()

	err = l.Restore(f)
	if err != nil {
		return &l, err
	}

	return &l, nil
}

func CacheFilePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "rssboat")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "data.json"), nil
}

func ConfigFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(dir, "rssboat")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	configFile := filepath.Join(appDir, "urls.yaml")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		f, err := os.Create(configFile)
		if err != nil {
			return "", err
		}
		defer f.Close()
	}

	return appDir, nil
}

func Clean(input string) string {
	p := bluemonday.StrictPolicy()
	clean := p.Sanitize(input)
	decoded := html.UnescapeString(clean)
	return normalizeSpaces(decoded)
}

func normalizeSpaces(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	return strings.Join(strings.Fields(s), " ")
}
