package rss

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestFeed(t *testing.T) {
	t.Run("Add feeds to feedList", func(t *testing.T) {
		feeds := make([]Feed, 3)

		var feedList FeedList

		feedList.Add(feeds)

		if len(feedList.All) != len(feeds) {
			t.Errorf("Wrong number of feeds added to list")
		}
	})

	t.Run("Update all feeds in list", func(t *testing.T) {
		server := Server(t)
		defer server.Close()

		feeds := make([]Feed, 3)
		for i := range feeds {
			feeds[i].Url = server.URL
		}
		feedList := FeedList{All: feeds}

		err := feedList.UpdateAll()
		if err != nil {
			t.Errorf("Error updating feeds: %q", err)
		}

		for _, feed := range feedList.All {
			if feed.Feed == nil {
				t.Errorf("Feed data empty after UpdateAll")
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

	t.Run("Get feed if url present", func(t *testing.T) {
		var feed Feed

		err := feed.GetFeed()
		assertError(t, err, ErrFeedHasNoUrl)
	})

	t.Run("Get and parse feed", func(t *testing.T) {
		server := Server(t)
		defer server.Close()

		feed := Feed{Url: server.URL}

		err := feed.GetFeed()
		if err != nil {
			t.Errorf("Error getting feed %q", err)
		}

		if feed.Error != "" {
			t.Error("Should unset error on feed")
		}

		if feed.Title != "NASA Space Station News" {
			t.Error("Error parsing feed")
		}
	})

	t.Run("Handle server error", func(t *testing.T) {
		server := ServerNotFound(t)
		defer server.Close()

		feed := Feed{Url: server.URL}

		err := feed.GetFeed()
		if err == nil {
			t.Errorf("Should return error on server error: %q", err)
		}

		if feed.Error == "" {
			t.Errorf("Should store error on feed: %s", feed.Error)
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

func Server(t *testing.T) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(rssData)
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

var yamlData = []byte(`
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

var rssData = []byte(`
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
