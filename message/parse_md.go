package message

import (
	"fmt"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/lakerszhy/rssx/rss"
)

func ParseMDCmd(i rss.FeedItem, wordWrap int) tea.Cmd {
	var cmds []tea.Cmd

	cmd := func() tea.Msg {
		return NewParseMDInProgress(i)
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		v := i.Content
		if v == "" {
			v = i.Description
		}

		// \u200B: ZERO WIDTH SPACE, can cause width not correct
		v = strings.ReplaceAll(v, "\u200B", "")
		v, err := md.ConvertString(v)
		if err != nil {
			return NewParseMDFailed(i, err)
		}

		r, err := glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(wordWrap),
			glamour.WithPreservedNewLines(),
		)
		if err != nil {
			return NewParseMDFailed(i, err)
		}

		var b strings.Builder
		b.WriteString(fmt.Sprintf("# %s\n", i.Title))
		b.WriteString(fmt.Sprintf("%s\n\n", i.PublishedAt.Format(time.DateTime)))
		b.WriteString("---\n")
		b.WriteString(v)

		v, err = r.Render(b.String())
		if err != nil {
			return NewParseMDFailed(i, err)
		}

		return NewParseMDSuccessful(i, v)
	}
	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}

type ParseMD struct {
	FeedItem rss.FeedItem
	MD       string
	status
	Err error
}

func NewParseMDInProgress(i rss.FeedItem) ParseMD {
	return ParseMD{
		FeedItem: i,
		status:   statusInProgress,
	}
}

func NewParseMDSuccessful(i rss.FeedItem, md string) ParseMD {
	return ParseMD{
		FeedItem: i,
		MD:       md,
		status:   statusSuccessful,
	}
}

func NewParseMDFailed(i rss.FeedItem, err error) ParseMD {
	return ParseMD{
		FeedItem: i,
		status:   statusFailed,
		Err:      err,
	}
}
