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

type SmartItem struct {
	theme      *config.AppTheme
	htmlPolicy *bluemonday.Policy
}

func NewSmartItem(theme *config.AppTheme) list.ItemDelegate {
	return SmartItem{
		theme:      theme,
		htmlPolicy: bluemonday.StrictPolicy(),
	}
}

func (d SmartItem) Render(w io.Writer, m list.Model, index int, item list.Item) {
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

	info := d.infoView(i, width, m.Index() == index)
	info = paddingStyle.Render(info)

	fmt.Fprintf(w, "%s\n%s\n%s", title, desc, info)
}

func (d SmartItem) descView(i rss.FeedItem, width int, isSelected bool) string {
	prompt := " "
	descStyle := lipgloss.NewStyle().Foreground(d.theme.ItemDesc)

	if isSelected {
		prompt = view.Border.Left
		prompt = lipgloss.NewStyle().Foreground(d.theme.ItemTitleActive).Render(prompt)
		descStyle = descStyle.Foreground(d.theme.ItemDescActive)
	}

	descWidth := width - ansi.StringWidth(prompt) - 1
	desc := i.PlainDescription(d.htmlPolicy)
	desc = ansi.Truncate(desc, descWidth, "...")
	desc = descStyle.Width(descWidth).Render(desc)

	return fmt.Sprintf("%s %s", prompt, desc)
}

func (d SmartItem) infoView(i rss.FeedItem, width int, isSelected bool) string {
	prompt := " "
	style := lipgloss.NewStyle().Bold(true).Foreground(d.theme.ItemDesc)

	if isSelected {
		prompt = view.Border.Left
		prompt = lipgloss.NewStyle().Foreground(d.theme.ItemTitleActive).Render(prompt)
		style = style.Foreground(d.theme.ItemDescActive)
	}

	date := i.PublishedAt.Format(time.DateOnly)
	date = style.Render(date)

	authorWidth := width - ansi.StringWidth(prompt) - ansi.StringWidth(date) - 2 //nolint:mnd // two space
	author := style.Width(authorWidth).Render(i.FeedName)
	author = ansi.Truncate(author, authorWidth, "...")
	return fmt.Sprintf("%s %s %s", prompt, author, date)
}

func (d SmartItem) Height() int {
	return 3 //nolint:mnd // title + desc
}

func (d SmartItem) Spacing() int {
	return 1
}

func (d SmartItem) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
