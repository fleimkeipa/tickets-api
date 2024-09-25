package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Function to load YAML file using Viper
func LoadEnv() error {
	// Set the file name and type
	viper.SetConfigName("config") // config.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // Look for the config file in the current directory

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	return nil
}
