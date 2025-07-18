package opml

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/lakerszhy/rssx/internal/rss"
)

// Specs: http://opml.org/spec2.opml
type opml struct {
	XMLName  xml.Name `xml:"opml"`
	Version  string   `xml:"version,attr"`
	Header   header   `xml:"head"`
	Outlines outlines `xml:"body>outline"`
}

func newOPML() *opml {
	return &opml{
		Version: "2.0",
		Header: header{
			Title:       "RssX",
			DateCreated: time.Now().Format(time.DateTime),
		},
	}
}

func (o opml) toFeeds() []rss.Feed {
	return outlinesToFeeds(o.Outlines)
}

type header struct {
	Title       string `xml:"title,omitempty"`
	DateCreated string `xml:"dateCreated,omitempty"`
	OwnerName   string `xml:"ownerName,omitempty"`
}

type outline struct {
	Title       string   `xml:"title,attr,omitempty"`
	Text        string   `xml:"text,attr"`
	FeedURL     string   `xml:"xmlUrl,attr,omitempty"`
	SiteURL     string   `xml:"htmlUrl,attr,omitempty"`
	Description string   `xml:"description,attr,omitempty"`
	Type        string   `xml:"type,attr,omitempty"`
	Outlines    outlines `xml:"outline,omitempty"`
}

func (o outline) isSubscriptions() bool {
	return strings.TrimSpace(o.FeedURL) != ""
}

func (o outline) toFeed() rss.Feed {
	name := o.Text
	if name == "" {
		name = o.Title
	}
	if name == "" {
		name = o.FeedURL
	}

	return rss.Feed{
		Name:        name,
		FeedURL:     o.FeedURL,
		HomePageURL: o.SiteURL,
	}
}

type outlines []outline

func outlinesToFeeds(o outlines) []rss.Feed {
	var feeds []rss.Feed

	for _, i := range o {
		if i.isSubscriptions() {
			feeds = append(feeds, i.toFeed())
		} else if len(i.Outlines) > 0 {
			feeds = append(feeds, outlinesToFeeds(i.Outlines)...)
		}
	}

	return feeds
}
