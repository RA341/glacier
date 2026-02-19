package assets

import (
	"context"

	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, gameId uint64, assetTypeStr ...string) ([]Asset, error)
	ListUnDownloaded(ctx context.Context) (map[uint][]Asset, error)
	ListUndownloadedByGame(ctx context.Context, gameId uint) ([]Asset, error)

	Get(ctx context.Context, id uint) (*Asset, error)
	GetById(ctx context.Context, id uint, assetType AssetType) (*Asset, error)
	GetByGame(ctx context.Context, id uint) ([]Asset, error)
	GetGameIds(ctx context.Context) ([]uint, error)

	Update(ctx context.Context, asset *Asset) error
	Delete(ctx context.Context, asset *Asset) error
	Save(ctx context.Context, asset *Asset, exists func(oldAsset *Asset) error) error
}

//go:generate go run github.com/dmarkham/enumer@latest -type=AssetType -output=enum_asset_type.go
type AssetType int

const (
	AssetUnknown AssetType = iota
	AssetThumbnail
	AssetBanner
	AssetTrailer
	AssetGameplayImage
	AssetGameplayVideo
	AssetArtwork
)

type Asset struct {
	gorm.Model
	GameID uint `gorm:"index"`

	// Type helps distinguish: "thumbnail", "banner", "gameplay_video"
	Type      AssetType `gorm:"index"`
	RemoteURL string
	LocalPath string // /GameID/<hash>
}

// todo detele local files when asset is deleted
//func (g *Game) BeforeDelete(tx *gorm.DB) (err error) {
//	var assets []Asset
//	// Find all assets associated with this game
//	tx.Model(&Asset{}).Where("game_id = ?", g.ID).Find(&assets)
//
//	for _, asset := range assets {
//		// logic to delete the file from your local storage
//		if asset.LocalPath != "" {
//			os.Remove(asset.LocalPath)
//		}
//	}
//
//	// If you are soft-deleting, GORM won't trigger the SQL Cascade.
//	// We manually tell GORM to soft-delete the associated assets.
//	tx.Where("game_id = ?", g.ID).Delete(&Asset{})
//
//	return nil
//}

func NewAsset(gameID uint, remoteUrl string) *Asset {
	return &Asset{
		GameID:    gameID,
		RemoteURL: remoteUrl,
	}
}

func NewRemoteAsset(remoteUrl string, assetType AssetType) Asset {
	return Asset{
		RemoteURL: remoteUrl,
		Type:      assetType,
	}
}
