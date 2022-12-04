package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

func (c *Config) parseFromOsArgs() {

	flag.StringVar(&c.servAPIAddr, "a", c.servAPIAddr, "grpc api server run address")
	flag.StringVar(&c.dirWithFiles, "d", c.dirWithFiles, "directory with files")
	flag.StringVar(&c.logLevel, "l", c.logLevel, "log level")


	flag.Parse()
}

func (c *Config) parseFromEnv() (err error) {

	envConfig := struct {
		ServAPIAddr               string        `env:"RUN_ADDRESS" toml:"RUN_ADDRESS"`
		LogLevel                  string        `env:"LOG_LEVEL" toml:"LOG_LEVEL"`
		DirWithFiles              string        `env:"DIR_WITH_FILES" toml:"DIR_WITH_FILES"`
	}{}

	if err = env.Parse(&envConfig); err != nil {
		return fmt.Errorf("parsing config from env: %w", err)
	}

	if envConfig.ServAPIAddr != "" {
		c.servAPIAddr = envConfig.ServAPIAddr
	}

	if envConfig.DirWithFiles != "" {
		c.dirWithFiles = envConfig.DirWithFiles
	}

	if envConfig.LogLevel != "" {
		c.logLevel = envConfig.LogLevel
	}

	return nil
}
