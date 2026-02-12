package config

type Config struct {
	Database
	GRPC
	Worker
}

func (c *Config) Prepare(prefix string) error {

	if err := c.Database.Prepare(prefix); err != nil {
		return err
	}

	if err := c.GRPC.Prepare(prefix); err != nil {
		return err
	}

	if err := c.Worker.Prepare(prefix); err != nil {
		return err
	}

	return nil
}
