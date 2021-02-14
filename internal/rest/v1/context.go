package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/project"
)

type Context struct {
	echo.Context
	ProjectService project.Service
	Log            *logrus.Entry
}
