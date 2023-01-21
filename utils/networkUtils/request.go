package networkUtils

import (
	"bytes"
	"github.com/goccy/go-json"
	"io/ioutil"
	"net/http"
)

func Post(url string, data map[string]any, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for key, header := range headers {
			req.Header.Set(key, header)
		}
	}
	resp, _err := client.Do(req)
	if _err != nil {
		return nil, _err
	}
	content, _ := ioutil.ReadAll(resp.Body)
	return content, nil
}

func Get(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for key, header := range headers {
			req.Header.Set(key, header)
		}
	}
	resp, _err := client.Do(req)
	if _err != nil {
		return nil, _err
	}
	content, _ := ioutil.ReadAll(resp.Body)
	return content, nil
}
