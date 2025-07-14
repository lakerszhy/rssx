package message

import "github.com/lakerszhy/rssx/rss"

type SelectFeed struct {
	Feed *rss.Feed
}

func NewSelectFeed(f *rss.Feed) SelectFeed {
	return SelectFeed{Feed: f}
}
