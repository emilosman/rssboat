package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emilosman/rssboat/internal/rss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
	feed        rss.Feed
}

func (i item) Title() string       { return i.title }
func (i item) Feed() rss.Feed      { return i.feed }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	feeds        rss.FeedList
	selectedFeed *rss.Feed
	list         list.Model
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
		feeds:     feedList,
		list:      list.New(items, list.NewDefaultDelegate(), 0, 0),
		itemsList: list.New(nil, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "rssboat"

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
			m.feeds.UpdateAll()
			list := BuildFeedList(m.feeds.All)
			m.list.SetItems(list)
		case "enter", "l":
			if m.selectedFeed == nil {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.selectedFeed = &i.feed
					list := BuildItemsList(m.selectedFeed)
					m.itemsList.Title = i.title
					m.itemsList.SetItems(list)
				}
			}
		case "h":
			m.selectedFeed = nil
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.selectedFeed != nil {
		return docStyle.Render(m.itemsList.View())
	}
	return docStyle.Render(m.list.View())
}

func BuildFeedList(feeds []*rss.Feed) []list.Item {
	var listItems []list.Item
	for _, feed := range feeds {
		title := feed.GetFields("Title")[0]
		description := feed.GetFields("Latest")[0]
		item := item{title: title, desc: description, feed: *feed}
		listItems = append(listItems, item)
	}
	return listItems
}

func BuildItemsList(feed *rss.Feed) []list.Item {
	var listItems []list.Item
	for _, rssItem := range feed.Feed.Items {
		item := item{title: rssItem.Title}
		listItems = append(listItems, item)
	}
	return listItems
}

func BuildApp() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
