package ui

import (
	"fmt"
	"strings"
	"time"

	lg "github.com/charmbracelet/lipgloss"
)

var (
	activeColor   = lg.Color("111")
	inactiveColor = lg.Color("240")
	unreadColor   = lg.Color("75")
	titleColor    = lg.Color("108")
	accentColor   = lg.Color("103")

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

func formatPublishDate(published time.Time) string {
	now := time.Now()
	diff := now.Sub(published)

	// <1 hour
	if diff < time.Hour {
		mins := int(diff.Minutes())
		if mins < 1 {
			return "just now"
		}
		return fmt.Sprintf("%dm ago", mins)
	}

	// <24 hours
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	}

	// <7 days
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%dd ago", days)
	}

	// < 1 year so show month and day
	if published.Year() == now.Year() {
		return published.Format("Jan 2")
	}

	// erry old
	return published.Format("Jan 2, 2006")
}

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

			// source and date
			if lineCount < availableHeight-2 {
				dateStr := formatPublishDate(article.Published)
				sourceLine := "    " + sourceStyle.Render("from "+item.feedTitle) + dimStyle.Render(" · "+dateStr)
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

func wrapLineWithIndent(line string, maxWidth int, leftMargin int) []string {
	// empty lines
	if strings.TrimSpace(line) == "" {
		return []string{""}
	}

	indent := strings.Repeat(" ", leftMargin)

	// dont wrap lines with ANSI escape codes (headers, code blocks)
	if strings.Contains(line, "\x1b[") {
		return []string{indent + line}
	}

	var wrapped []string
	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{""}
	}

	currentLine := ""
	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= maxWidth {
			currentLine = testLine
		} else {
			if currentLine != "" {
				wrapped = append(wrapped, indent+currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		wrapped = append(wrapped, indent+currentLine)
	}

	return wrapped
}

func (m Model) renderArticleView() string {
	var header []string

	// add top padding
	header = append(header, "")

	a := m.GetCurrentArticle()

	if m.loadingArticle {
		header = append(header, dimStyle.Render("Loading article..."))
		return strings.Join(header, "\n")
	}

	if m.cachedArticleURL != a.Link || m.cachedArticleContent == "" {
		header = append(header, dimStyle.Render("Article not loaded. Press Esc and Enter to reload."))
		return strings.Join(header, "\n")
	}

	readableWidth := 80
	if m.width < 100 {
		// for narrow screens use 90% of width
		readableWidth = int(float64(m.width) * 0.9)
	}
	if readableWidth < 40 {
		readableWidth = 40 // min width
	}

	// calc left margin
	leftMargin := max((m.width-readableWidth)/2, 2)

	// title
	titleLine := headerStyle.Render(a.Title)
	titleWidth := lg.Width(titleLine)
	titlePadding := (m.width - titleWidth) / 2
	if titlePadding > 0 {
		header = append(header, strings.Repeat(" ", titlePadding)+titleLine)
	} else {
		header = append(header, strings.Repeat(" ", leftMargin)+titleLine)
	}
	header = append(header, "")

	// split content into lines and wrap
	rawLines := strings.Split(m.cachedArticleContent, "\n")
	var wrappedLines []string
	for _, line := range rawLines {
		wrapped := wrapLineWithIndent(line, readableWidth, leftMargin)
		wrappedLines = append(wrappedLines, wrapped...)
	}

	// calcl height for content
	headerHeight := len(header)
	availableHeight := m.height - headerHeight - 1

	// clamp scroll
	scrollPos := m.articleScroll
	maxScroll := max(len(wrappedLines)-availableHeight, 0)
	if scrollPos > maxScroll {
		scrollPos = maxScroll
	}
	if scrollPos < 0 {
		scrollPos = 0
	}

	startLine := scrollPos
	endLine := min(startLine+availableHeight, len(wrappedLines))
	visibleContent := wrappedLines[startLine:endLine]

	// combine header and content
	result := strings.Join(header, "\n")
	if len(visibleContent) > 0 {
		result += "\n" + strings.Join(visibleContent, "\n")
	}

	return result
}

func (m Model) renderStatusBar() string {
	if m.statusMessage != "" {
		return statusStyle.Render(m.statusMessage)
	}

	return statusStyle.Render(m.getHelpText())
}

func (m Model) getHelpText() string {
	if m.searchInputTrap {
		return dimStyle.Render("Type to search | Enter: apply | Esc: cancel")
	}

	switch m.currentView {
	case articlesView:
		return dimStyle.Render("/: search | ↑↓/jk: nav | o: open | m: toggle-read | s: sources | r: refresh | q: quit")
	case articleView:
		return dimStyle.Render("↑↓/jk: scroll | PgDn/PgUp: page | g/G: top/bottom | o: open | Esc: back | q: quit")
	case sourcesView:
		return dimStyle.Render("s: back to articles | ↑↓/jk: navigate | enter: filter by source | a: show all | q: quit")
	}
	return ""
}
