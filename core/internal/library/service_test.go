package library

import (
	"context"
	"testing"

	"github.com/ra341/glacier/internal/downloader/types"
	metaTypes "github.com/ra341/glacier/internal/metadata/types"
	"github.com/stretchr/testify/require"
)

var testGame = Game{
	Meta: metaTypes.Meta{
		Name: "test",
	},
	Download: types.Download{
		State:        types.Complete,
		DownloadPath: "tmp/game",
	},
}

func TestMeta(t *testing.T) {
	sd := TestStore{}
	srv := New(sd, nil, nil)
	ctx := context.Background()

	_, err := srv.GetDownloadManifest(ctx, 1)
	require.NoError(t, err)

}

////////////////////////////////////////////////////////////////////////////////////////////////////////

type TestStore struct {
}

func (t TestStore) GetById(ctx context.Context, id uint) (Game, error) {
	return testGame, nil
}

func (t TestStore) Add(ctx context.Context, game *Game) error {
	//TODO implement me
	panic("implement me")
}

func (t TestStore) UpdateDownloadProgress(ctx context.Context, id uint, download types.Download) error {
	//TODO implement me
	panic("implement me")
}

func (t TestStore) DeleteGame(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}

func (t TestStore) List(ctx context.Context, limit uint, offset uint) ([]Game, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestStore) ListDownloadState(ctx context.Context, state types.DownloadState) ([]Game, error) {
	//TODO implement me
	panic("implement me")
}
