package packages

import (
	"bruce/exe"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

func updateApt() bool {
	return !exe.Run("/usr/bin/apt-get update -y", false).Failed()
}

func installAptPackage(pkg []string, isInstall bool) bool {
	action := "install"
	if !isInstall {
		action = "remove"
	}
	installCmd := fmt.Sprintf("/usr/bin/apt-get %s -y %s", action, strings.Join(pkg, " "))
	log.Debug().Msgf("apt install starting with: %s", installCmd)
	install := exe.Run(installCmd, false)
	if install.Failed() {
		if len(install.Get()) > 0 {
			strSplit := strings.Split(install.Get(), "\n")
			log.Error().Err(install.GetErr())
			for _, s := range strSplit {
				log.Info().Msg(s)
			}
		}
		return false
	}
	return true
}
