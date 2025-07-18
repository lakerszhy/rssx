package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/internal/rss"
)

func AddFeedCmd(v string, repo rss.Repo) tea.Cmd {
	var cmds []tea.Cmd

	cmd := func() tea.Msg {
		return NewAddFeedInProgress(v)
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		feed, err := rss.ParseURL(v)
		if err != nil {
			return NewAddFeedFailed(v, err)
		}

		feed, err = repo.AddFeed(feed)
		if err != nil {
			return NewAddFeedFailed(v, err)
		}

		return NewAddFeedSuccessful(feed)
	}
	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}

type AddFeed struct {
	Feed rss.Feed
	status
	Err error
}

func NewAddFeedInitial() AddFeed {
	return AddFeed{
		status: statusInitial,
	}
}

func NewAddFeedInProgress(url string) AddFeed {
	return AddFeed{
		Feed:   rss.Feed{FeedURL: url},
		status: statusInProgress,
	}
}

func NewAddFeedSuccessful(f rss.Feed) AddFeed {
	return AddFeed{
		Feed:   f,
		status: statusSuccessful,
	}
}

func NewAddFeedFailed(url string, err error) AddFeed {
	return AddFeed{
		Feed:   rss.Feed{FeedURL: url},
		status: statusFailed,
		Err:    err,
	}
}
