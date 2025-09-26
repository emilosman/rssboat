package rss

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mmcdole/gofeed"
)

var (
	invalidJson = bytes.NewBufferString(`
	{invalid
`)

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

func rssData(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/feed.xml")
	if err != nil {
		t.Fatal("Could not read test data")
	}
	return b
}

func yamlData(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/test_urls.yaml")
	if err != nil {
		t.Fatal("Could not read test data")
	}
	return b
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
