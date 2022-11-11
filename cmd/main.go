package main

import (
	"bruce/handlers"
	"bruce/packages"
	"bruce/system"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

var (
	version = "source"
)

func setLogger() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	if os.Getenv("BRUCE_DEBUG") != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	setLogger()
	system.InitSysInfo()
	packages.DoPackageManagerUpdate()
	app := &cli.App{
		Name:  "bruce",
		Usage: "By default will load config from /etc/bruce/config.yml",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "/etc/bruce/config.yml",
				Usage: "location where the config file will be example: https://s3.amazonaws.com/somebucket/my_install.yml",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "run package installs and template configuration & does daemon-reload & service restarts (systemd)",
				Action: func(cCtx *cli.Context) error {
					system.LoadConfig(cCtx.String("config"))
					handlers.Install(cCtx.Args().First())
					return nil
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "no package installs... run template updates & restarts only (optional)",
				/*Subcommands: []*cli.Command{
					{
						Name:  "template",
						Usage: "update a template",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("new task template: ", cCtx.Args().First())
							return nil
						},
					},
					{
						Name:  "restart",
						Usage: "restart a service",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("restarting service: ", cCtx.Args().First())
							return nil
						},
					},
				},*/
				Action: func(cCtx *cli.Context) error {
					fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
			},
			{
				Name:    "upgrade",
				Aliases: []string{"ug"},
				Usage:   "upgrade just packages but do not touch templates",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err)
	}
	//log.Info().Msgf("Starting Bruce (Version: %s)", version)
}
