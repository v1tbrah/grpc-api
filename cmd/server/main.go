package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"grpc-api/internal/server/api"
	"grpc-api/internal/server/config"
)

func main() {

	cfgOptions := []string{config.WithFlag, config.WithEnv}
	newCfg, err := config.New(cfgOptions...)
	if err != nil {
		log.Fatal().Err(err).Strs("cfg options", cfgOptions).Msg("creating new config")
	}

	logLevel, err := zerolog.ParseLevel(newCfg.LogLevel())
	if err != nil {
		log.Fatal().Err(err).Strs("cfg options", cfgOptions).Msg("parsing log level")
	}
	zerolog.SetGlobalLevel(logLevel)

	newAPI, err := api.New(newCfg)
	if err != nil {
		log.Fatal().Err(err).Str("cfg", newCfg.String()).Msg("creating new api")
	}

	newAPI.Run()

}
