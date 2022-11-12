package services

import (
	"bruce/exe"
	"bruce/system"
	"bruce/templates"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

var ()

func init() {

}

func StartOSServiceExecution() []string {
	var failedSvcs []string
	// TODO: Execute services and aggregate the list of ones that fail here
	if system.GetSysInfo().SystemType == "linux" {
		doDaemonReload := false
		for _, tpl := range templates.GetModifiedTemplates() {
			if strings.Contains(tpl, "systemd") {
				doDaemonReload = true
			}
		}
		if doDaemonReload {
			log.Debug().Msgf("daemon reload required due to service change")
			exe.Run("systemctl daemon-reload", system.GetSysInfo().TrySudo)
		}
		// We only support sytemd / systemctrl for right now...
		for _, svc := range system.GetSysInfo().Configuration.Services {
			status := exe.Run(fmt.Sprintf("systemctl is-active %s", svc.Name), system.GetSysInfo().TrySudo).Get()
			if strings.Contains(strings.ToLower(status), "could not be found") {
				log.Error().Err(fmt.Errorf("%s service not found", svc.Name)).Msg("service does not exist cannot manage state")
				continue
			}
			if svc.Enabled {
				// test if not enabled
				curState := exe.Run(fmt.Sprintf("systemctl is-enabled %s", svc.Name), system.GetSysInfo().TrySudo).Get()
				if curState != "enabled" {
					eno := exe.Run(fmt.Sprintf("systemctl enable %s --now", svc.Name), system.GetSysInfo().TrySudo).Get()
					log.Info().Str("output", eno).Msgf("set enabled for %s", svc.Name)
				}
			}

			if svc.State == "started" {
				if status != "active" {
					out := exe.Run(fmt.Sprintf("systemctl restart %s", svc.Name), system.GetSysInfo().TrySudo).Get()
					log.Info().Str("output", out).Msgf("issued restart to inactive service: %s", svc.Name)
				}
			}
			if svc.State == "stopped" {
				if status != "inactive" {
					out := exe.Run(fmt.Sprintf("systemctl stop %s", svc.Name), system.GetSysInfo().TrySudo).Get()
					log.Info().Str("output", out).Msgf("issued stop to active service: %s", svc.Name)
				}
			}
			if svc.RestartAlways {
				out := exe.Run(fmt.Sprintf("systemctl restart %s", svc.Name), system.GetSysInfo().TrySudo).Get()
				log.Info().Str("output", out).Msgf("issued restart (always) to service: %s", svc.Name)
			} else {
				for _, resTemp := range svc.RestartOnUpdate {
					for _, modT := range templates.GetModifiedTemplates() {
						if resTemp == modT {
							out := exe.Run(fmt.Sprintf("systemctl restart %s", svc.Name), system.GetSysInfo().TrySudo).Get()
							log.Info().Str("output", out).Msgf("issued restart (modified by template) to service: %s", svc.Name)
						}
					}
				}
			}
			// finally we recheck to see if it is started as we may have to revert
			status = exe.Run(fmt.Sprintf("systemctl is-active %s", svc.Name), system.GetSysInfo().TrySudo).Get()
			if svc.State == "started" {
				if status != "active" {
					log.Info().Str("status", status).Msgf("service is invalid state, need to revert: %s", svc.Name)
					failedSvcs = append(failedSvcs, svc.Name)
				}
			}
		}
	}
	return failedSvcs
}

func RestoreFailedServices(svcs []string) error {
	for _, svc := range svcs {
		for _, cs := range system.GetSysInfo().Configuration.Services {
			if svc == cs.Name {
				for _, srcName := range cs.RestartOnUpdate {
					log.Info().Msgf("restoring template %s", srcName)
					err := templates.RestoreBackupFile(srcName)
					if err != nil {
						log.Error().Err(err).Msg("could not restore template")
					}
				}
			}
		}
	}
	StartOSServiceExecution()
	return nil
}
