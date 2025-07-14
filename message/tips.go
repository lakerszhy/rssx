package message

import (
	tea "github.com/charmbracelet/bubbletea"
)

func TipsCmd(text string, autoDisappear bool) tea.Cmd {
	return func() tea.Msg {
		return Tips{
			Text:          text,
			AutoDisappear: autoDisappear,
		}
	}
}

func ErrTipsCmd(text string, err error, autoDisappear bool) tea.Cmd {
	return func() tea.Msg {
		return Tips{
			Text:          text,
			Err:           err,
			AutoDisappear: autoDisappear,
		}
	}
}

type Tips struct {
	Text          string
	Err           error
	AutoDisappear bool
}

func NewEmptyTips() Tips {
	return Tips{}
}

func NewTips(text string, autoDisappear bool) Tips {
	return Tips{
		Text:          text,
		AutoDisappear: autoDisappear,
	}
}

func NewErrTips(text string, err error, autoDisappear bool) Tips {
	return Tips{
		Text:          text,
		Err:           err,
		AutoDisappear: autoDisappear,
	}
}
