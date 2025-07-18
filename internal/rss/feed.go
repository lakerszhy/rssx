package rss

import (
	"strings"
)

const (
	smartFeedID = 0
)

type Feed struct {
	ID          int64
	Name        string
	FeedURL     string
	HomePageURL string
	Items       []FeedItem
}

func NewTodayFeed() Feed {
	return Feed{
		ID:   smartFeedID,
		Name: "⛭ Today",
	}
}

func NewUnreadFeed() Feed {
	return Feed{
		ID:   smartFeedID,
		Name: "⭘ Unread",
	}
}

func NewStarredFeed() Feed {
	return Feed{
		ID:   smartFeedID,
		Name: "⛤ Starred",
	}
}

func (f Feed) UnreadCount() int {
	count := 0
	for _, i := range f.Items {
		if !i.IsRead {
			count++
		}
	}
	return count
}

func (f *Feed) ToogleRead(itemID int64) *Feed {
	items := make([]FeedItem, 0, len(f.Items))
	for _, i := range f.Items {
		if i.ID == itemID {
			i.ToogleRead()
		}
		items = append(items, i)
	}
	f.Items = items
	return f
}

func (f *Feed) MarkAllRead(itemIDs []int64) *Feed {
	items := make([]FeedItem, 0, len(f.Items))
	for _, i := range f.Items {
		for _, id := range itemIDs {
			if i.ID == id {
				i.MarkRead()
				break
			}
		}
		items = append(items, i)
	}
	f.Items = items
	return f
}

func (f *Feed) ToogleStarred(itemID int64) *Feed {
	items := make([]FeedItem, 0, len(f.Items))
	for _, i := range f.Items {
		if i.ID == itemID {
			i.ToogleStarred()
		}
		items = append(items, i)
	}
	f.Items = items
	return f
}

func (f *Feed) Rename(v string) *Feed {
	f.Name = strings.TrimSpace(v)
	return f
}

func (f Feed) IsSmart() bool {
	return f.ID == smartFeedID
}

func (f Feed) FilterValue() string {
	return f.Name
}
