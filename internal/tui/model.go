package tui

import (
	"fmt"
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
	prog         *tea.Program
	l            *rss.List
	selectedFeed *rss.RssFeed
	lf           list.Model
	li           list.Model
	activeTab    int
	tabs         []string
}

func initialModel() *model {
	configFilePath, err := rss.ConfigFilePath()
	if err != nil {
		fmt.Println("Error opening config dir", err)
	}

	filesystem := os.DirFS(configFilePath)

	l, err := rss.LoadList(filesystem)

	t := getTabs(l)
	a := 0
	feeds := buildFeedList(l, t, a)

	m := &model{
		l:    l,
		lf:   list.New(feeds, list.NewDefaultDelegate(), 0, 0),
		li:   list.New(nil, list.NewDefaultDelegate(), 0, 0),
		tabs: t,
	}

	m.lf.DisableQuitKeybindings()
	m.li.DisableQuitKeybindings()
	m.lf.Title = activeTab(t, 0)

	if err != nil {
		m.lf.NewStatusMessage(err.Error())
	}

	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case feedUpdatedMsg:
		if msg.Err != nil {
			m.lf.NewStatusMessage(fmt.Sprintf("Error updating %s: %v", msg.Feed.Url, msg.Err))
		} else {
			rebuildFeedList(m)
			m.lf.NewStatusMessage(fmt.Sprintf("Updated %s", msg.Feed.Url))
		}
		return m, nil
	case feedsDoneMsg:
		m.lf.NewStatusMessage(MsgAllFeedsUpdated)
		m.li.NewStatusMessage(MsgAllFeedsUpdated)
		return m, nil
	case tea.KeyMsg:
		var handlers map[string]keyHandler
		if m.lf.FilterState().String() != "filtering" && m.li.FilterState().String() != "filtering" {
			if m.selectedFeed == nil {
				handlers = feedKeyHandlers
			} else {
				handlers = itemKeyHandlers
			}
		}

		if handler, ok := handlers[msg.String()]; ok {
			cmd := handler(m)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.lf.SetSize(msg.Width-h, msg.Height-v)
		m.li.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	if m.selectedFeed != nil {
		m.li, cmd = m.li.Update(msg)
	} else {
		m.lf, cmd = m.lf.Update(msg)
	}
	return m, cmd
}
