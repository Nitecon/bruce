package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

// BruceConfig will be marshalled from the provided config file that exists.
type BruceConfig struct {
	InstallPackages []string         `yaml:"packageList"`
	Templates       []ActionTemplate `yaml:"templates"`
	Services        []Services       `yaml:"services"`
}

// ActionTemplate provides the local and remote files to be used.
type ActionTemplate struct {
	LocalLocation  string `yaml:"localLocation"`
	RemoteLocation string `yaml:"remoteLocation"`
	Variables      Vars   `yaml:"vars"`
}

// Vars indicates the variables to replace in the template and how to replace them.
type Vars struct {
	Type     string `yaml:"type"`
	Action   string `yaml:"action"`
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
func LoadConfig(fileName string) (*BruceConfig, error) {
	if fileName == "" {
		return nil, fmt.Errorf("config file cannot be empty")
	}
	b, err := os.ReadFile(fileName) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	c := &BruceConfig{}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse config file")
	}
	return c, nil
}
