package ui

import (
	"github.com/pkg/errors"

	"lmm/api/http"
	"lmm/api/service/article/application"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/finder"
	"lmm/api/service/article/domain/model"
	"lmm/api/service/article/domain/repository"
	"lmm/api/service/article/domain/service"
	userModel "lmm/api/service/auth/domain/model"
)

var (
	errTitleRequired = errors.New("title required")
	errBodyRequired  = errors.New("body requried")
	errTagsRequired  = errors.New("tags requried")
)

// UI is the user interface to contact with network
type UI struct {
	appService *application.Service
}

// NewUI returns a new ui
func NewUI(
	articleFinder finder.ArticleFinder,
	articleRepository repository.ArticleRepository,
	authorService service.AuthorService,
) *UI {
	appService := application.NewService(
		application.NewArticleCommandService(articleRepository, authorService),
		application.NewArticleQueryService(articleFinder),
	)
	return &UI{appService: appService}
}

// PostNewArticle handles POST /1/articles
func (ui *UI) PostNewArticle(c http.Context) {
	user, ok := c.Value(http.StrCtxKey("user")).(*userModel.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	article := postArticleAdapter{}
	if err := c.Request().Bind(&article); err != nil {
		http.BadRequest(c)
		return
	}

	if err := ui.validatePostArticleAdaptor(&article); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(c,
		user.Name(),
		*article.Title,
		*article.Body,
		article.Tags,
	)
	switch errors.Cause(err) {
	case nil:
		c.Header("Location", "/v1/articles/"+articleID.String())
		c.String(http.StatusCreated, "Success")
	case domain.ErrArticleTitleTooLong, domain.ErrEmptyArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrInvalidArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrNoSuchUser:
		http.Unauthorized(c)
	default:
		http.Error(c, err.Error())
		http.ServiceUnavailable(c)
	}
}

// EditArticle handles PUT /1/article/:articleID
func (ui *UI) EditArticle(c http.Context) {
	user, ok := c.Value(http.StrCtxKey("user")).(*userModel.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	article := postArticleAdapter{}
	if err := c.Request().Bind(&article); err != nil {
		http.BadRequest(c)
		return
	}

	if err := ui.validatePostArticleAdaptor(&article); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err := ui.appService.ArticleCommandService().EditArticle(c,
		user.Name(),
		c.Request().PathParam("articleID"),
		*article.Title,
		*article.Body,
		article.Tags,
	)
	switch errors.Cause(err) {
	case nil:
		http.NoContent(c)
	case domain.ErrArticleTitleTooLong, domain.ErrEmptyArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrInvalidArticleID:
		c.String(http.StatusNotFound, domain.ErrNoSuchArticle.Error())
	case domain.ErrInvalidArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrNoSuchArticle:
		c.String(http.StatusNotFound, err.Error())
	case domain.ErrNoSuchUser:
		http.Unauthorized(c)
	case domain.ErrNotArticleAuthor:
		c.String(http.StatusForbidden, err.Error())
	default:
		http.Error(c, err.Error())
		http.ServiceUnavailable(c)
	}
}

func (ui *UI) validatePostArticleAdaptor(adaptor *postArticleAdapter) error {
	if adaptor.Title == nil {
		return errTitleRequired
	}
	if adaptor.Body == nil {
		return errBodyRequired
	}
	if adaptor.Tags == nil {
		return errTagsRequired
	}
	return nil
}

// ListArticles handles GET /v1/articles
func (ui *UI) ListArticles(c http.Context) {
	view, err := ui.appService.ArticleQueryService().ListArticlesByPage(c,
		c.Request().QueryParam("count"),
		c.Request().QueryParam("page"),
	)
	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, ui.articleListViewToJSON(view))
	case application.ErrInvalidCount, application.ErrInvalidPage:
		c.JSON(http.StatusBadRequest, err.Error())
	default:
		http.Error(c, err.Error())
		http.ServiceUnavailable(c)
	}
}

func (ui *UI) articleListViewToJSON(view *model.ArticleListView) *articleListAdapter {
	items := make([]articleListItem, len(view.Items()), len(view.Items()))
	for i, item := range view.Items() {
		items[i].ID = item.ID().String()
		items[i].Title = item.Title()
		items[i].PostAt = item.PostAt().UTC().String()
	}
	return &articleListAdapter{
		Articles:    items,
		HasNextPage: view.HasNextPage(),
	}
}

// GetArticle handles GET /v1/articles/:articleID
func (ui *UI) GetArticle(c http.Context) {
	view, err := ui.appService.ArticleQueryService().ArticleByID(c,
		c.Request().PathParam("articleID"),
	)
	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, ui.articleViewToJSON(view))
	case domain.ErrInvalidArticleID, domain.ErrNoSuchArticle:
		c.String(http.StatusNotFound, domain.ErrNoSuchArticle.Error())
	default:
		http.Error(c, err.Error())
		http.ServiceUnavailable(c)
	}
}

func (ui *UI) articleViewToJSON(view *model.ArticleView) *articleViewResponse {
	tags := make([]articleViewTag, len(view.Content().Tags()), len(view.Content().Tags()))
	for i, tag := range view.Content().Tags() {
		tags[i].Name = tag.Name()
	}
	return &articleViewResponse{
		ID:           view.ID().String(),
		Title:        view.Content().Text().Title(),
		Body:         view.Content().Text().Body(),
		PostAt:       view.PostAt().UTC().String(),
		LastEditedAt: view.LastEditedAt().UTC().String(),
		Tags:         tags,
	}
}

// GetAllArticleTags handles GET /v1/articleTags
func (ui *UI) GetAllArticleTags(c http.Context) {
	view, err := ui.appService.ArticleQueryService().AllArticleTags(c)

	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, ui.tagListViewToJSON(view))
	default:
		http.Error(c, err.Error())
		http.ServiceUnavailable(c)
	}
}

func (ui *UI) tagListViewToJSON(view model.TagListView) articleTagListView {
	tags := make([]articleTagListItemView, len(view), len(view))
	for i, tag := range view {
		tags[i].Name = tag.Name()
	}
	return tags
}