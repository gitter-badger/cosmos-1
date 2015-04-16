package util

import (
	"io/ioutil"
	"net/http"
)

func GetToken(req *http.Request) string {
	token := req.URL.Query().Get("token")
	if token == "" {
		token = "default"
	}
	return token
}

func GetBodyFromRequest(req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
