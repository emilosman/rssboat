package rss

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"

	yaml "github.com/goccy/go-yaml"
)

type List struct {
	Feeds     []*RssFeed
	FeedIndex map[string]*RssFeed
}

type FeedResult struct {
	Feed *RssFeed
	Err  error
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

func (l *List) Add(feeds ...*RssFeed) {
	l.Feeds = append(l.Feeds, feeds...)
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
			feeds = append(feeds, feed)
		}
	}

	l.Feeds = feeds

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

	if len(decoded.Feeds) == 0 {
		return ErrCacheEmpty
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
