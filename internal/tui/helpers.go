package tui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emilosman/rssboat/internal/rss"
)

var (
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ff4fff")).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#767676")).
				Padding(0, 1)
)

type feedUpdatedMsg struct {
	Feed *rss.RssFeed
	Err  error
}

func updateAllFeedsCmd(m *model) tea.Cmd {
	return func() tea.Msg {
		results, err := m.l.UpdateAllAsync()
		if err != nil {
			return feedUpdatedMsg{Feed: nil, Err: err}
		}

		go func() {
			for range m.l.Feeds {
				res := <-results
				m.prog.Send(feedUpdatedMsg{Feed: res.Feed, Err: res.Err})
			}
		}()

		return MsgUpdatingAllFeeds
	}
}

func rebuildFeedList(m *model) tea.Cmd {
	items := buildFeedList(m.l, m.tabs, m.activeTab)
	m.lf.SetItems(items)
	m.lf.Title = m.tabs[m.activeTab]
	return nil
}

func buildFeedList(l *rss.List, tabs []string, activeTab int) []list.Item {
	category := tabs[activeTab]
	feeds, err := l.GetCategory(category)
	if err != nil {
		feeds = l.Feeds
	}

	var listItems []list.Item
	for _, feed := range feeds {
		title := feed.Title()
		description := feed.Latest()
		listItems = append(listItems, feedItem{
			title:   title,
			desc:    description,
			rssFeed: feed,
		})
	}
	return listItems
}

func buildItemsList(feed *rss.RssFeed) []list.Item {
	listItems := make([]list.Item, 0, len(feed.RssItems))
	for idx := range feed.RssItems {
		ri := feed.RssItems[idx]
		title := ri.Title()
		description := ri.Item.Description
		listItems = append(listItems, rssListItem{
			title: title,
			desc:  description,
			item:  ri,
		})
	}
	return listItems
}

func getTabs(l *rss.List) []string {
	var tabs []string
	categories, err := l.GetAllCategories()
	if err != nil {
		return tabs
	}

	for category := range categories {
		tabs = append(tabs, category)
	}

	return tabs
}

func renderedTabs(m *model) string {
	var renderedTabs string
	for i, tab := range m.tabs {
		if i == m.activeTab {
			renderedTabs += activeTabStyle.Render(tab)
		} else {
			renderedTabs += inactiveTabStyle.Render(tab)
		}
	}

	return fmt.Sprintf("%s\n", renderedTabs)
}

func openInBrowser(url string) error {
	var browserCmd = map[string][]string{
		"darwin":  {"open"},
		"linux":   {"xdg-open"},
		"windows": {"rundll32", "url.dll,FileProtocolHandler"},
	}

	openCmd, ok := browserCmd[runtime.GOOS]
	if !ok {
		return fmt.Errorf("Unsuported platform: %s", runtime.GOOS)
	}

	cmd := exec.Command(openCmd[0], url)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
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
