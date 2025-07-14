package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/opml"
	"github.com/lakerszhy/rssx/rss"
)

func ImportCmd(fileName string, repo rss.Repo) tea.Cmd {
	var cmds []tea.Cmd

	cmd := func() tea.Msg {
		return NewImportInProgress()
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		feeds, err := opml.Import(fileName)
		if err != nil {
			return NewImportFailed(err)
		}

		feeds, err = repo.AddFeeds(feeds)
		if err != nil {
			return NewImportFailed(err)
		}

		return NewImportSuccessful(feeds)
	}
	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}

type Import struct {
	Feeds []rss.Feed
	status
	Err error
}

func NewImportInitial() Import {
	return Import{status: statusInitial}
}

func NewImportInProgress() Import {
	return Import{status: statusInProgress}
}

func NewImportSuccessful(feeds []rss.Feed) Import {
	return Import{
		status: statusSuccessful,
		Feeds:  feeds,
	}
}

func NewImportFailed(err error) Import {
	return Import{
		status: statusFailed,
		Err:    err,
	}
}
