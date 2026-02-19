package assets

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type StoreGorm struct {
	db *gorm.DB
}

func NewStoreGorm(db *gorm.DB) Store {
	return &StoreGorm{
		db: db,
	}
}

func (s *StoreGorm) GetGameIds(ctx context.Context) ([]uint, error) {
	var ids []uint
	err := s.db.WithContext(ctx).
		Model(&Asset{}).
		Distinct().
		Pluck("game_id", &ids).
		Error
	return ids, err
}

func (s *StoreGorm) List(ctx context.Context, gameId uint64, assetTypeStr ...string) ([]Asset, error) {
	var assets []Asset

	q := s.db.WithContext(ctx).Order("type DESC").Where("game_id = ?", gameId)
	if len(assetTypeStr) > 0 {
		var types []AssetType
		for _, str := range assetTypeStr {
			assetType, err := AssetTypeString(str)
			if err != nil {
				log.Warn().Err(err).Str("type", str).Msg("Failed to get asset type")
				continue
			}
			types = append(types, assetType)
		}
		if len(types) > 0 {
			q = q.Where("type IN ?", types)
		}
	}

	err := q.Find(&assets).Error
	return assets, err
}

func (s *StoreGorm) ListUnDownloaded(ctx context.Context) (map[uint][]Asset, error) {
	var assetList []Asset
	err := s.db.WithContext(ctx).
		Where("local_path IS NULL OR local_path = ''").
		Where("remote_url IS NOT NULL AND remote_url != ''").
		Find(&assetList).
		Error
	if err != nil {
		return nil, err
	}

	assets := make(map[uint][]Asset)
	for _, asset := range assetList {
		assets[asset.GameID] = append(assets[asset.GameID], asset)
	}

	return assets, nil
}

func (s *StoreGorm) Update(ctx context.Context, asset *Asset) error {
	return s.db.WithContext(ctx).Save(asset).Error
}

func (s *StoreGorm) Delete(ctx context.Context, asset *Asset) error {
	return s.db.WithContext(ctx).Unscoped().Delete(asset).Error
}

func (s *StoreGorm) Save(
	ctx context.Context,
	asset *Asset,
	removeIfExists func(oldAsset *Asset) error,
) error {
	return s.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			var oldAsset Asset

			err := tx.WithContext(ctx).
				Where("id = ?", asset.ID).
				First(&oldAsset).
				Error

			if err == nil {
				err := removeIfExists(&oldAsset)
				if err != nil {
					return err
				}
			}

			return tx.WithContext(ctx).Save(asset).Error
		})
}

func (s *StoreGorm) Get(ctx context.Context, id uint) (*Asset, error) {
	var asset Asset
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&asset).Error
	return &asset, err
}

func (s *StoreGorm) ListUndownloadedByGame(ctx context.Context, gameId uint) ([]Asset, error) {
	var gs []Asset

	err := s.db.WithContext(ctx).
		Where("game_id = ?", gameId).
		Where("local_path IS NULL OR local_path = ''").
		Where("remote_url IS NOT NULL AND remote_url != ''").
		Find(&gs).
		Error

	return gs, err
}

func (s *StoreGorm) GetById(ctx context.Context, id uint, assetType AssetType) (*Asset, error) {
	var asset Asset
	err := s.db.WithContext(ctx).
		Order("updated_at DESC").
		Where("game_id = ?", id).
		Where("type = ?", assetType).
		First(&asset).
		Error

	return &asset, err
}

func (s *StoreGorm) GetByGame(ctx context.Context, id uint) ([]Asset, error) {
	var assetList []Asset
	err := s.db.WithContext(ctx).Where("game_id = ?", id).Find(&assetList).Error
	return assetList, err
}
