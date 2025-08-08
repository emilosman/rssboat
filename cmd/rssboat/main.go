package main

import (
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
}
