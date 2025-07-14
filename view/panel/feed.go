package panel

import (
	"fmt"
	"log/slog"
	"slices"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/config"
	"github.com/lakerszhy/rssx/message"
	"github.com/lakerszhy/rssx/rss"
	"github.com/lakerszhy/rssx/view"
	"github.com/lakerszhy/rssx/view/delegate"
	"github.com/pkg/browser"
)

type Feed struct {
	cfg       *config.App
	logger    *slog.Logger
	listView  listView[rss.Feed]
	isFocused bool
}

func NewFeed(cfg *config.App, logger *slog.Logger) Feed {
	return Feed{
		cfg:       cfg,
		logger:    logger,
		listView:  newListView[rss.Feed](cfg, delegate.NewFeed(cfg.Theme)),
		isFocused: true,
	}
}

func (p Feed) Init() tea.Cmd {
	return nil
}

func (p Feed) Update(msg tea.Msg) (Feed, tea.Cmd) {
	switch msg := msg.(type) {
	case message.LoadFeeds:
		return p, p.setFeeds(msg.Feeds)
	case message.AddFeed:
		return p, p.addFeed(msg.Feed)
	case message.ToogleRead:
		return p, p.onToogleRead(msg)
	case message.ToogleStarred:
		return p, p.onToogleStarred(msg)
	case message.DeleteFeed:
		return p, p.onDeleteFeed(msg)
	case message.RenameFeed:
		return p, p.onRenameFeed(msg)
	case tea.KeyMsg:
		if key.Matches(msg, p.cfg.KeyMap.DeleteFeed) {
			return p, p.onDeleteFeedKeyMsg()
		}
		if key.Matches(msg, p.cfg.KeyMap.RenameFeed) {
			return p, p.onRenameFeedKeyMsg()
		}
		if key.Matches(msg, p.cfg.KeyMap.Open) {
			p.onOpenKeyMsg()
			return p, nil
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	p.listView, cmd = p.listView.Update(msg)
	cmds = append(cmds, cmd)

	cmd = func() tea.Msg {
		return message.NewSelectFeed(p.listView.selectedItem())
	}
	cmds = append(cmds, cmd)

	// if select smart feed, update smart feeds
	if i := p.listView.selectedItem(); i != nil && i.IsSmart() {
		p.updateSmartFeeds(p.listView.items())
	}

	return p, tea.Batch(cmds...)
}

func (p *Feed) UpdateFeed(f rss.Feed) tea.Cmd {
	feeds := p.listView.items()
	idx := slices.IndexFunc(feeds, func(i rss.Feed) bool {
		return f.ID == i.ID
	})
	if idx >= 0 {
		feeds[idx] = f
	}

	p.listView.setItems(feeds)
	p.updateSmartFeeds(feeds)
	return func() tea.Msg {
		return message.NewSelectFeed(p.listView.selectedItem())
	}
}

func (p *Feed) AddFeeds(feeds []rss.Feed) tea.Cmd {
	feeds = append(p.listView.items(), feeds...)
	return p.setFeeds(feeds)
}

func (p *Feed) setFeeds(feeds []rss.Feed) tea.Cmd {
	sort.SliceStable(feeds, func(i, j int) bool {
		return strings.ToLower(feeds[i].Name) < strings.ToLower(feeds[j].Name)
	})

	p.updateSmartFeeds(feeds)

	var cmd tea.Cmd
	if len(feeds) > 0 {
		p.listView.selectByIndex(3) //nolint:mnd // smart feeds length
		cmd = func() tea.Msg {
			return message.NewSelectFeed(&feeds[0])
		}
	}
	return cmd
}

func (p *Feed) addFeed(f rss.Feed) tea.Cmd {
	feeds := p.listView.items()
	feeds = append(feeds, f)
	p.setFeeds(feeds)
	// After setFeeds, smart feeds will be updated.
	// So have to get the correct index of the new feed.
	feeds = p.listView.items()

	idx := slices.IndexFunc(feeds, func(i rss.Feed) bool {
		return f.ID == i.ID
	})
	if idx >= 0 {
		p.listView.selectByIndex(idx)
	}

	cmd := func() tea.Msg {
		return message.NewSelectFeed(&f)
	}
	return cmd
}

func (p *Feed) onToogleRead(msg message.ToogleRead) tea.Cmd {
	cmd := p.update(func(f *rss.Feed) {
		f.ToogleRead(msg.ItemID)
	})

	// if select smart feed, don't update.
	// Otherwise, item in item panel will be gone.
	if i := p.listView.selectedItem(); i != nil && !i.IsSmart() {
		p.updateSmartFeeds(p.listView.items())
	}

	return cmd
}

func (p *Feed) onToogleStarred(msg message.ToogleStarred) tea.Cmd {
	cmd := p.update(func(f *rss.Feed) {
		f.ToogleStarred(msg.ItemID)
	})

	if i := p.listView.selectedItem(); i == nil && !i.IsSmart() {
		p.updateSmartFeeds(p.listView.items())
	}

	return cmd
}

func (p *Feed) onDeleteFeed(msg message.DeleteFeed) tea.Cmd {
	var feeds []rss.Feed
	for _, f := range p.listView.items() {
		if f.ID != msg.Feed.ID {
			feeds = append(feeds, f)
		}
	}
	p.updateSmartFeeds(feeds)

	if p.listView.selectedItem() == nil && len(feeds) > 0 {
		p.listView.selectByIndex(len(feeds) - 1)
	}

	cmd := func() tea.Msg {
		return message.NewSelectFeed(p.listView.selectedItem())
	}

	return cmd
}

func (p *Feed) onRenameFeed(msg message.RenameFeed) tea.Cmd {
	return p.update(func(f *rss.Feed) {
		if f.ID == msg.Feed.ID {
			f.Rename(msg.Feed.Name)
		}
	})
}

func (p Feed) onDeleteFeedKeyMsg() tea.Cmd {
	var cmd tea.Cmd
	if i := p.listView.selectedItem(); i != nil && !i.IsSmart() {
		cmd = func() tea.Msg {
			return message.NewDeleteFeedInitial(*i)
		}
	}
	return cmd
}

func (p Feed) onRenameFeedKeyMsg() tea.Cmd {
	var cmd tea.Cmd
	if i := p.listView.selectedItem(); i != nil && !i.IsSmart() {
		cmd = func() tea.Msg {
			return message.NewRenameFeedInitial(*i)
		}
	}
	return cmd
}

func (p Feed) onOpenKeyMsg() {
	i := p.listView.selectedItem()
	if i == nil || i.IsSmart() {
		return
	}

	if err := browser.OpenURL(i.HomePageURL); err != nil {
		p.logger.Error("open url", "link", i.HomePageURL, "err", err)
	}
}

func (p *Feed) update(fn func(f *rss.Feed)) tea.Cmd {
	feeds := make([]rss.Feed, 0, len(p.listView.items()))
	for _, f := range p.listView.items() {
		fn(&f)
		feeds = append(feeds, f)
	}
	p.listView.setItems(feeds)

	cmd := func() tea.Msg {
		return message.NewSelectFeed(p.listView.selectedItem())
	}
	return cmd
}

func (p *Feed) updateSmartFeeds(feeds []rss.Feed) {
	normalFeeds := []rss.Feed{}

	todayFeed := rss.NewTodayFeed()
	unreadFeed := rss.NewUnreadFeed()
	starredFeed := rss.NewStarredFeed()

	for _, f := range feeds {
		if f.IsSmart() {
			continue
		}
		normalFeeds = append(normalFeeds, f)
		for _, i := range f.Items {
			i.FeedName = f.Name
			if i.IsToday() {
				todayFeed.Items = append(todayFeed.Items, i)
			}
			if !i.IsRead {
				unreadFeed.Items = append(unreadFeed.Items, i)
			}
			if i.IsStarred {
				starredFeed.Items = append(starredFeed.Items, i)
			}
		}
	}

	slices.SortFunc(todayFeed.Items, func(a, b rss.FeedItem) int {
		return b.PublishedAt.Compare(a.PublishedAt)
	})
	slices.SortFunc(unreadFeed.Items, func(a, b rss.FeedItem) int {
		return b.PublishedAt.Compare(a.PublishedAt)
	})
	slices.SortFunc(starredFeed.Items, func(a, b rss.FeedItem) int {
		return b.PublishedAt.Compare(a.PublishedAt)
	})

	allFeeds := []rss.Feed{todayFeed, unreadFeed, starredFeed}
	allFeeds = append(allFeeds, normalFeeds...)
	p.listView.setItems(allFeeds)
}

func (p Feed) NormalFeeds() []rss.Feed {
	return slices.DeleteFunc(p.listView.items(), func(i rss.Feed) bool {
		return i.IsSmart()
	})
}

func (p *Feed) SetFocused(focused bool) {
	p.isFocused = focused
}

func (p *Feed) SetSize(width, height int) {
	p.listView.setSize(width, height)
}

func (p Feed) View() string {
	b := view.BorderWithFoot(p.footView(), p.cfg.FeedPanelWidth)
	style := view.BorderStyle(b, p.cfg.Theme, p.isFocused)
	return style.Render(p.listView.View())
}

func (p Feed) footView() string {
	total := len(p.NormalFeeds())
	if total == 0 {
		return ""
	}

	index := p.listView.index() - 3 //nolint:mnd // smart feeds length
	if index < 0 || index > total {
		return ""
	}
	return fmt.Sprintf("%d/%d", index+1, total)
}
