package rss

import (
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type FeedItem struct {
	ID          int64
	FeedName    string
	Title       string
	Description string
	Content     string
	Link        string
	IsRead      bool
	IsStarred   bool
	PublishedAt time.Time
}

func (i *FeedItem) ToogleRead() {
	i.IsRead = !i.IsRead
}

func (i *FeedItem) ToogleStarred() {
	i.IsStarred = !i.IsStarred
}

func (i FeedItem) IsToday() bool {
	return i.PublishedAt.After(time.Now().AddDate(0, 0, -1))
}

func (i FeedItem) FilterValue() string {
	return i.Title
}

func (i FeedItem) PlainDescription(policy *bluemonday.Policy) string {
	v := policy.Sanitize(i.Description)
	v = strings.ReplaceAll(v, "\n", " ")
	// \u200B: ZERO WIDTH SPACE, can cause width not correct
	v = strings.ReplaceAll(v, "\u200B", "")
	return strings.TrimSpace(v)
}
