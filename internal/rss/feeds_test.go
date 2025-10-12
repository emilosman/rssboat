package rss

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

func TestFeeds(t *testing.T) {
	t.Run("Get url instead of title when title not set", func(t *testing.T) {
		_, _, rssFeed, _, _, _ := newTestData()

		rssFeed.Feed.Title = ""
		field := rssFeed.Title()

		want := fmt.Sprintf("+ %s", rssFeed.Url)
		if field != want {
			t.Error("Feed title should be url when no title present")
		}

		rssFeed.Feed.Title = "Feed title"

		field = rssFeed.Title()
		if field != "+ Feed title" {
			t.Error("Unread feed title not returned")
		}

		rssFeed.MarkAllItemsRead()

		field = rssFeed.Title()
		if field != "Feed title" {
			t.Error("Read feed title not returned")
		}
	})

	t.Run("Should get feed description when feed does not have items", func(t *testing.T) {
		_, _, _, rssFeedWithoutItems, _, _ := newTestData()

		want := "Feed description"
		got := rssFeedWithoutItems.Latest()
		if got != want {
			t.Errorf("Did not get correct value, wanted %s, got %s", want, got)
		}
	})

	t.Run("Should get latest item title when items present", func(t *testing.T) {
		_, _, rssFeed, _, _, _ := newTestData()

		want := "Latest item title"
		got := rssFeed.Latest()
		if got != want {
			t.Errorf("Did not get correct value, wanted %s, got %s", want, got)
		}
	})

	t.Run("Should get error message if present", func(t *testing.T) {
		_, _, rssFeed, _, _, _ := newTestData()

		want := "Error happened"
		rssFeed.Error = want
		got := rssFeed.Latest()
		if got != want {
			t.Errorf("Did not get correct value, wanted %s, got %s", want, got)
		}
	})

	t.Run("Should get message when feed not loaded yet", func(t *testing.T) {
		_, _, _, _, rssFeedUnloaded, _ := newTestData()

		want := MsgFeedNotLoaded
		got := rssFeedUnloaded.Latest()
		if got != want {
			t.Errorf("Did not get latest feed item title, wanted %s, got %s", want, got)
		}
	})

	t.Run("Feed has unread item", func(t *testing.T) {
		unreadFeedItem := RssItem{
			Read: false,
			Item: &gofeed.Item{
				Title: "Latest item title",
			},
		}
		readFeedItem := RssItem{
			Read: true,
			Item: &gofeed.Item{
				Title: "Latest item title",
			},
		}
		rssFeed := RssFeed{
			Url:      "example.com",
			Category: "Fun",
			Feed: &gofeed.Feed{
				Title:       "Feed title",
				Description: "Feed description",
			},
			RssItems: []*RssItem{&unreadFeedItem, &readFeedItem},
		}

		if rssFeed.HasUnread() == false {
			t.Error("Feed should know there are unread items")
		}
	})

	t.Run("Mark all items read in feed", func(t *testing.T) {
		unreadFeedItem := RssItem{
			Read: false,
			Item: &gofeed.Item{
				Title: "Latest item title",
			},
		}
		readFeedItem := RssItem{
			Read: true,
			Item: &gofeed.Item{
				Title: "Latest item title",
			},
		}
		rssFeed := RssFeed{
			Url:      "example.com",
			Category: "Fun",
			Feed: &gofeed.Feed{
				Title:       "Feed title",
				Description: "Feed description",
			},
			RssItems: []*RssItem{&unreadFeedItem, &readFeedItem},
		}

		rssFeed.MarkAllItemsRead()

		if rssFeed.HasUnread() == true {
			t.Error("Error marking all items read in feed")
		}
	})

	t.Run("Mark feeds as read", func(t *testing.T) {
		feeds := []*RssFeed{
			{RssItems: []*RssItem{{Read: false}, {Read: false}}},
			{RssItems: []*RssItem{{Read: false}, {Read: false}}},
			{RssItems: []*RssItem{{Read: false}, {Read: false}}},
		}

		MarkFeedsAsRead(feeds...)

		for i := range feeds {
			if feeds[i].HasUnread() {
				t.Error("Feed not marked as read")
			}
		}
	})

	t.Run("Get feed if url present", func(t *testing.T) {
		var rssFeed RssFeed

		err := rssFeed.GetFeed()
		assertError(t, err, ErrFeedHasNoUrl)
	})

	t.Run("Get and parse feed", func(t *testing.T) {
		server := Server(t, testData(t, "feed.xml"))
		defer server.Close()

		rssFeed := RssFeed{Url: server.URL}

		err := rssFeed.GetFeed()
		if err != nil {
			t.Errorf("Error getting feed %q", err)
		}

		if rssFeed.Error != "" {
			t.Error("Should unset error on feed")
		}

		if rssFeed.Feed.Title != "NASA Space Station News" {
			t.Error("Error parsing feed")
		}

		rawItemCount := bytes.Count(testData(t, "feed.xml"), []byte(`<item>`))
		if len(rssFeed.RssItems) != rawItemCount {
			t.Errorf("Wrong number of feed items, wanted %d, got %d", rawItemCount, len(rssFeed.RssItems))
		}

		if rssFeed.RssItems[0].Item.Title != "Louisiana Students to Hear from NASA Astronauts Aboard Space Station" {
			t.Error("Wrong feed item title")
		}

		for _, item := range rssFeed.RssItems {
			if item.Item.Link == "" {
				t.Error("All items should have Link")
			}
		}
	})

	t.Run("Feed items should be sorted by pub date", func(t *testing.T) {
		server := Server(t, testData(t, "feed.xml"))
		defer server.Close()

		rssFeed := RssFeed{Url: server.URL}

		err := rssFeed.GetFeed()
		if err != nil {
			t.Errorf("Error getting feed %q", err)
		}

		prevItem := rssFeed.RssItems[0]
		for i := range rssFeed.RssItems[1:] {
			if rssFeed.RssItems[i].Item.PublishedParsed.After(*prevItem.Item.PublishedParsed) {
				t.Error("Wrong order of feed items")
			}
			prevItem = rssFeed.RssItems[i]
		}
	})

	t.Run("Should return the next unread item correctly", func(t *testing.T) {
		item1 := &RssItem{Read: true}
		item2 := &RssItem{Read: true}
		item3 := &RssItem{Read: false}

		rssFeed := RssFeed{
			RssItems: []*RssItem{item1, item2, item3},
		}

		tests := []struct {
			name     string
			prev     *RssItem
			expected *RssItem
		}{
			{"next unread after first", item1, item3},
			{"next unread after second", item2, item3},
			{"next unread after last", item3, nil},
			{"not in list", &RssItem{}, nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := rssFeed.NextUnread(tt.prev)
				if got != tt.expected {
					t.Errorf("NextAfter(%p) = %p, want %p", tt.prev, got, tt.expected)
				}
			})
		}
	})

	t.Run("Do not overwrite read state", func(t *testing.T) {
		server := Server(t, testData(t, "feed.xml"))
		defer server.Close()

		rssFeed := RssFeed{Url: server.URL}

		err := rssFeed.GetFeed()
		if err != nil {
			t.Errorf("Error getting feed %q", err)
		}

		rssFeed.MarkAllItemsRead()
		err = rssFeed.GetFeed()
		if err != nil {
			t.Errorf("Error getting feed %q", err)
		}

		if rssFeed.HasUnread() {
			t.Error("Unread state should not be overwritten")
		}
	})

	t.Run("Handle server error", func(t *testing.T) {
		server := ServerNotFound(t)
		defer server.Close()

		rssFeed := RssFeed{Url: server.URL}

		err := rssFeed.GetFeed()
		if err == nil {
			t.Errorf("Should return error on server error: %q", err)
		}

		if rssFeed.Error == "" {
			t.Errorf("Should store error on feed: %s", rssFeed.Error)
		}
	})

	t.Run("Update feeds", func(t *testing.T) {
		server := Server(t, testData(t, "feed.xml"))
		defer server.Close()

		feeds := []*RssFeed{
			{Url: server.URL},
			{Url: server.URL},
			{Url: ""},
		}

		results, err := UpdateFeeds(feeds...)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		received := 0
		successes := 0
		failures := 0

		for range feeds {
			select {
			case res := <-results:
				received++
				if res.Err != nil {
					failures++
				} else {
					successes++
				}
			case <-time.After(2 * time.Second):
				t.Fatalf("timeout waiting for feed results")
			}
		}

		if received != len(feeds) {
			t.Errorf("expected %d results, got %d", len(feeds), received)
		}
		if successes == 0 {
			t.Errorf("expected at least one successful feed")
		}
		if failures == 0 {
			t.Errorf("expected at least one failed feed")
		}
	})
}
