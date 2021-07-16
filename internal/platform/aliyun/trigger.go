// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package aliyun

import (
	"errors"
	"fmt"
	"net/http"
)

const TriggerName = "Raika_HTTPTrigger"

type CreateHTTPTriggerOptions struct {
	TriggerName  string
	ServiceName  string
	FunctionName string
}

type CreateHTTPTriggerRequest struct {
	Name   string `json:"triggerName"`
	Type   string `json:"triggerType"`
	Config struct {
		AuthType string `json:"authType"`
	} `json:"triggerConfig"`
	InvocationRole string `json:"invocationRole"`
	Qualifier      string `json:"qualifier"`
	SourceArn      string `json:"sourceArn"`
}

func (c *Client) CreateHTTPTrigger(opts CreateHTTPTriggerOptions) error {
	requestBody := CreateHTTPTriggerRequest{
		Name: opts.TriggerName,
		Type: "http",
		Config: struct {
			AuthType string `json:"authType"`
		}{
			AuthType: "anonymous",
		},
		InvocationRole: fmt.Sprintf("acs:ram::%s:role/aliyunfcdefaultrole", c.accountID),
		Qualifier:      "LATEST",
		SourceArn:      "anonymous",
	}
	resp, err := c.request(http.MethodPost, fmt.Sprintf("/services/%s/functions/%s/triggers", opts.ServiceName, opts.FunctionName), requestBody)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.ToString())
	}
	return nil
}

type GetHTTPTriggerOptions struct {
	TriggerName  string
	ServiceName  string
	FunctionName string
}

type HTTPTriggerResponse struct {
}

func (c *Client) GetHTTPTrigger(opts GetHTTPTriggerOptions) (*HTTPTriggerResponse, error) {
	resp, err := c.request(http.MethodGet, fmt.Sprintf("/services/%s/functions/%s/triggers/%s", opts.ServiceName, opts.FunctionName, opts.TriggerName))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.ToString())
	}
	return nil, nil
}
