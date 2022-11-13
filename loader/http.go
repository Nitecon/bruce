package loader

import (
	"io"
	"net/http"
)

// TODO: convert to interfaces some time soon.

func ReadFromHttp(fileName string) ([]byte, error) {
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

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		return nil, err
	}

	return b, nil
}
