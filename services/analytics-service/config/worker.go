package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Worker struct {
	Interval time.Duration `envconfig:"AGGREGATOR_WORKER_INTERVAL"`
}

func (c *Worker) Prepare(prefix string) error {
	return envconfig.Process(prefix, c)
}
