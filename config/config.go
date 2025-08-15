package config

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Port int
	}
}

var configPublic *Config
var once sync.Once

func LoadConfig() (config *Config) {
	once.Do(func() {
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			log.Fatal("CONFIG_PATH is not set")
		}
		//viper.SetConfigFile("config/config.toml")
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&config); err != nil {
			panic("ERROR load config file!")
		}
		configPublic = config
		log.Println("================ Loaded Configuration ================")
	})
	return
}

func GetConfig() *Config {
	return configPublic
}
