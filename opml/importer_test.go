package opml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	file := "./testdata/rssx.opml"
	feeds, err := Import(file)

	require.NoError(t, err)
	assert.Len(t, feeds, 113)

	for _, f := range feeds {
		assert.NotEmpty(t, f.Name)
		assert.NotEmpty(t, f.FeedURL)
	}
}
