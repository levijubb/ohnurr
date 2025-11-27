package ui

import (
	"fmt"
	"strings"

	"ohnurr/content"

	lg "github.com/charmbracelet/lipgloss"
)

var (
	// color
	activeColor   = lg.Color("205")
	inactiveColor = lg.Color("240")
	unreadColor   = lg.Color("81")
	titleColor    = lg.Color("86")
	accentColor   = lg.Color("147")

	// style
	articleTitleStyle = lg.NewStyle().
				Foreground(titleColor).
				Bold(true)

	articleTitleReadStyle = lg.NewStyle().
				Foreground(inactiveColor)

	selectedStyle = lg.NewStyle().
			Foreground(activeColor).
			Bold(true)

	dimStyle = lg.NewStyle().
			Foreground(inactiveColor)

	descriptionStyle = lg.NewStyle().
				Foreground(lg.Color("252"))

	sourceStyle = lg.NewStyle().
			Foreground(accentColor).
			Italic(true)

	headerStyle = lg.NewStyle().
			Foreground(titleColor).
			Bold(true).
			Padding(0, 1)

	statusStyle = lg.NewStyle().
			Foreground(lg.Color("205")).
			Padding(0, 1)

	unreadDotStyle = lg.NewStyle().
			Foreground(unreadColor).
			Bold(true)
)

func (m Model) View() string {
	if m.loading {
		return lg.Place(
			m.width, m.height,
			lg.Center, lg.Center,
			"Loading feeds...",
		)
	}

	if len(m.feeds) == 0 {
		return lg.Place(
			m.width, m.height,
			lg.Center, lg.Center,
			"No feeds loaded\nAdd feeds with: ohnurr add <url>",
		)
	}

	var content string
	switch m.currentView {
	case articlesView:
		content = m.renderArticlesView()
	case sourcesView:
		content = m.renderSourcesView()
	case articleView:
		content = m.renderArticleView()
	}

	statusBar := m.renderStatusBar()

	return lg.JoinVertical(lg.Left, content, statusBar)
}

func (m Model) renderArticlesView() string {
	var lines []string

	// header
	headerText := "📰 Articles"
	if m.filteredFeed != nil {
		headerText = fmt.Sprintf("📰 %s", m.filteredFeed.Title)
	}
	if m.searchInputTrap || m.searchQuery != "" {
		headerText += " " + dimStyle.Render(fmt.Sprintf("[search: %s", m.searchQuery))
		if m.searchInputTrap {
			headerText += selectedStyle.Render("_")
		}

		headerText += dimStyle.Render("]")
	}
	lines = append(lines, headerStyle.Render(headerText))
	lines = append(lines, "")

	// calc available height for articles
	// reserve space for header (2 lines) and status bar (1 line)
	availableHeight := m.height - 3

	visibleArticles := m.GetVisibleArticles()

	if len(visibleArticles) == 0 {
		if m.searchQuery != "" {
			lines = append(lines, dimStyle.Render("No articles match your search"))
		} else {
			lines = append(lines, dimStyle.Render("No articles available"))
		}
	} else {
		linesPerArticle := 4

		// determine articles to show based on scroll pos
		startIdx := min(m.selectedArticle, len(visibleArticles)-1)

		// try keep selected article in middle of view
		if startIdx > availableHeight/(linesPerArticle*2) {
			startIdx = startIdx - availableHeight/(linesPerArticle*2)
		} else {
			startIdx = 0
		}

		lineCount := 0
		for i := startIdx; i < len(visibleArticles) && lineCount < availableHeight-2; i++ {
			item := visibleArticles[i]
			article := item.article
			isRead := m.IsArticleRead(article)
			isSelected := i == m.selectedArticle

			// status indicator and title
			var titleLine string
			indicator := unreadDotStyle.Render("●")
			if isRead {
				indicator = dimStyle.Render("○")
			}

			titleText := article.Title
			if len(titleText) > m.width-6 {
				titleText = titleText[:m.width-9] + "..."
			}

			if isSelected {
				titleLine = selectedStyle.Render("▶ ") + indicator + " "
				if isRead {
					titleLine += articleTitleReadStyle.Render(titleText)
				} else {
					titleLine += articleTitleStyle.Render(titleText)
				}
			} else {
				titleLine = "  " + indicator + " "
				if isRead {
					titleLine += articleTitleReadStyle.Render(titleText)
				} else {
					titleLine += titleText
				}
			}

			lines = append(lines, titleLine)
			lineCount++

			// description
			if article.Description != "" && lineCount < availableHeight-2 {
				desc := article.Description
				maxDescWidth := m.width - 6
				if len(desc) > maxDescWidth {
					desc = desc[:maxDescWidth-3] + "..."
				}

				// skip description if it's too short
				if (len(strings.Split(desc, " "))) > 1 {
					descLine := "    " + descriptionStyle.Render(desc)
					lines = append(lines, descLine)
					lineCount++
				}
			}

			// source
			if lineCount < availableHeight-2 {
				sourceLine := "    " + sourceStyle.Render("from "+item.feedTitle)
				lines = append(lines, sourceLine)
				lineCount++
			}

			// blank lines between articles
			if lineCount < availableHeight-2 && i < len(visibleArticles)-1 {
				lines = append(lines, "")
				lineCount++
			}
		}
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderSourcesView() string {
	var lines []string

	// header
	lines = append(lines, headerStyle.Render("📚 Sources"))
	lines = append(lines, "")

	if m.filteredFeed != nil {
		lines = append(lines, dimStyle.Render("Press 'a' to show all feeds"))
		lines = append(lines, "")
	}

	for i, feed := range m.feeds {
		unreadCount := m.GetUnreadCount(feed)
		var line string

		if feed.Error != nil {
			line = fmt.Sprintf("❌ %s", feed.Title)
		} else {
			unreadIndicator := ""
			if unreadCount > 0 {
				unreadIndicator = unreadDotStyle.Render(fmt.Sprintf(" (%d unread)", unreadCount))
			}
			line = feed.Title + unreadIndicator
		}

		// truncate description
		if len(line) > m.width-6 {
			visualLen := lg.Width(line)
			if visualLen > m.width-6 {
				line = feed.Title
				if len(line) > m.width-9 {
					line = line[:m.width-9] + "..."
				}
			}
		}

		if i == m.selectedSource {
			line = selectedStyle.Render("▶ ") + line
		} else {
			line = "  " + line
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderArticleView() string {
	var l []string

	// header
	a := m.GetCurrentArticle()
	l = append(l, headerStyle.Render(a.Title))
	l = append(l, "")

	content, err := content.GetArticleContent(a.Link)
	if err != nil {
		l = append(l, dimStyle.Render(fmt.Sprintf("Error loading article: %v", err)))
	} else {
		l = append(l, content)
	}

	return strings.Join(l, "\n")
}

func (m Model) renderStatusBar() string {
	if m.statusMessage != "" {
		return statusStyle.Render(m.statusMessage)
	}

	var help string
	if m.searchInputTrap {
		help = dimStyle.Render("Type to search | Enter: apply | Esc: cancel")
	} else if m.currentView == articlesView {
		if m.searchQuery != "" {
			help = dimStyle.Render("/: search | Esc: clear search | ↑↓/jk: nav | o: open | m: toggle-read | q: quit")
		} else {
			help = dimStyle.Render("/: search | s: sources | ↑↓/jk: nav | o: open | m: toggle-read | r: refresh | ?: help | q: quit")
		}
	} else {
		help = dimStyle.Render("s: back to articles | ↑↓/jk: navigate | enter: filter by source | a: show all | q: quit")
	}
	return statusStyle.Render(help)
}
