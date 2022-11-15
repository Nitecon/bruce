package handlers

import (
	"bruce/packages"
	"bruce/services"
	"bruce/system"
	"bruce/templates"
	"github.com/rs/zerolog/log"
	"os"
)

func Install(arg string) error {

	log.Debug().Msg("starting install task")

	// Do initial cleanup of the backup dirs not after, so we have backups in case we need it!
	if _, err := os.Stat(system.GetSysInfo().Configuration.TempDir); os.IsExist(err) {
		err := os.RemoveAll(system.GetSysInfo().Configuration.TempDir)
		if err != nil {
			log.Info().Msgf("could not remove bruce temp directory, user removed?: %s", system.GetSysInfo().Configuration.BackupDir)
		}
	}

	StartPreExecCmds()
	CreateBackupLocation()

	// run package installers
	log.Debug().Msg("starting package installs")
	err := packages.RunPackageInstall()
	if err != nil {
		log.Error().Err(err).Msg("could not install packages")
	}
	log.Debug().Msg("package installs complete")
	BackupExistingTemplates()

	log.Debug().Msg("starting template setup")
	templates.RenderTemplates()
	log.Debug().Msg("complete template setup")

	log.Debug().Msg("starting post package installation commands")
	StartPostInstallCmds()
	log.Debug().Msg("completed post package installation commands")

	// run the systemd enablement / restarts etc
	svcs := services.StartOSServiceExecution()

	if len(svcs) > 0 {
		log.Info().Msgf("restoring services & associated templates: %s", svcs)
		rerr := services.RestoreFailedServices(svcs)
		if rerr != nil {
			log.Error().Err(err).Msg("could not install packages")
		}
	}
	// post execution commands are next
	StartPostExecCmds()
	return nil
}
