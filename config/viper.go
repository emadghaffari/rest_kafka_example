package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	// Viper variable instance of viperInterface
	Viper viperInterface = &viperStruct{}
)

type viperInterface interface {
	Configs()
}
type viperStruct struct{}

// InitViper func
func (v *viperStruct) Configs() {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Config file not found; ignore error if desired"))
		} else {
			panic(fmt.Errorf("fatal error config file: %s \n ", err))
		}
	}
	viper.ConfigFileUsed()
}
