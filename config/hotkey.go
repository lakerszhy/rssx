package config

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

type hotkey struct {
	Up            []string `toml:"up" comment:"Move up"` //nolint:golines
	Down          []string `toml:"down" comment:"Move down"`
	Start         []string `toml:"start" comment:"Go to start"`
	End           []string `toml:"end" comment:"Go to end"`
	PrevPage      []string `toml:"prev_page" comment:"\nMove to previous page"`
	NextPage      []string `toml:"next_page" comment:"Move to next page"`
	PrevFocus     []string `toml:"prev_focus" comment:"Focus on previous panel"`
	NextFocus     []string `toml:"next_focus" comment:"Focus on next panel"`
	AddFeed       []string `toml:"add_feed" comment:"\nAdd feed"`
	DeleteFeed    []string `toml:"delete_feed" comment:"Delete feed"`
	ToogleStarred []string `toml:"toogle_starred" comment:"Toogle starred status"`
	ToogleRead    []string `toml:"toogle_read" comment:"Toogle read status"` //nolint:golines
	MarkAllRead   []string `toml:"mark_all_read" comment:"Mark all items as read"`
	RenameFeed    []string `toml:"rename_feed" comment:"Rename feed"`
	Refresh       []string `toml:"refresh" comment:"Refresh feeds"`
	Open          []string `toml:"open" comment:"Open in browser"`
	Export        []string `toml:"export" comment:"Export OPML"`
	Import        []string `toml:"import" comment:"Import OPML"`
	Enter         []string `toml:"enter" comment:"\nEnter panel"`
	Esc           []string `toml:"esc" comment:"Cancel"`
	OpenDir       []string `toml:"open_dir" comment:"Open dir"`
	Help          []string `toml:"help" comment:"Show help"`
	Quit          []string `toml:"quit" comment:"Quit app"`
}

func (h hotkey) toApp() *keyMap {
	return &keyMap{
		Up:            newBinding(h.Up, "move up"),
		Down:          newBinding(h.Down, "move down"),
		Start:         newBinding(h.Start, "go to start"),
		End:           newBinding(h.End, "go to end"),
		PrevPage:      newBinding(h.PrevPage, "prev page"),
		NextPage:      newBinding(h.NextPage, "next page"),
		PrevFocus:     newBinding(h.PrevFocus, "prev focus"),
		NextFocus:     newBinding(h.NextFocus, "next focus"),
		AddFeed:       newBinding(h.AddFeed, "add feed"),
		DeleteFeed:    newBinding(h.DeleteFeed, "delete feed"),
		ToogleStarred: newBinding(h.ToogleStarred, "toogle starred"),
		ToogleRead:    newBinding(h.ToogleRead, "toogle read"),
		MarkAllRead:   newBinding(h.MarkAllRead, "mark all items as read"),
		RenameFeed:    newBinding(h.RenameFeed, "rename feed"),
		Refresh:       newBinding(h.Refresh, "refresh feed"),
		Open:          newBinding(h.Open, "open in browser"),
		Export:        newBinding(h.Export, "export OPML"),
		Import:        newBinding(h.Import, "import OPML"),
		Enter:         newBinding(h.Enter, "confirm"),
		Esc:           newBinding(h.Esc, "cancel"),
		OpenDir:       newBinding(h.OpenDir, "open dir"),
		Help:          newBinding(h.Help, "help"),
		Quit:          newBinding(h.Quit, "quit"),
	}
}

func newBinding(keys []string, desc string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(strings.Join(keys, "/"), desc),
	)
}
