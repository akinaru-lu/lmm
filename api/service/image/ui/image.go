package ui

import (
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"

	"lmm/api/http"
	account "lmm/api/service/account/domain/model"
	"lmm/api/service/image/application"
	"lmm/api/service/image/domain"
	"lmm/api/service/image/domain/repository"
)

const (
	maxFormDataSize                  = 32 << 20 // 32MB
	maxImageSize                     = 2 << 20  // 2MB
	errUploadImagesMaxNumberExceeded = "the number of uploaded images cannot be greater than 10 once"
)

var (
	errImageMaxSizeExceeded = errors.New("the size of image to upload is up to 2MB")
	errNotAllowedImageType  = errors.New("only gif, jpeg, png allowed")
)

type UI struct {
	app *application.AppService
}

func New(imageRepo repository.ImageRepository) *UI {
	app := application.NewAppService(imageRepo)
	return &UI{app: app}
}

func (ui *UI) Upload(c http.Context) {
	user, ok := c.Value(http.StrCtxKey("user")).(*account.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	if err := c.Request().ParseMultipartForm(maxFormDataSize); err != nil {
		log.Println(err.Error())
		http.BadRequest(c)
		return
	}

	imageSources := c.Request().MultipartForm.File["src"]
	if len(imageSources) > 10 {
		c.String(http.StatusBadRequest, errUploadImagesMaxNumberExceeded)
		return
	}

	for _, src := range imageSources {
		data, err := openImage(src)
		if err != nil {
			panic(err)
		}

		if err := ui.app.UploadImage(user, data); err != nil {
			panic(err.Error())
		}
	}
}

func openImage(fh *multipart.FileHeader) ([]byte, error) {
	// check type
	contentType := fh.Header.Get("Content-Type")
	switch contentType {
	case "image/gif", "image/jpeg", "image/png":
	default:
		return nil, errNotAllowedImageType
	}

	// check size
	if fh.Size > maxImageSize {
		return nil, errImageMaxSizeExceeded
	}

	// open file
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func (ui *UI) LoadImagesByPage(c http.Context) {
	models, hasNextPage, err := ui.app.FetchImagesByType(
		c.Request().QueryParam("type"),
		c.Request().QueryParam("count"),
		c.Request().QueryParam("page"),
	)
	switch err {
	case nil:
		images := make([]Image, len(models))
		for index, model := range models {
			images[index].Name = model.ID()
		}
		c.JSON(http.StatusOK, ImagesPage{
			Images:      images,
			HasNextPage: hasNextPage,
		})
	case domain.ErrInvalidPage, domain.ErrInvalidCount:
		c.String(http.StatusBadRequest, err.Error())
	default:
		panic(err.Error())
	}
}

func (ui *UI) MarkImage(c http.Context) {
	imageID := c.Request().PathParam("image")
	err := ui.app.MarkImageAs(imageID, c.Request().QueryParam("type"))
	switch err {
	case nil:
		c.String(http.StatusOK, "Success")
	case domain.ErrEmptyImageType:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrMarkImageFailed, domain.ErrNoSuchImage:
		c.String(http.StatusNotFound, err.Error())
	default:
		panic(err.Error())
	}
}