package rss

import (
	"errors"

	"github.com/mmcdole/gofeed"
)

var ErrFeedHasNoUrl = errors.New("Feed has no URL")

type Feed struct {
	url  string
	data *gofeed.Feed
}

type FeedList struct {
	All []Feed
}

func (f *Feed) GetFeed() error {
	if f.url == "" {
		return ErrFeedHasNoUrl
	}

	fp := gofeed.NewParser()
	data, err := fp.ParseURL(f.url)
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
		return errors.New("No feeds in list")
	}

	for _, feed := range l.All {
		err := feed.GetFeed()
		if err != nil {
			return err
		}
	}

	return nil
}
