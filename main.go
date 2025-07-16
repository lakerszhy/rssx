package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lakerszhy/rssx/config"
	"github.com/lakerszhy/rssx/store"
	_ "modernc.org/sqlite"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	dir, err := createDirs()
	if err != nil {
		return err
	}

	logFile, err := createLogFile(dir)
	if err != nil {
		return err
	}
	defer logFile.Close()
	logger := createLogger(logFile)

	cfg, err := config.Init(dir)
	if err != nil {
		return err
	}

	store, err := store.New(dir, logger)
	if err != nil {
		return err
	}
	defer store.Close()

	p := tea.NewProgram(newApp(dir, cfg, logger, store),
		tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err = p.Run(); err != nil {
		return err
	}

	return nil
}

func createDirs() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, "RssX")
	err = os.MkdirAll(dir, 0750)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func createLogFile(dir string) (*os.File, error) {
	logDir := filepath.Join(dir, "log")
	err := os.MkdirAll(logDir, 0750)
	if err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(filepath.Join(logDir, "rssx.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

func createLogger(w io.Writer) *slog.Logger {
	h := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(h)
}
