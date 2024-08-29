package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/tneuqole/heatmap/internal/util"
	"github.com/tneuqole/heatmap/template"
)

type HomeHandler struct{}

func (h HomeHandler) GetHome(c echo.Context) error {
	return util.Render(c, template.Home())
}
