package rss

type Repo interface {
	AddFeed(Feed) (Feed, error)
	AddFeeds([]Feed) ([]Feed, error)
	InsertItems(feedID int64, items []FeedItem) ([]FeedItem, error)
	GetAllFeeds() ([]Feed, error)
	DeleteFeed(id int64) error
	ToogleRead(id int64) error
	ToogleStarred(id int64) error
	RenameFeed(id int64, name string) error
}
