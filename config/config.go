package config

import (
	"bruce/loader"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"os/user"
	"runtime"
	"strings"
	"sync"
)

var (
	conf     *AppData
	confLock = new(sync.RWMutex)
)

type AppData struct {
	CurrentUser           string
	TrySudo               bool
	PackageHandler        string
	PackageHandlerPath    string
	PackageManagerUpdated bool
	SystemType            string
	Template              *TemplateData
	ServiceController     string
	ServiceControllerPath string
	CanUpdateServices     bool
	ChangedTemplates      []string
}

// TemplateData will be marshalled from the provided config file that exists.
type TemplateData struct {
	TempDir         string           `yaml:"tempDirectory"`
	PreExecCmds     []string         `yaml:"preExecCmds"`
	PostInstallCmds []string         `yaml:"postInstallerCmds"`
	InstallPackages []string         `yaml:"packageList"`
	Templates       []ActionTemplate `yaml:"templates"`
	OwnerShips      []OwnerShip      `yaml:"chowns"`
	Services        []Services       `yaml:"services"`
	PostExecCmds    []string         `yaml:"postExecCmds"`
	BackupDir       string
}

// OwnerShip provides a means to set the ownership of files or directories as needed.
type OwnerShip struct {
	Object    string `yaml:"type"`
	Path      string `yaml:"path"`
	Owner     string `yaml:"owner"`
	Group     string `yaml:"group"`
	Recursive bool   `yaml:"recursive"`
}

// ActionTemplate provides the local and remote files to be used.
type ActionTemplate struct {
	LocalLocation  string      `yaml:"localLocation"`
	RemoteLocation string      `yaml:"remoteLocation"`
	Variables      []Vars      `yaml:"vars"`
	Permissions    fs.FileMode `yaml:"perms"`
	Owner          string      `yaml:"owner"`
	Group          string      `yaml:"group"`
}

// Vars indicates the variables to replace in the template and how to replace them.
type Vars struct {
	Type     string `yaml:"type"`
	Output   string `yaml:"output"`
	Variable string `yaml:"variable"`
}

// Services are the list of services that must be set up as required.
type Services struct {
	Name            string   `yaml:"name"`
	Enabled         bool     `yaml:"setEnabled"`
	State           string   `yaml:"state"`
	RestartOnUpdate []string `yaml:"restartTrigger"`
	RestartAlways   bool     `yaml:"restartAlways"`
}

// LoadConfig attempts to load the user provided manifest.
func LoadConfig(fileName string) (*TemplateData, error) {
	ad := InitAppData()
	d, err := loader.ReadRemoteFile(fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read config file")
	}
	log.Debug().Bytes("rawConfig", d)
	c := &TemplateData{}

	err = yaml.Unmarshal(d, c)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse config file")
	}
	log.Debug().Interface("config", c)
	// setup some defaults
	for _, temps := range c.Templates {
		if temps.Permissions == 0 {
			temps.Permissions = 0664
		}
	}
	c.TempDir = fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, "bruce")
	ad.Template = c
	ad.Save()
	return c, nil
}

func InitAppData() *AppData {
	cfg := &AppData{}
	if runtime.GOOS == "linux" {
		u, err := user.Current()
		if err != nil {
			log.Fatal().Err(err).Msg("user should exist to operate")
		}
		if u.Username != "root" {
			cfg.TrySudo = true
			log.Debug().Msgf("attempting sudo, current user: %s", cfg.CurrentUser)
		}

		cfg.SystemType = "linux"
	}
	cfg.Save()
	return cfg
}

// Get function returns the currently set global system information to be used.
func Get() *AppData {
	confLock.RLock()
	defer confLock.RUnlock()
	return conf

}

// Save saves.
func (s *AppData) Save() {
	confLock.Lock()
	defer confLock.Unlock()
	conf = s
}

func GetValueForOSHandler(value string) string {
	log.Debug().Msgf("OS Handler value iteration: %#v", value)
	if Get().PackageHandler == "" {
		log.Error().Err(fmt.Errorf("cannot retrieve os handler value without a known package handler"))
		return ""
	}
	log.Debug().Msgf("testing for my package handler: %s", Get().PackageHandler)
	if strings.Contains(value, "|") {
		managerList := strings.Split(value, "|")
		var basePackage = ""
		var usablePackage = ""
		for _, mpkg := range managerList {
			log.Debug().Msgf("os handler iteration for manager: %#v", mpkg)
			if strings.Contains(mpkg, "=") {
				pmSplit := strings.Split(mpkg, "=")
				log.Debug().Msgf("handler [%s] specific value: %s", pmSplit[0], pmSplit[1])
				if pmSplit[0] == Get().PackageHandler {
					usablePackage = pmSplit[1]
				}
			} else {
				basePackage = mpkg
			}
		}
		if usablePackage != "" {
			log.Debug().Msgf("returning package manager value: %s", usablePackage)
			return usablePackage
		}
		log.Debug().Msgf("returning base value: %s", basePackage)
		return basePackage
	}
	log.Debug().Msgf("returning original value: %s", value)
	return value
}
