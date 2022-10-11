package config

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

// Config is used to get cofigurations of the project
type Config struct {
	Host                   string `envconfig:"HOST" required:"true"`
	Port                   string `envconfig:"PORT" required:"true"`
	PostgresHost           string `envconfig:"POSTGRES_HOST" required:"true"`
	PostgresPort           string `envconfig:"POSTGRES_PORT" required:"true"`
	PostgresUser           string `envconfig:"POSTGRES_USER" required:"true"`
	PostgresPassword       string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	PostgresDB             string `envconfig:"POSTGRES_DB" required:"true"`
	PostgresMigrationsPath string `envconfig:"POSTGRES_MIGRATIONS_PATH" required:"true"`
}

// Load loads configurations values
func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
