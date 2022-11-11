package handlers

import (
	"bruce/exe"
	"bruce/system"
	"github.com/rs/zerolog/log"
)

func StartPreExecCmds() {
	// start with pre execution cmds
	for _, v := range system.GetSysInfo().Configuration.PreExecCmds {
		pc := exe.Run(v, false)
		if pc.Failed() {
			log.Error().Err(pc.GetErr()).Msg(pc.Get())
		} else {
			log.Info().Msgf("completed: %s", v)
			log.Debug().Msgf("Output: %s", pc.Get())
		}
	}
}

func StartPostExecCmds() {
	// start with pre execution cmds

	for _, v := range system.GetSysInfo().Configuration.PostExecCmds {
		pc := exe.Run(v, false)
		if pc.Failed() {
			log.Error().Err(pc.GetErr()).Msg(pc.Get())
		} else {
			log.Info().Msgf("completed: %s", v)
			log.Debug().Msgf("Output: %s", pc.Get())
		}
	}
}
