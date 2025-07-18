package panel

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lakerszhy/rssx/internal/config"
	"github.com/lakerszhy/rssx/internal/message"
	"github.com/lakerszhy/rssx/internal/rss"
	"github.com/lakerszhy/rssx/internal/view"
	"github.com/pkg/browser"
)

type Preview struct {
	cfg       *config.App
	logger    *slog.Logger
	viewport  viewport.Model
	isFocused bool
	item      *rss.FeedItem
}

func NewPreview(cfg *config.App, logger *slog.Logger) Preview {
	vp := viewport.New(0, 0)
	vp.KeyMap = viewport.KeyMap{
		PageDown: cfg.KeyMap.NextPage,
		PageUp:   cfg.KeyMap.PrevPage,
		Down:     cfg.KeyMap.Down,
		Up:       cfg.KeyMap.Up,
	}
	return Preview{
		cfg:      cfg,
		logger:   logger,
		viewport: vp,
	}
}

func (p Preview) Init() tea.Cmd {
	return nil
}

func (p Preview) Update(msg tea.Msg) (Preview, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case message.SelectFeedItem:
		return p, p.onSelectFeedItemMsg(msg)
	case message.ParseMD:
		p.onParseMDMsg(msg)
		return p, nil
	case tea.KeyMsg:
		if key.Matches(msg, p.cfg.KeyMap.Open) {
			p.onOpenKeyMsg()
			return p, nil
		}
		if key.Matches(msg, p.cfg.KeyMap.Start) {
			p.viewport.SetYOffset(0)
			return p, nil
		}
		if key.Matches(msg, p.cfg.KeyMap.End) {
			p.viewport.SetYOffset(p.viewport.TotalLineCount())
			return p, nil
		}
	}

	p.viewport, cmd = p.viewport.Update(msg)
	return p, cmd
}

func (p *Preview) onSelectFeedItemMsg(msg message.SelectFeedItem) tea.Cmd {
	p.item = msg.FeedItem
	p.viewport.SetYOffset(0)

	var cmd tea.Cmd
	if p.item != nil {
		cmd = message.ParseMDCmd(*p.item, p.viewport.Width)
	}
	return cmd
}

func (p *Preview) onParseMDMsg(msg message.ParseMD) {
	if p.item == nil || p.item.ID != msg.FeedItem.ID {
		return
	}

	style := lipgloss.NewStyle().Padding(0, 1).Width(p.viewport.Width)

	if msg.IsFailed() {
		style = style.Foreground(p.cfg.Theme.Error)
		p.viewport.SetContent(style.Render(msg.Err.Error()))
		return
	}

	if msg.IsInProgress() {
		p.viewport.SetContent(style.Render("Loading..."))
		return
	}

	if msg.IsSuccessful() {
		p.viewport.SetContent(style.Render(msg.MD))
	}
}

func (p Preview) onOpenKeyMsg() {
	if p.item == nil {
		return
	}

	if err := browser.OpenURL(p.item.Link); err != nil {
		p.logger.Error("open url", "link", p.item.Link, "err", err)
	}
}

func (p Preview) View() string {
	if p.item == nil {
		style := view.BorderStyle(view.Border, p.cfg.Theme, p.isFocused)
		return style.Width(p.viewport.Width).Height(p.viewport.Height).
			AlignHorizontal(lipgloss.Center).Render("No content.")
	}

	foot := fmt.Sprintf("%.f%%", p.viewport.ScrollPercent()*100) //nolint:mnd // 100% is not a magic number
	b := view.BorderWithFoot(foot, p.viewport.Width)
	style := view.BorderStyle(b, p.cfg.Theme, p.isFocused)
	return style.Render(p.viewport.View())
}

func (p *Preview) SetSize(width, height int) {
	p.viewport.Width = width
	p.viewport.Height = height
}

func (p *Preview) SetFocused(focused bool) {
	p.isFocused = focused
}
