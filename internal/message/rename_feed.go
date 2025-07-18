package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/internal/rss"
)

func RenameFeedCmd(f rss.Feed, name string, repo rss.Repo) tea.Cmd {
	var cmds []tea.Cmd

	cmd := func() tea.Msg {
		return NewRenameFeedInProgress(f)
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		err := repo.RenameFeed(f.ID, name)
		if err != nil {
			return NewRenameFeedFailed(f, err)
		}
		f.Name = name
		return NewRenameFeedSuccessful(f)
	}
	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}

type RenameFeed struct {
	Feed rss.Feed
	status
	Err error
}

func NewRenameFeedInitial(f rss.Feed) RenameFeed {
	return RenameFeed{
		Feed:   f,
		status: statusInitial,
	}
}

func NewRenameFeedInProgress(f rss.Feed) RenameFeed {
	return RenameFeed{
		Feed:   f,
		status: statusInProgress,
	}
}

func NewRenameFeedSuccessful(f rss.Feed) RenameFeed {
	return RenameFeed{
		Feed:   f,
		status: statusSuccessful,
	}
}

func NewRenameFeedFailed(f rss.Feed, err error) RenameFeed {
	return RenameFeed{
		Feed:   f,
		status: statusFailed,
		Err:    err,
	}
}
