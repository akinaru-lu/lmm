package main

import (
	"lmm/api/articles"
	"lmm/api/db"
	"lmm/api/image"
	"lmm/api/profile"
	"log"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

func init() {
	db.Init()
}

func main() {
	router := elesion.Default("[api]")

	// /articles
	router.GET("/articles", articles.GetArticles)

	// /article
	router.GET("/article", articles.GetArticle)
	router.POST("/article", articles.PostArticle)
	router.PUT("/article", articles.UpdateArticle)

	// /articles/categories
	router.GET("/articles/categories", articles.GetCategories)

	// /articles/category
	router.POST("/articles/category", articles.NewCategory)

	// /articles/tags
	router.GET("/articles/tags", articles.GetTags)

	router.GET("/photos", image.Handler)

	router.GET("/profile", profile.Handler)

	log.Fatal(http.ListenAndServe(":8081", router))
}
