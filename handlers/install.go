package handlers

import (
	"bruce/config"
	"bruce/packages"
	"bruce/services"
	"bruce/templates"
	"github.com/rs/zerolog/log"
	"os"
)

func Install(arg string) error {

	log.Debug().Msg("starting install task")

	// Do initial cleanup of the backup dirs not after, so we have backups in case we need it!
	if _, err := os.Stat(config.Get().Configuration.TempDir); os.IsExist(err) {
		err := os.RemoveAll(config.Get().Configuration.TempDir)
		if err != nil {
			log.Info().Msgf("could not remove bruce temp directory, user removed?: %s", config.Get().Configuration.BackupDir)
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
	log.Debug().Msg("starting post package installation commands")
	StartPostInstallCmds()
	log.Debug().Msg("completed post package installation commands")

	BackupExistingTemplates()

	log.Debug().Msg("starting template setup")
	templates.RenderTemplates()
	log.Debug().Msg("complete template setup")

	// Now we set any ownership that must exist prior to service execution:
	StartOwnerships()

	// run the systemd enablement / restarts etc
	svcs := services.StartOSServiceExecution()

	/*if len(svcs) > 0 {
		log.Info().Msgf("restoring services & associated templates: %s", svcs)
		rerr := services.RestoreFailedServices(svcs)
		if rerr != nil {
			log.Error().Err(err).Msg("could not install packages")
		}
	}*/
	if len(svcs) > 0 {
		log.Info().Msgf("following services are in the wrong state: %s", svcs)
	} else {
		log.Info().Msg("all services are in the appropriate state")
	}

	// post execution commands are next
	StartPostExecCmds()
	return nil
}
