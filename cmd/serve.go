/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uniwise/parrot/internal/cache"
	"github.com/uniwise/parrot/internal/metrics"
	"github.com/uniwise/parrot/internal/project"
	privateRest "github.com/uniwise/parrot/internal/rest/v1/private"
	publicRest "github.com/uniwise/parrot/internal/rest/v1/public"
	"github.com/uniwise/parrot/internal/storage"
	"github.com/uniwise/parrot/pkg/poedit"
)

const (
	confServerPortPublic  = "server.port.public"
	confServerPortPrivate = "server.port.private"
	confServerGrace       = "server.gracePeriod"

	confLogLevel  = "log.level"
	confLogFormat = "log.format"

	confCacheType                  = "cache.type"
	confCacheTTL                   = "cache.ttl"
	confCacheRenewalThreshold      = "cache.renewalThreshold"
	confCacheFSDir                 = "cache.filesystem.dir"
	confCacheRedisMode             = "cache.redis.mode"
	confCacheRedisAddress          = "cache.redis.address"
	confCacheRedisUser             = "cache.redis.username"
	confCacheRedisPassword         = "cache.redis.password"
	confCacheRedisMaxRetries       = "cache.redis.maxRetries"
	confCacheRedisDB               = "cache.redis.db"
	confCacheRedisSentinelMaster   = "cache.redis.sentinel.master"
	confCacheRedisSentinelAddress  = "cache.redis.sentinel.addresses"
	confCacheRedisSentinelPassword = "cache.redis.sentinel.password" //nolint:gosec

	confPrometheusEnabled = "prometheus.enabled"
	confPrometheusPath    = "prometheus.path"
	confPrometheusPort    = "prometheus.port"

	confAPIToken = "api.token"

	confAWSRegion = "aws.region"
	confAWSBucket = "aws.bucket"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the parrot caching server",
	Long: `Start the parrot caching server.
This will take care of serving your translations,
by caching exports from poeditor`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := instantiateLogger()

		cacheInstance, err := instantiateCache(logger.WithField("subsystem", "cache"))
		if err != nil {
			logger.WithError(err).Fatal("Could not instantiate cache")
		}

		cli := poedit.NewClient(viper.GetString(confAPIToken), http.DefaultClient)

		s3Config := &storage.S3StorageConfig{
			Region: viper.GetString(confAWSRegion),
			Bucket: viper.GetString(confAWSBucket),
		}

		s3Client, err := storage.NewS3Client(context.Background(), *s3Config)
		if err != nil {
			logger.WithError(err).Fatal("Could not instantiate S3 client")
		}

		storageService := storage.NewService(context.Background(), s3Client)

		svc := project.NewService(cli, storageService, cacheInstance, viper.GetDuration(confCacheRenewalThreshold), logrus.NewEntry(logger))

		publicServer, err := publicRest.NewServer(logrus.NewEntry(logger), svc)
		if err != nil {
			logger.WithError(err).Fatal("Could not instantiate public server")
		}

		publicPort := viper.GetInt(confServerPortPublic)

		logger.Infof("Public server listening at :%d", publicPort)
		go func() {
			if err := publicServer.Start(publicPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal("shutting down public server")
			}
		}()

		privateServer, err := privateRest.NewServer(logrus.NewEntry(logger), svc, viper.GetBool(confPrometheusEnabled))
		if err != nil {
			logger.WithError(err).Fatal("Could not instantiate private server")
		}

		privatePort := viper.GetInt(confServerPortPrivate)

		logger.Infof("Private server listening at :%d", privatePort)
		go func() {
			if err := privateServer.Start(privatePort); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal("shutting down private server")
			}
		}()

		if viper.GetBool(confPrometheusEnabled) {
			logger.Infof("Prometheus metrics exposed at :%d%s", viper.GetInt(confPrometheusPort), viper.GetString(confPrometheusPath))
			go func() {
				logger.Fatal(metrics.Start(viper.GetString(confPrometheusPath), viper.GetInt(confPrometheusPort)))
			}()
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(confServerGrace))
		defer cancel()
		if err := privateServer.Shutdown(ctx); err != nil {
			logger.Fatal(err)
		}
		if err := publicServer.Shutdown(ctx); err != nil {
			logger.Fatal(err)
		}
	},
}

