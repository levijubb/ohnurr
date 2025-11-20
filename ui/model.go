package ui

import (
	"ohnurr/config"
	"ohnurr/rss"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type viewMode int

const (
	articlesView viewMode = iota
	sourcesView
)

type Model struct {
	config          *config.Config
	state           *config.State
	feeds           []*rss.Feed
	allArticles     []articleWithSource
	selectedArticle int
	selectedSource  int
	currentView     viewMode
	filteredFeed    *rss.Feed // nil means show all feeds
	width           int
	height          int
	loading         bool
	statusMessage   string
}

// combines an article with its source feed info
type articleWithSource struct {
	article   *rss.Article
	feedTitle string
}

type feedsLoadedMsg struct {
	feeds []*rss.Feed
}

func NewModel(cfg *config.Config, state *config.State) Model {
	return Model{
		config:          cfg,
		state:           state,
		feeds:           []*rss.Feed{},
		allArticles:     []articleWithSource{},
		selectedArticle: 0,
		selectedSource:  0,
		currentView:     articlesView,
		filteredFeed:    nil,
		loading:         true,
		statusMessage:   "Loading feeds...",
	}
}

func (m Model) Init() tea.Cmd {
	return loadFeeds(m.config.Feeds)
}

// creates a command to fetch all RSS feeds
func loadFeeds(urls []string) tea.Cmd {
	return func() tea.Msg {
		feeds := rss.FetchAllFeeds(urls)
		return feedsLoadedMsg{feeds: feeds}
	}
}

// creates a sorted list of all articles from all feeds
func (m *Model) buildAllArticles() {
	m.allArticles = []articleWithSource{}

	for _, feed := range m.feeds {
		if feed.Error != nil {
			continue
		}
		for i := range feed.Articles {
			m.allArticles = append(m.allArticles, articleWithSource{
				article:   &feed.Articles[i],
				feedTitle: feed.Title,
			})
		}
	}

	// sort by publication date
	sort.Slice(m.allArticles, func(i, j int) bool {
		return m.allArticles[i].article.Published.After(m.allArticles[j].article.Published)
	})
}

func (m Model) GetCurrentArticle() *rss.Article {
	if len(m.allArticles) == 0 || m.selectedArticle >= len(m.allArticles) {
		return nil
	}
	return m.allArticles[m.selectedArticle].article
}

// returns the currently selected feed in sources view
func (m Model) GetCurrentSource() *rss.Feed {
	if len(m.feeds) == 0 || m.selectedSource >= len(m.feeds) {
		return nil
	}
	return m.feeds[m.selectedSource]
}

func (m Model) IsArticleRead(article *rss.Article) bool {
	if article == nil {
		return false
	}
	return m.state.IsRead(article.GetArticleID())
}

func (m *Model) MarkCurrentArticleAsRead() {
	article := m.GetCurrentArticle()
	if article != nil {
		m.state.MarkAsRead(article.GetArticleID())
		m.state.Save()
		m.statusMessage = "Marked as read"
	}
}

// returns the number of unread articles in a feed
func (m Model) GetUnreadCount(feed *rss.Feed) int {
	if feed == nil {
		return 0
	}
	count := 0
	for _, article := range feed.Articles {
		if !m.state.IsRead(article.GetArticleID()) {
			count++
		}
	}
	return count
}

func (m *Model) RefreshFeeds() tea.Cmd {
	m.loading = true
	m.statusMessage = "Refreshing feeds..."
	m.selectedArticle = 0
	return loadFeeds(m.config.Feeds)
}

// sets a temporary status message
func (m *Model) SetStatusMessage(msg string) tea.Cmd {
	m.statusMessage = msg
	return tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

type clearStatusMsg struct{}
