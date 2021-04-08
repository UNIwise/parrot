package rest

import (
	"context"
	"fmt"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	eprom "github.com/paulfarver/echo-pack/middleware"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/project"
	v1 "github.com/uniwise/parrot/internal/rest/v1"
)

const (
	gzipCompressionLevel = 5
	requestIDLength      = 10
)

type Server struct {
	Echo *echo.Echo
}

func NewServer(projectService project.Service, entry *logrus.Entry, enablePrometheus bool) (*Server, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = NewValidator()

	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipCompressionLevel,
	}))
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return random.String(requestIDLength, random.Hex)
		},
	}))

	v1Group := e.Group("/v1", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &v1.Context{
				Context:        c,
				ProjectService: projectService,
				Log:            entry.WithField("requestID", c.Response().Header().Get(echo.HeaderXRequestID)),
			}

			return next(cc)
		}
	})
	if enablePrometheus {
		v1Group.Use(eprom.Prometheus())
	}
	v1.Register(v1Group)

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
