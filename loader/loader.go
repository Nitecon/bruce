package loader

import (
	"strings"
)

func ReadRemoteFile(remoteLoc string) ([]byte, error) {
	if strings.HasPrefix(remoteLoc, "http") {

	}
	// if no remote handlers can handle the reading of the file, lets try local
	return ReadFromLocal(remoteLoc)
}
