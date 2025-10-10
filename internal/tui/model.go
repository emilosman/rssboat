package tui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emilosman/rssboat/internal/rss"
	"github.com/muesli/reflow/wordwrap"
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
func (r rssListItem) FilterValue() string { return r.item.Content() }

type model struct {
	prog      *tea.Program
	ready     bool
	l         *rss.List
	f         *rss.RssFeed
	i         *rss.RssItem
	lf        list.Model
	li        list.Model
	v         viewport.Model
	tabs      []string
	activeTab int
}

func initialModel() *model {
	configFilePath, err := rss.ConfigFilePath()
	if err != nil {
		fmt.Println("Error opening config dir", err)
	}
	filesystem := os.DirFS(configFilePath)
	l, err := rss.LoadList(filesystem)
	t := getTabs(l)

	df := list.NewDefaultDelegate()
	df.ShortHelpFunc = listShortHelp
	df.FullHelpFunc = listFullHelp

	di := list.NewDefaultDelegate()
	di.ShortHelpFunc = itemsShortHelp
	di.FullHelpFunc = itemsFullHelp

	m := &model{
		l:         l,
		lf:        list.New(nil, df, 0, 0),
		li:        list.New(nil, di, 0, 0),
		tabs:      t,
		activeTab: 0,
		v:         viewport.New(10, 10),
	}

	rebuildFeedList(m)

	m.lf.DisableQuitKeybindings()
	m.li.DisableQuitKeybindings()
	m.lf.Title = activeTab(t, 0)

	if err != nil {
		m.lf.NewStatusMessage(err.Error())
	}

	if len(m.l.Feeds) == 0 {
		m.lf.NewStatusMessage(MsgNoFeedsInList)
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
			m.lf.NewStatusMessage(fmt.Sprintf("Error updating: %v", msg.Err))
			m.li.NewStatusMessage(fmt.Sprintf("Error updating: %v", msg.Err))
		} else {
			rebuildFeedList(m)
			m.lf.NewStatusMessage(fmt.Sprintf("Updated %s", msg.Feed.Url))
			m.li.NewStatusMessage(fmt.Sprintf("Updated %s", msg.Feed.Url))
		}
		return m, nil
	case feedsDoneMsg:
		m.lf.NewStatusMessage(MsgAllFeedsUpdated)
		m.li.NewStatusMessage(MsgAllFeedsUpdated)
		return m, nil
	case tea.KeyMsg:
		var handlers map[string]keyHandler
		lfState := m.lf.FilterState().String()
		liState := m.li.FilterState().String()

		if lfState != "filtering" && liState != "filtering" {
			if m.f != nil {
				handlers = itemKeyHandlers
			} else {
				handlers = feedKeyHandlers

				i, err := strconv.Atoi(msg.String())
				if err == nil {
					cmd := handleTabNumber(m, i)
					return m, cmd
				}
			}

			if handler, ok := handlers[msg.String()]; ok {
				cmd := handler(m)
				return m, cmd
			}
		}

	case tea.WindowSizeMsg:
		lh, lv := listStyle.GetFrameSize()
		m.lf.SetSize(msg.Width-lh, msg.Height-lv)
		m.li.SetSize(msg.Width-lh, msg.Height-lv)

		vh, vv := viewStyle.GetFrameSize()
		m.v.Width = msg.Width - vh
		m.v.Height = msg.Height - vv
		if m.i != nil {
			m.v.SetContent(wordwrap.String(m.i.Content(), m.v.Width))
		}
	}

	var cmd tea.Cmd

	switch {
	case m.i != nil:
		m.v, cmd = m.v.Update(msg)
	case m.f != nil:
		m.li, cmd = m.li.Update(msg)
	default:
		m.lf, cmd = m.lf.Update(msg)
	}

	return m, cmd
}
