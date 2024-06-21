package main

import (
	"context"
	"database/sql"
	"os"
	"path"

	"github.com/carlosarguelles/jellygo/internal/dashboard"
	libinfra "github.com/carlosarguelles/jellygo/internal/library/infrastructure"
	"github.com/carlosarguelles/jellygo/internal/libraryhtml"
	"github.com/carlosarguelles/jellygo/internal/librarymoviehtml"
	"github.com/carlosarguelles/jellygo/internal/media"
	movieinfra "github.com/carlosarguelles/jellygo/internal/movie/infrastructure"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryanbradynd05/go-tmdb"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()

	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync()

	db, err := sql.Open("sqlite3", "file:./db/dev.db")
	if err != nil {
		panic("error connecting to db")
	}
	defer db.Close()

	tmdbAPI := tmdb.Init(tmdb.Config{
		APIKey: os.Getenv("TMDB_API_KEY"),
	})

	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Info("req", zap.String("uri", v.URI), zap.Int("status", v.Status), zap.String("method", v.Method))
			} else {
				logger.Error("req error", zap.String("uri", v.URI), zap.Int("status", v.Status), zap.Error(v.Error))
			}
			return nil
		},
	}))

	e.Static("/static", "public")

	libraryRepository := libinfra.NewSqliteLibraryRepository(db)

	movieAPI := movieinfra.NewTmdbMovieAPI(tmdbAPI)
	movieRepository := movieinfra.NewSqliteMovieRepository(db)

	imageDownloader := media.NewImageManagerService(path.Join("public", "images"), logger)

	libraryQueue := libinfra.NewLibraryQueue(movieAPI, movieRepository, imageDownloader, logger)
	go libraryQueue.Run(ctx)
	defer libraryQueue.Close()

	libraryController := libinfra.NewLibraryController(libraryRepository, libraryQueue)

	e.GET("/api/libraries", libraryController.Index)
	e.POST("/api/libraries", libraryController.Create)

	htmlLibraryController := libraryhtml.NewHtmlLibraryController(libraryRepository, libraryQueue)

	e.GET("/libraries", htmlLibraryController.Index)
	e.GET("/libraries/:libraryID", htmlLibraryController.Show)
	e.GET("/libraries/new", htmlLibraryController.Create)
	e.POST("/libraries", htmlLibraryController.Save)
	e.PATCH("/libraries/:libraryID", htmlLibraryController.Refresh)

	htmlLibraryMovieController := librarymoviehtml.NewHtmlLibraryMovieController(libraryRepository, movieRepository)

	e.GET("/libraries/:libraryID/movies", htmlLibraryMovieController.Index)
	e.GET("/libraries/:libraryID/movies/:movieID", htmlLibraryMovieController.Show)

	dashboardController := dashboard.NewDashboardController()

	e.GET("/", dashboardController.Index)

	e.Logger.Fatal(e.Start(":6969"))
}
