package librarymoviehtml

import (
	"strconv"

	libraryapp "github.com/carlosarguelles/jellygo/internal/library/application"
	"github.com/labstack/echo/v4"
)

type HtmlLibraryMovieController struct {
	libraryRepository libraryapp.LibraryRepository
}

func NewHtmlLibraryMovieController(libraryRepository libraryapp.LibraryRepository) *HtmlLibraryMovieController {
	return &HtmlLibraryMovieController{libraryRepository}
}

func (lmc *HtmlLibraryMovieController) Show(c echo.Context) error {
	libraryIDStr := c.Param("libraryID")
	libraryID, _ := strconv.Atoi(libraryIDStr)
	library, err := lmc.libraryRepository.GetByIDWithMovies(c.Request().Context(), libraryID)
	if err != nil {
		return err
	}
	return Index(library.Movies).Render(c.Request().Context(), c.Response().Writer)
}
