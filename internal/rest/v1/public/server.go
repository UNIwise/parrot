package public

import (
	"context"
	"fmt"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/project"
	controllers "github.com/uniwise/parrot/internal/rest/v1/public/controllers"
	"github.com/uniwise/parrot/internal/rest/v1/helpers"
)

const (
	gzipCompressionLevel = 5
)

type Server struct {
	Echo *echo.Echo
}

func NewServer(l *logrus.Entry, projectService project.Service) (*Server, error) {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.Validator = helpers.NewValidator()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipCompressionLevel,
	}))

	controllers.Register(e, l, projectService)

	h := gosundheit.New()

	if err := projectService.RegisterChecks(h); err != nil {
		return nil, errors.Wrap(err, "Failed to register healthchecks")
	}

	e.GET("/health", echo.WrapHandler(healthhttp.HandleHealthJSON(h)))

	return &Server{
		Echo: e,
	}, nil
}

func (s *Server) Start(port int) error {
	return s.Echo.Start(fmt.Sprintf(":%d", port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Echo.Shutdown(ctx)
}
