package main

import (
	"io/ioutil"
	"net/http"
)

func GetBodyFromRequest(req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
