package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/emilosman/rssboat/internal/rss"
)

type keyHandler func(*model) tea.Cmd

var (
	feedKeyHandlers = map[string]keyHandler{
		"A":      handleMarkFeedRead,
		"C":      handleMarkAllFeedsRead,
		"E":      handleEdit,
		"h":      handlePrevTab,
		"l":      handleNextTab,
		"o":      handleOpenFeed,
		"r":      handleUpdateFeed,
		"R":      handleUpdateAllFeeds,
		"q":      handleQuit,
		"enter":  handleEnterFeed,
		"esc":    handleQuit,
		"ctrl+c": handleInterrupt,
		"tab":    handleNextTab,
	}

	itemKeyHandlers = map[string]keyHandler{
		"a":     handleToggleRead,
		"o":     handleOpenItem,
		"b":     handleBack,
		"q":     handleBack,
		"r":     handleUpdateFeed,
		"R":     handleUpdateAllFeeds,
		"enter": handleEnterItem,
	}
)

func handleEdit(m *model) tea.Cmd {
	configFilePath, err := rss.ConfigFilePath()
	if err != nil {
		fmt.Println("Error opening config dir", err)
		return nil
	}
	configFile := filepath.Join(configFilePath, "urls.yaml")
	cmd := exec.Command("vi", configFile)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		m.lf.NewStatusMessage(err.Error())
		return nil
	}

	filesystem := os.DirFS(configFilePath)
	l, err := rss.LoadList(filesystem)
	if err != nil {
		m.lf.NewStatusMessage(err.Error())
		return nil
	}

	m.l = l
	m.activeTab = 0
	m.tabs = getTabs(l)
	m.lf.NewStatusMessage("URLs file edited")

	return rebuildFeedList(m)
}

func handleNextTab(m *model) tea.Cmd {
	m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
	return rebuildFeedList(m)
}

func handlePrevTab(m *model) tea.Cmd {
	m.activeTab = max(m.activeTab-1, 0)
	return rebuildFeedList(m)
}

func handleToggleRead(m *model) tea.Cmd {
	i, ok := m.li.SelectedItem().(rssListItem)
	if ok {
		i.item.ToggleRead()
		items := buildItemsList(m.selectedFeed)
		m.li.SetItems(items)
		m.li.NewStatusMessage(MsgItemReadToggled)
	}
	return nil
}

func handleMarkFeedRead(m *model) tea.Cmd {
	if i, ok := m.lf.SelectedItem().(feedItem); ok {
		f := i.rssFeed
		f.MarkAllItemsRead()
		rebuildFeedList(m)
		m.lf.NewStatusMessage(MsgMarkFeedRead)
	}
	return nil
}

func handleBack(m *model) tea.Cmd {
	rebuildFeedList(m)
	m.lf.ResetFilter()
	m.selectedFeed = nil
	return nil
}

func handleMarkAllFeedsRead(m *model) tea.Cmd {
	m.l.MarkAllFeedsRead()
	rebuildFeedList(m)
	m.lf.NewStatusMessage(MsgMarkAllFeedsRead)
	m.SaveState()
	return nil
}

func handleOpenFeed(m *model) tea.Cmd {
	i, ok := m.lf.SelectedItem().(feedItem)
	if ok {
		rssFeed := i.rssFeed
		if rssFeed.Feed != nil {
			err := openInBrowser(rssFeed.Feed.Link)
			if err != nil {
				errorMessage := fmt.Sprintf("Error opening item, %q", err)
				m.li.NewStatusMessage(errorMessage)
			}
		}
	}
	return nil
}

func handleOpenItem(m *model) tea.Cmd {
	i, ok := m.li.SelectedItem().(rssListItem)
	if ok {
		rssItem := i.item
		if rssItem.Item != nil {
			err := openInBrowser(rssItem.Link())
			if err != nil {
				errorMessage := fmt.Sprintf("Error opening item, %q", err)
				m.li.NewStatusMessage(errorMessage)
			}
			rssItem.MarkRead()
			items := buildItemsList(m.selectedFeed)
			m.li.SetItems(items)
		}
	}
	return nil
}

func handleUpdateFeed(m *model) tea.Cmd {
	if i, ok := m.lf.SelectedItem().(feedItem); ok {
		f := i.rssFeed
		statusMsg := fmt.Sprintf("Updating feed %s", f.Url)
		m.lf.NewStatusMessage(statusMsg)
		go func(m *model) {
			err := f.GetFeed()
			if err != nil {
				m.lf.NewStatusMessage(ErrUpdatingFeed)
			}
			rebuildFeedList(m)
			m.lf.NewStatusMessage(MsgFeedUpdated)
			m.SaveState()
		}(m)
	}
	return nil
}

func handleUpdateAllFeeds(m *model) tea.Cmd {
	m.lf.NewStatusMessage(MsgUpdatingAllFeeds)
	m.li.NewStatusMessage(MsgUpdatingAllFeeds)

	return updateAllFeedsCmd(m)
}

func handleQuit(m *model) tea.Cmd {
	m.SaveState()
	return tea.Quit
}

func handleEnterFeed(m *model) tea.Cmd {
	if m.lf.FilterState().String() != "filtering" {
		if i, ok := m.lf.SelectedItem().(feedItem); ok {
			if i.rssFeed.Feed != nil && i.rssFeed.Error == "" {
				m.selectedFeed = i.rssFeed
				items := buildItemsList(m.selectedFeed)
				m.li.Title = i.title
				m.li.SetItems(items)
			}
		}
	}
	return nil
}

func handleEnterItem(m *model) tea.Cmd {
	if m.li.FilterState().String() != "filtering" {
		i, ok := m.li.SelectedItem().(rssListItem)
		if ok {
			rssItem := i.item
			if rssItem.Item != nil {
				err := openInBrowser(rssItem.Link())
				if err != nil {
					errorMessage := fmt.Sprintf("Error opening item, %q", err)
					m.li.NewStatusMessage(errorMessage)
				}
				rssItem.MarkRead()
				items := buildItemsList(m.selectedFeed)
				m.li.SetItems(items)
			}
		}
	}
	return nil
}

func handleInterrupt(m *model) tea.Cmd {
	m.SaveState()
	return tea.Quit
}
