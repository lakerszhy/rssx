package config

import (
	"testing"

	"dario.cat/mergo"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	c1 := config{
		ThemeName:      "dark",
		FeedPanelWidth: 10,
	}
	c2 := config{
		ThemeName:      "light",
		ItemPanelWidth: 100,
	}
	err := mergo.Merge(&c1, c2,
		mergo.WithOverride)
	require.NoError(t, err)
	require.Equal(t, "light", c1.ThemeName)
	require.Equal(t, 100, c1.ItemPanelWidth)
	require.Equal(t, 10, c1.FeedPanelWidth)
}
