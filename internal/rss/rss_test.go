package rss

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/mmcdole/gofeed"
)

var (
	unreadFeedItem = RssItem{
		Read: false,
		Item: &gofeed.Item{
			Title: "Latest item title",
		},
	}

	readFeedItem = RssItem{
		Read: true,
		Item: &gofeed.Item{
			Title: "Latest item title",
		},
	}

	rssFeed = RssFeed{
		Url:      "example.com",
		Category: "Fun",
		Feed: &gofeed.Feed{
			Title:       "Feed title",
			Description: "Feed description",
		},
		RssItems: []RssItem{unreadFeedItem, readFeedItem},
	}

	rssFeedWithoutItems = RssFeed{
		Url:      "example.com",
		Category: "Fun",
		Feed: &gofeed.Feed{
			Title:       "Feed title",
			Description: "Feed description",
		},
	}

	rssFeedUnloaded = RssFeed{
		Url:      "example.com",
		Category: "Fun",
	}

	feedList = FeedList{
		All: []*RssFeed{
			&rssFeed, &rssFeedUnloaded, &rssFeedWithoutItems,
		},
	}
)

func TestFeed(t *testing.T) {
	t.Run("Should marshal feed list to JSON", func(t *testing.T) {
		_, err := feedList.ToJson()
		if err != nil {
			t.Errorf("Error marshaling feed list to JSON: %q", err)
		}
	})

	t.Run("Should write JSON to file", func(t *testing.T) {
		var buf bytes.Buffer

		err := feedList.Save(&buf)
		if err != nil {
			t.Fatalf("Unexpected error: %q", err)
		}

		got := buf.String()
		if !bytes.Contains([]byte(got), []byte("Latest item title")) {
			t.Errorf("JSON output does not contain expected feeds: %s", got)
		}
	})

	t.Run("Should restore feeds from JSON file", func(t *testing.T) {
		var buf bytes.Buffer

		err := feedList.Save(&buf)
		if err != nil {
			t.Fatalf("Unexpected error: %q", err)
		}

		got, err := Restore(&buf)
		if err != nil {
			t.Fatalf("Unexpected error restoring: %q", err)
		}

		if len(got.All) != len(feedList.All) {
			t.Errorf("Expected %d feeds, got %d", len(feedList.All), len(got.All))
		}

		if got.All[0].Feed.Title != feedList.All[0].Feed.Title {
			t.Errorf("Expected first feed title %q, got %q", feedList.All[0].Feed.Title, got.All[0].Feed.Title)
		}
	})

	t.Run("Should handle invalid JSON file", func(t *testing.T) {
		_, err := Restore(invalidJson)
		if err == nil {
			t.Error("Should handle invalid JSON file")
		}
	})

	t.Run("Get url instead of title when title not set", func(t *testing.T) {
		columnName := "Title"
		rssFeed.Feed.Title = ""
		field := rssFeed.GetField(columnName)

		if field != rssFeed.Url {
			t.Error("Feed title should be url when no title present")
		}

		rssFeed.Feed.Title = "Feed title"

		field = rssFeed.GetField(columnName)
		if field != "🟢 Feed title" {
			t.Error("Unread feed title not returned")
		}

		rssFeed.MarkAllItemsRead()

		field = rssFeed.GetField(columnName)
		if field != "Feed title" {
			t.Error("Read feed title not returned")
		}
	})

	t.Run("Test feed fields", func(t *testing.T) {
		fieldsTest := []struct {
			field string
			want  string
		}{
			{"Url", rssFeed.Url},
			{"Category", rssFeed.Category},
		}
		for _, tt := range fieldsTest {
			got := rssFeed.GetField(tt.field)
			if got != tt.want {
				t.Error("Should return correct field")
			}
		}
	})

	t.Run("Should get feed description when feed does not have items", func(t *testing.T) {
		field := "Latest"
		want := "Feed description"
		got := rssFeedWithoutItems.GetField(field)
		if got != want {
			t.Errorf("Did not get correct value, wanted %s, got %s", want, got)
		}
	})

	t.Run("Should get latest item title when items present", func(t *testing.T) {
		field := "Latest"
		want := "Latest item title"
		got := rssFeed.GetField(field)
		if got != want {
			t.Errorf("Did not get correct value, wanted %s, got %s", want, got)
		}
	})

	t.Run("Should get error message if present", func(t *testing.T) {
		field := "Latest"
		want := "Error happened"
		rssFeed.Error = want
		got := rssFeed.GetField(field)
		if got != want {
			t.Errorf("Did not get correct value, wanted %s, got %s", want, got)
		}
	})

	t.Run("Should get message when feed not loaded yet", func(t *testing.T) {
		field := "Latest"
		want := MsgFeedNotLoaded
		rssFeed.Error = want
		got := rssFeedUnloaded.GetField(field)
		if got != want {
			t.Errorf("Did not get latest feed item title, wanted %s, got %s", want, got)
		}
	})

	t.Run("Should handle when no field name given", func(t *testing.T) {
		field := "XYZ"
		want := ""
		got := rssFeed.GetField(field)
		if got != want {
			t.Errorf("Did not get default field value")
		}
	})

	t.Run("Should get read status of unread feed item", func(t *testing.T) {
		field := "Title"
		want := "🟢 Latest item title"
		got := unreadFeedItem.GetField(field)
		if got != want {
			t.Errorf("Did not get correct field value, want %s, got %s", want, got)
		}
	})

	t.Run("Should get title of read feed item", func(t *testing.T) {
		field := "Title"
		want := "Latest item title"
		got := readFeedItem.GetField(field)
		if got != want {
			t.Errorf("Did not get correct field value, want %s, got %s", want, got)
		}
	})
	t.Run("Should handle when no field name given", func(t *testing.T) {
		field := "XYZ"
		want := ""
		got := readFeedItem.GetField(field)
		if got != want {
			t.Errorf("Did not get default field value")
		}
	})

	t.Run("Toggle feed item read flag", func(t *testing.T) {
		var feedItem RssItem

		feedItem.ToggleRead()

		if feedItem.Read != true {
			t.Error("Error toggling feed item read flag")
		}

		feedItem.ToggleRead()

		if feedItem.Read != false {
			t.Error("Error toggling feed item read flag")
		}
	})

	t.Run("Mark item read", func(t *testing.T) {
		var feedItem RssItem

		feedItem.MarkRead()

		if !feedItem.Read {
			t.Error("Item should be read")
		}
	})

	t.Run("Add feeds to feedList", func(t *testing.T) {
		feeds := make([]*RssFeed, 3)
		for i := range feeds {
			feeds[i] = &RssFeed{}
		}

		var feedList FeedList

		feedList.Add(feeds...)

		if len(feedList.All) != len(feeds) {
			t.Errorf("Wrong number of feeds added to list")
		}
	})

	t.Run("Update all feeds in list", func(t *testing.T) {
		var feedList FeedList
		server := Server(t, rssData)
		defer server.Close()

		feeds := make([]*RssFeed, 3)
		for i := range feeds {
			feeds[i] = &RssFeed{}
			feeds[i].Url = server.URL
		}

		feedList.Add(feeds...)

		err := feedList.UpdateAll()
		if err != nil {
			t.Errorf("Error updating feeds: %q", err)
		}

		for _, feed := range feedList.All {
			if feed.Feed == nil {
				t.Errorf("Feed data empty after UpdateAll")
			}
		}

		feedList.MarkAllFeedsRead()
		err = feedList.UpdateAll()
		if err != nil {
			t.Errorf("Error updating feeds: %q", err)
		}

		for _, feed := range feedList.All {
			if feed.HasUnread() == true {
				t.Errorf("Unread state should not be overwritten")
			}
		}
	})

	t.Run("Update all only when feeds in list", func(t *testing.T) {
		var feedList FeedList

		err := feedList.UpdateAll()
		assertError(t, err, ErrNoFeedsInList)

		for _, feed := range feedList.All {
			if feed.Feed == nil {
				t.Errorf("Feed data empty for feed %s", feed.Url)
			}
		}
	})

	t.Run("Feed has unread item", func(t *testing.T) {
		var rssFeed RssFeed

		readItems := make([]RssItem, 3)
		for i := range readItems {
			readItems[i].Read = true
		}

		unreadItem := RssItem{Read: false}

		items := append(readItems, unreadItem)

		rssFeed.RssItems = items

		if rssFeed.HasUnread() == false {
			t.Error("Feed should know there are unread items")
		}
	})

	t.Run("Mark all items read in feed", func(t *testing.T) {
		var rssFeed RssFeed

		unreadItems := make([]RssItem, 3)
		for i := range unreadItems {
			unreadItems[i].Read = false
		}

		rssFeed.RssItems = unreadItems

		rssFeed.MarkAllItemsRead()

		if rssFeed.HasUnread() == true {
			t.Error("Error marking all items read in feed")
		}
	})

	t.Run("Mark all feeds read in feedList", func(t *testing.T) {
		var rssFeed RssFeed
		var feedList FeedList

		unreadItems := make([]RssItem, 3)
		for i := range unreadItems {
			unreadItems[i].Read = false
		}

		rssFeed.RssItems = unreadItems

		feedList.Add(&rssFeed)

		feedList.MarkAllFeedsRead()

		for _, feed := range feedList.All {
			if feed.HasUnread() == true {
				t.Error("Error marking all feeds read in feedList")
			}
		}
	})

	t.Run("Get feed if url present", func(t *testing.T) {
		var rssFeed RssFeed

		err := rssFeed.GetFeed()
		assertError(t, err, ErrFeedHasNoUrl)
	})

	t.Run("Get and parse feed", func(t *testing.T) {
		server := Server(t, rssData)
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

		if len(rssFeed.RssItems) != 5 {
			t.Errorf("Wrong number of feed items, wanted %d, got %d", 5, len(rssFeed.RssItems))
		}

		if rssFeed.RssItems[0].Item.Title != "Louisiana Students to Hear from NASA Astronauts Aboard Space Station" {
			t.Error("Wrong feed item title")
		}
	})

	t.Run("Do not overwrite read state", func(t *testing.T) {
		server := Server(t, rssData)
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

	t.Run("Create feeds from FS", func(t *testing.T) {
		fs := fstest.MapFS{
			"feeds.yml": {Data: yamlData},
		}

		feeds, err := CreateFeedsFromFS(fs)
		if err != nil {
			t.Errorf("Error reading file: %q", err)
		}

		if len(feeds) != 7 {
			t.Errorf("Wrong number of feeds created, wanted %d, get %d", 7, len(feeds))
		}

		for _, feed := range feeds {
			if feed.Url == "" {
				t.Error("Feed URL not set when creating from file")
			}
		}
	})

	t.Run("Handle missing feeds file", func(t *testing.T) {
		fs := fstest.MapFS{
			"other.yml": {Data: []byte(``)},
		}

		_, err := CreateFeedsFromFS(fs)
		if err == nil {
			t.Error("Should raise error when file not found")
		}
	})

	t.Run("Handle invalid feeds file", func(t *testing.T) {
		fs := fstest.MapFS{
			"feeds.yml": {Data: []byte("invalid: [unbalanced")},
		}

		_, err := CreateFeedsFromFS(fs)
		if err == nil {
			t.Error("Should raise error when file invalid")
		}
	})
}

func Server(t *testing.T, data []byte) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	return server
}

func ServerNotFound(t *testing.T) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	return server
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("wanted error but did not get one")
	}

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

