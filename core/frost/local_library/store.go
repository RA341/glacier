package download

import (
	"context"

	"github.com/ra341/glacier/frost/local_library/download"
	v1 "github.com/ra341/glacier/generated/frost_library/v1"
	metadata "github.com/ra341/glacier/internal/metadata/types"
	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, query string, limit, offset int) ([]LocalGame, error)
	ListWithState(ctx context.Context, status download.Status) ([]LocalGame, error)

	Get(ctx context.Context, id int) (LocalGame, error)
	Add(ctx context.Context, game *LocalGame) error
	Edit(ctx context.Context, id int, game *LocalGame) error
	EditStatus(ctx context.Context, id int, Status download.Status, StatusMessage string) error
	Delete(ctx context.Context, id int) error
}

type LocalGame struct {
	gorm.Model

	Game          metadata.Meta `gorm:"embedded"`
	GameId        int
	DownloadPath  string
	InstallerPath string
	ExePath       string

	Status        download.Status
	StatusMessage string
}

func (g *LocalGame) ToProto() *v1.LocalGame {
	return &v1.LocalGame{
		ID:            uint64(g.ID),
		DownloadPath:  g.DownloadPath,
		InstallerPath: g.InstallerPath,
		ExePath:       g.ExePath,
		Status:        g.Status.String(),
		StatusMessage: g.StatusMessage,
	}
}
