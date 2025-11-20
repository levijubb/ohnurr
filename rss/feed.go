package rss

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

type Feed struct {
	URL      string
	Title    string
	Articles []Article
	Error    error
}

type Article struct {
	Title       string
	Link        string
	Description string
	Published   time.Time
	GUID        string
	FeedTitle   string
}

// fetches and parses an RSS feed from the given URL
func FetchFeed(url string) (*Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return &Feed{
			URL:   url,
			Error: err,
		}, err
	}

	articles := make([]Article, 0, len(feed.Items))
	for _, item := range feed.Items {
		guid := item.GUID
		if guid == "" {
			guid = item.Link
		}

		var published time.Time
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			published = *item.UpdatedParsed
		}

		articles = append(articles, Article{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Published:   published,
			GUID:        guid,
			FeedTitle:   feed.Title,
		})
	}

	return &Feed{
		URL:      url,
		Title:    feed.Title,
		Articles: articles,
	}, nil
}

// fetch multiple RSS feeds concurrently
func FetchAllFeeds(urls []string) []*Feed {
	results := make([]*Feed, len(urls))
	done := make(chan bool)

	for i, url := range urls {
		go func(index int, url string) {
			feed, err := FetchFeed(url)
			if err != nil {
				results[index] = &Feed{
					URL:   url,
					Title: fmt.Sprintf("Error loading feed: %s", url),
					Error: err,
				}
			} else {
				results[index] = feed
			}
			done <- true
		}(i, url)
	}

	// wait for all feeds to complete
	for range urls {
		<-done
	}

	return results
}

// returns a unique identifier for an article
func (a *Article) GetArticleID() string {
	if a.GUID != "" {
		return a.GUID
	}
	return a.Link
}
