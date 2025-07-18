package rss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToggleItemReadInFeed(t *testing.T) {
	cases := []struct {
		name   string
		itemID int64
		isRead bool
	}{
		{name: "Item is read", itemID: 1, isRead: true},
		{name: "Item is unread", itemID: 1, isRead: false},
	}

	for _, c := range cases {
		f := Feed{
			Items: []FeedItem{
				{ID: c.itemID, IsRead: c.isRead},
				{ID: c.itemID + 1, IsRead: c.isRead},
			},
		}
		f.ToogleRead(c.itemID)
		for _, i := range f.Items {
			if i.ID == c.itemID {
				assert.Equal(t, c.isRead, !i.IsRead)
			} else {
				assert.Equal(t, c.isRead, i.IsRead)
			}
		}
	}
}

func TestToggleItemStarInFeed(t *testing.T) {
	cases := []struct {
		name   string
		itemID int64
		isStar bool
	}{
		{name: "Item is starred", itemID: 1, isStar: true},
		{name: "Item is unstarred", itemID: 1, isStar: false},
	}
	for _, c := range cases {
		f := Feed{
			Items: []FeedItem{
				{ID: c.itemID, IsStarred: c.isStar},
				{ID: c.itemID + 1, IsStarred: c.isStar},
			},
		}
		f.ToogleStarred(c.itemID)
		for _, i := range f.Items {
			if i.ID == c.itemID {
				assert.Equal(t, c.isStar, !i.IsStarred)
			} else {
				assert.Equal(t, c.isStar, i.IsStarred)
			}
		}
	}
}
