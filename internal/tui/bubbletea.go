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

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type feedItem struct {
	title, desc string
	feed        *rss.Feed
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
	selectedFeed *rss.Feed
	feedsList    list.Model
	itemsList    list.Model
}

func initialModel() *model {
	filesystem := os.DirFS(".")
	feeds, err := rss.CreateFeedsFromFS(filesystem)
	if err != nil {
		log.Fatal(err)
	}

	var feedList rss.FeedList
	feedList.Add(feeds...)

	items := BuildFeedList(feedList.All)

	m := model{
		feedList:  feedList,
		feedsList: list.New(items, list.NewDefaultDelegate(), 0, 0),
		itemsList: list.New(nil, list.NewDefaultDelegate(), 0, 0),
	}
	m.feedsList.Title = "rssboat"

	return &m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "R":
			m.feedsList.NewStatusMessage("Updating feeds...")
			go func(m *model) {
				m.feedList.UpdateAll()
				list := BuildFeedList(m.feedList.All)
				m.feedsList.SetItems(list)
				m.feedsList.NewStatusMessage("Updated all.")
			}(m)
		case "enter":
			if m.selectedFeed == nil {
				if i, ok := m.feedsList.SelectedItem().(feedItem); ok {
					if i.feed.Feed != nil && i.feed.Error == "" {
						m.selectedFeed = i.feed
						items := BuildItemsList(m.selectedFeed)
						m.itemsList.Title = i.title
						m.itemsList.SetItems(items)
					}
				}
			} else {
				i, ok := m.itemsList.SelectedItem().(rssListItem)
				if ok {
					rssItem := i.item
					cmd := exec.Command("open", rssItem.Link)
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}
				}
			}
		case "b":
			m.selectedFeed = nil
		case "o":
			if m.selectedFeed == nil {
				i, ok := m.feedsList.SelectedItem().(feedItem)
				if ok {
					feed := i.feed
					if feed.Feed != nil {
						cmd := exec.Command("open", feed.Link)
						if err := cmd.Run(); err != nil {
							log.Fatal(err)
						}
					}
				}
			} else {
				i, ok := m.itemsList.SelectedItem().(rssListItem)
				if ok {
					rssItem := i.item
					cmd := exec.Command("open", rssItem.Link)
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}
				}
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

func BuildFeedList(feeds []*rss.Feed) []list.Item {
	var listItems []list.Item
	for _, feed := range feeds {
		title := feed.GetField("Title")
		description := feed.GetField("Latest")
		listItems = append(listItems, feedItem{
			title: title,
			desc:  description,
			feed:  feed,
		})
	}
	return listItems
}

func BuildItemsList(feed *rss.Feed) []list.Item {
	var listItems []list.Item
	for _, rssItem := range feed.Feed.Items {
		title := rssItem.Title
		description := rssItem.Description
		listItems = append(listItems, rssListItem{
			title: title,
			desc:  description,
			item:  &rss.RssItem{Item: rssItem},
		})
	}
	return listItems
}

func BuildApp() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
