package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/emilosman/rssboat/internal/rss"
	"github.com/muesli/reflow/wordwrap"
)

type keyHandler func(*model) tea.Cmd

var (
	feedKeyHandlers = map[string]keyHandler{
		"A":      handleMarkFeedRead,
		"C":      handleMarkAllFeedsRead,
		"E":      handleEdit,
		"h":      handlePrevTab,
		"left":   handlePrevTab,
		"l":      handleNextTab,
		"right":  handleNextTab,
		"n":      handleNextUnreadFeed,
		"o":      handleOpenFeed,
		"p":      handlePrevUnreadFeed,
		"r":      handleUpdateFeed,
		"R":      handleUpdateAllFeeds,
		"q":      handleQuit,
		"ctrl+a": handleMarkTabAsRead,
		"ctrl+c": handleInterrupt,
		"ctrl+r": handleTabUpdate,
		"enter":  handleEnterFeed,
		"esc":    handleQuit,
		"tab":    handleNextTab,
	}

	itemKeyHandlers = map[string]keyHandler{
		"a":     handleToggleRead,
		"A":     handleMarkItemsRead,
		"b":     handleBack,
		"n":     handleNextUnreadItem,
		"o":     handleOpenItem,
		"p":     handlePrevUnreadItem,
		"q":     handleBack,
		"esc":   handleBack,
		"r":     handleUpdateFeed,
		"R":     handleUpdateAllFeeds,
		"enter": handleViewItem,
	}

	viewKeyHandlers = map[string]keyHandler{
		"b":     handleBack,
		"l":     handleViewNext,
		"right": handleViewNext,
		"h":     handleViewPrev,
		"left":  handleViewPrev,
		"o":     handleOpenItem,
		"q":     handleBack,
		"esc":   handleBack,
	}
)

func handleEdit(m *model) tea.Cmd {
	configFilePath, err := rss.ConfigFilePath()
	if err != nil {
		fmt.Println("Error opening config dir", err)
		return nil
	}
	configFile := filepath.Join(configFilePath, "urls.yaml")

	editor := os.Getenv("EDITOR")
	if editor == "" {
		switch runtime.GOOS {
		case "windows":
			editor = "notepad"
		default:
			editor = "vi"
		}
	}

	cmd := exec.Command(editor, configFile)

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
		rebuildItemsList(m)
		m.li.NewStatusMessage(MsgItemReadToggled)
	}
	return nil
}

func handleNextUnreadItem(m *model) tea.Cmd {
	i, ok := m.li.SelectedItem().(rssListItem)
	if ok {
		prev := i.item
		index, next := m.f.NextUnreadItem(prev)
		if next != nil {
			m.li.Select(index)
		}
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

func handleMarkItemsRead(m *model) tea.Cmd {
	if m.f != nil {
		m.f.MarkAllItemsRead()
		rebuildItemsList(m)
		m.li.NewStatusMessage(MsgMarkFeedRead)
	}
	return nil
}

func handleMarkTabAsRead(m *model) tea.Cmd {
	feeds, err := m.l.GetCategory(activeTab(m.tabs, m.activeTab))
	if err != nil {
		m.lf.NewStatusMessage(err.Error())
	}

	rss.MarkFeedsAsRead(feeds...)
	rebuildFeedList(m)
	m.lf.NewStatusMessage(MsgMakrTabAsRead)

	return nil
}

func handleBack(m *model) tea.Cmd {
	if m.i != nil {
		m.i = nil
	} else {
		m.lf.ResetFilter()
		m.li.ResetFilter()
		rebuildFeedList(m)
		m.f = nil
	}
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
			rebuildItemsList(m)
		}
	}
	return nil
}

func handlePrevUnreadItem(m *model) tea.Cmd {
	i, ok := m.li.SelectedItem().(rssListItem)
	if ok {
		next := i.item
		index, prev := m.f.PrevUnreadItem(next)
		if prev != nil {
			m.li.Select(index)
		}
	}
	return nil
}

