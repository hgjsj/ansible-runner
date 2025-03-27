package utils

import (
	"path/filepath"

	"github.com/spf13/viper"
)

// LoadConfig loads the application configuration.
func LoadConfig(p string) {
	configDir, configFileName := filepath.Split(p)

	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}
}
