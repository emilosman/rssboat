package rss

import "testing"

var feedUrls = []string{
	"https://www.reddit.com/r/golang.rss",
	"https://cprss.s3.amazonaws.com/golangweekly.com.xml",
	"https://go.dev/blog/feed.atom",
	"https://commandcenter.blogspot.com/feeds/posts/default?alt=rss",
	"https://research.swtch.com/feed.atom",
	"https://www.americanexpress.io/feed.xml",
	"https://golang.cafe/rss",
}

func TestFeed(t *testing.T) {
	t.Run("Get feeds", func(t *testing.T) {
		feeds, err := GetFeeds(feedUrls)
		if err != nil {
			t.Errorf("Error getting feeds %q", err)
		}

		if len(feeds) != len(feedUrls) {
			t.Errorf("Wrong number of feeds returned, wanted %d, got %d", len(feedUrls), len(feeds))
		}
	})
}
