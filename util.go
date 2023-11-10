package credhub

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func closeResponse(response *http.Response) {
	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}
}

func getJson[T any](client HttpClient, url string) (*T, error) {
	request, e := http.NewRequest(http.MethodGet, url, nil)
	if e != nil {
		return nil, e
	}
	request.Header.Set("Accept", "application/json")

	response, e := client.Do(request)
	defer closeResponse(response)
	if e != nil {
		return nil, e
	}

	result := new(T)
	if e = json.NewDecoder(response.Body).Decode(result); e != nil {
		return nil, e
	}

	return result, nil
}

func putAndGetJsons[T any](client HttpClient, url string, content any) (*T, error) {

	body, e := json.Marshal(content)
	if e != nil {
		return nil, e
	}

	request, e := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if e != nil {
		return nil, e
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, e := client.Do(request)
	defer closeResponse(response)
	if e != nil {
		return nil, e
	}

	result := new(T)
	if e = json.NewDecoder(response.Body).Decode(result); e != nil {
		return nil, e
	}

	return result, nil
}
