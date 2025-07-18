package dialog

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/lakerszhy/rssx/internal/config"
	"github.com/lakerszhy/rssx/internal/view"
)

const (
	dialogWidth  = 60
	buttonWidth  = 9
	buttonMargin = 2
)

func render(name, content string, theme *config.AppTheme) string {
	b := view.BorderWithTitle(name, dialogWidth)
	return view.BorderStyle(b, theme, true).
		Width(dialogWidth).Padding(1, 2). //nolint:mnd // horizontal padding
		Render(content)
}

func newTextInput(theme *config.AppTheme, placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(theme.TextInputPlaceholder)
	ti.Width = dialogWidth - 7 //nolint:mnd // horizontal padding + prompt length
	ti.Focus()
	ti.TextStyle = lipgloss.NewStyle().Foreground(theme.TextInput)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(theme.Cursor)
	ti.PromptStyle = lipgloss.NewStyle().Foreground(theme.TextInputPrompt)
	return ti
}

func inputView(ti textinput.Model, theme *config.AppTheme) string {
	return view.BorderStyle(view.Border, theme, false).
		Border(view.Border, false, false, true).
		Render(ti.View())
}

func actionsView(theme *config.AppTheme, isDanger bool) string {
	esc := lipgloss.NewStyle().
		Background(theme.CancelButtonBackground).Foreground(theme.CancelButton).
		Align(lipgloss.Center).MarginRight(buttonMargin).Width(buttonWidth).
		Render("Esc")

	enter := enterButton(theme, isDanger)

	actions := lipgloss.JoinHorizontal(lipgloss.Center, esc, enter)
	return lipgloss.NewStyle().Width(dialogWidth - buttonMargin*2).
		Align(lipgloss.Center).Render(actions)
}

func enterButton(theme *config.AppTheme, isDanger bool) string {
	background := theme.ConfirmButtonBackground
	foreground := theme.ConfirmButton
	if isDanger {
		background = theme.DangerButtonBackground
		foreground = theme.DangerButton
	}
	return lipgloss.NewStyle().
		Background(background).Foreground(foreground).
		Align(lipgloss.Center).MarginLeft(buttonMargin).Width(buttonWidth).
		Render("Enter")
}
