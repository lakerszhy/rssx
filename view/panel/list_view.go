package panel

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lakerszhy/rssx/config"
)

type listView[T list.Item] struct {
	model list.Model
}

func newListView[T list.Item](
	cfg *config.App,
	delegate list.ItemDelegate,
) listView[T] {
	model := list.New([]list.Item{}, delegate, 0, 0)
	model.SetShowTitle(false)
	model.SetShowStatusBar(false)
	model.SetFilteringEnabled(false)
	model.SetShowPagination(false)
	model.SetShowHelp(false)
	model.KeyMap = list.KeyMap{
		CursorUp:   cfg.KeyMap.Up,
		CursorDown: cfg.KeyMap.Down,
		GoToStart:  cfg.KeyMap.Start,
		GoToEnd:    cfg.KeyMap.End,
		PrevPage:   cfg.KeyMap.PrevPage,
		NextPage:   cfg.KeyMap.NextPage,
	}
	return listView[T]{
		model: model,
	}
}

func (l listView[T]) Init() tea.Cmd {
	return nil
}

func (l listView[T]) Update(msg tea.Msg) (listView[T], tea.Cmd) {
	var cmd tea.Cmd
	l.model, cmd = l.model.Update(msg)
	return l, cmd
}

func (l listView[T]) selectedItem() *T {
	if i, ok := l.model.SelectedItem().(T); ok {
		return &i
	}
	return nil
}

func (l listView[T]) index() int {
	return l.model.Index()
}

func (l *listView[T]) selectByIndex(i int) {
	l.model.Select(i)
}

func (l *listView[T]) setItems(v []T) {
	items := make([]list.Item, 0, len(v))
	for _, item := range v {
		items = append(items, item)
	}
	l.model.SetItems(items)
}

func (l listView[T]) items() []T {
	items := make([]T, 0, len(l.model.Items()))
	for _, item := range l.model.Items() {
		if i, ok := item.(T); ok {
			items = append(items, i)
		}
	}
	return items
}

func (l listView[T]) View() string {
	if len(l.items()) == 0 {
		return lipgloss.NewStyle().Width(l.model.Width()).
			Height(l.model.Height()).AlignHorizontal(lipgloss.Center).
			Render("No items.")
	}
	return l.model.View()
}

func (l listView[T]) footView() string {
	if len(l.model.Items()) > 0 {
		return fmt.Sprintf("%d/%d", l.model.Index()+1, len(l.model.Items()))
	}
	return ""
}

func (l *listView[T]) setSize(width, height int) {
	l.model.SetSize(width, height)
}

func (l *listView[T]) setDelegate(delegate list.ItemDelegate) {
	l.model.SetDelegate(delegate)
}
