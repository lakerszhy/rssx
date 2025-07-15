package config

import (
	"embed"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"

	"dario.cat/mergo"
	"github.com/pelletier/go-toml/v2"
)

//go:embed config.toml
var configFS embed.FS

//go:embed hotkey.toml
var hotkeyFS embed.FS

//go:embed theme/*.toml
var themeFS embed.FS

const (
	hotkeyFileName = "hotkey.toml"
	configFileName = "config.toml"
)

func Init(dir string) (*App, error) {
	hotKey, err := loadX[hotkey](dir, hotkeyFS, hotkeyFileName)
	if err != nil {
		return nil, err
	}

	cfg, err := loadX[config](dir, configFS, configFileName)
	if err != nil {
		return nil, err
	}
	cfg.Hotkey = hotKey

	// embed path use / for all os
	themeFileName := path.Join("theme", cfg.ThemeName+".toml")
	theme, err := loadX[theme](dir, themeFS, themeFileName)
	if err != nil {
		return nil, err
	}
	cfg.Theme = theme

	return cfg.toApp(), nil
}

func loadX[T any](dir string, embedFS embed.FS, filename string) (*T, error) {
	// Copy .toml to user dir
	if err := os.CopyFS(dir, embedFS); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			return nil, err
		}
	}

	// Load user config
	userData, err := os.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		return nil, err
	}

	var userConfig T
	err = toml.Unmarshal(userData, &userConfig)
	if err != nil {
		return nil, err
	}

	// Load default config
	defaultData, err := embedFS.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var defaultConfig T
	err = toml.Unmarshal(defaultData, &defaultConfig)
	if err != nil {
		return nil, err
	}

	// Merge user config to default config
	err = mergo.Merge(&defaultConfig, userConfig, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	b, err := toml.Marshal(defaultConfig)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(filepath.Join(dir, filename), b, 0600)
	if err != nil {
		return nil, err
	}

	return &defaultConfig, nil
}

type config struct {
	ThemeName       string  `toml:"theme" comment:"Theme name"` //nolint:golines
	FeedPanelWidth  int     `toml:"feed_panel_width" comment:"\nWidth of feed panel"`
	ItemPanelWidth  int     `toml:"item_panel_width" comment:"\nWidth of item panel"`
	RefreshInterval int     `toml:"refresh_interval" comment:"\nAuto refresh interval in minutes"`
	Hotkey          *hotkey `toml:"-"`
	Theme           *theme  `toml:"-"`
}

func (c config) toApp() *App {
	return &App{
		RefreshInterval: time.Duration(c.RefreshInterval) * time.Minute,
		FeedPanelWidth:  c.FeedPanelWidth,
		ItemPanelWidth:  c.ItemPanelWidth,
		Theme:           c.Theme.toApp(),
		KeyMap:          c.Hotkey.toApp(),
	}
}
