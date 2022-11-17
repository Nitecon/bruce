package packages

import (
	"bruce/config"
	"bruce/exe"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

func updateApt() bool {
	return !exe.Run("apt update -y", config.Get().TrySudo).Failed()
}

func installAptPackage(pkg []string) bool {
	installCmd := fmt.Sprintf("apt install -y %s", strings.Join(pkg, " "))
	log.Debug().Msgf("apt install starting with: %s", installCmd)
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
