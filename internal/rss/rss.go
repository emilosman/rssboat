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
	Url      string
	Category string
	Error    string
	Data     *gofeed.Feed
}

type FeedList struct {
	All []Feed
}

func (f *Feed) GetFeed() error {
	if f.Url == "" {
		return ErrFeedHasNoUrl
	}

	fp := gofeed.NewParser()
	data, err := fp.ParseURL(f.Url)
	if err != nil {
		f.Error = err.Error()
		return err
	}

	f.Data = data
	f.Error = ""

	return nil
}

func (l *FeedList) Add(feeds []Feed) {
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

func CreateFeedsFromFS(filesystem fs.FS) ([]Feed, error) {
	var feeds []Feed
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

	for category, urls := range raw {
		for _, u := range urls {
			feeds = append(feeds, Feed{
				Url:      u,
				Category: category,
				Data:     nil,
			})
		}
	}

	return feeds, nil
}
