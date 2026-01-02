package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/browser"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case feedsLoadedMsg:
		m.feeds = msg.feeds
		m.buildArticles()
		m.loading = false
		m.statusMessage = ""
		// reset selections if out of bounds
		if m.selectedArticle >= len(m.allArticles) {
			m.selectedArticle = 0
		}
		if m.selectedSource >= len(m.feeds) {
			m.selectedSource = 0
		}
		return m, nil

	case articleContentLoadedMsg:
		m.loadingArticle = false
		if msg.err != nil {
			m.statusMessage = fmt.Sprintf("Error loading article: %v", msg.err)
		} else {
			m.cachedArticleURL = msg.url
			m.cachedArticleContent = msg.content
		}
		return m, nil

	case clearStatusMsg:
		m.statusMessage = ""
		return m, nil

	case tea.KeyMsg:
		if m.searchInputTrap {
			return m.handleSearchInput(msg)
		}

		// global keybindings (when not in search mode)
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "s":
			// toggle sources and articles view
			switch m.currentView {
			case articlesView:
				m.currentView = sourcesView
			case sourcesView:
				m.currentView = articlesView
			}
			return m, nil

		case "r":
			// refresh
			return m, m.RefreshFeeds()

		case "/":
			// enter search mode (only in articles view)
			if m.currentView == articlesView {
				m.searchInputTrap = true
				m.searchQuery = ""
				m.selectedArticle = 0
				return m, nil
			}

		case "esc":
			// clear search if search query is active
			if m.searchQuery != "" && m.currentView == articlesView {
				m.searchQuery = ""
				m.selectedArticle = 0
				return m, nil
			}
		}

		switch m.currentView {
		case articlesView:
			return m.handleArticlesViewKeys(msg)
		case articleView:
			return m.handleArticleViewKeys(msg)
		case sourcesView:
			return m.handleSourcesViewKeys(msg)
		}
	}

	return m, nil
}

func (m Model) handleArticlesViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visibleArticles := m.GetVisibleArticles()

	switch msg.String() {
	case "up", "k":
		if m.selectedArticle > 0 {
			m.selectedArticle--
		}

	case "down", "j":
		if len(visibleArticles) > 0 && m.selectedArticle < len(visibleArticles)-1 {
			m.selectedArticle++
		}

	case "m":
		// manually update read status
		m.ToggleCurrentArticleReadStatus()

	case "o":
		// open article in browser
		article := m.GetCurrentArticle()
		if article != nil && article.Link != "" {
			m.MarkCurrentArticleAsRead()
			err := browser.OpenURL(article.Link)
			if err != nil {
				return m, m.SetStatusMessage("Failed to open browser")
			}
			return m, m.SetStatusMessage("Opened in browser")
		}

	case "enter":
		// view article
		article := m.GetCurrentArticle()
		if article == nil {
			return m, m.SetStatusMessage("Could not get article")
		}

		m.currentView = articleView
		m.articleScroll = 0 // reset scroll when entering article

		m.MarkCurrentArticleAsRead()

		// check if article is cached
		if m.cachedArticleURL != article.Link {
			m.loadingArticle = true
			m.cachedArticleURL = ""
			m.cachedArticleContent = ""
			return m, loadArticleContent(article.Link)
		}
	}

	return m, nil
}

func (m Model) handleArticleViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// return to articles list
		m.currentView = articlesView
		return m, nil

	case "o":
		// open article in browser
		article := m.GetCurrentArticle()
		if article != nil && article.Link != "" {
			err := browser.OpenURL(article.Link)
			if err != nil {
				return m, m.SetStatusMessage("Failed to open browser")
			}
			return m, m.SetStatusMessage("Opened in browser")
		}

	case "up", "k":
		if m.articleScroll > 0 {
			m.articleScroll--
		}

	case "down", "j":
		m.articleScroll++

	case "pgup":
		// scroll up by page
		pageSize := m.height - 4 // account for header and status bar
		m.articleScroll -= pageSize
		if m.articleScroll < 0 {
			m.articleScroll = 0
		}

	case "pgdown":
		// scroll down by page
		pageSize := m.height - 4
		m.articleScroll += pageSize

	case "g":
		m.articleScroll = 0

	case "G":
		// scroll to bottom
		m.articleScroll = m.height
	}

	return m, nil
}

// handleSearchInput processes keyboard input in search mode
func (m Model) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		// exit search mode and clear query
		m.searchInputTrap = false
		m.searchQuery = ""
		m.selectedArticle = 0
		return m, nil

	case tea.KeyEnter:
		// exit search mode but keep query active
		m.searchInputTrap = false
		return m, nil

	case tea.KeyBackspace:
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.selectedArticle = 0 // reset selection when query changes
		} else {
			// if query is empty, exit search mode
			m.searchInputTrap = false
		}
		return m, nil

	case tea.KeyRunes:
		// add typed character to query
		m.searchQuery += string(msg.Runes)
		m.selectedArticle = 0 // reset selection when query changes
		return m, nil
	}

	return m, nil
}

func (m Model) handleSourcesViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedSource > 0 {
			m.selectedSource--
		}

	case "down", "j":
		if len(m.feeds) > 0 && m.selectedSource < len(m.feeds)-1 {
			m.selectedSource++
		}

	case "enter":
		// filter by selected source
		selectedFeed := m.GetCurrentSource()
		if selectedFeed != nil {
			m.filteredFeed = selectedFeed
			m.currentView = articlesView
			m.selectedArticle = 0
			// rebuild articles list
			m.allArticles = []articleWithSource{}
			for i := range selectedFeed.Articles {
				m.allArticles = append(m.allArticles, articleWithSource{
					article:   &selectedFeed.Articles[i],
					feedTitle: selectedFeed.Title,
				})
			}
			return m, m.SetStatusMessage("Filtered by: " + selectedFeed.Title)
		}

	case "a":
		// show all feeds
		if m.filteredFeed != nil {
			m.filteredFeed = nil
			m.buildArticles()
			m.selectedArticle = 0
			m.currentView = articlesView
			return m, m.SetStatusMessage("Showing all feeds")
		}
	}

	return m, nil
}
