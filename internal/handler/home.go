package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/tneuqole/habitmap/internal/template"
	"github.com/tneuqole/habitmap/internal/util"
)

type HomeHandler struct{}

func (h HomeHandler) GetHome(c echo.Context) error {
	return util.Render(c, template.Home())
}
