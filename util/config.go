package util

import (
	"time"

	"github.com/spf13/viper"
)

// stores all configurations of the application
// values are read by viper from the config file or environment variables
type Config struct {
	AppEnv               string        `mapstructure:"APP_ENV"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ProdDBSource         string        `mapstructure:"PROD_DB_SOURCE"`
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
	viper.SetEnvPrefix("app")
	
	// makes env variables provided in the terminal have more priority than those in .env file
	viper.AutomaticEnv()

	viper.SetDefault("PROD_DB_SOURCE", "")
	viper.SetDefault("PROD_DB_DRIVER", "")
	viper.SetDefault("APP_ENV", "")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	viper.Unmarshal(&config)
	return
}

func (config *Config) FetchDBSource() string {
	env := viper.GetString("APP_ENV")
	remoteDB := viper.GetString("PROD_DB_SOURCE")
	if env == "production" {
		return remoteDB
	} else {
		return config.DBSource
	}
}

func (config *Config) FetchDBDriver() string {
	env := viper.GetString("APP_ENV")
	remoteDBDRIVER := viper.GetString("PROD_DB_DRIVER")
	if env == "production" {
		return remoteDBDRIVER
	} else {
		return config.DBDriver
	}
}
