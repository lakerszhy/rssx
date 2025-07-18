package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/lakerszhy/rssx/internal/config"
)

var (
	Border               = lipgloss.RoundedBorder()
	BorderHorizontalSize = Border.GetLeftSize() + Border.GetRightSize()
	BorderVerticalSize   = Border.GetTopSize() + Border.GetBottomSize()
)

func BorderStyle(b lipgloss.Border, theme *config.AppTheme, isFocused bool) lipgloss.Style {
	color := theme.Border
	if isFocused {
		color = theme.BorderActive
	}
	return lipgloss.NewStyle().Border(b).BorderForeground(color)
}

func BorderWithTitle(title string, width int) lipgloss.Border {
	info := fmt.Sprintf("%s%s%s", Border.MiddleRight, title, Border.MiddleLeft)
	repeatCount := max(width-ansi.StringWidth(info)-ansi.StringWidth(Border.Top), 0)
	end := strings.Repeat(Border.Top, repeatCount)
	top := Border.Top + info + end

	b := Border
	b.Top = top
	return b
}

func BorderWithFoot(foot string, width int) lipgloss.Border {
	if foot == "" {
		return Border
	}

	info := fmt.Sprintf("%s%s%s", Border.MiddleRight, foot, Border.MiddleLeft)
	repeatCount := max(width-ansi.StringWidth(info)-ansi.StringWidth(Border.Bottom), 0)
	start := strings.Repeat(Border.Bottom, repeatCount)
	bottom := start + info + Border.Bottom

	b := Border
	b.Bottom = bottom
	return b
}
