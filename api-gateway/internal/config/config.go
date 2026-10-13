package config

import (
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

type service struct {
	Url       string   `yaml:"url"`
	Endpoints []string `yaml:"endpoints"`
}

type Config struct {
	Services map[string]service `yaml:"services"`
}

func NewConfig() (Config, error) {
	file, err := os.Open("config.yaml")
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	var c Config

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
