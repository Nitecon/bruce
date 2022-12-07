package loader

import (
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

func ReaderFromLocal(fileName string) (io.ReadCloser, error) {
	log.Debug().Msgf("starting local read of %s", fileName)
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		log.Info().Msgf("local reader engine: (file does not exist): %s", fileName)
		return nil, err
	}
	return os.Open(fileName)
}
