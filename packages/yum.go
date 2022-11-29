package packages

import (
	"bruce/exe"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

func updateYum() bool {
	return !exe.Run("/usr/bin/yum update -y", false).Failed()
}

func installYumPackage(pkg []string) bool {
	installCmd := fmt.Sprintf("/usr/bin/yum install -y %s", strings.Join(pkg, " "))
	log.Debug().Msgf("/usr/bin/yum install starting with: %s", installCmd)
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
