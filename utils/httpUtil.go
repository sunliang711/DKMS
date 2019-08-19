package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Post TODO
// 2019/08/16 10:45:16
func Post(url string, headers map[string]string, body interface{}) (*http.Response, error) {
	bs, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", url, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}

	client := &http.Client{}
	return client.Do(request)
}
