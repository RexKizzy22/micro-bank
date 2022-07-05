package util

import (
	"time"

	"github.com/spf13/viper"
)

// stores all configurations of the application
// values are read by viper from the config file or environment variables
type Config struct {
	DBSource          string        `mapstructure:"DB_SOURCE"`
	DBDriver          string        `mapstructure:"DB_DRIVER"`
	ServerAddress     string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenDuration     time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	viper.Unmarshal(&config)
	return
}
