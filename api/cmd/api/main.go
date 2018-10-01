package main

import (
	"fmt"
	"lmm/api/storage/uploader"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"lmm/api/http"
	"lmm/api/log"
	"lmm/api/middleware"
	"lmm/api/storage/db"

	// user
	userApp "lmm/api/service/user/application"
	userStorage "lmm/api/service/user/infra/persistence"
	userUI "lmm/api/service/user/ui"

	// auth
	authApp "lmm/api/service/auth/application"
	authStorage "lmm/api/service/auth/infra/persistence"
	authUI "lmm/api/service/auth/ui"

	// article
	articleFetcher "lmm/api/service/article/infra/fetcher"
	articleStorage "lmm/api/service/article/infra/persistence"
	authorService "lmm/api/service/article/infra/service"
	articleUI "lmm/api/service/article/ui"

	// image
	imageStorage "lmm/api/service/image/infra/persistence"
	imageService "lmm/api/service/image/infra/service"
	image "lmm/api/service/image/ui"
)

var (
	dbSrcName  = "root:@tcp(lmm-mysql:3306)/"
	dbName     = os.Getenv("DATABASE_NAME")
	connParams = "parseTime=true"
)

func main() {
	logger := globalRecorder()
	defer logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	mysql := db.NewMySQL(fmt.Sprintf("%s%s?%s", dbSrcName, dbName, connParams))
	defer mysql.Close()

	router := http.NewRouter()

	// middleware
	accessRecorder := accessLogRecorder()
	defer accessRecorder.Sync()
	router.Use(middleware.NewAccessLog(accessRecorder))

	recvRecoder := recoveryRecorder()
	defer recvRecoder.Sync()
	router.Use(middleware.NewRecovery(recvRecoder))

	// user
	userRepo := userStorage.NewUserStorage(mysql)
	userAppService := userApp.NewService(userRepo)
	userUI := userUI.NewUI(userAppService)
	router.POST("/v1/users", userUI.SignUp)

	// auth
	authRepo := authStorage.NewUserStorage(mysql)
	authAppService := authApp.NewService(authRepo)
	authUI := authUI.NewUI(authAppService)
	router.POST("/v1/auth/login", authUI.Login)

	// article
	authorAdapter := authorService.NewAuthorAdapter(mysql)
	articleRepo := articleStorage.NewArticleStorage(mysql, authorAdapter)
	articleFinder := articleFetcher.NewArticleFetcher(mysql)
	articleUI := articleUI.NewUI(articleFinder, articleRepo, authorAdapter)
	router.POST("/v1/articles", articleUI.PostNewArticle)
	router.PUT("/v1/articles", articleUI.EditArticle)
	router.GET("/v1/articles", articleUI.ListArticles)
	router.GET("/v1/articles/:articleID", articleUI.GetArticle)
	router.GET("/v1/articleTags", articleUI.GetAllArticleTags)

	// image
	imageRepo := imageStorage.NewImageStorage(mysql, uploader.NewLocalImageUploader())
	image := image.New(imageRepo, imageService.NewUserAdapter(mysql))
	router.POST("/v1/images", authUI.BearerAuth(image.Upload))

	server := http.NewServer(":8002", router)
	server.Run()
}

func globalRecorder() *zap.Logger {
	cfg := log.DefaultZapConfig()

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Named("global")
}

func recoveryRecorder() *zap.Logger {
	cfg := log.DefaultZapConfig()
	cfg.DisableCaller = true
	cfg.DisableStacktrace = false
	cfg.EncoderConfig.StacktraceKey = "stacktrace"

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Named("recovery")
}

func accessLogRecorder() *zap.Logger {
	cfg := log.DefaultZapConfig()
	cfg.DisableCaller = true

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Named("access_log")
}
