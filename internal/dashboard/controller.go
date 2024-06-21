package dashboard

import (
	"github.com/carlosarguelles/jellygo/internal/layouts"
	"github.com/labstack/echo/v4"
)

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

func (dc *DashboardController) Index(c echo.Context) error {
	return layouts.Root(Dashboard()).Render(c.Request().Context(), c.Response().Writer)
}