func handleUpdateFeed(m *model) tea.Cmd {
	if len(m.l.Feeds) == 0 {
		m.lf.NewStatusMessage(ErrUpdatingFeed)
		return nil
	}

	feed := m.f
	if m.f == nil {
		if i, ok := m.lf.SelectedItem().(feedItem); ok {
			feed = i.rssFeed
		}
	}

	message := fmt.Sprintf("%s %s", MsgUpdatingFeed, feed.Url)
	m.lf.NewStatusMessage(message)
	m.li.NewStatusMessage(message)

	return updateFeedCmd(m, feed)
}

func handleUpdateAllFeeds(m *model) tea.Cmd {
	m.lf.NewStatusMessage(MsgUpdatingAllFeeds)
	m.li.NewStatusMessage(MsgUpdatingAllFeeds)

	return updateAllFeedsCmd(m)
}

func handleTabUpdate(m *model) tea.Cmd {
	m.lf.NewStatusMessage(MsgUpdatingAllFeeds)
	m.li.NewStatusMessage(MsgUpdatingAllFeeds)

	return updateTabFeedsCmd(m)
}

func handleQuit(m *model) tea.Cmd {
	m.SaveState()
	return tea.Quit
}

func handleEnterFeed(m *model) tea.Cmd {
	if i, ok := m.lf.SelectedItem().(feedItem); ok {
		if i.rssFeed.Feed != nil {
			m.f = i.rssFeed
			m.li.Title = i.title
			m.li.Select(0)
			rebuildItemsList(m)
		}
	}
	return nil
}

func handleViewItem(m *model) tea.Cmd {
	i, ok := m.li.SelectedItem().(rssListItem)
	if ok {
		m.i = i.item
		if m.i.Item != nil {
			m.v.YOffset = 0
			m.v.SetContent(wordwrap.String(m.i.Content(), m.v.Width))
			m.i.MarkRead()
			rebuildItemsList(m)
		}
	}
	return nil
}

func handleInterrupt(m *model) tea.Cmd {
	m.SaveState()
	return tea.Quit
}

func handleTabNumber(m *model, i int) tea.Cmd {
	if i > len(m.tabs) {
		return nil
	}

	if i == 0 {
		return nil
	}

	m.activeTab = i - 1

	return rebuildFeedList(m)
}

func handleViewNext(m *model) tea.Cmd {
	index, next := m.f.NextAfter(m.i)
	if next != nil {
		m.i = next
		m.li.Select(index)
		m.v.YOffset = 0
		m.v.SetContent(wordwrap.String(next.Content(), m.v.Width))
		next.MarkRead()
		rebuildItemsList(m)
	}
	return nil
}

func handleViewPrev(m *model) tea.Cmd {
	index, prev := m.f.PrevBefore(m.i)
	if prev != nil {
		m.i = prev
		m.li.Select(index)
		m.v.YOffset = 0
		m.v.SetContent(wordwrap.String(prev.Content(), m.v.Width))
		prev.MarkRead()
		rebuildItemsList(m)
	}
	return nil
}

func handleNextUnreadFeed(m *model) tea.Cmd {
	i, ok := m.lf.SelectedItem().(feedItem)
	if ok {
		prev := i.rssFeed
		feeds, err := m.l.GetCategory(activeTab(m.tabs, m.activeTab))
		if err != nil {
			return nil
		}
		index, next := rss.NextUnreadFeed(feeds, prev)
		if next != nil {
			m.lf.Select(index)
		}
	}
	return nil
}

func handlePrevUnreadFeed(m *model) tea.Cmd {
	i, ok := m.lf.SelectedItem().(feedItem)
	if ok {
		next := i.rssFeed
		feeds, err := m.l.GetCategory(activeTab(m.tabs, m.activeTab))
		if err != nil {
			return nil
		}
		index, prev := rss.PrevUnreadFeed(feeds, next)
		if prev != nil {
			m.lf.Select(index)
		}
	}
	return nil
}
