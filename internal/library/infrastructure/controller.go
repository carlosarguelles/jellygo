package infrastructure

import (
	"net/http"

	"github.com/carlosarguelles/jellygo/internal/library/application"
	"github.com/carlosarguelles/jellygo/internal/library/domain"
	"github.com/labstack/echo/v4"
)

type LibraryController struct {
	repository   application.LibraryRepository
	libraryQueue *LibraryQueue
}

func (lc *LibraryController) Index(c echo.Context) error {
	libs, err := lc.repository.GetAll(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "err")
	}
	return c.JSON(http.StatusOK, libs)
}

func (lc *LibraryController) Create(c echo.Context) error {
	path := c.FormValue("path")
	if path == "" {
		return c.String(http.StatusBadRequest, "empty path")
	}
	lib := domain.Library{Path: path}
	err := lc.repository.Create(c.Request().Context(), &lib)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error creating library")
	}
	lc.libraryQueue.Add(lib)
	return c.String(http.StatusOK, "ok")
}

func NewLibraryController(repository application.LibraryRepository, libraryQueue *LibraryQueue) *LibraryController {
	return &LibraryController{repository, libraryQueue}
}
