package configs

import (
	"ansible-runner/utils"

	"github.com/spf13/viper"
)

// LoadConfig loads the application configuration.
func LoadConfig(p string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(p)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		utils.Log.Error(err.Error())
	}
}


