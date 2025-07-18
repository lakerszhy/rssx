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

type Import struct {
	cfg       *config.App
	repo      rss.Repo
	ti        textinput.Model
	importMsg message.Import
}

func NewImport(cfg *config.App, repo rss.Repo) tea.Model {
	return Import{
		cfg:       cfg,
		repo:      repo,
		ti:        newTextInput(cfg.Theme, "File Path"),
		importMsg: message.NewImportInitial(),
	}
}

func (d Import) Init() tea.Cmd {
	return textinput.Blink
}

func (d Import) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case message.Import:
		return d.onImportMsg(msg)
	case tea.KeyMsg:
		if key.Matches(msg, d.cfg.KeyMap.Enter) {
			return d.onEnterKeyMsg()
		}
	}

	d.ti, cmd = d.ti.Update(msg)
	return d, cmd
}

func (d Import) onImportMsg(msg message.Import) (tea.Model, tea.Cmd) {
	d.importMsg = msg

	var cmd tea.Cmd
	if msg.IsInProgress() {
		d.ti.Blur()
	} else {
		cmd = d.ti.Focus()
	}

	return d, cmd
}

func (d Import) onEnterKeyMsg() (tea.Model, tea.Cmd) {
	if d.importMsg.IsInProgress() {
		return d, nil
	}

	v := strings.TrimSpace(d.ti.Value())
	if v == "" {
		return d, nil
	}

	return d, message.ImportCmd(v, d.repo)
}

func (d Import) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		inputView(d.ti, d.cfg.Theme),
		fmt.Sprintf("%s\n", d.msgView()),
		actionsView(d.cfg.Theme, false),
	)
	return render("Import Feeds", content, d.cfg.Theme)
}

func (d Import) msgView() string {
	style := lipgloss.NewStyle().Width(dialogWidth)
	if d.importMsg.IsInProgress() {
		return style.Foreground(d.cfg.Theme.DialogMsg).Render("Importing...")
	}
	if d.importMsg.IsFailed() {
		return style.Foreground(d.cfg.Theme.Error).Render(d.importMsg.Err.Error())
	}
	return ""
}
