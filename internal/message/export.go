package message

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/internal/opml"
	"github.com/lakerszhy/rssx/internal/rss"
)

func ExportCmd(feeds []rss.Feed, dir string) tea.Cmd {
	var cmds []tea.Cmd

	cmd := func() tea.Msg {
		return NewExportInProgress()
	}
	cmds = append(cmds, cmd)

	p := filepath.Join(dir, "rssx.opml")
	cmd = func() tea.Msg {
		err := opml.Export(feeds, p)
		if err != nil {
			return NewExportFailed(err)
		}
		return NewExportSuccessful(p)
	}
	cmds = append(cmds, cmd)

	return tea.Sequence(cmds...)
}

type Export struct {
	FilePath string
	status
	Err error
}

func NewExportInProgress() Export {
	return Export{
		status: statusInProgress,
	}
}

func NewExportSuccessful(filePath string) Export {
	return Export{
		FilePath: filePath,
		status:   statusSuccessful,
	}
}

func NewExportFailed(err error) Export {
	return Export{
		status: statusFailed,
		Err:    err,
	}
}
