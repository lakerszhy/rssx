package app

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lakerszhy/rssx/config"
	"github.com/lakerszhy/rssx/message"
	"github.com/lakerszhy/rssx/rss"
	"github.com/lakerszhy/rssx/view"
	"github.com/lakerszhy/rssx/view/dialog"
	"github.com/lakerszhy/rssx/view/panel"
	"github.com/pkg/browser"
)

type app struct {
	windowWidth  int
	windowHeight int

	dir     string
	cfg     *config.App
	logger  *slog.Logger
	repo    rss.Repo
	version string

	feedPanel    panel.Feed
	itemPanel    panel.Item
	previewPanel panel.Preview
	dialog       tea.Model
	statusBar    view.StatusBar

	focus focus

	loadFeedsMsg message.LoadFeeds
	refreshMsg   message.Refresh
}

func New(dir string, cfg *config.App, logger *slog.Logger,
	repo rss.Repo, version string) tea.Model {
	return app{
		dir:          dir,
		cfg:          cfg,
		logger:       logger,
		repo:         repo,
		focus:        focusFeed,
		loadFeedsMsg: message.NewLoadFeedsInProgress(),
		feedPanel:    panel.NewFeed(cfg, logger, repo),
		itemPanel:    panel.NewItem(cfg, logger, repo),
		previewPanel: panel.NewPreview(cfg, logger),
		statusBar:    view.NewStatusBar(cfg, logger, version),
	}
}

func (a app) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("RssX"),
		message.LoadFeedsCmd(a.repo),
	)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case message.AddFeed:
		return a.onAddFeedMsg(msg)
	case message.DeleteFeed:
		return a.onDeleteFeedMsg(msg)
	case message.RenameFeed:
		return a.onRenameFeedMsg(msg)
	case message.LoadFeeds:
		return a.onLoadFeedsMsg(msg)
	case message.SelectFeed:
		return a.onSelectFeedMsg(msg)
	case message.SelectFeedItem:
		return a.onSelectFeedItemMsg(msg)
	case message.ToogleRead:
		return a.onToogleReadMsg(msg)
	case message.MarkAllRead:
		return a.onMarkAllReadMsg(msg)
	case message.ToogleStarred:
		return a.onToogleStarredMsg(msg)
	case message.ParseMD:
		return a.onParseMDMsg(msg)
	case message.Refresh:
		return a.onRefreshMsg(msg)
	case message.RefreshTick:
		return a.onRefreshTickMsg(msg)
	case message.Export:
		return a.onExportMsg(msg)
	case message.Import:
		return a.onImportMsg(msg)
	case message.Tips:
		a.statusBar, cmd = a.statusBar.Update(msg)
		return a, cmd
	case tea.KeyMsg:
		return a, a.onKeyMsg(msg)
	case tea.WindowSizeMsg:
		a.onWindowSizeMsg(msg)
		return a, nil
	}

	if a.dialog != nil {
		a.dialog, cmd = a.dialog.Update(msg)
		return a, cmd
	}

	return a, nil
}

