package delegate

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/lakerszhy/rssx/config"
	"github.com/lakerszhy/rssx/rss"
	"github.com/lakerszhy/rssx/view"
	"github.com/microcosm-cc/bluemonday"
)

type item struct {
	theme      *config.AppTheme
	htmlPolicy *bluemonday.Policy
}

func NewItem(theme *config.AppTheme) list.ItemDelegate {
	return item{
		theme:      theme,
		htmlPolicy: bluemonday.StrictPolicy(),
	}
}

func (d item) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if m.Width() <= 0 {
		return
	}

	i, ok := item.(rss.FeedItem)
	if !ok {
		return
	}

	paddingStyle := lipgloss.NewStyle().Padding(0, 1)
	width := m.Width() - paddingStyle.GetHorizontalPadding()

	title := titleView(i, width, m.Index() == index, d.theme)
	title = paddingStyle.Render(title)

	desc := d.descView(i, width, m.Index() == index)
	desc = paddingStyle.Render(desc)

	fmt.Fprintf(w, "%s\n%s", title, desc)
}

func (d item) descView(i rss.FeedItem, width int, isSelected bool) string {
	prompt := " "
	descStyle := lipgloss.NewStyle().Foreground(d.theme.ItemDesc)

	if isSelected {
		prompt = view.Border.Left
		prompt = lipgloss.NewStyle().Foreground(d.theme.ItemTitleActive).Render(prompt)
		descStyle = descStyle.Foreground(d.theme.ItemDescActive)
	}

	date := i.PublishedAt.Format(time.DateOnly)
	date = descStyle.Bold(true).Render(date)

	descWidth := width - ansi.StringWidth(prompt) - ansi.StringWidth(date) - 2 //nolint:mnd // two space
	desc := i.PlainDescription(d.htmlPolicy)
	desc = ansi.Truncate(desc, descWidth, "...")
	desc = descStyle.Width(descWidth).Render(desc)

	return fmt.Sprintf("%s %s %s", prompt, desc, date)
}

func (d item) Height() int {
	return 2 //nolint:mnd // title + desc
}

func (d item) Spacing() int {
	return 1
}

func (d item) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func titleView(i rss.FeedItem, width int, isSelected bool, theme *config.AppTheme) string {
	prompt := " "
	unread := ""
	starred := ""
	titleStyle := lipgloss.NewStyle().Foreground(theme.ItemTitle)

	if !i.IsRead {
		unread = lipgloss.NewStyle().Foreground(theme.Unread).Render("⏺")
	}
	if i.IsStarred {
		starred = lipgloss.NewStyle().Foreground(theme.Starred).Render("⭑")
	}
	if isSelected {
		prompt = view.Border.Left
		titleStyle = titleStyle.Foreground(theme.ItemTitleActive)
	}

	suffix := starred
	if len(unread) > 0 {
		suffix = fmt.Sprintf("%s %s", starred, unread)
	}
	titleWidth := width - ansi.StringWidth(suffix)
	if len(suffix) > 0 {
		titleWidth--
	}

	title := fmt.Sprintf("%s %s", prompt, i.Title)
	title = ansi.Truncate(title, titleWidth, "...")
	title = titleStyle.Width(titleWidth).Render(title)

	if len(suffix) > 0 {
		title = fmt.Sprintf("%s %s", title, suffix)
	}

	return title
}
