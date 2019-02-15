package application

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"lmm/api/http"
	"lmm/api/service/asset/application/command"
	"lmm/api/service/asset/domain"
	"lmm/api/service/asset/domain/model"
	"lmm/api/service/asset/domain/repository"
	"lmm/api/service/asset/domain/service"
	"lmm/api/util/stringutil"
)

var (
	ErrInvalidPage    = errors.New("invalid page")
	ErrInvalidPerPage = errors.New("invalid perPage")
)

// Service struct
type Service struct {
	uploaderService service.UploaderService
	assetRepository repository.AssetRepository
	imageService    service.ImageService
	imageEncoder    service.ImageEncoder
	assetFinder     service.AssetFinder
	cacheService    CacheService
}

// NewService creates a new image application service
func NewService(
	assetFinder service.AssetFinder,
	assetRepository repository.AssetRepository,
	cacheService CacheService,
	imageService service.ImageService,
	imageEncoder service.ImageEncoder,
	uploaderService service.UploaderService,
) *Service {
	return &Service{
		assetFinder:     assetFinder,
		assetRepository: assetRepository,
		cacheService:    cacheService,
		imageService:    imageService,
		imageEncoder:    imageEncoder,
		uploaderService: uploaderService,
	}
}

// UploadAsset handles upload asset command
func (app *Service) UploadAsset(c context.Context, cmd *command.UploadAsset) error {
	uploader, err := app.uploaderService.FromUserID(c, cmd.UserID())
	if err != nil {
		return errors.Wrap(err, cmd.UserID())
	}

	t := cmd.Type()
	switch t {
	case model.Image, model.Photo:
		if err := app.uploadImage(c, uploader, t, cmd.Data()); err != nil {
			return err
		}
	default:
		return errors.Wrap(domain.ErrUnsupportedAssetType, t.String())
	}

	if err := app.cacheService.ClearPhotos(c); err != nil {
		http.Log().Warn(c, err.Error())
	}
	return nil
}

func (app *Service) uploadImage(c context.Context, uploader *model.Uploader, assetType model.AssetType, data []byte) error {
	dst, ext, err := app.imageEncoder.Encode(c, data)
	if err != nil {
		if err == domain.ErrUnsupportedImageFormat {
			return errors.Wrap(err, ext)
		}
		return err
	}

	name := base64.URLEncoding.EncodeToString([]byte(uuid.NewMD5(uuid.New(), dst).String()))
	asset := model.NewAsset(assetType, name+"."+ext, uploader, dst)

	return app.assetRepository.Save(c, asset)
}

// ListImages lists images by given page and perPage
func (app *Service) ListImages(c context.Context, pageStr, perPageStr string) (*model.ImageCollection, error) {
	page, perPage, err := app.parseLimitAndCursorOrDefault(pageStr, perPageStr)
	if err != nil {
		return nil, err
	}

	return app.assetFinder.FindAllImages(c, page, perPage)
}

// ListPhotos lists images by given page and perPage
func (app *Service) ListPhotos(c context.Context, pageStr, perPageStr string) (*model.PhotoCollection, error) {
	page, perPage, err := app.parseLimitAndCursorOrDefault(pageStr, perPageStr)
	if err != nil {
		return nil, err
	}

	if photos, ok := app.cacheService.FetchPhotos(c, page, perPage); ok {
		return photos, nil
	}

	photos, err := app.assetFinder.FindAllPhotos(c, page, perPage*perPage+1)
	if err != nil {
		return nil, err
	}

	if len(photos.List()) == 0 {
		return photos, nil
	}

	if err := app.cacheService.StorePhotos(c, page, perPage, photos.List()); err != nil {
		http.Log().Warn(c, err.Error())
	}

	hasNextPage := len(photos.List()) > int(perPage)
	if hasNextPage {
		photos = model.NewPhotoCollection(photos.List()[:perPage], hasNextPage)
	}
	return photos, nil
}

func (app *Service) parseLimitAndCursorOrDefault(pageStr, perPageStr string) (uint, uint, error) {
	page, err := stringutil.ParseUint(pageStr)
	if err != nil {
		return 0, 0, errors.Wrap(ErrInvalidPage, err.Error())
	}
	if page < 1 {
		return 0, 0, errors.Wrap(ErrInvalidPage, "page can not be less than 1")
	}

	perPage, err := stringutil.ParseUint(perPageStr)
	if err != nil {
		return 0, 0, errors.Wrap(ErrInvalidPerPage, err.Error())
	}

	return page, perPage, nil
}

func (app *Service) SetPhotoAlternateTexts(c context.Context, cmd *command.SetImageAlternateTexts) error {
	asset, err := app.assetRepository.FindAssetByName(c, cmd.ImageName())
	if err != nil {
		return errors.Wrap(domain.ErrNoSuchAsset, err.Error())
	}

	if asset.Type() != model.Photo {
		return domain.ErrInvalidTypeNotAPhoto
	}

	alts := make([]*model.Alt, len(cmd.AltNames()))
	for i, name := range cmd.AltNames() {
		alts[i] = model.NewAlt(asset.Name(), name)
	}

	if err := app.imageService.SetAlt(c, asset, alts); err != nil {
		return err
	}

	if err := app.cacheService.ClearPhotos(c); err != nil {
		http.Log().Warn(c, err.Error())
	}

	return nil
}

// GetPhotoDescription get photo's description
func (app *Service) GetPhotoDescription(c context.Context, name string) (*model.PhotoDescriptor, error) {
	photo, err := app.assetRepository.FindPhotoByName(c, name)
	if err != nil {
		return nil, errors.Wrap(domain.ErrNoSuchPhoto, err.Error())
	}
	return photo, nil
}