var (
	invalidJson = bytes.NewBufferString(`
	{invalid
`)

	yamlData = []byte(`
golang:
  - https://www.reddit.com/r/golang.rss
  - https://cprss.s3.amazonaws.com/golangweekly.com.xml
  - https://go.dev/blog/feed.atom
  - https://commandcenter.blogspot.com/feeds/posts/default?alt=rss
  - https://research.swtch.com/feed.atom
  - https://www.americanexpress.io/feed.xml

jobs:
  - https://golang.cafe/rss
`)

	rssData = []byte(`
<?xml version="1.0"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
   <channel>
      <title>NASA Space Station News</title>
      <link>http://www.nasa.gov/</link>
      <description>A RSS news feed containing the latest NASA press releases on the International Space Station.</description>
      <language>en-us</language>
      <pubDate>Tue, 10 Jun 2003 04:00:00 GMT</pubDate>
      <lastBuildDate>Fri, 21 Jul 2023 09:04 EDT</lastBuildDate>
      <docs>https://www.rssboard.org/rss-specification</docs>
      <generator>Blosxom 2.1.2</generator>
      <managingEditor>neil.armstrong@example.com (Neil Armstrong)</managingEditor>
      <webMaster>sally.ride@example.com (Sally Ride)</webMaster>
      <atom:link href="https://www.rssboard.org/files/sample-rss-2.xml" rel="self" type="application/rss+xml" />
      <item>
         <title>Louisiana Students to Hear from NASA Astronauts Aboard Space Station</title>
         <link>http://www.nasa.gov/press-release/louisiana-students-to-hear-from-nasa-astronauts-aboard-space-station</link>
         <description>As part of the state's first Earth-to-space call, students from Louisiana will have an opportunity soon to hear from NASA astronauts aboard the International Space Station.</description>
         <pubDate>Fri, 21 Jul 2023 09:04 EDT</pubDate>
         <guid>http://www.nasa.gov/press-release/louisiana-students-to-hear-from-nasa-astronauts-aboard-space-station</guid>
      </item>
      <item>
         <description>NASA has selected KBR Wyle Services, LLC, of Fulton, Maryland, to provide mission and flight crew operations support for the International Space Station and future human space exploration.</description>
         <link>http://www.nasa.gov/press-release/nasa-awards-integrated-mission-operations-contract-iii</link>
         <pubDate>Thu, 20 Jul 2023 15:05 EDT</pubDate>
         <guid>http://www.nasa.gov/press-release/nasa-awards-integrated-mission-operations-contract-iii</guid>
      </item>
      <item>
         <title>NASA Expands Options for Spacewalking, Moonwalking Suits</title>
         <link>http://www.nasa.gov/press-release/nasa-expands-options-for-spacewalking-moonwalking-suits-services</link>
         <description>NASA has awarded Axiom Space and Collins Aerospace task orders under existing contracts to advance spacewalking capabilities in low Earth orbit, as well as moonwalking services for Artemis missions.</description>
         <enclosure url="http://www.nasa.gov/sites/default/files/styles/1x1_cardfeed/public/thumbnails/image/iss068e027836orig.jpg?itok=ucNUaaGx" length="1032272" type="image/jpeg" />
         <pubDate>Mon, 10 Jul 2023 14:14 EDT</pubDate>
         <guid>http://www.nasa.gov/press-release/nasa-expands-options-for-spacewalking-moonwalking-suits-services</guid>
      </item>
      <item>
         <title>NASA to Provide Coverage as Dragon Departs Station</title>
         <link>http://www.nasa.gov/press-release/nasa-to-provide-coverage-as-dragon-departs-station-with-science</link>
         <description>NASA is set to receive scientific research samples and hardware as a SpaceX Dragon cargo resupply spacecraft departs the International Space Station on Thursday, June 29.</description>
         <pubDate>Tue, 20 May 2003 08:56:02 GMT</pubDate>
         <guid>http://www.nasa.gov/press-release/nasa-to-provide-coverage-as-dragon-departs-station-with-science</guid>
      </item>
      <item>
         <title>NASA Plans Coverage of Roscosmos Spacewalk Outside Space Station</title>
         <link>http://liftoff.msfc.nasa.gov/news/2003/news-laundry.asp</link>
         <description>Compared to earlier spacecraft, the International Space Station has many luxuries, but laundry facilities are not one of them.  Instead, astronauts have other options.</description>
         <enclosure url="http://www.nasa.gov/sites/default/files/styles/1x1_cardfeed/public/thumbnails/image/spacex_dragon_june_29.jpg?itok=nIYlBLme" length="269866" type="image/jpeg" />
         <pubDate>Mon, 26 Jun 2023 12:45 EDT</pubDate>
         <guid>http://liftoff.msfc.nasa.gov/2003/05/20.html#item570</guid>
      </item>
   </channel>
</rss>
`)
)
