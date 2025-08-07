package rss

import "github.com/mmcdole/gofeed"

func GetFeeds(feedUrls []string) ([]gofeed.Feed, error) {
	var feeds []gofeed.Feed

	for _, feedUrl := range feedUrls {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(feedUrl)
		if err != nil {
			return feeds, err
		}
		feeds = append(feeds, *feed)
	}

	return feeds, nil
}
