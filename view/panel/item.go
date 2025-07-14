package panel

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/lakerszhy/rssx/config"
	"github.com/lakerszhy/rssx/message"
	"github.com/lakerszhy/rssx/rss"
	"github.com/lakerszhy/rssx/view"
	"github.com/lakerszhy/rssx/view/delegate"
	"github.com/pkg/browser"
)

type Item struct {
	height    int
	cfg       *config.App
	logger    *slog.Logger
	repo      rss.Repo
	feed      *rss.Feed
	listView  listView[rss.FeedItem]
	isFocused bool
}

func NewItem(cfg *config.App, logger *slog.Logger, repo rss.Repo) Item {
	return Item{
		cfg:      cfg,
		logger:   logger,
		repo:     repo,
		listView: newListView[rss.FeedItem](cfg, delegate.NewItem(cfg.Theme)),
	}
}

func (p Item) Init() tea.Cmd {
	return nil
}

func (p Item) Update(msg tea.Msg) (Item, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case message.SelectFeed:
		return p, p.onSelectFeed(msg)
	case tea.KeyMsg:
		if key.Matches(msg, p.cfg.KeyMap.ToogleRead) {
			return p, p.sendToogleReadCmd()
		}
		if key.Matches(msg, p.cfg.KeyMap.ToogleStarred) {
			return p, p.sendToogleStarredCmd()
		}
		if key.Matches(msg, p.cfg.KeyMap.Open) {
			p.onOpenKeyMsg()
			return p, nil
		}
	}

	p.listView, cmd = p.listView.Update(msg)
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		return message.NewSelectFeedItem(p.listView.selectedItem())
	}
	cmds = append(cmds, cmd)

	// When selected item is not read, send toogle read message
	cmds = append(cmds, p.sendReadCmd())

	return p, tea.Batch(cmds...)
}

func (p *Item) onSelectFeed(msg message.SelectFeed) tea.Cmd {
	var cmd tea.Cmd

	if msg.Feed == nil {
		p.feed = msg.Feed
		p.listView.setItems([]rss.FeedItem{})
		return cmd
	}

	slices.SortFunc(msg.Feed.Items, func(a, b rss.FeedItem) int {
		return b.PublishedAt.Compare(a.PublishedAt)
	})
	p.listView.setItems(msg.Feed.Items)

	// When selected feed is changed, unselect item
	if p.feed == nil || msg.Feed.ID != p.feed.ID {
		p.listView.selectByIndex(-1)
		cmd = func() tea.Msg {
			return message.NewSelectFeedItem(nil)
		}
	}

	p.feed = msg.Feed

	if p.feed != nil {
		if p.feed.IsSmart() {
			p.listView.setDelegate(delegate.NewSmartItem(p.cfg.Theme))
		} else {
			p.listView.setDelegate(delegate.NewItem(p.cfg.Theme))
		}
	}

	return cmd
}

func (p *Item) SetFocused(focused bool) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	p.isFocused = focused
	if p.isFocused && p.listView.selectedItem() == nil && len(p.listView.items()) > 0 {
		p.listView.selectByIndex(0)
		cmd = func() tea.Msg {
			return message.NewSelectFeedItem(p.listView.selectedItem())
		}
		cmds = append(cmds, cmd)
	}

	// When focused, if selected item is not read, send toogle read message.
	if p.isFocused {
		cmd = p.sendReadCmd()
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (p Item) sendReadCmd() tea.Cmd {
	var cmd tea.Cmd
	i := p.listView.selectedItem()
	if p.feed != nil && i != nil && !i.IsRead {
		cmd = message.ToogleReadCmd(i.ID, p.repo)
	}
	return cmd
}

func (p Item) sendToogleReadCmd() tea.Cmd {
	var cmd tea.Cmd
	i := p.listView.selectedItem()
	if p.feed != nil && i != nil {
		cmd = message.ToogleReadCmd(i.ID, p.repo)
	}
	return cmd
}

func (p Item) sendToogleStarredCmd() tea.Cmd {
	var cmd tea.Cmd
	i := p.listView.selectedItem()
	if p.feed != nil && i != nil {
		cmd = message.ToogleStarredCmd(i.ID, p.repo)
	}
	return cmd
}

func (p Item) onOpenKeyMsg() {
	i := p.listView.selectedItem()
	if i == nil {
		return
	}

	if err := browser.OpenURL(i.Link); err != nil {
		p.logger.Error("open url", "link", i.Link, "err", err)
	}
}

func (p *Item) SetSize(width, height int) {
	p.height = height
	p.listView.setSize(width, height-2) //nolint:mnd // title+divider
}

func (p Item) View() string {
	if p.feed == nil {
		return ""
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		p.titleView(),
		p.divider(),
		p.listView.View(),
	)

	style := view.BorderStyle(p.border(), p.cfg.Theme, p.isFocused)
	return style.Render(content)
}

func (p Item) titleView() string {
	title := p.feed.Name

	unread := ""
	if p.feed.UnreadCount() > 0 {
		unread = fmt.Sprintf(" [%d]", p.feed.UnreadCount())
	}

	titleStyle := lipgloss.NewStyle().Foreground(p.cfg.Theme.FeedTitleActive).
		Padding(0, 0, 0, 1)
	if p.feed.IsSmart() {
		titleStyle = titleStyle.Foreground(p.cfg.Theme.SmartFeedActive)
	}

	titleWidth := p.cfg.ItemPanelWidth - ansi.StringWidth(unread) -
		titleStyle.GetHorizontalPadding()

	title = ansi.Truncate(title, titleWidth, "...")
	title = lipgloss.NewStyle().Width(titleWidth).Render(title)

	return titleStyle.Render(title + unread)
}

func (p Item) border() lipgloss.Border {
	b := view.BorderWithFoot(p.listView.footView(), p.cfg.ItemPanelWidth)

	left := b.Left + b.MiddleLeft
	left = fmt.Sprintf("%s%s", left, strings.Repeat(b.Left, p.height))
	b.Left = left

	right := b.Right + b.MiddleRight
	right = fmt.Sprintf("%s%s", right, strings.Repeat(b.Right, p.height))
	b.Right = right

	return b
}

func (p Item) divider() string {
	color := p.cfg.Theme.Border
	if p.isFocused {
		color = p.cfg.Theme.BorderActive
	}

	return lipgloss.NewStyle().Foreground(color).
		Render(strings.Repeat(view.Border.Top, p.cfg.ItemPanelWidth))
}
