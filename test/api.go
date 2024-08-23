package test

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

func NewFormRequest(method string, target string, data map[string]any) (*http.Request, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	req := httptest.NewRequest(method, target, formFromMap(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func JsonToMap(data []byte) (map[string]any, error) {
	var res map[string]any
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func formFromMap(data map[string]any) io.Reader {
	form := url.Values{}

	for key, val := range maps.All(data) {
		form.Add(key, fmt.Sprintf("%v", val))
	}

	return strings.NewReader(form.Encode())
}
