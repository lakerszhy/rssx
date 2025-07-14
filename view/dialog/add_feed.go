package dialog

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lakerszhy/rssx/config"
	"github.com/lakerszhy/rssx/message"
	"github.com/lakerszhy/rssx/rss"
)

type AddFeed struct {
	cfg        *config.App
	repo       rss.Repo
	ti         textinput.Model
	addFeedMsg message.AddFeed
}

func NewAddFeed(cfg *config.App, repo rss.Repo) tea.Model {
	return AddFeed{
		cfg:  cfg,
		repo: repo,
		ti:   newTextInput(cfg.Theme, "Feed URL"),
	}
}

func (d AddFeed) Init() tea.Cmd {
	return textinput.Blink
}

func (d AddFeed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case message.AddFeed:
		return d.onAddFeedMsg(msg)
	case tea.KeyMsg:
		if key.Matches(msg, d.cfg.KeyMap.Enter) {
			return d.onEnterKeyMsg()
		}
	}

	d.ti, cmd = d.ti.Update(msg)
	return d, cmd
}

func (d AddFeed) onAddFeedMsg(msg message.AddFeed) (tea.Model, tea.Cmd) {
	d.addFeedMsg = msg

	var cmd tea.Cmd
	if msg.IsInProgress() {
		d.ti.Blur()
	} else {
		cmd = d.ti.Focus()
	}

	return d, cmd
}

func (d AddFeed) onEnterKeyMsg() (tea.Model, tea.Cmd) {
	if d.addFeedMsg.IsInProgress() {
		return d, nil
	}

	v := strings.TrimSpace(d.ti.Value())
	if v == "" {
		return d, nil
	}

	return d, message.AddFeedCmd(v, d.repo)
}

func (d AddFeed) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		inputView(d.ti, d.cfg.Theme),
		fmt.Sprintf("%s\n", d.msgView()),
		actionsView(d.cfg.Theme, false),
	)
	return render("Add Feed", content, d.cfg.Theme)
}

func (d AddFeed) msgView() string {
	style := lipgloss.NewStyle().Width(dialogWidth)
	if d.addFeedMsg.IsInProgress() {
		return style.Foreground(d.cfg.Theme.DialogMsg).Render("Adding...")
	}
	if d.addFeedMsg.IsFailed() {
		return style.Foreground(d.cfg.Theme.Error).Render(d.addFeedMsg.Err.Error())
	}
	return ""
}
