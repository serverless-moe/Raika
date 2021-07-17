// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tencentcloud

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/wuhan005/Raika/internal/platform"
)

var _ platform.Cloud = (*Client)(nil)

type Client struct {
	regionID            string
	secretID, secretKey string
}

func New(opts platform.AuthenticateOptions) *Client {
	return &Client{
		regionID:  opts[RegionIDField],
		secretID:  opts[SecretIDField],
		secretKey: opts[SecretKeyField],
	}
}

func (c *Client) String() string {
	return "tencentcloud"
}

func (c *Client) Authenticate() error {
	resp, err := c.request(http.MethodGet, "GetAccount")
	if err != nil {
		return err
	}

	var respJSON struct {
		Response struct {
			Error struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			} `json:"Error"`
		} `json:"Response"`
	}
	if err := resp.ToJSON(&respJSON); err != nil {
		return err
	}

	if respJSON.Response.Error.Code != "" {
		return errors.Errorf("[%s] %s", respJSON.Response.Error.Code, respJSON.Response.Error.Message)
	}
	return nil
}

func (c *Client) request(method, action string, requestBody ...interface{}) (*response, error) {
	u := "https://scf.tencentcloudapi.com/"

	var err error
	var body io.Reader
	var reqBody []byte
	var query url.Values
	if len(requestBody) == 1 {
		if method == http.MethodPost {
			reqBody, err = json.Marshal(requestBody[0])
			if err != nil {
				return nil, errors.Wrap(err, "JSON encode")
			}
			body = bytes.NewReader(reqBody)
		} else if method == http.MethodGet {
			if v, ok := requestBody[0].(url.Values); ok {
				query = v
			}
		} else {
			body = bytes.NewReader(nil)
		}
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	req.Header.Set("x-tc-action", action)
	req.Header.Set("x-tc-region", c.regionID)
	req.Header.Set("x-tc-version", "2018-04-16")
	if req.Method == http.MethodGet {
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("content-type", "application/json")
	}
	req.Header.Set("host", "scf.tencentcloudapi.com")
	req.Header.Set("Authorization", c.GetAuthorizationHeader(req, reqBody))

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
