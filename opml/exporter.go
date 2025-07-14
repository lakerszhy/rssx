package opml

import (
	"encoding/xml"
	"os"

	"github.com/lakerszhy/rssx/rss"
)

func Export(feeds []rss.Feed, filePath string) error {
	doc := newOPML()

	for _, f := range feeds {
		o := outline{
			Title:       f.Name,
			Text:        f.Name,
			FeedURL:     f.FeedURL,
			SiteURL:     f.HomePageURL,
			Description: "",
			Type:        "rss",
		}
		doc.Outlines = append(doc.Outlines, o)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(xml.Header)
	if err != nil {
		return err
	}

	encoder := xml.NewEncoder(f)
	encoder.Indent("", "    ")
	return encoder.Encode(doc)
}
