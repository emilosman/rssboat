package rss

import (
	"github.com/mmcdole/gofeed"
)

type Feed struct {
	url  string
	data *gofeed.Feed
}

type FeedList struct {
	All []Feed
}

func (f *Feed) GetFeed() error {
	fp := gofeed.NewParser()
	data, err := fp.ParseURL(f.url)
	f.data = data
	if err != nil {
		return err
	}

	return nil
}
