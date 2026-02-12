package config

type Config struct {
	GRPC
	Redis
	Worker
	Memory
}

func (c *Config) Prepare(prefix string) error {
	if err := c.Redis.Prepare(prefix); err != nil {
		return err
	}

	if err := c.GRPC.Prepare(prefix); err != nil {
		return err
	}

	if err := c.Worker.Prepare(prefix); err != nil {
		return err
	}

	if err := c.Memory.Prepare(prefix); err != nil {
		return err
	}

	return nil
}
