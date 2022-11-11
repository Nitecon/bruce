package packages

import (
	"bruce/exe"
	"bruce/system"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

func updateDnf() bool {
	return !exe.Run("dnf update -y", system.GetSysInfo().TrySudo).Failed()
}

func installDnfPackage(pkg []string) bool {
	install := exe.Run(fmt.Sprintf("apt install -y %s", strings.Join(pkg, " ")), system.GetSysInfo().TrySudo)
	if install.Failed() {
		log.Error().Err(install.GetErr()).Msg("error installing")
		return false
	}
	return true
}
