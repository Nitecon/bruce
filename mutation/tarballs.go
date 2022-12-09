package mutation

import (
	"archive/tar"
	"bruce/loader"
	"compress/gzip"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func useGzipReader(filename string, fileReader io.ReadCloser) io.ReadCloser {
	if strings.HasSuffix(filename, ".tgz") || strings.HasSuffix(filename, ".tar.gz") {
		gzr, err := gzip.NewReader(fileReader)
		if err != nil {
			log.Error().Err(err).Msg("could not instantiate gzip reader returning original")
			return fileReader
		}
		return gzr
	}
	return fileReader
}

func ExtractTarball(src, dst string, force bool) error {
	// We just check dest currently as we will read from multiple source locations and they may fail by time we cleaned up so worthless to check upfront.
	if _, err := os.Stat(dst); err == nil {
		if !force {
			log.Info().Msgf("%s already exists cannot extract tarball to location", dst)
			return nil
		}
		// Destination exists and we are forcing new extraction so...
		err = os.RemoveAll(dst)
		if err != nil {
			return fmt.Errorf("failed to remove destination prior to extraction: %s", dst)
		}
	}
	rr, err := loader.GetRemoteReader(src)
	if err != nil {
		log.Error().Err(err).Msgf("cannot read tarball at src: %s", src)
		return err
	}
	tr := tar.NewReader(useGzipReader(src, rr))
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
					return err
				}
			}
		// create file with existing file mode
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			// save contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			f.Close()
		}
	}
}
