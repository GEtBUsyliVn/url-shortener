package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Memory struct {
	CacheTtl time.Duration `envconfig:"MEMORY_CACHE_TTL"`
}

func (c *Memory) Prepare(prefix string) error {
	return envconfig.Process(prefix, c)
}
