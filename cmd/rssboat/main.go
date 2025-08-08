package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/emilosman/rssboat/internal/rss"
)

func main() {
	feeds, err := rss.CreateFeedsFromFS(os.DirFS("."))
	if err != nil {
		slog.Error("Error", "error", err)
		os.Exit(1)
	}

	var feedList rss.FeedList
	feedList.Add(feeds)

	feedList.UpdateAll()

	for _, feed := range feedList.All {
		if feed.Data != nil {
			fmt.Println(feed.Data.Title)
		} else {
			fmt.Println(feed.Error)
		}
	}
}
