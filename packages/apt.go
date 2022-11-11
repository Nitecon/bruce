package packages

import (
	"bruce/exe"
	"bruce/system"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

func updateApt() bool {
	return !exe.Run("apt update -y", system.GetSysInfo().TrySudo).Failed()
}

func installAptPackage(pkg []string) bool {
	install := exe.Run(fmt.Sprintf("apt install -y %s", strings.Join(pkg, " ")), system.GetSysInfo().TrySudo)
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
