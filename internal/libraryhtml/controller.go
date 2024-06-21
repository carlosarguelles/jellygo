package libraryhtml

import (
	"fmt"
	"strconv"

	"github.com/carlosarguelles/jellygo/internal/layouts"
	"github.com/carlosarguelles/jellygo/internal/library/application"
	"github.com/carlosarguelles/jellygo/internal/library/domain"
	"github.com/carlosarguelles/jellygo/internal/library/infrastructure"
	"github.com/labstack/echo/v4"
)

type HtmlLibraryController struct {
	libraryRepository application.LibraryRepository
	libraryQueue      *infrastructure.LibraryQueue
}

func (lc *HtmlLibraryController) Index(c echo.Context) error {
	libs, _ := lc.libraryRepository.GetAll(c.Request().Context())
	return layouts.Root(Index(libs)).Render(c.Request().Context(), c.Response().Writer)
}

func (lc *HtmlLibraryController) Show(c echo.Context) error {
	libraryIDStr := c.Param("libraryID")
	libraryID, _ := strconv.Atoi(libraryIDStr)
	lib, _ := lc.libraryRepository.GetByID(c.Request().Context(), libraryID)
	return layouts.Root(Show(*lib)).Render(c.Request().Context(), c.Response().Writer)
}

func (lc *HtmlLibraryController) Create(c echo.Context) error {
	return layouts.Root(Create()).Render(c.Request().Context(), c.Response().Writer)
}

func (lc *HtmlLibraryController) Save(c echo.Context) error {
	path := c.FormValue("path")
	name := c.FormValue("name")
	typeStr := c.FormValue("type")
	lib := domain.NewLibrary(path, typeStr, name)
	lc.libraryRepository.Create(c.Request().Context(), lib)
	lc.libraryQueue.Add(*lib)
	c.Response().Header().Add("HX-Location", fmt.Sprintf("/libraries/%d", lib.ID))
	return nil
}

func (lc *HtmlLibraryController) Refresh(c echo.Context) error {
	libraryIDStr := c.Param("libraryID")
	libraryID, _ := strconv.Atoi(libraryIDStr)
	lib, _ := lc.libraryRepository.GetByID(c.Request().Context(), libraryID)
	lc.libraryQueue.Add(*lib)
	return Show(*lib).Render(c.Request().Context(), c.Response().Writer)
}

func NewHtmlLibraryController(repository application.LibraryRepository, libraryQueue *infrastructure.LibraryQueue) *HtmlLibraryController {
	return &HtmlLibraryController{repository, libraryQueue}
}
