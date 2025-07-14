package view

import (
	"log/slog"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/lakerszhy/rssx/config"
	"github.com/lakerszhy/rssx/message"
)

type StatusBar struct {
	model   help.Model
	width   int
	cfg     *config.App
	logger  *slog.Logger
	tipsMsg message.Tips
}

func NewStatusBar(cfg *config.App, logger *slog.Logger) StatusBar {
	model := help.New()
	model.Styles.FullKey = lipgloss.NewStyle().Foreground(cfg.Theme.HelpKey)
	model.Styles.FullDesc = lipgloss.NewStyle().Foreground(cfg.Theme.HelpKeyDesc)

	return StatusBar{
		model:  model,
		cfg:    cfg,
		logger: logger,
	}
}

func (s StatusBar) Init() tea.Cmd {
	return nil
}

func (s StatusBar) Update(msg tea.Msg) (StatusBar, tea.Cmd) {
	switch msg := msg.(type) {
	case message.Tips:
		return s, s.onTipsMsg(msg)
	}

	return s, nil
}

func (s *StatusBar) onTipsMsg(msg message.Tips) tea.Cmd {
	s.tipsMsg = msg
	if msg.Text == "" || !msg.AutoDisappear {
		return nil
	}

	//nolint:mnd // 5s
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return message.NewEmptyTips()
	})
}

func (s *StatusBar) Toggle() {
	s.model.ShowAll = !s.model.ShowAll
}

func (s *StatusBar) Hide() {
	s.model.ShowAll = false
}

func (s StatusBar) View() string {
	help := lipgloss.NewStyle().Padding(0, 1).
		Background(s.cfg.Theme.HelpBackground).
		Foreground(s.cfg.Theme.Help).Render("? help")

	msgWidth := s.width - lipgloss.Width(help)
	msg := s.tipsView(msgWidth)

	bar := lipgloss.NewStyle().Width(s.width).Render(msg + help)

	if s.model.ShowAll {
		fullHelp := s.model.View(s.cfg.KeyMap)
		return lipgloss.JoinVertical(lipgloss.Left, bar, fullHelp)
	}
	return bar
}

func (s *StatusBar) tipsView(width int) string {
	style := lipgloss.NewStyle().Padding(0, 1).Width(width).
		Background(s.cfg.Theme.StatusBarBackground).Foreground(s.cfg.Theme.Tips)

	if s.tipsMsg.Text == "" {
		return style.Render("")
	}

	msg := ansi.Truncate(s.tipsMsg.Text, width-2, "...") //nolint:mnd // padding

	if s.tipsMsg.Err != nil {
		style = style.Foreground(s.cfg.Theme.Error)
	}

	return style.Render(msg)
}

func (s *StatusBar) SetWidth(v int) {
	s.width = v
}
