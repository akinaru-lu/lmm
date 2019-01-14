package persistence

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"lmm/api/service/asset/domain/model"
	"lmm/api/storage/db"
	"lmm/api/storage/uploader"
)

// AssetStorage s a AssetRepository implementation
type AssetStorage struct {
	db       db.DB
	uploader uploader.Uploader
}

// NewAssetStorage creates new AssetStorage
func NewAssetStorage(db db.DB, uploader uploader.Uploader) *AssetStorage {
	return &AssetStorage{db: db, uploader: uploader}
}

// FindByName implementation
func (s *AssetStorage) FindByName(c context.Context, name string) (*model.AssetDescriptor, error) {
	stmt := s.db.Prepare(c, "select name, type from asset where name = ? ")

	var (
		assetName     string
		assetTypeCode uint8
	)
	if err := stmt.QueryRow(c, name).Scan(&assetName, &assetTypeCode); err != nil {
		return nil, err
	}

	return model.NewAssetDescriptor(assetName, s.decodeAssetType(assetTypeCode)), nil
}

// Save implementation
func (s *AssetStorage) Save(c context.Context, asset *model.Asset) error {
	tx, err := s.db.Begin(c, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(c, `insert into asset (name, type, user, created_at) values (?, ?, ?, ?)`)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	if _, err := stmt.Exec(c,
		asset.Name(),
		s.encodeAssetType(asset.Type()),
		asset.Uploader().ID(),
		time.Now(),
	); err != nil {
		if e := tx.Rollback(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	if err := s.uploader.Upload(c,
		asset.Name(),
		asset.Data(),
		uploader.ImageUploaderConfig{
			Type: asset.Type().String(),
		},
	); err != nil {
		if e := tx.Rollback(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	return tx.Commit()
}

// Remove implementation
func (s *AssetStorage) Remove(c context.Context, asset *model.Asset) error {
	panic("not implemented")
}

func (s *AssetStorage) encodeAssetType(assetType model.AssetType) uint8 {
	switch assetType.String() {
	case "image":
		return 0
	case "photo":
		return 1
	default:
		panic("invalid asset type: '" + assetType.String() + "'")
	}
}

func (s *AssetStorage) decodeAssetType(code uint8) model.AssetType {
	switch code {
	case 0:
		return model.Image
	case 1:
		return model.Photo
	default:
		panic("invalid asset code: '" + string(code) + "'")
	}
}
