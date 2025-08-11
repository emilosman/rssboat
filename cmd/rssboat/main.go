package main

import (
	"fmt"
	"os"

	"github.com/emilosman/rssboat/internal/rss"
)

func main() {
	feeds, err := rss.CreateFeedsFromFS(os.DirFS("."))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var feedList rss.FeedList
	feedList.Add(feeds...)

	feedList.UpdateAll()

	for _, feed := range feedList.All {
		if feed.Feed != nil {
			fmt.Println(feed.Title)
		} else {
			fmt.Println(feed.Error)
		}
	}
}
