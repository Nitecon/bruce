package packages

import (
	"bruce/exe"
	"bruce/system"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

func updateYum() bool {
	return !exe.Run("yum update -y", system.GetSysInfo().TrySudo).Failed()
}

func installYumPackage(pkg []string) bool {

	install := exe.Run(fmt.Sprintf("yum install -y %s", strings.Join(pkg, " ")), system.GetSysInfo().TrySudo)
	if install.Failed() {
		log.Error().Err(install.GetErr()).Msg("error installing")
		return false
	}
	return true
}
