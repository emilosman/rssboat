package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emilosman/rssboat/internal/rss"
)

type feedItem struct {
	title, desc string
	rssFeed     *rss.RssFeed
}

func (f feedItem) Title() string       { return f.title }
func (f feedItem) Description() string { return f.desc }
func (f feedItem) FilterValue() string { return f.title }

type rssListItem struct {
	title, desc string
	item        *rss.RssItem
}

func (r rssListItem) Title() string       { return r.title }
func (r rssListItem) Description() string { return r.desc }
func (r rssListItem) FilterValue() string { return r.title }

type model struct {
	feedList     rss.FeedList
	selectedFeed *rss.RssFeed
	feedsList    list.Model
	itemsList    list.Model
}

func initialModel() *model {
	var feedList rss.FeedList
	var initialStatusMsg string
	f, err := os.Open("./data.json")
	if err != nil {
		fmt.Println("Error opening data file:", err)
		filesystem := os.DirFS(".")

		feeds, err := rss.CreateFeedsFromFS(filesystem)
		if err != nil {
			log.Fatal(err)
		}

		feedList.Add(feeds...)
		initialStatusMsg = "Feeds loaded from YAML file"
	} else {
		defer f.Close()

		feedList, err = rss.Restore(f)
		if err != nil {
			log.Fatalf("failed to restore feeds: %v", err)
		}
		initialStatusMsg = "Feeds restored from JSON file"
	}

	all := buildFeedList(feedList.All)

	m := model{
		feedList:  feedList,
		feedsList: list.New(all, list.NewDefaultDelegate(), 0, 0),
		itemsList: list.New(nil, list.NewDefaultDelegate(), 0, 0),
	}
	m.feedsList.DisableQuitKeybindings()
	m.itemsList.DisableQuitKeybindings()
	m.feedsList.Title = "rssboat"
	m.feedsList.NewStatusMessage(initialStatusMsg)

	return &m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		var handlers map[string]keyHandler
		if m.selectedFeed == nil {
			handlers = feedKeyHandlers
		} else {
			handlers = itemKeyHandlers
		}

		if handler, ok := handlers[msg.String()]; ok {
			cmd := handler(m)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.feedsList.SetSize(msg.Width-h, msg.Height-v)
		m.itemsList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	if m.selectedFeed != nil {
		m.itemsList, cmd = m.itemsList.Update(msg)
	} else {
		m.feedsList, cmd = m.feedsList.Update(msg)
	}
	return m, cmd
}
