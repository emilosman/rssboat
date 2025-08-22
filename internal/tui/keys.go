package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type keyHandler func(*model) tea.Cmd

var (
	feedKeyHandlers = map[string]keyHandler{
		"A":      handleMarkFeedRead,
		"b":      handleBack,
		"C":      handleMarkAllFeedsRead,
		"o":      handleOpenFeed,
		"r":      handleUpdateFeed,
		"R":      handleUpdateAllFeeds,
		"q":      handleQuit,
		"enter":  handleEnterFeed,
		"esc":    handleQuit,
		"ctrl+c": handleInterrupt,
	}

	itemKeyHandlers = map[string]keyHandler{
		"a":     handleToggleRead,
		"o":     handleOpenItem,
		"q":     handleBack,
		"r":     handleUpdateFeed,
		"R":     handleUpdateAllFeeds,
		"enter": handleEnterItem,
	}
)

func handleToggleRead(m *model) tea.Cmd {
	i, ok := m.itemsList.SelectedItem().(rssListItem)
	if ok {
		i.item.ToggleRead()
		items := buildItemsList(m.selectedFeed)
		m.itemsList.SetItems(items)
		m.itemsList.NewStatusMessage(MsgItemReadToggled)
	}
	return nil
}

func handleMarkFeedRead(m *model) tea.Cmd {
	if i, ok := m.feedsList.SelectedItem().(feedItem); ok {
		f := i.rssFeed
		f.MarkAllItemsRead()
		all := buildFeedList(m.feedList.All)
		m.feedsList.SetItems(all)
		m.feedsList.NewStatusMessage(MsgMarkFeedRead)
	}
	return nil
}

func handleBack(m *model) tea.Cmd {
	all := buildFeedList(m.feedList.All)
	m.feedsList.SetItems(all)
	m.selectedFeed = nil
	return nil
}

func handleMarkAllFeedsRead(m *model) tea.Cmd {
	m.feedList.MarkAllFeedsRead()
	all := buildFeedList(m.feedList.All)
	m.feedsList.SetItems(all)
	m.feedsList.NewStatusMessage(MsgMarkAllFeedsRead)
	m.SaveState()
	return nil
}

func handleOpenFeed(m *model) tea.Cmd {
	i, ok := m.feedsList.SelectedItem().(feedItem)
	if ok {
		rssFeed := i.rssFeed
		if rssFeed.Feed != nil {
			err := openInBrowser(rssFeed.Feed.Link)
			if err != nil {
				errorMessage := fmt.Sprintf("Error opening item, %q", err)
				m.itemsList.NewStatusMessage(errorMessage)
			}
		}
	}
	return nil
}

func handleOpenItem(m *model) tea.Cmd {
	i, ok := m.itemsList.SelectedItem().(rssListItem)
	if ok {
		rssItem := i.item
		if rssItem.Item != nil {
			err := openInBrowser(rssItem.Item.Link)
			if err != nil {
				errorMessage := fmt.Sprintf("Error opening item, %q", err)
				m.itemsList.NewStatusMessage(errorMessage)
			}
			rssItem.ToggleRead()
			items := buildItemsList(m.selectedFeed)
			m.itemsList.SetItems(items)
		}
	}
	return nil
}

func handleUpdateFeed(m *model) tea.Cmd {
	if i, ok := m.feedsList.SelectedItem().(feedItem); ok {
		f := i.rssFeed
		statusMsg := fmt.Sprintf("Updating feed %s", f.Url)
		m.feedsList.NewStatusMessage(statusMsg)
		go func(m *model) {
			err := f.GetFeed()
			if err != nil {
				m.feedsList.NewStatusMessage(ErrUpdatingFeed)
			}
			all := buildFeedList(m.feedList.All)
			m.feedsList.SetItems(all)
			m.feedsList.NewStatusMessage(MsgFeedUpdated)
			m.SaveState()
		}(m)
	}
	return nil
}

func handleUpdateAllFeeds(m *model) tea.Cmd {
	m.feedsList.NewStatusMessage(MsgUpdatingAllFeeds)
	m.itemsList.NewStatusMessage(MsgUpdatingAllFeeds)
	go func(m *model) {
		err := m.feedList.UpdateAll()
		if err != nil {
			m.feedsList.NewStatusMessage(ErrUpdatingFeeds)
		}
		all := buildFeedList(m.feedList.All)
		m.feedsList.SetItems(all)
		m.feedsList.NewStatusMessage(MsgAllFeedsUpdated)
		m.itemsList.NewStatusMessage(MsgAllFeedsUpdated)
		m.SaveState()
	}(m)
	return nil
}

func handleQuit(m *model) tea.Cmd {
	m.SaveState()
	return tea.Quit
}

func handleEnterFeed(m *model) tea.Cmd {
	if m.feedsList.FilterState().String() != "filtering" {
		if i, ok := m.feedsList.SelectedItem().(feedItem); ok {
			if i.rssFeed.Feed != nil && i.rssFeed.Error == "" {
				m.selectedFeed = i.rssFeed
				items := buildItemsList(m.selectedFeed)
				m.itemsList.Title = i.title
				m.itemsList.SetItems(items)
			}
		}
	}
	return nil
}

func handleEnterItem(m *model) tea.Cmd {
	if m.itemsList.FilterState().String() != "filtering" {
		i, ok := m.itemsList.SelectedItem().(rssListItem)
		if ok {
			rssItem := i.item
			if rssItem.Item != nil {
				err := openInBrowser(rssItem.Item.Link)
				if err != nil {
					errorMessage := fmt.Sprintf("Error opening item, %q", err)
					m.itemsList.NewStatusMessage(errorMessage)
				}
				rssItem.ToggleRead()
				items := buildItemsList(m.selectedFeed)
				m.itemsList.SetItems(items)
			}
		}
	}
	return nil
}

func handleInterrupt(m *model) tea.Cmd {
	m.SaveState()
	return tea.Quit
}
