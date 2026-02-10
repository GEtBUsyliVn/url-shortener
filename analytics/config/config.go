package config

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

func InitConfig() *Config {
	data, err := os.ReadFile("/Users/Alex/IdeaProjects/url-shortener/analytics/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}
