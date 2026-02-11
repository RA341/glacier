package updater

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSrv(t *testing.T) {
	ctx := context.Background()

	srv := New()
	releases, err := srv.FetchReleases(ctx)
	require.NoError(t, err)

	t.Log(releases)
}
