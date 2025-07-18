package dialog

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lakerszhy/rssx/internal/config"
	"github.com/lakerszhy/rssx/internal/message"
	"github.com/lakerszhy/rssx/internal/rss"
)

type RenameFeed struct {
	cfg           *config.App
	repo          rss.Repo
	ti            textinput.Model
	renameFeedMsg message.RenameFeed
}

func NewRenameFeed(cfg *config.App, repo rss.Repo) tea.Model {
	return RenameFeed{
		cfg:  cfg,
		repo: repo,
		ti:   newTextInput(cfg.Theme, "Feed Name"),
	}
}

func (d RenameFeed) Init() tea.Cmd {
	return textinput.Blink
}

func (d RenameFeed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case message.RenameFeed:
		return d.onRenameFeedMsg(msg)
	case tea.KeyMsg:
		if key.Matches(msg, d.cfg.KeyMap.Enter) {
			return d.onEnterKeyMsg()
		}
	}

	d.ti, cmd = d.ti.Update(msg)
	return d, cmd
}

func (d RenameFeed) onRenameFeedMsg(msg message.RenameFeed) (tea.Model, tea.Cmd) {
	d.renameFeedMsg = msg

	var cmd tea.Cmd

	switch {
	case msg.IsInitial():
		d.ti.SetValue(msg.Feed.Name)
		cmd = d.ti.Focus()
	case msg.IsInProgress():
		d.ti.Blur()
	case msg.IsFailed():
		cmd = d.ti.Focus()
	}

	return d, cmd
}

func (d RenameFeed) onEnterKeyMsg() (tea.Model, tea.Cmd) {
	if d.renameFeedMsg.IsInProgress() {
		return d, nil
	}

	v := strings.TrimSpace(d.ti.Value())
	if v == "" {
		return d, nil
	}

	return d, message.RenameFeedCmd(d.renameFeedMsg.Feed, v, d.repo)
}

func (d RenameFeed) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		inputView(d.ti, d.cfg.Theme),
		fmt.Sprintf("%s\n", d.msgView()),
		actionsView(d.cfg.Theme, false),
	)

	return render("Rename Feed", content, d.cfg.Theme)
}

func (d RenameFeed) msgView() string {
	style := lipgloss.NewStyle().Width(dialogWidth)
	if d.renameFeedMsg.IsInProgress() {
		return style.Foreground(d.cfg.Theme.DialogMsg).Render("Renaming...")
	}
	if d.renameFeedMsg.IsFailed() {
		return style.Foreground(d.cfg.Theme.Error).Render(d.renameFeedMsg.Err.Error())
	}
	return ""
}
