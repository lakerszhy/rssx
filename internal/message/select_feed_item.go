package message

import "github.com/lakerszhy/rssx/internal/rss"

type SelectFeedItem struct {
	FeedItem *rss.FeedItem
}

func NewSelectFeedItem(i *rss.FeedItem) SelectFeedItem {
	return SelectFeedItem{FeedItem: i}
}
