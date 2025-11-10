package tui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emilosman/rssboat/internal/rss"
)

type feedUpdatedMsg struct {
	Feed *rss.RssFeed
	Err  error
}

type feedsDoneMsg struct{}
type statusClearMsg struct{}

func updateAllFeedsCmd(m *model) tea.Cmd {
	return func() tea.Msg {
		results, err := m.l.UpdateAllFeeds()
		if err != nil {
			return feedUpdatedMsg{Feed: nil, Err: err}
		}

		go func() {
			for res := range results {
				m.prog.Send(feedUpdatedMsg{Feed: res.Feed, Err: res.Err})
			}
			m.prog.Send(feedsDoneMsg{})
		}()

		return MsgUpdatingAllFeeds
	}
}

func updateTabFeedsCmd(m *model) tea.Cmd {
	return func() tea.Msg {
		feeds, err := m.l.GetCategory(activeTab(m.tabs, m.activeTab))
		if err != nil {
			return feedUpdatedMsg{Feed: nil, Err: err}
		}

		results, err := rss.UpdateFeeds(feeds...)
		if err != nil {
			return feedUpdatedMsg{Feed: nil, Err: err}
		}

		go func() {
			for res := range results {
				m.prog.Send(feedUpdatedMsg{Feed: res.Feed, Err: res.Err})
			}
			m.prog.Send(feedsDoneMsg{})
		}()

		return MsgUpdatingAllFeeds
	}
}

func updateFeedCmd(m *model, feed *rss.RssFeed) tea.Cmd {
	return func() tea.Msg {
		results, err := rss.UpdateFeeds(feed)
		if err != nil {
			return feedUpdatedMsg{Feed: nil, Err: err}
		}

		go func() {
			for res := range results {
				m.prog.Send(feedUpdatedMsg{Feed: res.Feed, Err: res.Err})
			}
			m.prog.Send(feedsDoneMsg{})
		}()

		return fmt.Sprintf("%s %s", MsgUpdatingFeed, feed.Url)
	}
}

// Builds the feed list and sets the items
func rebuildFeedList(m *model) tea.Cmd {
	items := buildFeedList(m.l, m.tabs, m.activeTab)
	m.lf.SetItems(items)
	m.lf.Select(0)
	return nil
}

func rebuildItemsList(m *model) tea.Cmd {
	if m.li.FilterState().String() != "filter applied" {
		items := buildItemsList(m.f)
		m.li.SetItems(items)
		m.li.Select(0)
	}
	return nil
}

// Builds the feed list
func buildFeedList(l *rss.List, t []string, a int) []list.Item {
	var listItems []list.Item

	feeds, err := l.GetCategory(activeTab(t, a))
	if err != nil {
		feeds = l.Feeds
	}

	if len(feeds) != 0 {
		for _, feed := range feeds {
			title := feed.Title()
			description := feed.Latest()

			if feed.HasUnread() {
				title = unreadStyle.Render(title)
			}

			if feed.Error != "" {
				description = errorStyle.Render(description)
			}

			listItems = append(listItems, feedItem{
				title:   title,
				desc:    description,
				rssFeed: feed,
			})
		}
	}

	return listItems
}

func activeTab(t []string, a int) string {
	var activeTab string
	if len(t) != 0 {
		activeTab = t[a]
	}
	return activeTab
}

func buildItemsList(feed *rss.RssFeed) []list.Item {
	listItems := make([]list.Item, 0, len(feed.RssItems))
	for idx := range feed.RssItems {
		ri := feed.RssItems[idx]
		title := ri.Title()
		description := ri.Description()

		if ri.Bookmark {
			title = bookmarkStyle.Render(title)
		}

		if !ri.Read {
			title = unreadStyle.Render(title)
		}

		listItems = append(listItems, rssListItem{
			title: title,
			desc:  description,
			item:  ri,
		})
	}
	return listItems
}

func renderedTabs(m *model) string {
	var renderedTabs string
	for i, tab := range m.tabs {
		if i == m.activeTab {
			renderedTabs += activeTabStyle.Render(tab)
		} else {
			feeds, _ := m.l.GetCategory(tab)
			hasUnread := false
			for _, f := range feeds {
				if f.HasUnread() {
					hasUnread = true
					break
				}
			}
			if hasUnread {
				renderedTabs += unreadTabStyle.Render(tab)
			} else {
				renderedTabs += inactiveTabStyle.Render(tab)
			}
		}
	}

	return renderedTabs
}

func renderedTitle(m *model) string {
	return titleStyle.Render(m.title)
}

func renderedStatus(m *model) string {
	return statusStyle.Render(m.status)
}

func openInBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

func (m *model) UpdateTitle(title string) {
	m.title = title
}

func (m *model) UpdateStatus(msg string) {
	m.status = msg

	if m.clearTimer != nil {
		m.clearTimer.Stop()
	}

	m.clearTimer = time.AfterFunc(8*time.Second, func() {
		m.prog.Send(statusClearMsg{})
	})
}

func (m *model) SaveState() error {
	cacheFilePath, err := rss.CacheFilePath()
	if err != nil {
		return err
	}

	f, err := os.Create(cacheFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return m.l.Save(f)
}

func BuildApp() {
	m := initialModel()
	p := tea.NewProgram(m)
	m.prog = p

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
