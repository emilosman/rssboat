package tui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emilosman/rssboat/internal/rss"
)

func buildFeedList(feeds []*rss.RssFeed) []list.Item {
	var listItems []list.Item
	for _, feed := range feeds {
		title := feed.GetField("Title")
		description := feed.GetField("Latest")
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
		ri := &feed.RssItems[idx]
		title := ri.GetField("Title")
		description := ri.Item.Description

		listItems = append(listItems, rssListItem{
			title: title,
			desc:  description,
			item:  ri,
		})
	}
	return listItems
}

func openInBrowser(url string) error {
	cmd := exec.Command("open", url)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (m *model) SaveState() error {
	f, err := os.Create("./data.json")
	if err != nil {
		return err
	}
	defer f.Close()
	return m.feedList.Save(f)
}

func BuildApp() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
