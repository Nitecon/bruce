package handlers

import (
	"bruce/config"
	"bruce/services"
	"bruce/templates"
	"github.com/rs/zerolog/log"
	"os"
)

func Update(arg string) error {

	log.Debug().Msg("starting Update task")
	cfg := config.Get()
	if _, err := os.Stat(cfg.Template.TempDir); os.IsExist(err) {
		err := os.RemoveAll(cfg.Template.TempDir)
		if err != nil {
			log.Info().Msgf("could not remove bruce temp directory, user removed?: %s", cfg.Template.BackupDir)
		}
	}
	CreateBackupLocation()
	RunCLICmds(cfg.Template.PreUpdateCmds)

	BackupExistingTemplates(cfg.Template.UpdateTemplates)

	log.Debug().Msg("starting template setup")
	templates.RenderTemplates(cfg.Template.UpdateTemplates)
	log.Debug().Msg("complete template setup")
	RunCLICmds(cfg.Template.PostTplUpdateCmds)
	// Now we set any ownership that must exist prior to service execution:
	StartOwnerships()

	// run the systemd enablement / restarts etc
	svcs := services.StartOSServiceReloads()
	if len(svcs) > 0 {
		log.Info().Msgf("following services are in the wrong state: %s", svcs)
	} else {
		log.Info().Msg("all services are in the appropriate state")
	}

	// post execution commands are next
	RunCLICmds(cfg.Template.PostUpdateExecCmds)
	return nil
}
