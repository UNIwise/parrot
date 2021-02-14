package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/cache"
	"github.com/uniwise/parrot/internal/poedit"
)

type Context struct {
	echo.Context
	Client poedit.Client
	Cacher cache.Cache
	Log    *logrus.Entry
}
