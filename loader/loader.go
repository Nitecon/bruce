package loader

import "strings"

func ReadRemoteFile(remoteLoc string) ([]byte, error) {
	if strings.ToLower(remoteLoc[0:4]) == "http" {
		return ReadFromHttp(remoteLoc)
	}
	if strings.ToLower(remoteLoc[0:5]) == "s3://" {
		return ReadFromS3(remoteLoc)
	}
	// if no remote handlers can handle the reading of the file, lets try local
	return ReadFromLocal(remoteLoc)
}
