package download

import (
	"context"

	"github.com/ra341/glacier/frost/local_library/download"
	"github.com/ra341/glacier/internal/library"
	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, query string, limit, offset int) ([]LocalGame, error)
	ListWithState(ctx context.Context, status ...download.Status) ([]LocalGame, error)

	Get(ctx context.Context, id int) (LocalGame, error)
	Add(ctx context.Context, game *LocalGame) error
	Edit(ctx context.Context, game *LocalGame) error
	EditStatus(ctx context.Context, id int, down *download.LocalDownload) error
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, id int) error
	GetByGameId(ctx context.Context, id uint64, localDownload bool) (*LocalGame, error)
}

type GamePlay struct {
	// this is the final game exe that actually runs
	LaunchExe string
}

type LocalGame struct {
	gorm.Model

	GameId int
	Game   library.Game `gorm:"embedded"`

	Download download.LocalDownload `gorm:"embedded"`
	Play     GamePlay               `gorm:"embedded"`
}
