package rss

import (
	"bytes"
	"testing"
	"testing/fstest"
	"time"

	"github.com/mmcdole/gofeed"
)

func newList() List {
	unreadRssItem := RssItem{
		Read: false,
		Item: &gofeed.Item{Title: "Latest item title"},
	}

	readRssItem := RssItem{
		Read: true,
		Item: &gofeed.Item{Title: "Latest item title"},
	}

	rssFeedWithoutItems := RssFeed{
		Url:      "example.com",
		Category: "Serious",
		Feed: &gofeed.Feed{
			Title:       "Feed title",
			Description: "Feed description",
		},
	}

	rssFeedUnloaded := RssFeed{Url: "example.com"}

	rssFeed := RssFeed{
		Url:      "example.com",
		Category: "Fun",
		Feed: &gofeed.Feed{
			Title:       "Feed title",
			Description: "Feed description",
		},
		RssItems: []*RssItem{&unreadRssItem, &readRssItem},
	}

	return List{
		Feeds:     []*RssFeed{&rssFeed, &rssFeedUnloaded, &rssFeedWithoutItems},
		FeedIndex: make(map[string]*RssFeed),
	}
}

func TestLists(t *testing.T) {
	t.Run("Should marshal feed list to JSON", func(t *testing.T) {
		l := newList()

		_, err := l.ToJson()
		if err != nil {
			t.Errorf("Error marshaling feed list to JSON: %q", err)
		}
	})

	t.Run("Should write JSON to file", func(t *testing.T) {
		l := newList()
		var buf bytes.Buffer

		err := l.Save(&buf)
		if err != nil {
			t.Fatalf("Unexpected error: %q", err)
		}

		got := buf.String()
		if !bytes.Contains([]byte(got), []byte("Latest item title")) {
			t.Errorf("JSON output does not contain expected feeds: %s", got)
		}
	})

	t.Run("Should return specified category", func(t *testing.T) {
		l := newList()

		category := "Fun"
		feeds, err := l.GetCategory(category)
		if err != nil {
			t.Errorf("Error returning category: %q", err)
		}

		if len(feeds) == 0 {
			t.Error("No feeds returned")
		}

		for _, feed := range feeds {
			if feed.Category != category {
				t.Error("Wrong category returned")
			}
		}
	})

	t.Run("Should handle unspecified category", func(t *testing.T) {
		l := newList()

		var category string
		_, err := l.GetCategory(category)
		assertError(t, err, ErrNoCategoryGiven)
	})

	t.Run("Should return all categories", func(t *testing.T) {
		l := newList()

		categories, err := l.GetAllCategories()
		if err != nil {
			t.Errorf("Error getting categories: %q", err)
		}

		control := []string{"Fun", "Serious"}
		for _, category := range control {
			feeds, ok := categories[category]
			if !ok {
				t.Errorf("Category not returned: %s", category)
			}
			for _, feed := range feeds {
				if feed.Category != category {
					t.Errorf("Feed has wrong category: got %s, want %s", feed.Category, category)
				}
			}
		}
	})

	t.Run("Should return all categories", func(t *testing.T) {
		l := newList()

		categories, err := l.GetAllCategories()
		if err != nil {
			t.Errorf("Error getting categories: %q", err)
		}

		feeds, ok := categories["Uncategorized"]
		if !ok {
			t.Error("Uncategorized feeds not returned")
		}
		for _, feed := range feeds {
			if feed.Category != "" {
				t.Errorf("Feed has wrong category: got %s, want Uncategorized", feed.Category)
			}
		}
	})

	t.Run("Should restore feeds from JSON file", func(t *testing.T) {
		l := newList()

		var buf bytes.Buffer

		err := l.Save(&buf)
		if err != nil {
			t.Fatalf("Unexpected error: %q", err)
		}

		err = l.Restore(&buf)
		if err != nil {
			t.Fatalf("Unexpected error restoring: %q", err)
		}
	})

	t.Run("Should handle restore feeds from empty JSON file", func(t *testing.T) {
		var l List

		var buf bytes.Buffer

		err := l.Save(&buf)
		if err != nil {
			t.Fatalf("Unexpected error: %q", err)
		}

		err = l.Restore(&buf)
		assertError(t, err, ErrChacheEmpty)
		if err == nil {
			t.Fatalf("")
		}
	})

	t.Run("Should load list", func(t *testing.T) {
		fs := fstest.MapFS{
			"urls.yaml": {Data: testData(t, "test_urls.yaml")},
		}

		l, err := LoadList(fs)
		if err != nil {
			t.Errorf("Error loading list: %q", err)
		}

		if l == nil {
			t.Error("List should have been returned")
		}
	})

	t.Run("Should handle urls.yaml not existing", func(t *testing.T) {
		fs := fstest.MapFS{}

		l, err := LoadList(fs)
		if err == nil {
			t.Errorf("Should throw error")
		}

		if l == nil {
			t.Error("List should have been returned")
		}
	})

	t.Run("Should handle empty urls.yaml file", func(t *testing.T) {
		fs := fstest.MapFS{
			"urls.yaml": {Data: []byte(``)},
		}

		l, err := LoadList(fs)
		if err != nil {
			t.Errorf("Error loading list: %q", err)
		}

		if l == nil {
			t.Error("List should have been returned")
		}
	})

	t.Run("Should handle invalid urls.yaml file", func(t *testing.T) {
		fs := fstest.MapFS{
			"urls.yaml": {Data: testData(t, "feed.xml")},
		}

		l, err := LoadList(fs)
		if err == nil {
			t.Errorf("Should have thrown error")
		}

		if l == nil {
			t.Error("List should have been returned")
		}
	})

	t.Run("Should handle invalid JSON file", func(t *testing.T) {
		l := newList()
		err := l.Restore(invalidJson)
		if err == nil {
			t.Error("Should handle invalid JSON file")
		}
	})

	t.Run("Add feeds to list", func(t *testing.T) {
		feeds := make([]*RssFeed, 3)
		for i := range feeds {
			feeds[i] = &RssFeed{}
		}

		var l List

		l.Add(feeds...)

		if len(l.Feeds) != len(feeds) {
			t.Errorf("Wrong number of feeds added to list")
		}
	})

	t.Run("Update all feeds in list", func(t *testing.T) {
		server := Server(t, testData(t, "feed.xml"))
		defer server.Close()

		feeds := []*RssFeed{
			{Url: server.URL},
			{Url: server.URL},
			{Url: ""},
		}

		var l List
		l.Add(feeds...)

		results, err := l.UpdateAllFeeds()
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

	t.Run("Update all only when feeds in list", func(t *testing.T) {
		var l List

		_, err := l.UpdateAllFeeds()
		assertError(t, err, ErrNoFeedsInList)
	})

	t.Run("Mark all feeds read in list", func(t *testing.T) {
		var rssFeed RssFeed
		var l List

		l.Add(&rssFeed)

		l.MarkAllFeedsRead()

		for _, feed := range l.Feeds {
			if feed.HasUnread() == true {
				t.Error("Error marking all feeds read in l")
			}
		}
	})

	t.Run("Create feeds from YAML", func(t *testing.T) {
		_, _, _, _, _, l := newTestData()
		fs := fstest.MapFS{
			"urls.yaml": {Data: testData(t, "test_urls.yaml")},
		}

		err := l.CreateFeedsFromYaml(fs, "urls.yaml")
		if err != nil {
			t.Errorf("Error reading file: %q", err)
		}

		rawItemCount := bytes.Count(testData(t, "test_urls.yaml"), []byte(`http`))
		if len(l.Feeds) != rawItemCount {
			t.Errorf("Wrong number of feeds created, wanted %d, got %d", rawItemCount, len(l.Feeds))
		}

		for _, feed := range l.Feeds {
			if feed.Url == "" {
				t.Error("Feed URL not set when creating from file")
			}
		}
	})

	t.Run("Handle missing feeds file", func(t *testing.T) {
		_, _, _, _, _, l := newTestData()

		fs := fstest.MapFS{}

		err := l.CreateFeedsFromYaml(fs, "urls.yaml")
		if err == nil {
			t.Error("Should raise error when file not found")
		}
	})

	t.Run("Handle invalid feeds file", func(t *testing.T) {
		_, _, _, _, _, l := newTestData()

		fs := fstest.MapFS{
			"urls.yaml": {Data: []byte("invalid: [unbalanced")},
		}

		err := l.CreateFeedsFromYaml(fs, "urls.yaml")
		if err == nil {
			t.Error("Should raise error when file invalid")
		}
	})

	t.Run("Should return the next unread feed correctly", func(t *testing.T) {
		feed1 := &RssFeed{
			RssItems: []*RssItem{
				{Read: true},
			},
		}

		feed2 := &RssFeed{
			RssItems: []*RssItem{
				{Read: true},
			},
		}

		feed3 := &RssFeed{
			RssItems: []*RssItem{
				{Read: false},
			},
		}

		l := &List{
			Feeds: []*RssFeed{feed1, feed2, feed3},
		}

		tests := []struct {
			name     string
			prev     *RssFeed
			index    int
			expected *RssFeed
		}{
			{"next unread after first", feed1, 2, feed3},
			{"next unread after second", feed2, 2, feed3},
			{"next unread after last", feed3, -1, nil},
			{"not in list", &RssFeed{}, -1, nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				index, got := l.NextUnreadFeed(tt.prev)
				if got != tt.expected {
					t.Errorf("NextAfter(%p) = %p, want %p", tt.prev, got, tt.expected)
				}
				if index != tt.index {
					t.Errorf("Wrong index returned, want %d, got %d", tt.index, index)
				}
			})
		}
	})
}
