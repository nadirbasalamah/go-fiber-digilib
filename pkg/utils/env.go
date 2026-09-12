package utils

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func GetConfig(key string) string {
	isDevelop := os.Getenv("APP_MODE") != "production"

	if isDevelop {
		viper.AddConfigPath(".")
		viper.SetConfigFile(".env")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("error when reading configuration file: %s\n", err)
		}

		return viper.GetString(key)
	}

	return os.Getenv(key)
}
