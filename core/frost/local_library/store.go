package download

import (
	"context"

	"github.com/ra341/glacier/frost/local_library/download"
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

	Game metadata.Meta `gorm:"embedded"`

	DownloadPath  string
	InstallerPath string
	ExePath       string

	Status        download.Status
	StatusMessage string
}
