package rss

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mmcdole/gofeed"
)

var (
	sanitationTests = []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "strips HTML tags",
			input:    `Hello <b>World</b>`,
			expected: "Hello World",
		},
		{
			name:     "removes scripts",
			input:    `Test <script>alert("x")</script> Done`,
			expected: "Test Done",
		},
		{
			name:     "decodes entities",
			input:    `I&#39;M fine`,
			expected: "I'M fine",
		},
		{
			name:     "normalizes spaces and newlines",
			input:    "Line1\nLine2   Line3\r\nLine4",
			expected: "Line1 Line2 Line3 Line4",
		},
		{
			name:     "mixed case",
			input:    "COASTAL REGION:\n thu - partly cloudy; <sup><i>Run time: Wednesday ,",
			expected: "COASTAL REGION: thu - partly cloudy; Run time: Wednesday ,",
		},
	}
)

func newTestData() (RssItem, RssItem, RssFeed, RssFeed, RssFeed, List) {
	unreadRssItem := RssItem{
		Read: false,
		Item: &gofeed.Item{Title: "Latest item title"},
	}

	readRssItem := RssItem{
		Read: true,
		Item: &gofeed.Item{Title: "Latest item title"},
	}

	rssFeed := RssFeed{
		Url:      "example.com",
		Category: "Fun",
		Feed: &gofeed.Feed{
			Title:       "Feed title",
			Description: "Feed description",
		},
		RssItems: []*RssItem{&unreadRssItem, &readRssItem},
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

	l := List{
		Feeds:     []*RssFeed{&rssFeed, &rssFeedUnloaded, &rssFeedWithoutItems},
		FeedIndex: make(map[string]*RssFeed),
	}

	return unreadRssItem, readRssItem, rssFeed, rssFeedWithoutItems, rssFeedUnloaded, l
}

func TestRss(t *testing.T) {
	t.Run("Should sanitize text", func(t *testing.T) {
		for _, tt := range sanitationTests {
			t.Run(tt.name, func(t *testing.T) {
				got := Clean(tt.input)
				if got != tt.expected {
					t.Errorf("Clean(%q) = %q, want %q", tt.input, got, tt.expected)
				}
			})
		}

		for _, tt := range sanitationTests {
			t.Run(tt.name, func(t *testing.T) {
				item := RssItem{Item: &gofeed.Item{Description: tt.input}}

				got := item.Description()
				if got != tt.expected {
					t.Errorf("Clean(%q) = %q, want %q", tt.input, got, tt.expected)
				}
			})
		}

		for _, tt := range sanitationTests {
			t.Run(tt.name, func(t *testing.T) {
				item := RssItem{
					Read: true,
					Item: &gofeed.Item{Title: tt.input},
				}

				got := item.Title()
				if got != tt.expected {
					t.Errorf("Clean(%q) = %q, want %q", tt.input, got, tt.expected)
				}
			})
		}

		for _, tt := range sanitationTests {
			t.Run(tt.name, func(t *testing.T) {
				feed := RssFeed{
					Feed: &gofeed.Feed{
						Title: tt.input,
					},
				}

				got := feed.Title()
				if got != tt.expected {
					t.Errorf("Clean(%q) = %q, want %q", tt.input, got, tt.expected)
				}
			})
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
         <title>NASA to Provide Coverage as Dragon Departs Station</title>
         <link>http://www.nasa.gov/press-release/nasa-to-provide-coverage-as-dragon-departs-station-with-science</link>
         <description>NASA is set to receive scientific research samples and hardware as a SpaceX Dragon cargo resupply spacecraft departs the International Space Station on Thursday, June 29.</description>
         <pubDate>Tue, 20 May 2003 08:56:02 GMT</pubDate>
         <guid>http://www.nasa.gov/press-release/nasa-to-provide-coverage-as-dragon-departs-station-with-science</guid>
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
         <title>NASA Plans Coverage of Roscosmos Spacewalk Outside Space Station</title>
         <link>http://liftoff.msfc.nasa.gov/news/2003/news-laundry.asp</link>
         <description>Compared to earlier spacecraft, the International Space Station has many luxuries, but laundry facilities are not one of them.  Instead, astronauts have other options.</description>
         <enclosure url="http://www.nasa.gov/sites/default/files/styles/1x1_cardfeed/public/thumbnails/image/spacex_dragon_june_29.jpg?itok=nIYlBLme" length="269866" type="image/jpeg" />
         <pubDate>Mon, 26 Jun 2023 12:45 EDT</pubDate>
         <guid>http://liftoff.msfc.nasa.gov/2003/05/20.html#item570</guid>
      </item>
      <item>
         <title>NASA Plans Coverage of Roscosmos Spacewalk Outside Space Station</title>
         <link>http://liftoff.msfc.nasa.gov/news/2003/news-laundry.asp</link>
         <description>Compared to earlier spacecraft, the International Space Station has many luxuries, but laundry facilities are not one of them.  Instead, astronauts have other options.</description>
         <enclosure url="http://www.nasa.gov/sites/default/files/styles/1x1_cardfeed/public/thumbnails/image/spacex_dragon_june_29.jpg?itok=nIYlBLme" length="269866" type="image/jpeg" />
         <pubDate>Mon, 27 Jun 2023 12:45 EDT</pubDate>
         <guid></guid>
      </item>
      <item>
         <title>NASA Plans Coverage of Roscosmos Spacewalk Outside Space Station</title>
         <link>http://liftoff.msfc.nasa.gov/news/2003/news-laundry.asp</link>
         <description>Compared to earlier spacecraft, the International Space Station has many luxuries, but laundry facilities are not one of them.  Instead, astronauts have other options.</description>
         <enclosure url="http://www.nasa.gov/sites/default/files/styles/1x1_cardfeed/public/thumbnails/image/spacex_dragon_june_29.jpg?itok=nIYlBLme" length="269866" type="image/jpeg" />
         <pubDate></pubDate>
         <guid>http://liftoff.msfc.nasa.gov/2003/05/20.html#item571</guid>
      </item>
   </channel>
</rss>
`)
)
