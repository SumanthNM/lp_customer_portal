package common

import (
	"net/http"

	"github.com/go-chassis/foundation/httputil"
	"github.com/go-chassis/openlog"
)

func GETReqeust(url string, queryParms, headers map[string]string) (string, error) {
	openlog.Debug("Making GET request to: " + url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		openlog.Error("Error making GET request to: " + url)
		return "", err
	}
	q := req.URL.Query()
	for key, value := range queryParms {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		openlog.Error("Error making GET request to: " + url)
		openlog.Error(err.Error())
		body := httputil.ReadBody(resp)
		openlog.Error(string(body))
		return "", err
	}
	defer resp.Body.Close()
	body := httputil.ReadBody(resp)
	return string(body), nil
}
