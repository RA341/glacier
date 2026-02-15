package assets

import (
	"context"

	"gorm.io/gorm"
)

type StoreGorm struct {
	db *gorm.DB
}

func NewStoreGorm(db *gorm.DB) *StoreGorm {
	return &StoreGorm{
		db: db,
	}
}

func (s *StoreGorm) List(ctx context.Context) ([]Asset, error) {
	var assets []Asset
	err := s.db.WithContext(ctx).Find(&assets).Error
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (s *StoreGorm) ListUnDownloaded(ctx context.Context) ([]Asset, error) {
	var assets []Asset
	err := s.db.WithContext(ctx).
		Where("local_path IS NULL OR local_path = ''").
		Where("remote_url IS NOT NULL AND remote_url != ''").
		Find(&assets).
		Error

	return assets, err
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
