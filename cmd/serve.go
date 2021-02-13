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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		server, err := rest.NewServer(viper.GetInt("server.port"), logrus.NewEntry(logger))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Fatal(server.Start())
	},
}

func init() {
	viper.SetDefault("server.port", 80)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	serveCmd.PersistentFlags().Int32("port", 0, "Port for the server to listen on")
	serveCmd.PersistentFlags().String("loglevel", "", "Log level")
	serveCmd.PersistentFlags().String("logformat", "", "Formatter for the logs")

	viper.BindPFlag("server.port", serveCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("log.level", serveCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("log.format", serveCmd.PersistentFlags().Lookup("logformat"))

	rootCmd.AddCommand(serveCmd)
}
