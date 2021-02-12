package rest

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v1 "github.com/uniwise/parrot/internal/rest/v1"
)

const (
	gzipCompressionLevel = 5
)

type Server struct {
	Config *Config
	Echo   *echo.Echo
}

func NewServer(conf *Config) (*Server, error) {
	if conf == nil {
		return nil, errors.New("Config can't be nil")
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = NewValidator()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipCompressionLevel,
	}))

	// Routes
	v1.Register(e.Group("/v1"))

	// Health endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"health": "ok",
		})
	})

	return &Server{
		Config: conf,
		Echo:   e,
	}, nil
}

func (s *Server) Start() error {
	return s.Echo.Start(fmt.Sprintf(":%d", s.Config.Port))
}
