package loader

import (
	"io"
	"strings"
)

func ReadRemoteFile(remoteLoc string) ([]byte, error) {
	r, err := GetRemoteReader(remoteLoc)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

func GetRemoteReader(remoteLoc string) (io.ReadCloser, error) {
	if strings.ToLower(remoteLoc[0:4]) == "http" {
		return ReaderFromHttp(remoteLoc)
	}
	if strings.ToLower(remoteLoc[0:5]) == "s3://" {
		return ReaderFromS3(remoteLoc)
	}
	// if no remote handlers can handle the reading of the file, lets try local
	return ReaderFromLocal(remoteLoc)
}
