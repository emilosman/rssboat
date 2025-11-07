package tui

import (
	"testing"

	"github.com/emilosman/rssboat/internal/rss"
	"github.com/mmcdole/gofeed"
)

func newList() rss.List {
	unreadRssItem := rss.RssItem{
		Read: false,
		Item: &gofeed.Item{Title: "Latest item title"},
	}

	readRssItem := rss.RssItem{
		Read: true,
		Item: &gofeed.Item{Title: "Latest item title"},
	}

	rssFeedWithoutItems := rss.RssFeed{
		Url:      "example.com",
		Category: "Serious",
		Feed: &gofeed.Feed{
			Title:       "Feed title",
			Description: "Feed description",
		},
	}

	rssFeedUnloaded := rss.RssFeed{Url: "example.com"}

	rssFeed := rss.RssFeed{
		Url:      "example.com",
		Category: "Fun",
		Feed: &gofeed.Feed{
			Title:       "Feed title",
			Description: "Feed description",
		},
		RssItems: []*rss.RssItem{&unreadRssItem, &readRssItem},
	}

	return rss.List{
		Feeds:     []*rss.RssFeed{&rssFeed, &rssFeedUnloaded, &rssFeedWithoutItems},
		FeedIndex: make(map[string]*rss.RssFeed),
		CategoryIndex: map[string][]*rss.RssFeed{
			"Fun":     {&rssFeed},
			"Serious": {&rssFeedWithoutItems},
		},
	}
}

func TestHelpers(t *testing.T) {
	t.Run("Should build feeds list", func(t *testing.T) {
		l := newList()
		tabs := []string{"Fun"}
		m := model{
			l:         &l,
			tabs:      tabs,
			activeTab: 0,
		}
		listItems := buildFeedList(m.l, m.tabs, m.activeTab)

		if len(listItems) == 0 {
			t.Errorf("No list items returned")
		}
	})
}
