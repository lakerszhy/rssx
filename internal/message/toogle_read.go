package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/internal/rss"
)

func ToogleReadCmd(itemID int64, repo rss.Repo) tea.Cmd {
	var cmds []tea.Cmd

	cmd := func() tea.Msg {
		return NewToogleReadSuccessful(itemID)
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		return NewToogleReadInProgress(itemID)
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		err := repo.ToogleRead(itemID)
		if err != nil {
			return NewToogleReadFailed(itemID, err)
		}
		return nil
	}
	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}

type ToogleRead struct {
	ItemID int64
	Err    error
	status
}

func NewToogleReadInProgress(itemID int64) ToogleRead {
	return ToogleRead{
		ItemID: itemID,
		status: statusInProgress,
	}
}

func NewToogleReadSuccessful(itemID int64) ToogleRead {
	return ToogleRead{
		ItemID: itemID,
		status: statusSuccessful,
	}
}

func NewToogleReadFailed(itemID int64, err error) ToogleRead {
	return ToogleRead{
		ItemID: itemID,
		Err:    err,
		status: statusFailed,
	}
}
