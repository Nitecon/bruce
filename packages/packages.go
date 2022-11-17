package packages

import (
	"bruce/config"
	"bruce/exe"
	"fmt"
	"github.com/rs/zerolog/log"
	"path"
	"strings"
)

func InstallOSPackage(pkgs []string) bool {
	cfg := config.Get()
	cfg.PackageHandlerPath = GetLinuxPackageHandler()
	cfg.PackageHandler = path.Base(cfg.PackageHandlerPath)
	svcInfo, err := GetLinuxServiceController()
	if err != nil {
		cfg.CanUpdateServices = false
	} else {
		cfg.ServiceControllerPath = svcInfo
		cfg.ServiceController = path.Base(svcInfo)
	}
	DoPackageManagerUpdate(cfg)
	cfg.Save()
	if len(pkgs) < 1 {
		log.Error().Err(fmt.Errorf("can't install nothing"))
		return false
	}
	switch cfg.PackageHandler {
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

func GetLinuxServiceController() (string, error) {
	// We only support systemctl for now
	sysPath := exe.HasExecInPath("systemctl")
	if sysPath == "" {
		return "", fmt.Errorf("systemctl not found on this system")
	}
	return sysPath, nil
}

func GetLinuxPackageHandler() string {
	packageHandler := "/usr/bin/yum"
	if exe.FileExists(packageHandler) {
		log.Debug().Msgf("using package manager: %s", packageHandler)
		return packageHandler
	}
	packageHandler = "/usr/bin/apt"
	if exe.FileExists(packageHandler) {
		log.Debug().Msgf("using package manager: %s", packageHandler)
		return packageHandler
	}
	packageHandler = "/usr/bin/dnf"
	if exe.FileExists(packageHandler) {
		log.Debug().Msgf("using package manager: %s", packageHandler)
		return packageHandler
	}
	log.Error().Err(fmt.Errorf("no package handler")).Msg("could not find a supported package handler for this system")
	return ""
}

func GetManagerPackages(pkgs []string, manager string) []string {
	var newList []string
	for _, pkg := range pkgs {
		log.Debug().Msgf("package iteration: %#v", pkg)
		if strings.Contains(pkg, "|") {
			managerList := strings.Split(pkg, "|")
			var basePackage = ""
			var usablePackage = ""
			for _, mpkg := range managerList {
				log.Debug().Msgf("package iteration for manager: %#v", mpkg)
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
		// no manager substitutes so just add it
		newList = append(newList, pkg)
	}
	return newList
}

func DoPackageManagerUpdate(cfg *config.SysInfo) bool {
	updateComplete := false
	switch cfg.PackageHandler {
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
	cfg.PackageManagerUpdated = updateComplete
	if !updateComplete {
		log.Info().Msg("no package manager to check for installed package, during packaging update")
	}
	return false
}

func RunPackageInstall() error {

	pkgs := config.Get().Configuration.InstallPackages
	if InstallOSPackage(pkgs) {
		log.Info().Msgf("[%s] installed", pkgs)
		return nil
	}

	log.Error().Err(fmt.Errorf("failed to install [%s]", pkgs))
	return nil
}
