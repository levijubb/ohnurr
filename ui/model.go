package ui

import (
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"ohnurr/config"
	"ohnurr/content"
	"ohnurr/rss"
)

type viewMode int

const (
	articlesView viewMode = iota
	articleView
	sourcesView
)

type Model struct {
	config               *config.Config
	state                *config.State
	feeds                []*rss.Feed
	allArticles          []articleWithSource
	selectedArticle      int
	selectedSource       int
	currentView          viewMode
	filteredFeed         *rss.Feed // nil == show all feeds
	searchInputTrap      bool
	searchQuery          string
	width                int
	height               int
	loading              bool
	statusMessage        string
	articleScroll        int // scroll position in article view
	cachedArticleURL     string
	cachedArticleContent string // TODO: cache more than one article at a time
	loadingArticle       bool
}

// combines an article with its source feed info
type articleWithSource struct {
	article   *rss.Article
	feedTitle string
}

type feedsLoadedMsg struct {
	feeds []*rss.Feed
}

type articleContentLoadedMsg struct {
	url     string
	content string
	err     error
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
		searchInputTrap: false,
		searchQuery:     "",
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

// creates a command to fetch article content
func loadArticleContent(url string) tea.Cmd {
	return func() tea.Msg {
		articleContent, err := content.GetArticleContent(url)
		return articleContentLoadedMsg{
			url:     url,
			content: articleContent,
			err:     err,
		}
	}
}

// creates a sorted list of all articles from all feeds
func (m *Model) buildArticles() {
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

// returns articles filtered by search query if active
func (m Model) GetVisibleArticles() []articleWithSource {
	if m.searchQuery != "" {
		return FilterArticles(m.allArticles, m.searchQuery)
	}
	return m.allArticles
}

func (m Model) GetCurrentArticle() *rss.Article {
	visibleArticles := m.GetVisibleArticles()
	if len(visibleArticles) == 0 || m.selectedArticle >= len(visibleArticles) {
		return nil
	}
	return visibleArticles[m.selectedArticle].article
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

	if article == nil {
		return
	}

	m.state.MarkAsRead(article.GetArticleID())
	_ = m.state.Save()
}

func (m *Model) MarkCurrentArticleAsUnread() {
	article := m.GetCurrentArticle()
	if article == nil {
		return
	}

	m.state.UnmarkAsRead(article.GetArticleID())
	_ = m.state.Save()
}

func (m *Model) ToggleCurrentArticleReadStatus() {
	article := m.GetCurrentArticle()

	if m.IsArticleRead(article) {
		m.MarkCurrentArticleAsUnread()
	} else {
		m.MarkCurrentArticleAsRead()
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
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

type clearStatusMsg struct{}
