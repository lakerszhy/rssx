package message

import "github.com/lakerszhy/rssx/internal/rss"

type SelectFeed struct {
	Feed *rss.Feed
}

func NewSelectFeed(f *rss.Feed) SelectFeed {
	return SelectFeed{Feed: f}
}
