package opml

import (
	"encoding/xml"
	"os"

	"github.com/lakerszhy/rssx/rss"
)

func Import(name string) ([]rss.Feed, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var doc opml
	if err = xml.NewDecoder(f).Decode(&doc); err != nil {
		return nil, err
	}

	return doc.toFeeds(), nil
}
