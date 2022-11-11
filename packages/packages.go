package packages

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
	"runtime"
	"strings"
)

var (
	packageHandler    = ""
	hasAlreadyUpdated = false
	curUser           = ""
	pfxCmd            = ""
)

func init() {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("whoami")
		d, err := cmd.CombinedOutput()
		if err == nil {
			curUser = strings.TrimSpace(string(d))
			log.Debug().Msgf("using package manager: %s", packageHandler)
		}
		if curUser != "root" {
			pfxCmd = "sudo"
			log.Debug().Msgf("package installs may fail, you're not root: %s", curUser)
		}
		cmd = exec.Command("which", "yum")
		d, err = cmd.CombinedOutput()
		if err == nil {
			packageHandler = strings.TrimSpace(string(d))
			log.Debug().Msgf("using package manager: %s", packageHandler)
			return
		}
		cmd = exec.Command("which", "apt")
		d, err = cmd.CombinedOutput()
		if err == nil {
			packageHandler = strings.TrimSpace(string(d))
			log.Debug().Msgf("using package manager: %s", packageHandler)
			return
		}
		cmd = exec.Command("which", "dnf")
		d, err = cmd.CombinedOutput()
		if err == nil {
			packageHandler = strings.TrimSpace(string(d))
			log.Debug().Msgf("using package manager: %s", packageHandler)
			return
		}
	}
}
func checkDebPackageInstalled(pkg string) bool {
	if pkg == "" {
		log.Debug().Msg("can't check for nothing")
		return false
	}
	log.Debug().Msgf("looking for %s package", pkg)
	cmd := exec.Command("/usr/bin/dpkg-query", "-s", pkg)
	d, err := cmd.CombinedOutput()
	if err != nil {
		log.Debug().Err(err).Msgf("error getting info on: %s", string(d))
		return false
	}
	if strings.Contains(strings.ToLower(string(d)), "is not installed") {
		return false
	}
	return true
}

func RunPackageInstall(pkg string) error {
	if packageHandler != "" {
		if !hasAlreadyUpdated {
			d, err := exec.Command(pfxCmd, packageHandler, "update", "-y").CombinedOutput()
			if err != nil {
				log.Debug().Msgf("error updating packages with (%s): %s", packageHandler, string(d))
				return err
			}
			hasAlreadyUpdated = true
			log.Debug().Msgf("successfully updated %s", packageHandler)
		}
		if checkDebPackageInstalled(pkg) {
			log.Info().Msgf("%s is already installed...", pkg)
			return nil
		}
		log.Debug().Msgf("package [%s] not installed... installing now", pkg)
		d, err := exec.Command(pfxCmd, packageHandler, "install", "-y", pkg).CombinedOutput()
		if err == nil {
			log.Debug().Msgf("package installed: %s", string(d))
		}
		if strings.Contains(strings.ToLower(string(d)), "unable to locate") {
			perr := fmt.Errorf("no install candidate for %s on %s", pkg, packageHandler)
			log.Error().Err(perr).Msgf("package %s does not exist with %s", pkg, packageHandler)
			return perr
		}
		log.Debug().Msgf("install output: %s", d)
	}
	return nil
}