func (a app) onAddFeedMsg(msg message.AddFeed) (app, tea.Cmd) {
	// When add feed is in progress, user close add feed dialog,
	// we should not update dialog.
	if _, ok := a.dialog.(dialog.AddFeed); !ok {
		return a, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if msg.IsSuccessful() {
		a.dialog = nil
		a.feedPanel, cmd = a.feedPanel.Update(msg)
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	}

	a.dialog, cmd = a.dialog.Update(msg)
	return a, cmd
}

func (a app) onDeleteFeedMsg(msg message.DeleteFeed) (app, tea.Cmd) {
	var cmd tea.Cmd

	if a.dialog == nil && msg.IsInitial() {
		a.dialog = dialog.NewDeleteFeed(a.cfg, a.repo)
		a.dialog, cmd = a.dialog.Update(msg)
		return a, cmd
	}

	// When add feed is in progress, user close add feed dialog,
	// we should not update dialog.
	if _, ok := a.dialog.(dialog.DeleteFeed); !ok {
		return a, nil
	}

	var cmds []tea.Cmd

	if msg.IsSuccessful() {
		a.dialog = nil
		a.feedPanel, cmd = a.feedPanel.Update(msg)
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	}

	a.dialog, cmd = a.dialog.Update(msg)
	return a, cmd
}

func (a app) onRenameFeedMsg(msg message.RenameFeed) (app, tea.Cmd) {
	var cmd tea.Cmd

	if a.dialog == nil && msg.IsInitial() {
		a.dialog = dialog.NewRenameFeed(a.cfg, a.repo)
		a.dialog, cmd = a.dialog.Update(msg)
		return a, cmd
	}

	// When add feed is in progress, user close add feed dialog,
	// we should not update dialog.
	if _, ok := a.dialog.(dialog.RenameFeed); !ok {
		return a, nil
	}

	var cmds []tea.Cmd

	if msg.IsSuccessful() {
		a.dialog = nil
		a.feedPanel, cmd = a.feedPanel.Update(msg)
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	}

	a.dialog, cmd = a.dialog.Update(msg)
	return a, cmd
}

func (a app) onLoadFeedsMsg(msg message.LoadFeeds) (app, tea.Cmd) {
	a.loadFeedsMsg = msg

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if msg.IsSuccessful() {
		a.feedPanel, cmd = a.feedPanel.Update(msg)
		cmds = append(cmds, cmd)

		cmd = message.RefreshTickCmd(a.cfg.RefreshInterval)
		cmds = append(cmds, cmd)

		cmd = message.RefreshCmd(msg.Feeds, a.repo)
		cmds = append(cmds, cmd)
	}
	return a, tea.Batch(cmds...)
}

func (a app) onSelectFeedMsg(msg message.SelectFeed) (app, tea.Cmd) {
	var cmd tea.Cmd
	a.itemPanel, cmd = a.itemPanel.Update(msg)
	return a, cmd
}

func (a app) onSelectFeedItemMsg(msg message.SelectFeedItem) (app, tea.Cmd) {
	var cmd tea.Cmd
	a.previewPanel, cmd = a.previewPanel.Update(msg)
	return a, cmd
}

func (a app) onToogleReadMsg(msg message.ToogleRead) (app, tea.Cmd) {
	if msg.IsFailed() {
		a.logger.Error("toogle read item failed",
			"item id", msg.ItemID, "err", msg.Err)
		return a, nil
	}

	if msg.IsSuccessful() {
		var cmd tea.Cmd
		a.feedPanel, cmd = a.feedPanel.Update(msg)
		return a, cmd
	}

	return a, nil
}

func (a app) onMarkAllReadMsg(msg message.MarkAllRead) (app, tea.Cmd) {
	if msg.IsFailed() {
		a.logger.Error("mark all read failed",
			"item ids", msg.ItemIDs, "err", msg.Err)
		return a, nil
	}

	if msg.IsSuccessful() {
		var cmd tea.Cmd
		a.feedPanel, cmd = a.feedPanel.Update(msg)
		return a, cmd
	}

	return a, nil
}

func (a app) onToogleStarredMsg(msg message.ToogleStarred) (app, tea.Cmd) {
	if msg.IsFailed() {
		a.logger.Error("toogle starred item failed",
			"item id", msg.ItemID, "err", msg.Err)
		return a, nil
	}

	if msg.IsSuccessful() {
		var cmd tea.Cmd
		a.feedPanel, cmd = a.feedPanel.Update(msg)
		return a, cmd
	}

	return a, nil
}

func (a app) onParseMDMsg(msg message.ParseMD) (app, tea.Cmd) {
	var cmd tea.Cmd
	a.previewPanel, cmd = a.previewPanel.Update(msg)
	return a, cmd
}

func (a app) onRefreshMsg(msg message.Refresh) (app, tea.Cmd) {
	a.refreshMsg = msg

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if msg.IsInProgress() {
		v := fmt.Sprintf("Refreshing %d of %d ...", len(msg.Results)+1, msg.Total)
		cmds = append(cmds, message.TipsCmd(v, false))
	}

	if len(msg.Results) == 0 {
		return a, tea.Batch(cmds...)
	}

	ret := msg.Results[len(msg.Results)-1]
	if ret.IsFailed() {
		a.logger.Error("refresh failed", "feed id", ret.Feed.ID, "err", ret.Err)
	}

	if ret.IsSuccessful() {
		cmd = a.feedPanel.UpdateFeed(ret.Feed)
		cmds = append(cmds, cmd)
	}

	if len(msg.Results) == msg.Total {
		a.refreshMsg = message.NewRefreshSuccessful(msg.Total, msg.Results)
		cmd = message.TipsCmd("Refresh finished", true)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a app) onRefreshTickMsg(_ message.RefreshTick) (app, tea.Cmd) {
	return a, message.RefreshTickCmd(a.cfg.RefreshInterval)
}

func (a app) onExportMsg(msg message.Export) (app, tea.Cmd) {
	if msg.IsInProgress() {
		return a, message.TipsCmd("Exporting...", false)
	}

	if msg.IsSuccessful() {
		return a, message.TipsCmd("Export successful: "+msg.FilePath, true)
	}

	if msg.IsFailed() {
		return a, message.ErrTipsCmd("Export failed", msg.Err, true)
	}

	return a, nil
}

func (a app) onImportMsg(msg message.Import) (app, tea.Cmd) {
	// When add feed is in progress, user close add feed dialog,
	// we should not update dialog.
	if _, ok := a.dialog.(dialog.Import); !ok {
		return a, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	if msg.IsSuccessful() {
		a.dialog = nil

		cmd = a.feedPanel.AddFeeds(msg.Feeds)
		cmds = append(cmds, cmd)

		cmd = message.RefreshCmd(msg.Feeds, a.repo)
		cmds = append(cmds, cmd)

		return a, tea.Batch(cmds...)
	}

	a.dialog, cmd = a.dialog.Update(msg)
	return a, cmd
}

func (a *app) onKeyMsg(msg tea.KeyMsg) tea.Cmd {
	if key.Matches(msg, a.cfg.KeyMap.Quit) {
		return tea.Quit
	}

	if key.Matches(msg, a.cfg.KeyMap.Esc) {
		a.dialog = nil
		a.statusBar.Hide()
		a.setSizes()
		return nil
	}

	// If feeds are loading, don't handle key msgs.
	if a.loadFeedsMsg.IsInProgress() {
		return nil
	}

	var cmd tea.Cmd

	// If dialog is displayed, key msgs should be sent to dialog.
	if a.dialog != nil {
		a.dialog, cmd = a.dialog.Update(msg)
		return cmd
	}

	if key.Matches(msg, a.cfg.KeyMap.OpenDir) {
		err := browser.OpenFile(a.dir)
		if err != nil {
			a.logger.Error("open dir failed", "err", err)
			return message.ErrTipsCmd("open dir failed", err, true)
		}
	}

	if key.Matches(msg, a.cfg.KeyMap.Help) {
		a.statusBar.Toggle()
		a.setSizes()
		return nil
	}

	if key.Matches(msg, a.cfg.KeyMap.PrevFocus) {
		a.focus = a.focus.prev()
		return a.updateFocus()
	}
	if key.Matches(msg, a.cfg.KeyMap.NextFocus) {
		a.focus = a.focus.next()
		return a.updateFocus()
	}

	if key.Matches(msg, a.cfg.KeyMap.AddFeed) {
		a.dialog = dialog.NewAddFeed(a.cfg, a.repo)
		return a.dialog.Init()
	}

	if key.Matches(msg, a.cfg.KeyMap.Refresh) {
		return a.onRefreshKeyMsg()
	}

	if key.Matches(msg, a.cfg.KeyMap.Export) {
		return a.onExportKeyMsg()
	}

	if key.Matches(msg, a.cfg.KeyMap.Import) {
		a.dialog = dialog.NewImport(a.cfg, a.repo)
		return a.dialog.Init()
	}

	switch a.focus {
	case focusFeed:
		a.feedPanel, cmd = a.feedPanel.Update(msg)
	case focusItem:
		a.itemPanel, cmd = a.itemPanel.Update(msg)
	case focusPreview:
		a.previewPanel, cmd = a.previewPanel.Update(msg)
	}
	return cmd
}

func (a *app) onRefreshKeyMsg() tea.Cmd {
	if a.refreshMsg.IsInProgress() {
		return nil
	}

	feeds := a.feedPanel.NormalFeeds()
	if len(feeds) == 0 {
		return nil
	}

	return message.RefreshCmd(feeds, a.repo)
}

func (a *app) onExportKeyMsg() tea.Cmd {
	feeds := a.feedPanel.NormalFeeds()
	if len(feeds) == 0 {
		return message.TipsCmd("No feeds to export", true)
	}

	return message.ExportCmd(feeds, a.dir)
}

func (a *app) updateFocus() tea.Cmd {
	a.feedPanel.SetFocused(a.focus == focusFeed)
	cmd := a.itemPanel.SetFocused(a.focus == focusItem)
	a.previewPanel.SetFocused(a.focus == focusPreview)
	return cmd
}

func (a *app) onWindowSizeMsg(msg tea.WindowSizeMsg) {
	a.windowWidth = msg.Width
	a.windowHeight = msg.Height
	a.setSizes()
}

func (a *app) setSizes() {
	height := a.windowHeight - lipgloss.Height(a.statusBar.View()) - view.BorderVerticalSize

	previewWidth := a.windowWidth - a.cfg.FeedPanelWidth - a.cfg.ItemPanelWidth -
		view.BorderHorizontalSize*3 //nolint:mnd // panels count

	a.feedPanel.SetSize(a.cfg.FeedPanelWidth, height)
	a.itemPanel.SetSize(a.cfg.ItemPanelWidth, height)
	a.previewPanel.SetSize(previewWidth, height)

	a.statusBar.SetWidth(a.windowWidth)
}

func (a app) View() string {
	if a.dialog != nil {
		return a.dialogView()
	}

	if a.loadFeedsMsg.IsInProgress() {
		return a.loadingView("Loading feeds...")
	}

	if a.loadFeedsMsg.IsFailed() {
		v := fmt.Sprintf("Load feeds failed: %v", a.loadFeedsMsg.Err)
		return a.loadingView(v)
	}

	if len(a.feedPanel.NormalFeeds()) == 0 {
		return a.loadingView("No Feeds")
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		a.feedPanel.View(),
		a.itemPanel.View(),
		a.previewPanel.View(),
	)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		a.statusBar.View(),
	)
}

func (a app) dialogView() string {
	return lipgloss.Place(
		a.windowWidth,
		a.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		a.dialog.View(),
	)
}

func (a app) loadingView(v ...string) string {
	l := lipgloss.NewStyle().Width(a.windowWidth).
		Foreground(a.cfg.Theme.Logo).
		Align(lipgloss.Center).Render(logo)
	v = append([]string{l}, v...)

	h := a.windowHeight - lipgloss.Height(a.statusBar.View())
	view := lipgloss.JoinVertical(lipgloss.Center, v...)
	view = lipgloss.NewStyle().Height(h).Width(a.windowWidth).Render(view)

	return lipgloss.JoinVertical(lipgloss.Center,
		view, a.statusBar.View())
}
