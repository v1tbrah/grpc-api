package config

type Config struct {
	servAPIAddr  string
	logLevel     string
	dirWithFiles string
}

func New(options ...string) (newCfg *Config, err error) {

	newCfg = &Config{}

	for _, opt := range options {
		switch opt {
		case WithFlag:
			newCfg.parseFromOsArgs()
		case WithEnv:
			if err = newCfg.parseFromEnv(); err != nil {
				return nil, err
			}
		}
	}

	newCfg.setDefaultIfNotConfigured()

	return newCfg, nil
}

func (c *Config) setDefaultIfNotConfigured() {

	if c.servAPIAddr == "" {
		c.servAPIAddr = ":8080"
	}

	if c.dirWithFiles == "" {
		c.dirWithFiles = "filesSavedInGRPCServer"
	}

	if c.logLevel == "" {
		c.logLevel = "info"
	}

}

func (c *Config) ServAPIAddr() string {
	return c.servAPIAddr
}

func (c *Config) LogLevel() string {
	return c.logLevel
}

func (c *Config) DirWithFiles() string {
	return c.dirWithFiles
}

func (c *Config) String() string {
	return "servAPIAddr: " + c.servAPIAddr +
		"dirWithFiles: " + c.dirWithFiles +
		" logLevel" + c.LogLevel()
}
