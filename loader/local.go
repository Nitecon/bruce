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

func WriterFromLocal(fileName string) (io.WriteCloser, error) {
	w, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		log.Error().Err(err).Msgf("could not open file for writing: %s", fileName)
		return nil, err
	}
	return w, nil
}
