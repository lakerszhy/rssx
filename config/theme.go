package config

import (
	"github.com/charmbracelet/lipgloss"
)

type theme struct {
	Logo string `toml:"logo"`

	Border       string `toml:"border" comment:"\nBorder"`
	BorderActive string `toml:"border_active"`

	FeedTitle       string `toml:"feed_title" comment:"\nFeed Title"`
	FeedTitleActive string `toml:"feed_title_active"`
	SmartFeed       string `toml:"smart_feed"`
	SmartFeedActive string `toml:"smart_feed_active"`

	ItemTitle       string `toml:"item_title" comment:"\nFeed Item Title"`
	ItemTitleActive string `toml:"item_title_active"`
	ItemDesc        string `toml:"item_desc"`
	ItemDescActive  string `toml:"item_desc_active"`

	Starred string `toml:"starred" comment:"\nStarred Feed Item"`
	Unread  string `toml:"unread" comment:"Unread Feed Item"`

	TextInput            string `toml:"text_input" comment:"\nText Input"`
	TextInputPlaceholder string `toml:"text_input_placeholder"`
	TextInputPrompt      string `toml:"text_input_prompt"`
	Cursor               string `toml:"cursor"`
	DialogMsg            string `toml:"dialog_msg"`

	CancelButton            string `toml:"cancel_button" comment:"\nButtons"`
	CancelButtonBackground  string `toml:"cancel_button_background"`
	ConfirmButton           string `toml:"confirm_button"`
	ConfirmButtonBackground string `toml:"confirm_button_background"`
	DangerButton            string `toml:"danger_button"`
	DangerButtonBackground  string `toml:"danger_button_background"`

	Error string `toml:"error" comment:"\n"`

	StatusBarBackground string `toml:"status_bar_background" comment:"\nStatus Bar"`
	HelpKey             string `toml:"help_key"`
	HelpKeyDesc         string `toml:"help_key_desc"`
	Tips                string `toml:"tips"`
	Help                string `toml:"help"`
	HelpBackground      string `toml:"help_background"`
}

func (t theme) toApp() *AppTheme {
	return &AppTheme{
		Cursor:                  lipgloss.Color(t.Cursor),
		DialogMsg:               lipgloss.Color(t.DialogMsg),
		Border:                  lipgloss.Color(t.Border),
		BorderActive:            lipgloss.Color(t.BorderActive),
		Starred:                 lipgloss.Color(t.Starred),
		Unread:                  lipgloss.Color(t.Unread),
		Error:                   lipgloss.Color(t.Error),
		CancelButton:            lipgloss.Color(t.CancelButton),
		CancelButtonBackground:  lipgloss.Color(t.CancelButtonBackground),
		ConfirmButton:           lipgloss.Color(t.ConfirmButton),
		ConfirmButtonBackground: lipgloss.Color(t.ConfirmButtonBackground),
		DangerButton:            lipgloss.Color(t.DangerButton),
		DangerButtonBackground:  lipgloss.Color(t.DangerButtonBackground),
		HelpKey:                 lipgloss.Color(t.HelpKey),
		HelpKeyDesc:             lipgloss.Color(t.HelpKeyDesc),
		StatusBarBackground:     lipgloss.Color(t.StatusBarBackground),
		Help:                    lipgloss.Color(t.Help),
		HelpBackground:          lipgloss.Color(t.HelpBackground),
		Tips:                    lipgloss.Color(t.Tips),
		TextInput:               lipgloss.Color(t.TextInput),
		TextInputPlaceholder:    lipgloss.Color(t.TextInputPlaceholder),
		TextInputPrompt:         lipgloss.Color(t.TextInputPrompt),
		FeedTitle:               lipgloss.Color(t.FeedTitle),
		FeedTitleActive:         lipgloss.Color(t.FeedTitleActive),
		SmartFeed:               lipgloss.Color(t.SmartFeed),
		SmartFeedActive:         lipgloss.Color(t.SmartFeedActive),
		ItemTitle:               lipgloss.Color(t.ItemTitle),
		ItemTitleActive:         lipgloss.Color(t.ItemTitleActive),
		ItemDesc:                lipgloss.Color(t.ItemDesc),
		ItemDescActive:          lipgloss.Color(t.ItemDescActive),
		Logo:                    lipgloss.Color(t.Logo),
	}
}
