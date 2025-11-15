package rss

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"slices"
	"sort"

	yaml "github.com/goccy/go-yaml"
)

type List struct {
	Feeds         []*RssFeed
	FeedIndex     map[string]*RssFeed
	CategoryIndex map[string][]*RssFeed
}

type FeedResult struct {
	Feed *RssFeed
	Err  error
}

func (l *List) Categories() []string {
	var categories []string
	for category := range l.CategoryIndex {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

func (l *List) GetCategory(category string) ([]*RssFeed, error) {
	var feeds []*RssFeed

	if category == "" {
		return feeds, ErrNoCategoryGiven
	}

	return l.CategoryIndex[category], nil
}

func (l *List) Add(feeds ...*RssFeed) {
	l.Feeds = append(l.Feeds, feeds...)
}

func (l *List) bookmarks() *RssFeed {
	return l.FeedIndex["Bookmarks"]
}

func (l *List) BookmarkItem(i *RssItem) (string, error) {
	bookmarks := l.bookmarks()
	if bookmarks != nil {
		i.ToggleBookmark()
		if i.Bookmark {
			bookmarks.RssItems = append(bookmarks.RssItems, i)
			return MsgBookmarkItem, nil
		} else {
			i := slices.Index(bookmarks.RssItems, i)
			if i != -1 {
				bookmarks.RssItems = append(bookmarks.RssItems[:i], bookmarks.RssItems[i+1:]...)
			}
			return MsgUnBookmarkItem, nil
		}
	}
	return "", ErrNoBookmarkFeed
}

func (l *List) UpdateAllFeeds() (<-chan FeedResult, error) {
	return UpdateFeeds(l.Feeds...)
}

func (l *List) CreateFeedsFromYaml(filesystem fs.FS, filename string) error {
	file, err := filesystem.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, _ := io.ReadAll(file)

	var raw map[string][]string
	if err := yaml.Unmarshal(data, &raw); err != nil {
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
			l.CategoryIndex[category] = append(l.CategoryIndex[category], feed)
			feeds = append(feeds, feed)
		}
	}

	l.Feeds = append(l.Feeds, feeds...)

	return nil
}

func (l *List) MarkAllFeedsRead() {
	for _, feed := range l.Feeds {
		feed.MarkAllItemsRead()
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

	for _, decodedFeed := range decoded.Feeds {
		feed := l.FeedIndex[decodedFeed.Url]
		if feed != nil {
			feed.Error = decodedFeed.Error
			feed.Feed = decodedFeed.Feed
			feed.RssItems = decodedFeed.RssItems
		}
	}

	return nil
}

func NewListWithDefaults() *List {
	bookmarks := &RssFeed{Url: "Bookmarks"}
	return &List{
		Feeds: []*RssFeed{bookmarks},
		FeedIndex: map[string]*RssFeed{
			"Bookmarks": bookmarks,
		},
		CategoryIndex: map[string][]*RssFeed{},
	}
}

func LoadList(filesystem fs.FS) (*List, error) {
	l := NewListWithDefaults()

	err := l.CreateFeedsFromYaml(filesystem, "urls.yaml")
	if err != nil {
		return l, err
	}

	cacheFilePath, err := CacheFilePath()
	if err != nil {
		return l, err
	}

	f, err := os.Open(cacheFilePath)
	if err != nil {
		return l, err
	}
	defer f.Close()

	err = l.Restore(f)
	if err != nil {
		return l, err
	}

	return l, nil
}
