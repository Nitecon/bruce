package loader

import (
	"io"
	"net/http"
)

func ReaderFromHttp(fileName string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", fileName, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
