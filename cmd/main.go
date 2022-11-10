package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
	"os"
)

var (
	version = "source"
)

func setLogger() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	if os.Getenv("BRUCE_DEBUG") != ""{
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
			return
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	setLogger()
	log.Info().Msgf("Starting Bruce (Version: %s)", version)
}