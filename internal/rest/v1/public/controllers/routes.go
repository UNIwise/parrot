package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/project"
)

type Handlers struct {
	ProjectService project.Service
}

type HandlerFunction func(ctx echo.Context, l *logrus.Entry) error

func Register(e *echo.Echo, l *logrus.Entry, projectService project.Service) {
	h := &Handlers{
		ProjectService: projectService,
	}

	g := e.Group("/v1")

	g.GET("/project/:project/language/:language", wrap(h.getProjectLanguage, l))
}

func wrap(fn HandlerFunction, logger *logrus.Entry) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		l := logger.WithFields(logrus.Fields{
			"method":    ctx.Request().Method,
			"path":      ctx.Request().URL.Path,
			"requestId": ctx.Response().Header().Get(echo.HeaderXRequestID),
		})

		return fn(ctx, l)
	}
}
