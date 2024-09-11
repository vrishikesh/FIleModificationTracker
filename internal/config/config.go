package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Directory   string `validate:"required"`
	CheckFreq   int    `validate:"required,gte=1"`
	APIEndpoint string `validate:"required,url"`
}

func loadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	config := &Config{
		Directory:   viper.GetString("directory"),
		CheckFreq:   viper.GetInt("check_freq"),
		APIEndpoint: viper.GetString("api_endpoint"),
	}

	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
