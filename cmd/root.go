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
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const asciiart string = "ICAgICAoXA0KICAgICggIFwgIC8obylcDQogICAgKCAgIFwvICAoKS8gLykNCiAgICAgKCAgIGA7LikpJyIuKQ0KICAgICAgYCgvLy8vLy4tJw0KICAgPT09PT0pKT0pKT09PSgpDQogICAgIC8vLycNCiAgICAvLw0KICAgJw=="

func ascii() string {
	b, _ := base64.StdEncoding.DecodeString(asciiart)

	return string(b)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "parrot",
	Short: "Your friendly neighborhood Poeditor pull-through-cache",
	Long: fmt.Sprintf(`Parrot is a pull-through-cache for poeditor.
%s
Serve translations from poeditor instead of baking them into your frontend!
Parrot is designed to act as a wrapper for poeditor, so you can update
translations without rebuilding your frontend.
`, ascii()),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.parrot.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".parrot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".parrot")
	}

	godotenv.Load(".env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
