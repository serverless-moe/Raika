// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

var host = "http://127.0.0.1:3000"

func request(method, baseURL string, requestBody ...interface{}) (*response, error) {
	var body io.Reader
	if len(requestBody) == 1 {
		repBody, err := json.Marshal(requestBody[0])
		if err != nil {
			return nil, errors.Wrap(err, "JSON encode")
		}
		body = bytes.NewReader(repBody)
	}

	req, err := http.NewRequest(method, host+baseURL, body)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}

	return &response{
		Response: resp,
	}, nil
}

type response struct {
	*http.Response
}

func (r *response) ToJSON(v interface{}) error {
	defer func() { _ = r.Body.Close() }()
	return json.NewDecoder(r.Body).Decode(v)
}

func (r *response) ToString() string {
	defer func() { _ = r.Body.Close() }()
	resp, _ := io.ReadAll(r.Body)
	return string(resp)
}
