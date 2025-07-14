package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/rss"
)

func ToogleStarredCmd(itemID int64, repo rss.Repo) tea.Cmd {
	var cmds []tea.Cmd

	cmd := func() tea.Msg {
		return NewToogleStarredSuccessful(itemID)
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		return NewToogleStarredInProgress(itemID)
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		err := repo.ToogleStarred(itemID)
		if err != nil {
			return NewToogleStarredFailed(itemID, err)
		}
		return nil
	}
	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}

type ToogleStarred struct {
	ItemID int64
	Err    error
	status
}

func NewToogleStarredInProgress(itemID int64) ToogleStarred {
	return ToogleStarred{
		ItemID: itemID,
		status: statusInProgress,
	}
}

func NewToogleStarredSuccessful(itemID int64) ToogleStarred {
	return ToogleStarred{
		ItemID: itemID,
		status: statusSuccessful,
	}
}

func NewToogleStarredFailed(itemID int64, err error) ToogleStarred {
	return ToogleStarred{
		ItemID: itemID,
		Err:    err,
		status: statusFailed,
	}
}
