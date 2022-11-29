package operators

import (
	"bruce/packages"
	"bruce/system"
	"fmt"
	"github.com/rs/zerolog/log"
)

type Packages struct {
	Name        string   `yaml:"name"`
	PackageList []string `yaml:"packageList"`
	OsLimits    string   `yaml:"osLimits"`
}

func (p *Packages) Execute() error {
	if system.Get().CanExecOnOs(p.OsLimits) {
		log.Info().Msgf("starting package installs for %s", system.Get().PackageHandler)
		success := packages.InstallOSPackage(p.PackageList, system.Get().PackageHandler)
		if !success {
			err := fmt.Errorf("cannot install packages: %s", p.PackageList)
			log.Error().Err(err).Msg("package install failed")
			return err
		}
		return nil
	} else {
		si := system.Get()
		log.Debug().Msgf("System (%s|%s) limited execution of installs for: %s", si.OSID, si.OSVersionID, p.OsLimits)
	}
	return nil
}