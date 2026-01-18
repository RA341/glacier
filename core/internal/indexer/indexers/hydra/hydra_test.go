package hydra

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHydra(t *testing.T) {
	hydra, err := newRaw(map[string]any{
		"cacheDir":       "./tmp/cache",
		"updateInterval": "24h",
		"sources": map[string]string{
			"fitgirl": "https://hydralinks.pages.dev/sources/fitgirl.json",
		},
		"debug": true,
	})
	require.NoError(t, err)

	hydra.updateIndexes()

	search, err := hydra.Search("warhammer")
	require.NoError(t, err)

	t.Log(search)
}
