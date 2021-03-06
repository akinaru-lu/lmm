package presentation

import (
	"net/http"

	httpUtil "lmm/api/pkg/http"
	"lmm/api/service/asset/usecase"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	ErrUnsupportType = errors.New("unsupport type")
)

type GinRouterProvider struct {
	usecase *usecase.Usecase
}

func NewGinRouterProvider(app *usecase.Usecase) *GinRouterProvider {
	return &GinRouterProvider{usecase: app}
}

func (p *GinRouterProvider) Provide(router *gin.Engine) {
	router.POST("/v1/photos", p.PostV1Photos)
	router.PUT("/v1/photos/:photo/tags", p.PutV1PhotoTags)
	router.GET("/v1/photos", p.GetV1Photos)
	router.GET("/v1/photos/:photo", p.GetV1Photo)
}

// PostV1Photos handles POST /v1/photos
func (p *GinRouterProvider) PostV1Photos(c *gin.Context) {
	user, ok := httpUtil.AuthFromGinContext(c)
	if !ok {
		httpUtil.Unauthorized(c)
		return
	}

	f, fh, err := c.Request.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			c.String(http.StatusBadRequest, "photo required")
			return
		}
		httpUtil.LogWarn(c, "failed to get file from request", err)
		httpUtil.BadRequest(c)
		return
	}

	contentType := fh.Header.Get("Content-Type")
	switch contentType {
	case "image/gif", "image/jpeg", "image/jpg", "image/png":
	default:
		c.String(http.StatusBadRequest, ErrUnsupportType.Error())
	}

	url, err := p.usecase.UploadPhoto(c, &usecase.AssetToUpload{
		ContentType: contentType,
		DataSource:  f,
		Filename:    fh.Filename,
		UserID:      user.ID,
	})

	switch err {
	case nil:
		c.Header("Location", url)
		httpUtil.Response(c, http.StatusCreated, "Success")
	default:
		httpUtil.LogPanic(c, "unexpected error", err)
	}
}

type tagList struct {
	Tags []string `json:"tags"`
}

type photo struct {
	ID string `uri:"photo" binding:"required"`
}

// GetV1Photo handles GET /v1/photos/:photo
func (p *GinRouterProvider) GetV1Photo(c *gin.Context) {
	var photo photo
	if err := c.ShouldBindUri(&photo); err != nil {
		httpUtil.LogWarn(c, "bind uri error", err)
		httpUtil.BadRequest(c)
		return
	}

	json, err := p.usecase.GetPhotoInfo(c, photo.ID)
	if err != nil {
		httpUtil.LogWarn(c, "error on getting photo info", err)
		httpUtil.NotFound(c)
		return
	}

	c.JSON(http.StatusOK, json)
}

// PutV1PhotoTags handles PUT /v1/photos/:photo/tags
func (p *GinRouterProvider) PutV1PhotoTags(c *gin.Context) {
	user, ok := httpUtil.AuthFromGinContext(c)
	if !ok {
		httpUtil.Unauthorized(c)
		return
	}

	var tags tagList
	if err := c.ShouldBindJSON(&tags); err != nil {
		httpUtil.LogWarn(c, "bind json error", err)
		httpUtil.BadRequest(c)
		return
	}

	var photo photo
	if err := c.ShouldBindUri(&photo); err != nil {
		httpUtil.LogWarn(c, "bind uri error", err)
		httpUtil.BadRequest(c)
		return
	}

	err := p.usecase.SetPhotoTags(c, user.ID, photo.ID, tags.Tags)
	if err != nil {
		httpUtil.LogWarn(c, "error on setting photo tags", err)
	}

	original := errors.Cause(err)

	switch original {
	case nil:
		httpUtil.Response(c, http.StatusOK, "Success")
	case usecase.ErrNoSuchPhoto:
		httpUtil.NotFound(c)
	case usecase.ErrForbidden:
		httpUtil.Forbidden(c)
	default:
		httpUtil.LogPanic(c, "unexpected error", err)
	}
}

type photoList struct {
	Items      []*usecase.Photo `json:"items"`
	NextCursor string           `json:"next_cursor"`
}

// GetV1Photos handles GET /v1/photos
func (p *GinRouterProvider) GetV1Photos(c *gin.Context) {
	photos, cursor, err := p.usecase.ListPhotos(c,
		c.DefaultQuery("count", "10"),
		c.DefaultQuery("cursor", ""),
	)

	if err != nil {
		httpUtil.LogWarn(c, "error on getting photo list", err)
		httpUtil.NotFound(c)
		return
	}

	c.JSON(http.StatusOK, &photoList{
		Items:      photos,
		NextCursor: cursor,
	})
}

// PatchV1Photos handles PATCH /v1/photos/:photo
func (p *GinRouterProvider) PatchV1Photos(c *gin.Context) {
}

// PostV1Assets handles POST /v1/assets
// This endpoint is to upload common assets
func (p *GinRouterProvider) PostV1Assets(c *gin.Context) {
}
