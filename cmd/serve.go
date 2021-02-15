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
	"net/http"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uniwise/parrot/internal/cache"
	"github.com/uniwise/parrot/internal/poedit"
	"github.com/uniwise/parrot/internal/project"
	"github.com/uniwise/parrot/internal/rest"
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

		c, err := instantiateCache()
		if err != nil {
			logger.Fatal(err)
		}

		cli := poedit.NewClient(viper.GetString("api.token"), http.DefaultClient)

		svc := project.NewService(cli, c)

		server, err := rest.NewServer(svc, logrus.NewEntry(logger))
		if err != nil {
			logger.Fatal(err)
		}
		port := viper.GetInt("server.port")
		logger.Infof("Server listening at :%d", port)
		logger.Fatal(server.Start(port))
	},
}

func init() {
	cDir, err := os.UserCacheDir()
	if err != nil {
		cDir = "/tmp"
	}

	viper.SetDefault("server.port", 80)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("cache.type", "filesystem")
	viper.SetDefault("cache.ttl", time.Hour)
	viper.SetDefault("cache.filesystem.dir", path.Join(cDir, "parrot"))
	viper.SetDefault("api.token", "")

	serveCmd.PersistentFlags().Int32("port", 0, "Port for the server to listen on")
	serveCmd.PersistentFlags().String("loglevel", "", "Log level")
	serveCmd.PersistentFlags().String("logformat", "", "Formatter for the logs")
	serveCmd.PersistentFlags().String("cache-type", "", "Which system to use for the underlying cache")
	serveCmd.PersistentFlags().Duration("cache-ttl", 0, "Time to live for the cache")
	serveCmd.PersistentFlags().String("cache-dir", "", "Directory where the filesystem cache lives")
	serveCmd.PersistentFlags().String("api-token", "", "API token to authenticate against poeditor")

	viper.BindPFlag("server.port", serveCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("log.level", serveCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("log.format", serveCmd.PersistentFlags().Lookup("logformat"))
	viper.BindPFlag("cache.type", serveCmd.PersistentFlags().Lookup("cache-type"))
	viper.BindPFlag("cache.ttl", serveCmd.PersistentFlags().Lookup("cache-ttl"))
	viper.BindPFlag("api.token", serveCmd.PersistentFlags().Lookup("api-token"))
	viper.BindPFlag("cache.filesystem.dir", serveCmd.PersistentFlags().Lookup("cache-dir"))

	rootCmd.AddCommand(serveCmd)
}

func instantiateLogger() *logrus.Logger {
	logger := logrus.New()

	lvl, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		logger.WithError(err).Warnf("Could not parse log level '%s' defaulting to INFO", viper.GetString("log.level"))
		lvl = logrus.InfoLevel
	}
	logger.SetLevel(lvl)

	switch viper.GetString("log.format") {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	default:
		logger.Warnf("Did not understand log format '%s'. Defaulting to json format", viper.GetString("log.format"))
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return logger
}

func instantiateCache() (cache.Cache, error) {
	cType := viper.GetString("cache.type")
	ttl := viper.GetDuration("cache.ttl")
	switch cType {
	case "filesystem":
		c, err := cache.NewFilesystemCache(viper.GetString("cache.filesystem.dir"), ttl)
		if err != nil {
			return nil, errors.Wrap(err, "Could not instantiate filesystem cache")
		}
		return c, nil
	case "redis":
		return nil, errors.Errorf("'%s' cache type is not yet implemented", cType)
	default:
		return nil, errors.Errorf("'%s' cache type is not yet implemented", cType)
	}
}
