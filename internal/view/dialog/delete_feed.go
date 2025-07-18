package dialog

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lakerszhy/rssx/internal/config"
	"github.com/lakerszhy/rssx/internal/message"
	"github.com/lakerszhy/rssx/internal/rss"
)

type DeleteFeed struct {
	cfg           *config.App
	repo          rss.Repo
	deleteFeedMsg message.DeleteFeed
}

func NewDeleteFeed(cfg *config.App, repo rss.Repo) tea.Model {
	return DeleteFeed{
		cfg:  cfg,
		repo: repo,
	}
}

func (d DeleteFeed) Init() tea.Cmd {
	return nil
}

func (d DeleteFeed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case message.DeleteFeed:
		d.deleteFeedMsg = msg
		return d, nil
	case tea.KeyMsg:
		if key.Matches(msg, d.cfg.KeyMap.Enter) {
			return d.onEnterKeyMsg()
		}
	}

	return d, nil
}

func (d DeleteFeed) onEnterKeyMsg() (tea.Model, tea.Cmd) {
	if d.deleteFeedMsg.IsInProgress() {
		return d, nil
	}

	return d, message.DeleteFeedCmd(d.deleteFeedMsg.Feed, d.repo)
}

func (d DeleteFeed) View() string {
	prompt := lipgloss.NewStyle().Foreground(d.cfg.Theme.BorderActive).Bold(true).
		Render("Are you sure to delete the feed ?")
	name := lipgloss.NewStyle().Width(dialogWidth).Align(lipgloss.Center).
		Foreground(d.cfg.Theme.FeedTitle).Render(d.deleteFeedMsg.Feed.Name)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		prompt,
		name,
		fmt.Sprintf("%s\n", d.msgView()),
		actionsView(d.cfg.Theme, true),
	)
	return render("Delete Feed", content, d.cfg.Theme)
}

func (d DeleteFeed) msgView() string {
	style := lipgloss.NewStyle().Width(dialogWidth)
	if d.deleteFeedMsg.IsInProgress() {
		return style.Foreground(d.cfg.Theme.DialogMsg).Render("Deleting...")
	}
	if d.deleteFeedMsg.IsFailed() {
		return style.Foreground(d.cfg.Theme.Error).
			Render(d.deleteFeedMsg.Err.Error())
	}
	return ""
}
