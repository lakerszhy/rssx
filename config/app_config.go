package config

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	FeedPanelWidth  int
	ItemPanelWidth  int
	RefreshInterval time.Duration
	Theme           *AppTheme
	KeyMap          *keyMap
}

type keyMap struct {
	Up            key.Binding
	Down          key.Binding
	Start         key.Binding
	End           key.Binding
	PrevPage      key.Binding
	NextPage      key.Binding
	PrevFocus     key.Binding
	NextFocus     key.Binding
	AddFeed       key.Binding
	DeleteFeed    key.Binding
	ToogleStarred key.Binding
	ToogleRead    key.Binding
	MarkAllRead   key.Binding
	RenameFeed    key.Binding
	Refresh       key.Binding
	Open          key.Binding
	Export        key.Binding
	Import        key.Binding
	Enter         key.Binding
	Esc           key.Binding
	OpenDir       key.Binding
	Help          key.Binding
	Quit          key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PrevPage, k.NextPage, k.Start, k.End},
		{k.PrevFocus, k.NextFocus, k.ToogleStarred, k.ToogleRead, k.MarkAllRead, k.Refresh},
		{k.AddFeed, k.DeleteFeed, k.RenameFeed, k.Open, k.Export, k.Import},
		{k.Enter, k.Esc, k.OpenDir, k.Help, k.Quit},
	}
}

type AppTheme struct {
	Cursor                  lipgloss.Color
	DialogMsg               lipgloss.Color
	Border                  lipgloss.Color
	BorderActive            lipgloss.Color
	Starred                 lipgloss.Color
	Unread                  lipgloss.Color
	Error                   lipgloss.Color
	CancelButton            lipgloss.Color
	CancelButtonBackground  lipgloss.Color
	ConfirmButton           lipgloss.Color
	ConfirmButtonBackground lipgloss.Color
	DangerButton            lipgloss.Color
	DangerButtonBackground  lipgloss.Color
	HelpKey                 lipgloss.Color
	HelpKeyDesc             lipgloss.Color
	StatusBarBackground     lipgloss.Color
	Tips                    lipgloss.Color
	Help                    lipgloss.Color
	HelpBackground          lipgloss.Color
	TextInput               lipgloss.Color
	TextInputPlaceholder    lipgloss.Color
	TextInputPrompt         lipgloss.Color
	FeedTitle               lipgloss.Color
	FeedTitleActive         lipgloss.Color
	SmartFeed               lipgloss.Color
	SmartFeedActive         lipgloss.Color
	ItemTitle               lipgloss.Color
	ItemTitleActive         lipgloss.Color
	ItemDesc                lipgloss.Color
	ItemDescActive          lipgloss.Color
	Logo                    lipgloss.Color
}
