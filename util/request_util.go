package util

import (
	"io/ioutil"
	"net/http"
)

func GetQueryParam(req *http.Request, key, defaultVal string) string {
	param := req.URL.Query().Get(key)
	if param == "" {
		return defaultVal
	}
	return param
}

func GetBodyFromRequest(req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
