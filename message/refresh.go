package message

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/rss"
)

type RefreshTick time.Time

func RefreshTickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return RefreshTick(t)
	})
}

func RefreshCmd(feeds []rss.Feed, repo rss.Repo) tea.Cmd {
	total := len(feeds)
	if total == 0 {
		return nil
	}

	results := make([]FeedRefreshResult, 0, len(feeds))

	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmd = func() tea.Msg {
		return NewRefreshInProgress(total, results)
	}
	cmds = append(cmds, cmd)

	for _, f := range feeds {
		cmd = func() tea.Msg {
			results = append(results, refreshFeed(f, repo))
			return NewRefreshInProgress(total, results)
		}
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func refreshFeed(f rss.Feed, repo rss.Repo) FeedRefreshResult {
	newFeed, err := rss.ParseURL(f.FeedURL)
	if err != nil {
		return newFeedRefreshResultFailed(f, err)
	}

	// only insert new items.
	var newItems []rss.FeedItem
	for _, i := range newFeed.Items {
		found := false
		for _, j := range f.Items {
			if j.Link == i.Link {
				found = true
				break
			}
		}
		if !found {
			newItems = append(newItems, i)
		}
	}

	newItems, err = repo.InsertItems(f.ID, newItems)
	if err != nil {
		return newFeedRefreshResultFailed(f, err)
	}

	f.Items = append(f.Items, newItems...)
	return newFeedRefreshResultSuccessful(f)
}

type Refresh struct {
	Results []FeedRefreshResult
	Total   int
	status
}

func NewRefreshInProgress(total int, results []FeedRefreshResult) Refresh {
	return Refresh{
		Results: results,
		Total:   total,
		status:  statusInProgress,
	}
}

func NewRefreshSuccessful(total int, results []FeedRefreshResult) Refresh {
	return Refresh{
		Results: results,
		Total:   total,
		status:  statusSuccessful,
	}
}

type FeedRefreshResult struct {
	Feed rss.Feed
	Err  error
	status
}

func newFeedRefreshResultSuccessful(f rss.Feed) FeedRefreshResult {
	return FeedRefreshResult{
		Feed:   f,
		status: statusSuccessful,
	}
}

func newFeedRefreshResultFailed(f rss.Feed, err error) FeedRefreshResult {
	return FeedRefreshResult{
		Feed:   f,
		status: statusFailed,
		Err:    err,
	}
}
