package system

import (
	"bruce/loader"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
)

// BruceConfig will be marshalled from the provided config file that exists.
type BruceConfig struct {
	TempDir         string           `yaml:"tempDirectory"`
	PreExecCmds     []string         `yaml:"preExecCmds"`
	InstallPackages []string         `yaml:"packageList"`
	Templates       []ActionTemplate `yaml:"templates"`
	Services        []Services       `yaml:"services"`
	PostExecCmds    []string         `yaml:"postExecCmds"`
	BackupDir       string
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
	Name             string   `yaml:"name"`
	Enabled          bool     `yaml:"setEnabled"`
	State            string   `yaml:"state"`
	RestartOnUpdate  []string `yaml:"restartTrigger"`
	OnFailRevertTPLs []string `yaml:"failRevertTemplate"`
	RestartAlways    bool     `yaml:"restartAlways"`
}

// LoadConfig attempts to load the user provided manifest.
func LoadConfig(fileName string) (*BruceConfig, error) {
	d, err := loader.ReadRemoteFile(fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read config file")
	}
	log.Debug().Bytes("rawConfig", d)
	c := &BruceConfig{}

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
	info := GetSysInfo()
	info.Configuration = c
	SetSysInfo(info)
	return c, nil
}
