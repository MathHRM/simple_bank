package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBdriver      string `mapstructure:"DB_DRIVER"`
	DBsource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	TokenSimetricKey string `mapstructure:"TOKEN_SIMETRIC_KEY"`
	AcessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`

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

	err = viper.Unmarshal(&config)

	return
}