// nolint:gochecknoinits
func init() {
	cDir, err := os.UserCacheDir()
	if err != nil {
		cDir = "/tmp"
	}

	viper.SetDefault(confServerPortPublic, 8000)
	viper.SetDefault(confServerPortPrivate, 8001)
	viper.SetDefault(confServerGrace, time.Second*10)

	viper.SetDefault(confLogLevel, "info")
	viper.SetDefault(confLogFormat, "json")

	viper.SetDefault(confCacheType, "filesystem")
	viper.SetDefault(confCacheTTL, time.Hour)
	viper.SetDefault(confCacheRenewalThreshold, time.Minute*30)
	viper.SetDefault(confCacheFSDir, path.Join(cDir, "parrot"))
	viper.SetDefault(confCacheRedisMode, "single")
	viper.SetDefault(confCacheRedisMaxRetries, -1)
	viper.SetDefault(confCacheRedisDB, 1)

	viper.SetDefault(confPrometheusEnabled, true)
	viper.SetDefault(confPrometheusPort, 9090)
	viper.SetDefault(confPrometheusPath, "/metrics")
	viper.SetDefault(confAWSRegion, "eu-west-1")

	rootCmd.AddCommand(serveCmd)
}

func instantiateLogger() *logrus.Logger {
	logger := logrus.New()

	lvl, err := logrus.ParseLevel(viper.GetString(confLogLevel))
	if err != nil {
		logger.WithError(err).Warnf("Could not parse log level '%s' defaulting to INFO", viper.GetString(confLogLevel))
		lvl = logrus.InfoLevel
	}
	logger.SetLevel(lvl)

	switch viper.GetString(confLogFormat) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	default:
		logger.Warnf("Did not understand log format '%s'. Defaulting to json format", viper.GetString(confLogFormat))
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	logger.WithFields(logrus.Fields{
		"level":  viper.GetString(confLogLevel),
		"format": viper.GetString(confLogFormat),
	}).Info("Logger initialized")

	return logger
}

func instantiateCache(l *logrus.Entry) (cache.Cache, error) {
	cType := viper.GetString(confCacheType)
	switch cType {
	case "filesystem":
		return instantiateFilesystemCache()
	case "redis":
		redis.SetLogger(&cache.RedisLogger{Entry: l})

		return instantiateRedisCache()
	default:
		return nil, errors.Errorf("'%s' cache type is not yet implemented", cType)
	}
}

func instantiateRedisCache() (*cache.RedisCache, error) {
	switch viper.GetString(confCacheRedisMode) {
	case "sentinel":
		return cache.NewRedisCache(redis.NewFailoverClient(&redis.FailoverOptions{
			Username:   viper.GetString(confCacheRedisUser),
			Password:   viper.GetString(confCacheRedisPassword),
			MaxRetries: viper.GetInt(confCacheRedisMaxRetries),
			DB:         viper.GetInt(confCacheRedisDB),

			MasterName:       viper.GetString(confCacheRedisSentinelMaster),
			SentinelAddrs:    viper.GetStringSlice(confCacheRedisSentinelAddress),
			SentinelPassword: viper.GetString(confCacheRedisSentinelPassword),
		}), viper.GetDuration(confCacheTTL)), nil
	case "single":
		return cache.NewRedisCache(redis.NewClient(&redis.Options{
			Username:   viper.GetString(confCacheRedisUser),
			Password:   viper.GetString(confCacheRedisPassword),
			MaxRetries: viper.GetInt(confCacheRedisMaxRetries),
			DB:         viper.GetInt(confCacheRedisDB),

			Addr: viper.GetString(confCacheRedisAddress),
		}), viper.GetDuration(confCacheTTL)), nil
	case "cluster":
		return nil, errors.New("Cluster mode is not supported")
	case "ring":
		return nil, errors.New("Ring mode is not supported")
	default:
		return nil, errors.Errorf("Did not understand redis mode '%s'", viper.GetString(confCacheRedisMode))
	}
}

func instantiateFilesystemCache() (*cache.FilesystemCache, error) {
	return cache.NewFilesystemCache(viper.GetString(confCacheFSDir), viper.GetDuration(confCacheTTL))
}
