package v1

import (
	"github.com/labstack/echo/v4"
	eprom "github.com/paulfarver/echo-pack/middleware"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/project"
)

type Handlers struct {
	ProjectService project.Service
}

type HandlerFunction func(ctx echo.Context, l *logrus.Entry) error

func Register(e *echo.Echo, l *logrus.Entry, projectService project.Service, enablePrometheus bool) {
	h := &Handlers{
		ProjectService: projectService,
	}

	g := e.Group("/v1")

	if enablePrometheus {
		g.Use(eprom.Prometheus())
	}

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
