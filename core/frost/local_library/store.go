package download

import (
	"context"

	"github.com/ra341/glacier/frost/local_library/download"
	v1 "github.com/ra341/glacier/generated/frost_library/v1"
	"github.com/ra341/glacier/internal/library"
	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, query string, limit, offset int) ([]LocalGame, error)
	ListWithState(ctx context.Context, status ...download.Status) ([]LocalGame, error)

	Get(ctx context.Context, id int) (LocalGame, error)
	Add(ctx context.Context, game *LocalGame) error
	Edit(ctx context.Context, id int, game *LocalGame) error
	EditStatus(ctx context.Context, id int, down *download.Info) error
	Delete(ctx context.Context, id int) error
}

type GamePlay struct {
	InstallerPath string
	ExePath       string
}

type LocalGame struct {
	gorm.Model

	GameId int
	Game   library.Game `gorm:"embedded"`

	Download download.Info `gorm:"embedded"`
	Play     GamePlay      `gorm:"embedded"`
}

func (g *LocalGame) ToProto() *v1.LocalGame {
	return &v1.LocalGame{
		ID:            uint64(g.ID),
		DownloadPath:  g.Download.DownloadPath,
		Status:        g.Download.Status.String(),
		StatusMessage: g.Download.StatusMessage,
	}
}
