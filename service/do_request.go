package service

import (
	"errors"
	"net/http"
)

func DoRequest(req *http.Request) (map[string]interface{}, int, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, res.StatusCode, err
	}

	data, err := CheckResponse(res)
	if err != nil {
		return nil, res.StatusCode, err
	}

	if _, ok := data["error"]; ok {
		return data, res.StatusCode, errors.New("Error")
	}

	return data, res.StatusCode, nil
}
