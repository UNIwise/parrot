package rest

import (
	// Used to load .env files for environment variables
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

// Config contains environment configurable variables
type Config struct {
	Port         uint32 `envconfig:"API_PORT" default:"8080"`
	GithubSecret string `envconfig:"GITHUB_SECRET" required:"true"`
}

// NewEnvConfig is the constructor for EnvConfig
func NewEnvConfig() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	return &c, err
}
