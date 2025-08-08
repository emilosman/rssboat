package rss

import (
	"errors"
	"io"
	"io/fs"

	yaml "github.com/goccy/go-yaml"
	"github.com/mmcdole/gofeed"
)

var ErrFeedHasNoUrl = errors.New("Feed has no URL")
var ErrNoFeedsInList = errors.New("No feeds in list")

type Feed struct {
	Url      string
	Category string
	data     *gofeed.Feed
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
		return err
	}
	f.data = data

	return nil
}

func (l *FeedList) Add(feeds []Feed) {
	l.All = append(l.All, feeds...)
}

func (l *FeedList) UpdateAll() error {
	if len(l.All) == 0 {
		return ErrNoFeedsInList
	}

	for _, feed := range l.All {
		err := feed.GetFeed()
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateFeedsFromFS(filesystem fs.FS) ([]Feed, error) {
	var feeds []Feed
	file, err := filesystem.Open("feeds.yml")
	if err != nil {
		return feeds, err
	}
	defer file.Close()

	yamlData, err := io.ReadAll(file)
	if err != nil {
		return feeds, err
	}

	var raw map[string][]string
	if err := yaml.Unmarshal(yamlData, &raw); err != nil {
		return feeds, err
	}

	for category, urls := range raw {
		for _, u := range urls {
			feeds = append(feeds, Feed{
				Url:      u,
				Category: category,
				data:     nil,
			})
		}
	}

	return feeds, nil
}
