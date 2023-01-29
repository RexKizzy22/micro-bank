package util

import (
	"time"

	"github.com/spf13/viper"
)

// stores all configurations of the application
// values are read by viper from the config file or environment variables
type Config struct {
	DBSource             string        `mapstructure:"DB_SOURCE"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	HTTP_ServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPC_ServerAddress   string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// makes env variables provided in the terminal have more priority than those in .env file
	viper.AutomaticEnv()
	viper.BindEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	viper.Unmarshal(&config)
	return
}
