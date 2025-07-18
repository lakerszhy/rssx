package delegate

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/lakerszhy/rssx/internal/config"
	"github.com/lakerszhy/rssx/internal/rss"
)

type feed struct {
	theme *config.AppTheme
}

func NewFeed(theme *config.AppTheme) list.ItemDelegate {
	return feed{
		theme: theme,
	}
}

func (f feed) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(rss.Feed)
	if !ok {
		return
	}

	style := f.titleStyle(i, index == m.Index())

	prompt := " "
	if index == m.Index() {
		prompt = ">"
	}

	unreadCount := i.UnreadCount()
	unreadStr := ""
	if unreadCount > 0 {
		unreadStr = fmt.Sprintf(" [%d]", unreadCount)
		unreadStr = style.Render(unreadStr)
	}

	nameWidth := m.Width() - lipgloss.Width(unreadStr)
	name := ansi.Truncate(fmt.Sprintf("%s %s", prompt, i.Name),
		nameWidth, "...")
	name = style.Width(nameWidth).Render(name)

	text := lipgloss.NewStyle().Render(name + unreadStr)
	fmt.Fprint(w, text)
}

func (f feed) titleStyle(i rss.Feed, isSelected bool) lipgloss.Style {
	style := lipgloss.NewStyle().Foreground(f.theme.FeedTitle)

	if i.IsSmart() {
		style = style.Foreground(f.theme.SmartFeed)
	}

	if isSelected {
		if i.IsSmart() {
			style = style.Foreground(f.theme.SmartFeedActive).Bold(true)
		} else {
			style = style.Foreground(f.theme.FeedTitleActive)
		}
	}

	return style
}

func (f feed) Height() int {
	return 1
}

func (f feed) Spacing() int {
	return 0
}

func (f feed) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
