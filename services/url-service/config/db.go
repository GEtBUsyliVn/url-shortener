package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Host     string `envconfig:"DB_HOST"`
	Port     int    `envconfig:"DB_PORT"`
	DataBase string `envconfig:"DATABASE"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
}

func (c *Database) Prepare(prefix string) error {
	return envconfig.Process(prefix, c)
}
