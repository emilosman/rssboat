package core

import "github.com/mmcdole/gofeed"

type Feed struct {
	url  string
	data gofeed.Feed
}

type FeedList struct {
	All []Feed
}
