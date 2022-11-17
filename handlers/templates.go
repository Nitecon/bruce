package handlers

import (
	"bruce/config"
	"bruce/random"
	"bruce/templates"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

func CreateBackupLocation() {
	// First create a temporary backup directory where we will store existing templates
	cfg := config.Get()
	backupDir := fmt.Sprintf("%s%c%s", cfg.Configuration.TempDir, os.PathSeparator, random.String(16))
	err := os.MkdirAll(backupDir, 0775)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot continue must have a backup dir for templates, please specify a temp directory that the user has access to create under")
	}
	cfg.Configuration.BackupDir = backupDir
	cfg.Save()

	log.Debug().Msgf("created backup directory: %s", backupDir)
}

func BackupExistingTemplates() {
	// back up existing templates to be updated
	cfg := config.Get()
	err := templates.BackupLocal(cfg.Configuration.BackupDir, cfg.Configuration.Templates)
	if err != nil {
		log.Fatal().Err(err).Msg("backup failed... cannot continue")
	}
}
