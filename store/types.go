package store

import (
	"database/sql"
	"time"

	"github.com/lakerszhy/rssx/rss"
)

type feed struct {
	id          int64
	name        string
	feedURL     string
	homePageURL string
}

func (f feed) toFeed() rss.Feed {
	return rss.Feed{
		ID:          f.id,
		Name:        f.name,
		FeedURL:     f.feedURL,
		HomePageURL: f.homePageURL,
	}
}

type feedItem struct {
	id          int64
	feedID      int64
	title       string
	description sql.NullString
	content     sql.NullString
	link        string
	isRead      bool
	isStarred   bool
	publishedAt sql.NullInt64
}

func (i feedItem) toItem() rss.FeedItem {
	return rss.FeedItem{
		ID:          i.id,
		Title:       i.title,
		Description: i.description.String,
		Content:     i.content.String,
		Link:        i.link,
		IsRead:      i.isRead,
		IsStarred:   i.isStarred,
		PublishedAt: time.Unix(i.publishedAt.Int64, 0),
	}
}
