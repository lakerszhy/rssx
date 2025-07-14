package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/rss"
)

func LoadFeedsCmd(repo rss.Repo) tea.Cmd {
	var cmds []tea.Cmd

	cmd := func() tea.Msg {
		return NewLoadFeedsInProgress()
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		feeds, err := repo.GetAllFeeds()
		if err != nil {
			return NewLoadFeedsFailed(err)
		}
		return NewLoadFeedsSuccessful(feeds)
	}
	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}

type LoadFeeds struct {
	Feeds []rss.Feed
	status
	Err error
}

func NewLoadFeedsInProgress() LoadFeeds {
	return LoadFeeds{
		status: statusInProgress,
	}
}

func NewLoadFeedsSuccessful(feeds []rss.Feed) LoadFeeds {
	return LoadFeeds{
		Feeds:  feeds,
		status: statusSuccessful,
	}
}

func NewLoadFeedsFailed(err error) LoadFeeds {
	return LoadFeeds{
		status: statusFailed,
		Err:    err,
	}
}
