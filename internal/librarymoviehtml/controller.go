package librarymoviehtml

import (
	"strconv"

	"github.com/carlosarguelles/jellygo/internal/layouts"
	libapp "github.com/carlosarguelles/jellygo/internal/library/application"
	movapp "github.com/carlosarguelles/jellygo/internal/movie/application"
	"github.com/labstack/echo/v4"
)

type HtmlLibraryMovieController struct {
	libraryRepository libapp.LibraryRepository
	movieRepository   movapp.MovieRepository
}

func NewHtmlLibraryMovieController(libraryRepository libapp.LibraryRepository, movieRepository movapp.MovieRepository) *HtmlLibraryMovieController {
	return &HtmlLibraryMovieController{libraryRepository, movieRepository}
}

func (lmc *HtmlLibraryMovieController) Index(c echo.Context) error {
	libraryIDStr := c.Param("libraryID")
	libraryID, _ := strconv.Atoi(libraryIDStr)
	library, err := lmc.libraryRepository.GetByIDWithMovies(c.Request().Context(), libraryID)
	if err != nil {
		return err
	}
	return Index(library.Movies).Render(c.Request().Context(), c.Response().Writer)
}

func (lmc *HtmlLibraryMovieController) Show(c echo.Context) error {
	movieIDStr := c.Param("movieID")
	movieID, _ := strconv.Atoi(movieIDStr)
	movie, err := lmc.movieRepository.GetByID(c.Request().Context(), movieID)
	if err != nil {
		return err
	}
	return layouts.Root(Show(*movie)).Render(c.Request().Context(), c.Response().Writer)
}
