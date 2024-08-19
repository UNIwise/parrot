package database

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DSN   string `mapstructure:"dsn" validate:"required"`
	Debug bool   `mapstructure:"debug" default:"false"`
}

func NewClient(ctx context.Context, conf Config) (*gorm.DB, error) {
	dialect := mysql.Open(conf.DSN)
	
	gormConfig := &gorm.Config{}
	if(conf.Debug) {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(dialect, gormConfig)
	if(err != nil) {
		return nil, errors.Wrap(err, "failed to open database connection")
	}

	sqlDB, err := db.DB()
	if(err != nil) {
		return nil, errors.Wrap(err, "failed to get database connection for pinging")
	}

	if err = sqlDB.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}

	return db, nil
}
