package handlers

import (
	"bruce/system"
	"bruce/templates"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

func CreateBackupLocation() {
	// First create a temporary backup directory where we will store existing templates
	backupDir := fmt.Sprintf("%s%c%s", system.GetSysInfo().Configuration.TempDir, os.PathSeparator, RandDirName(16))
	err := os.MkdirAll(backupDir, 0775)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot continue must have a backup dir for templates, please specify a temp directory that the user has access to create under")
	}
	si := system.GetSysInfo()
	si.Configuration.BackupDir = backupDir
	system.SetSysInfo(si)

	log.Debug().Msgf("created backup directory: %s", backupDir)
}

func BackupExistingTemplates() {
	// back up existing templates to be updated
	err := templates.BackupLocal(system.GetSysInfo().Configuration.BackupDir, system.GetSysInfo().Configuration.Templates)
	if err != nil {
		log.Fatal().Err(err).Msg("backup failed... cannot continue")
	}
}
