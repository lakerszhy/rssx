package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/internal/rss"
)

func DeleteFeedCmd(f rss.Feed, repo rss.Repo) tea.Cmd {
	var cmds []tea.Cmd

	cmd := func() tea.Msg {
		return NewDeleteFeedInProgress(f)
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		err := repo.DeleteFeed(f.ID)
		if err != nil {
			return NewDeleteFeedFailed(f, err)
		}
		return NewDeleteFeedSuccessful(f)
	}
	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}

type DeleteFeed struct {
	Feed rss.Feed
	status
	Err error
}

func NewDeleteFeedInitial(feed rss.Feed) DeleteFeed {
	return DeleteFeed{
		status: statusInitial,
		Feed:   feed,
	}
}

func NewDeleteFeedInProgress(feed rss.Feed) DeleteFeed {
	return DeleteFeed{
		status: statusInProgress,
		Feed:   feed,
	}
}

func NewDeleteFeedSuccessful(feed rss.Feed) DeleteFeed {
	return DeleteFeed{
		status: statusSuccessful,
		Feed:   feed,
	}
}

func NewDeleteFeedFailed(feed rss.Feed, err error) DeleteFeed {
	return DeleteFeed{
		status: statusFailed,
		Feed:   feed,
		Err:    err,
	}
}
