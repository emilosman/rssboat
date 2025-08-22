package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type keyHandler func(*model) tea.Cmd

var keyHandlers = map[string]keyHandler{
	"a":      handleToggleItemRead,
	"A":      handleMarkAllFeedRead,
	"b":      handleBack,
	"C":      handleMarkAllFeedsRead,
	"o":      handleOpen,
	"r":      handleUpdateFeed,
	"R":      handleUpdateAllFeeds,
	"q":      handleQuit,
	"enter":  handleEnter,
	"esc":    handleQuit,
	"ctrl+c": handleInterrupt,
}

func handleToggleItemRead(m *model) tea.Cmd {
	if m.selectedFeed != nil {
		i, ok := m.itemsList.SelectedItem().(rssListItem)
		if ok {
			i.item.ToggleRead()
			items := buildItemsList(m.selectedFeed)
			m.itemsList.SetItems(items)
			m.itemsList.NewStatusMessage("Item read state toggled")
		}
	}
	return nil
}

func handleMarkAllFeedRead(m *model) tea.Cmd {
	if m.selectedFeed == nil {
		if i, ok := m.feedsList.SelectedItem().(feedItem); ok {
			f := i.rssFeed
			f.MarkAllItemsRead()
			all := buildFeedList(m.feedList.All)
			m.feedsList.SetItems(all)
			m.feedsList.NewStatusMessage("Marked all feed items read")
		}
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
	m.feedsList.NewStatusMessage("All feeds marked read")
	m.SaveState()
	return nil
}

func handleOpen(m *model) tea.Cmd {
	if m.selectedFeed == nil {
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
	} else {
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

func handleUpdateFeed(m *model) tea.Cmd {
	if m.selectedFeed == nil {
		if i, ok := m.feedsList.SelectedItem().(feedItem); ok {
			f := i.rssFeed
			statusMsg := fmt.Sprintf("Updating feed %s", f.Url)
			m.feedsList.NewStatusMessage(statusMsg)
			go func(m *model) {
				err := f.GetFeed()
				if err != nil {
					m.feedsList.NewStatusMessage("Error updating feed")
				}
				all := buildFeedList(m.feedList.All)
				m.feedsList.SetItems(all)
				m.feedsList.NewStatusMessage("Feed updated")
				m.SaveState()
			}(m)
		}
	}
	return nil
}

func handleUpdateAllFeeds(m *model) tea.Cmd {
	m.feedsList.NewStatusMessage(MsgUpdatingAllFeeds)
	m.itemsList.NewStatusMessage(MsgUpdatingAllFeeds)
	go func(m *model) {
		err := m.feedList.UpdateAll()
		if err != nil {
			m.feedsList.NewStatusMessage("Error updating feeds")
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
	if m.selectedFeed != nil {
		all := buildFeedList(m.feedList.All)
		m.feedsList.SetItems(all)
		m.selectedFeed = nil
	} else {
		m.SaveState()
	}
	return tea.Quit
}

func handleEnter(m *model) tea.Cmd {
	if m.selectedFeed == nil && m.feedsList.FilterState().String() != "filtering" {
		if i, ok := m.feedsList.SelectedItem().(feedItem); ok {
			if i.rssFeed.Feed != nil && i.rssFeed.Error == "" {
				m.selectedFeed = i.rssFeed
				items := buildItemsList(m.selectedFeed)
				m.itemsList.Title = i.title
				m.itemsList.SetItems(items)
			}
		}
	} else {
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
	}
	return nil
}

func handleInterrupt(m *model) tea.Cmd {
	m.SaveState()
	return tea.Quit
}
