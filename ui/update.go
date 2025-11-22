package ui

import (
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
		m.buildAllArticles()
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

	case clearStatusMsg:
		m.statusMessage = ""
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "s":
			// toggle sources and articles view
			if m.currentView == articlesView {
				m.currentView = sourcesView
			} else {
				m.currentView = articlesView
			}
			return m, nil

		case "r":
			// refresh
			return m, m.RefreshFeeds()

		case "?":
			// help
			return m, m.SetStatusMessage("↑/k:up ↓/j:down s:sources m:mark-read o:open r:refresh ?:help q:quit")
		}

		if m.currentView == articlesView {
			return m.handleArticlesViewKeys(msg)
		} else {
			return m.handleSourcesViewKeys(msg)
		}
	}

	return m, nil
}

func (m Model) handleArticlesViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedArticle > 0 {
			m.selectedArticle--
		}

	case "down", "j":
		if len(m.allArticles) > 0 && m.selectedArticle < len(m.allArticles)-1 {
			m.selectedArticle++
		}

	case "m":
		// manually update read status
		m.ToggleCurrentArticleReadStatus()

	case "o", "enter":
		// open article in browser
		// TODO: Implement scraping and reading within the TUI on enter press
		article := m.GetCurrentArticle()
		if article != nil && article.Link != "" {
			m.MarkCurrentArticleAsRead()
			err := browser.OpenURL(article.Link)
			if err != nil {
				return m, m.SetStatusMessage("Failed to open browser")
			}
			return m, m.SetStatusMessage("Opened in browser")
		}
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
			m.buildAllArticles()
			m.selectedArticle = 0
			m.currentView = articlesView
			return m, m.SetStatusMessage("Showing all feeds")
		}
	}

	return m, nil
}
