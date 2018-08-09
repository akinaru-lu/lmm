package main

import (
	accountInfra "lmm/api/context/account/infra"
	account "lmm/api/context/account/ui"
	blogInfra "lmm/api/context/blog/infra"
	blog "lmm/api/context/blog/ui"
	imageInfra "lmm/api/context/image/infra"
	img "lmm/api/context/image/ui"

	"lmm/api/http"
	"lmm/api/storage"
)

var (
	accountUI *account.UI
	blogUI    *blog.UI
	imageUI   *img.UI
)

func initUIs(db *storage.DB) {
	userRepo := accountInfra.NewUserStorage(db)
	accountUI = account.New(userRepo)

	blogRepo := blogInfra.NewBlogStorage(db)
	categoryRepo := blogInfra.NewCategoryStorage(db)
	tagRepo := blogInfra.NewTagStorage(db)
	blogUI = blog.New(blogRepo, categoryRepo, tagRepo)

	imgRepo := imageInfra.NewImageStorage(db)
	imageUI = img.New(imgRepo)
}

func main() {
	db := storage.NewDB()
	initUIs(db)

	router := http.NewRouter()

	// account
	router.POST("/v1/signup", accountUI.SignUp)
	router.POST("/v1/signin", accountUI.SignIn)
	router.GET("/v1/verify", accountUI.BearerAuth(accountUI.Verify))

	// blog
	router.GET("/v1/blog", blogUI.GetAllBlog)
	router.GET("/v1/blog/:blog", blogUI.GetBlog)
	router.POST("/v1/blog", accountUI.BearerAuth(blogUI.PostBlog))
	router.PUT("/v1/blog/:blog", accountUI.BearerAuth(blogUI.UpdateBlog))
	// blog category
	router.GET("/v1/blog/:blog/category", blogUI.GetBlogCagetory)
	router.PUT("/v1/blog/:blog/category", accountUI.BearerAuth(blogUI.SetBlogCategory))
	// blog tag
	router.GET("/v1/blog/:blog/tags", blogUI.GetAllTagsOfBlog)
	router.POST("/v1/blog/:blog/tags", accountUI.BearerAuth(blogUI.NewBlogTag))

	// category
	router.GET("/v1/categories", blogUI.GetAllCategoris)
	router.POST("/v1/categories", accountUI.BearerAuth(blogUI.PostCategory))
	router.PUT("/v1/categories/:category", accountUI.BearerAuth(blogUI.UpdateCategory))
	router.DELETE("/v1/categories/:category", accountUI.BearerAuth(blogUI.DeleteCategory))

	// tag
	router.GET("/v1/tags", blogUI.GetAllTags)
	router.PUT("/v1/tags/:tag", accountUI.BearerAuth(blogUI.UpdateTag))
	router.DELETE("/v1/tags/:tag", accountUI.BearerAuth(blogUI.DeleteTag))

	// image
	router.POST("/v1/images", accountUI.BearerAuth(imageUI.Upload))
	router.GET("/v1/images", imageUI.LoadImagesByPage)
	router.PUT("/v1/images/:image/photo", accountUI.BearerAuth(imageUI.SetAsPhoto))
	router.DELETE("/v1/images/:image/photo", accountUI.BearerAuth(imageUI.SetAsPhoto))

	http.Serve(":8002", router)
}
