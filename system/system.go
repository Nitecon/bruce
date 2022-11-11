package system

import (
	"bruce/exe"
	"fmt"
	"github.com/rs/zerolog/log"
	"path"
	"runtime"
	"sync"
)

var (
	sysinfo     *SysInfo
	sysinfoLock = new(sync.RWMutex)
)

type SysInfo struct {
	CurrentUser           string
	TrySudo               bool
	PackageHandler        string
	PackageHandlerPath    string
	PackageManagerUpdated bool
	SystemType            string
}

func SetPackageMangerUpdated(isUpdated bool) {
	cfg := GetSysInfo()
	cfg.PackageManagerUpdated = isUpdated
	SetSysInfo(cfg)
}

// GetSysInfo function returns the currently set global system information to be used.
func GetSysInfo() *SysInfo {
	sysinfoLock.RLock()
	defer sysinfoLock.RUnlock()
	return sysinfo

}

// SetSysInfo sets the global system configuration
func SetSysInfo(cfg *SysInfo) *SysInfo {
	sysinfoLock.Lock()
	defer sysinfoLock.Unlock()
	sysinfo = cfg
	return sysinfo
}

func InitSysInfo() *SysInfo {
	cfg := &SysInfo{}
	if runtime.GOOS == "linux" {
		cfg.CurrentUser = exe.Run("whoami", false).Get()
		if cfg.CurrentUser != "root" {
			cfg.TrySudo = true
			log.Debug().Msgf("attempting sudo, current user: %s", cfg.CurrentUser)
		}
		cfg.PackageHandlerPath = GetLinuxPackageHandler()
		cfg.PackageHandler = path.Base(cfg.PackageHandlerPath)
		cfg.SystemType = "linux"
	}
	return SetSysInfo(cfg)
}

func GetLinuxPackageHandler() string {
	packageHandler := exe.HasExecInPath("yum")
	if packageHandler != "" {
		log.Debug().Msgf("using package manager: %s", packageHandler)
		return packageHandler
	}
	packageHandler = exe.HasExecInPath("apt")
	if packageHandler != "" {
		log.Debug().Msgf("using package manager: %s", packageHandler)
		return packageHandler
	}
	packageHandler = exe.HasExecInPath("dnf")
	if packageHandler != "" {
		log.Debug().Msgf("using package manager: %s", packageHandler)
		return packageHandler
	}
	log.Error().Err(fmt.Errorf("no package handler")).Msg("could not find a supported package handler for this system")
	return ""
}
