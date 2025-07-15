package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/rss"
)

func MarkAllReadCmd(itemIDs []int64, repo rss.Repo) tea.Cmd {
	var cmds []tea.Cmd

	cmd := func() tea.Msg {
		return NewMarkAllReadSuccessful(itemIDs)
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		return NewMarkAllReadInProgress(itemIDs)
	}
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		err := repo.MarkAllRead(itemIDs)
		if err != nil {
			return NewMarkAllReadFailed(itemIDs, err)
		}
		return nil
	}
	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}

type MarkAllRead struct {
	ItemIDs []int64
	Err     error
	status
}

func NewMarkAllReadInProgress(itemIDs []int64) MarkAllRead {
	return MarkAllRead{
		ItemIDs: itemIDs,
		status:  statusInProgress,
	}
}

func NewMarkAllReadSuccessful(itemIDs []int64) MarkAllRead {
	return MarkAllRead{
		ItemIDs: itemIDs,
		status:  statusSuccessful,
	}
}

func NewMarkAllReadFailed(itemIDs []int64, err error) MarkAllRead {
	return MarkAllRead{
		ItemIDs: itemIDs,
		Err:     err,
		status:  statusFailed,
	}
}
