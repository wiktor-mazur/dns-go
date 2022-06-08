package utils

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	UDP_ENABLED bool   `default:"true"`
	UDP_IP      string `default:"0.0.0.0"`
	UDP_PORT    int    `default:"8053"`

	INTERNET_ROOT_SERVER string
}

func LoadConfig() (*Config, error) {
	result := &Config{}

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	err = envconfig.Process("", result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
