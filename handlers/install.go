package handlers

import (
	"bruce/config"
	"fmt"
	"github.com/rs/zerolog/log"
)

func Install(cfgf, arg string) error {
	cfg, err := config.LoadConfig(cfgf)
	if err != nil {
		log.Fatal().Err(err).Msg("install cannot continue without config")
	}
	fmt.Printf("%#v", cfg)
	return nil
}
