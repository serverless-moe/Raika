// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package aliyun

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

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

type HTTPTriggerResponse struct {
}

func (c *Client) GetHTTPTrigger(serviceName, functionName, triggerName string) (*HTTPTriggerResponse, error) {
	resp, err := c.request(http.MethodGet, fmt.Sprintf("/services/%s/functions/%s/triggers/%s", serviceName, functionName, triggerName))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.ToString())
	}
	return nil, nil
}

type CreateCronTriggerOptions struct {
	TriggerName  string
	ServiceName  string
	FunctionName string
	CronString   string
}

type CreateCronTriggerRequest struct {
	Name   string `json:"triggerName"`
	Type   string `json:"triggerType"`
	Config struct {
		AuthType       string `json:"authType"`
		CronExpression string `json:"cronExpression"`
		Enable         bool   `json:"enable"`
	} `json:"triggerConfig"`
	InvocationRole string `json:"invocationRole"`
	Qualifier      string `json:"qualifier"`
	SourceArn      string `json:"sourceArn"`
}

func (c *Client) CreateCronTrigger(opts CreateCronTriggerOptions) error {
	requestBody := CreateCronTriggerRequest{
		Name: opts.TriggerName,
		Type: "timer",
		Config: struct {
			AuthType       string `json:"authType"`
			CronExpression string `json:"cronExpression"`
			Enable         bool   `json:"enable"`
		}{
			AuthType:       "anonymous",
			CronExpression: opts.CronString,
			Enable:         true,
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

type ListTriggersResponse struct {
	Triggers []struct {
		TriggerName    string      `json:"triggerName"`
		Description    string      `json:"description"`
		TriggerId      string      `json:"triggerId"`
		SourceArn      interface{} `json:"sourceArn"`
		TriggerType    string      `json:"triggerType"`
		InvocationRole string      `json:"invocationRole"`
		Qualifier      string      `json:"qualifier"`
		TriggerConfig  struct {
			Methods  []string `json:"methods"`
			AuthType string   `json:"authType"`
		} `json:"triggerConfig"`
		CreatedTime      time.Time `json:"createdTime"`
		LastModifiedTime time.Time `json:"lastModifiedTime"`
	} `json:"triggers"`
}

func (c *Client) ListTriggers(serviceName, functionName string) (*ListTriggersResponse, error) {
	// TODO support nextToken
	resp, err := c.request(http.MethodGet, fmt.Sprintf("/services/%s/functions/%s/triggers?limit=100", serviceName, functionName))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.ToString())
	}

	var respJSON ListTriggersResponse
	return &respJSON, resp.ToJSON(&respJSON)
}

func (c *Client) DeleteTrigger(serviceName, functionName, triggerName string) error {
	resp, err := c.request(http.MethodDelete, fmt.Sprintf("/services/%s/functions/%s/triggers/%s", serviceName, functionName, triggerName))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New(resp.ToString())
	}
	return nil
}
