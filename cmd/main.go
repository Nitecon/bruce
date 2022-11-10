package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
	"fmt"
	"os"
	"github.com/urfave/cli/v2"
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
	
	app := &cli.App{
		Name:  "bruce",
		Usage: "specify a configuration file to be used.",
		Action: func(*cli.Context) error {
			fmt.Println("Batman has nothing on me!")
			return nil
		},
	}
	
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err)
	}
	//log.Info().Msgf("Starting Bruce (Version: %s)", version)
}