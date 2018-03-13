package main

import (
	"log"
	"net/http"

	"github.com/akinaru-lu/elesion"

	"lmm/api/controller/blog"
	"lmm/api/controller/category"
	"lmm/api/controller/image"
	"lmm/api/controller/tag"
	"lmm/api/controller/user"
)

func main() {
	router := elesion.Default("[api]")

	// user
	router.POST("/v1/signup", user.SignUp)
	router.POST("/v1/signin", user.SignIn)
	router.GET("/v1/verify", user.Verify)

	// blog
	router.GET("/v1/blog/:blog", blog.Get)
	router.GET("/v1/users/:user/blog", blog.GetList)
	router.POST("/v1/blog", blog.Post)
	router.PUT("/v1/blog/:blog", blog.Update)
	router.DELETE("/v1/blog/:blog", blog.Delete)
	// blog category
	router.GET("/v1/blog/:blog/category", blog.GetCategory)
	router.PUT("/v1/blog/:blog/category", blog.SetCategory)
	router.DELETE("/v1/blog/:blog/category", blog.DeleteCategory)

	// category
	router.GET("/v1/users/:user/categories", category.GetByUser)
	router.POST("/v1/categories", category.Register)
	router.PUT("/v1/categories/:category", category.Update)
	router.DELETE("/v1/categories/:category", category.Delete)

	// tag
	router.GET("/v1/users/:user/tags", tag.GetByUser)
	router.GET("/v1/blog/:blog/tags", tag.GetByBlog)
	router.POST("/v1/blog/:blog/tags", tag.Register)
	router.PUT("/v1/blog/:blog/tags/:tag", tag.Update)
	router.DELETE("/v1/blog/:blog/tags/:tag", tag.Delete)

	// image
	router.GET("/v1/users/:user/images", image.GetAllImages)
	router.GET("/v1/users/:user/images/photos", image.GetAllPhotos)
	router.POST("/v1/images", image.Upload)

	log.Fatal(http.ListenAndServe(":8081", router))
}
