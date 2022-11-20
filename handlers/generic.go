package handlers

import (
	"bruce/config"
	"bruce/exe"
	"github.com/rs/zerolog/log"
	"os"
)

func RunCLICmds(cmds []string) {
	// start with pre execution cmds
	for _, v := range cmds {
		fileName := exe.EchoToFile(v)
		err := os.Chmod(fileName, 0775)
		if err != nil {
			log.Fatal().Err(err).Msg("temp file must exist to continue")
		}
		log.Debug().Str("command", v).Msgf("executing local file: %s", fileName)
		pc := exe.Run(fileName, false)
		if pc.Failed() {
			log.Error().Err(pc.GetErr()).Msg(pc.Get())
		} else {
			log.Info().Msgf("completed executing: %s", fileName)
			log.Debug().Msgf("Output: %s", pc.Get())
			os.Remove(fileName)
		}
	}
}

func StartPreExecCmds() {
	// start with pre execution cmds
	for _, v := range config.Get().Template.PreExecCmds {
		fileName := exe.EchoToFile(v)
		err := os.Chmod(fileName, 0775)
		if err != nil {
			log.Fatal().Err(err).Msg("temp file must exist to continue")
		}
		log.Debug().Str("command", v).Msgf("executing local file: %s", fileName)
		pc := exe.Run(fileName, false)
		if pc.Failed() {
			log.Error().Err(pc.GetErr()).Msg(pc.Get())
		} else {
			log.Info().Msgf("completed executing: %s", fileName)
			log.Debug().Msgf("Output: %s", pc.Get())
			os.Remove(fileName)
		}
	}
}

func StartPostInstallCmds() {
	// start the post installation commands
	for _, v := range config.Get().Template.PostInstallCmds {
		fileName := exe.EchoToFile(v)
		err := os.Chmod(fileName, 0775)
		if err != nil {
			log.Fatal().Err(err).Msg("temp file must exist to continue")
		}
		log.Debug().Msgf("executing local file: %s", fileName)
		pc := exe.Run(fileName, false)
		if pc.Failed() {
			log.Error().Err(pc.GetErr()).Msg(pc.Get())
		} else {
			log.Info().Msgf("completed executing: %s", fileName)
			log.Debug().Msgf("Output: %s", pc.Get())
			os.Remove(fileName)
		}
	}
}

func StartOwnerships() {
	for _, v := range config.Get().Template.OwnerShips {
		err := exe.SetOwnership(v)
		if err != nil {
			log.Error().Err(err).Msg("could not set ownership")
		}
	}
}

func StartPostExecCmds() {
	// start with pre execution cmds
	for _, v := range config.Get().Template.PostExecCmds {
		fileName := exe.EchoToFile(v)
		err := os.Chmod(fileName, 0775)
		if err != nil {
			log.Fatal().Err(err).Msg("temp file must exist to continue")
		}
		log.Debug().Msgf("executing local file: %s", fileName)
		pc := exe.Run(fileName, config.Get().TrySudo)
		if pc.Failed() {
			log.Fatal().Err(pc.GetErr()).Msg(pc.Get())
		} else {
			log.Info().Msgf("completed executing: %s", fileName)
			log.Debug().Msgf("Output: %s", pc.Get())
			os.Remove(fileName)
		}
	}
}
