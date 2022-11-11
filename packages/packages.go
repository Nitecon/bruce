package packages

import (
	"bruce/system"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

func InstallOSPackage(pkgs []string) bool {
	if len(pkgs) < 1 {
		log.Error().Err(fmt.Errorf("can't install nothing"))
		return false
	}
	switch system.GetSysInfo().PackageHandler {
	case "apt":
		return installAptPackage(GetManagerPackages(pkgs, "apt"))
	case "yum":
		return installYumPackage(GetManagerPackages(pkgs, "yum"))
	case "dnf":
		return installDnfPackage(GetManagerPackages(pkgs, "dnf"))
	}
	log.Info().Msg("no package manager to check for installed package")
	return false
}

func GetManagerPackages(pkgs []string, manager string) []string {
	var newList []string
	for _, pkg := range pkgs {
		if strings.Contains(pkg, "|") {
			managerList := strings.Split(pkg, "|")
			var basePackage = ""
			var usablePackage = ""
			for _, mpkg := range managerList {
				if strings.Contains(mpkg, "=") {
					pmSplit := strings.Split(mpkg, "=")
					if pmSplit[0] == manager {
						usablePackage = pmSplit[1]
					}
				} else {
					basePackage = mpkg
				}
			}
			if usablePackage != "" {
				newList = append(newList, usablePackage)
			} else {
				newList = append(newList, basePackage)
			}
		}
	}
	return newList
}

func DoPackageManagerUpdate() bool {
	updateComplete := false
	switch system.GetSysInfo().PackageHandler {
	case "apt":
		updateComplete = updateApt()
		break
	case "yum":
		updateComplete = updateYum()
		break
	case "dnf":
		updateComplete = updateDnf()
		break
	}
	system.SetPackageMangerUpdated(updateComplete)
	if !updateComplete {
		log.Info().Msg("no package manager to check for installed package")
	}
	return false
}

func RunPackageInstall(pkgs []string) error {
	if InstallOSPackage(pkgs) {
		log.Info().Msgf("[%s] installed", pkgs)
	}

	log.Error().Err(fmt.Errorf("failed to install [%s]", pkgs))
	return nil
}
