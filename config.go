package api

type Config struct {
	Port int
}

func (c *Config) Validate() bool {
	return c.Port != 0
}
