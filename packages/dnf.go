package packages

import (
	"bruce/config"
	"bruce/exe"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

func updateDnf() bool {
	return !exe.Run("dnf update -y", config.Get().TrySudo).Failed()
}

func installDnfPackage(pkg []string) bool {
	installCmd := fmt.Sprintf("dnf install -y %s", strings.Join(pkg, " "))
	log.Debug().Msgf("dnf install starting with: %s", installCmd)
	install := exe.Run(installCmd, config.Get().TrySudo)
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
