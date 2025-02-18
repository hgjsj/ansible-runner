package configs

import (
	"github.com/spf13/viper"
	"log"
)

// LoadConfig loads the application configuration.
func LoadConfig(p string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(p)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}
