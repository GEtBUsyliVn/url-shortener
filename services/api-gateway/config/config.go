package config

type Config struct {
	GRPC
}

func (c *Config) Prepare(prefix string) error {
	if err := c.GRPC.Prepare(prefix); err != nil {
		return err
	}
	return nil
}
