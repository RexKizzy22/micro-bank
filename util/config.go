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
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	HTTP_ServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPC_ServerAddress   string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_Address"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.SetEnvPrefix("app")

	// makes env variables provided in the terminal have more priority than those in .env file
	viper.AutomaticEnv()

	viper.SetDefault("PROD_DB_SOURCE", "")
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
	remoteDBConnectionString := viper.GetString("PROD_DB_SOURCE")

	if env == "production" {
		return remoteDBConnectionString
	} else {
		return config.DBSource
	}
}
