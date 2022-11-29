package operators

import (
	"bruce/exe"
	"github.com/rs/zerolog/log"
)

// OwnerShip provides a means to set the ownership of files or directories as needed.
type OwnerShip struct {
	Name      string `yaml:"name"`
	ObType    string `yaml:"type"`
	Path      string `yaml:"path"`
	Owner     string `yaml:"owner"`
	Group     string `yaml:"group"`
	Recursive bool   `yaml:"recursive"`
}

func (o *OwnerShip) Execute() error {
	log.Debug().Msgf("starting chown on: %s", o.Path)
	err := exe.SetOwnership(o.ObType, o.Path, o.Owner, o.Group, o.Recursive)
	if err != nil {
		log.Error().Err(err).Msg("could not set ownership")
		return err
	}
	return nil
}
