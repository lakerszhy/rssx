package rss

import (
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

func ParseURL(u string) (Feed, error) {
	fp := gofeed.NewParser()
	fp.Client = &http.Client{
		Timeout: 10 * time.Second, //nolint:mnd // timeout
	}
	f, err := fp.ParseURL(u)
	if err != nil {
		return Feed{FeedURL: u}, err
	}

	feed := toFeed(f)
	if feed.FeedURL == "" {
		feed.FeedURL = u
	}
	return feed, nil
}

func toFeed(f *gofeed.Feed) Feed {
	items := make([]FeedItem, 0, len(f.Items))
	for _, v := range f.Items {
		desc := strings.TrimSpace(v.Description)
		publishedAt := time.Now()
		if v.PublishedParsed != nil {
			publishedAt = *v.PublishedParsed
		}
		items = append(items, FeedItem{
			Title:       v.Title,
			Description: desc,
			Content:     v.Content,
			Link:        v.Link,
			PublishedAt: publishedAt,
		})
	}

	return Feed{
		Name:        f.Title,
		FeedURL:     f.FeedLink,
		HomePageURL: f.Link,
		Items:       items,
	}
}
