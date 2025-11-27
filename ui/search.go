package ui

import (
	"strings"

	"ohnurr/rss"
)

// returns articles matching the search query
func FilterArticles(articles []articleWithSource, query string) []articleWithSource {
	if query == "" {
		return articles
	}

	query = strings.ToLower(query)
	filtered := make([]articleWithSource, 0)

	for _, item := range articles {
		if matchesSearch(item, query) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// checks if an article matches the search query
func FilterArticlesByContent(article *rss.Article, query string) bool {
	if query == "" {
		return true
	}

	query = strings.ToLower(query)

	// title
	if strings.Contains(strings.ToLower(article.Title), query) {
		return true
	}

	// description
	if strings.Contains(strings.ToLower(article.Description), query) {
		return true
	}

	return false
}

// checks if an article matches the search query
func matchesSearch(item articleWithSource, query string) bool {
	article := item.article

	// title
	if strings.Contains(strings.ToLower(article.Title), query) {
		return true
	}

	// description
	if strings.Contains(strings.ToLower(article.Description), query) {
		return true
	}

	// feed title
	if strings.Contains(strings.ToLower(item.feedTitle), query) {
		return true
	}

	return false
}
