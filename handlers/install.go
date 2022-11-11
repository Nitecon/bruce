package handlers

import (
	"bruce/config"
	"bruce/templates"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

func Install(cfgf, arg string) error {
	cfg, err := config.LoadConfig(cfgf)
	if err != nil {
		log.Fatal().Err(err).Msg("install cannot continue without config")
	}
	log.Debug().Msg("starting install task")
	// start with pre execution cmds
	for _, v := range cfg.PreExecCmds {
		log.Info().Msgf("executing: %s", v)
	}
	// First create a temporary backup directory where we will store existing templates
	backupDir := fmt.Sprintf("%s%c%s", cfg.TempDir, os.PathSeparator, RandDirName(16))
	err = os.MkdirAll(backupDir, 0775)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot continue must have a backup dir for templates, please specify a temp directory that the user has access to create under")
	}
	log.Debug().Msgf("created backup directory: %s", backupDir)
	// concurrently write out the templates (make backups for sys restarts)
	err = templates.BackupLocal(backupDir, cfg.Templates)
	if err != nil {
		log.Info().Err(err).Msg("backup failed... okay to continue?")
	}
	// now install the list of packages
	templates.RenderTemplates(cfg.Templates)

	// run the systemd enablement / restarts etc

	// post execution commands are next

	// now we do cleanup of our backup directories if everything went well!
	err = os.RemoveAll(backupDir)
	if err != nil {
		log.Info().Msgf("could not remove backup directory, user removed?: ", backupDir)
	}
	return nil
}
