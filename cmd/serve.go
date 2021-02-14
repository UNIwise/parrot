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
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uniwise/parrot/internal/cache"
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

		server, err := rest.NewServer(nil, c, logrus.NewEntry(logger))
		if err != nil {
			logger.Fatal(err)
		}
		port := viper.GetInt("server.port")
		logger.Infof("Server listening at :%d", port)
		logger.Fatal(server.Start(port))
	},
}

func init() {
	viper.SetDefault("server.port", 80)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("cache.type", "filesystem")
	viper.SetDefault("cache.ttl", time.Hour)

	serveCmd.PersistentFlags().Int32("port", 0, "Port for the server to listen on")
	serveCmd.PersistentFlags().String("loglevel", "", "Log level")
	serveCmd.PersistentFlags().String("logformat", "", "Formatter for the logs")
	serveCmd.PersistentFlags().String("cache-type", "", "Which system to use for the underlying cache")
	serveCmd.PersistentFlags().Duration("cache-ttl", 0, "Time to live for the cache")

	viper.BindPFlag("server.port", serveCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("log.level", serveCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("log.format", serveCmd.PersistentFlags().Lookup("logformat"))
	viper.BindPFlag("cache.type", serveCmd.PersistentFlags().Lookup("cache-type"))
	viper.BindPFlag("cache.ttl", serveCmd.PersistentFlags().Lookup("cache-ttl"))

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
		c, err := cache.NewFileCache(ttl)
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
