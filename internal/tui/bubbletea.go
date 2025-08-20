package tui

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emilosman/rssboat/internal/rss"
)

var (
	MsgUpdatingAllFeeds = "Updating all feeds..."
	MsgAllFeedsUpdated  = "All feeds updated."
	docStyle            = lipgloss.NewStyle().Margin(1, 2)
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
		f.Close()
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
		switch msg.String() {
		case "A":
			if m.selectedFeed == nil {
				if i, ok := m.feedsList.SelectedItem().(feedItem); ok {
					f := i.rssFeed
					f.MarkAllItemsRead()
					all := buildFeedList(m.feedList.All)
					m.feedsList.SetItems(all)
					m.feedsList.NewStatusMessage("Marked all feed items read")
				}
			}
		case "a":
			if m.selectedFeed != nil {
				i, ok := m.itemsList.SelectedItem().(rssListItem)
				if ok {
					rssItem := i.item
					rssItem.ToggleRead()
					items := buildItemsList(m.selectedFeed)
					m.itemsList.SetItems(items)
					m.itemsList.NewStatusMessage("Item read state toggled")
				}
			}
		case "r":
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
					}(m)
				}
			}
		case "R":
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
				f, _ := os.Create("./data.json")
				defer f.Close()
				m.feedList.Save(f)
			}(m)
		case "enter":
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
						}
					}
				}
			}
		case "b":
			all := buildFeedList(m.feedList.All)
			m.feedsList.SetItems(all)
			m.selectedFeed = nil
		case "o":
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
					}
				}
			}
		case "q", "esc":
			if m.selectedFeed != nil {
				all := buildFeedList(m.feedList.All)
				m.feedsList.SetItems(all)
				m.selectedFeed = nil
			} else {
				return m, tea.Quit
			}
		case "ctrl+c":
			return m, tea.Quit
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

func (m *model) View() string {
	if m.selectedFeed != nil {
		return docStyle.Render(m.itemsList.View())
	}
	return docStyle.Render(m.feedsList.View())
}

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
	var listItems []list.Item
	for _, rssItem := range feed.RssItems {
		title := rssItem.GetField("Title")
		description := rssItem.Item.Description
		listItems = append(listItems, rssListItem{
			title: title,
			desc:  description,
			item:  &rssItem,
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

func BuildApp() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
