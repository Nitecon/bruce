package loader

import (
	"github.com/rs/zerolog/log"
	"os"
)

func ReadFromLocal(fileName string) ([]byte, error) {
	log.Debug().Msgf("starting local read of %s", fileName)
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		log.Info().Msgf("local reader engine: (file does not exist): %s", fileName)
		return nil, err
	}
	return os.ReadFile(fileName)
}
