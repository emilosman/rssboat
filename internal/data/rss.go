package rss

import (
	"github.com/mmcdole/gofeed"
)

func GetFeed(feedUrl string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedUrl)
	if err != nil {
		return feed, err
	}

	return feed, nil
}
