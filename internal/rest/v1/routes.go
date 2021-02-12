package v1

import (
	"github.com/labstack/echo/v4"
)

func Register(g *echo.Group) {
	g.GET("/project/:project/language/:language", getProjectLanguage)
}
