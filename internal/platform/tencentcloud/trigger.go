// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tencentcloud

import (
	"net/http"

	"github.com/pkg/errors"
)

type CreateHTTPTriggerOptions struct {
	TriggerName  string
	FunctionName string
}

type CreateHTTPTriggerRequest struct {
	FunctionName string `json:"FunctionName"`
	TriggerName  string `json:"TriggerName"`
	Type         string `json:"Type"`
	TriggerDesc  string `json:"TriggerDesc"`
}

type CreateTriggerResponse struct {
	Response struct {
		RequestId   string `json:"RequestId"`
		TriggerInfo struct {
			AddTime          string `json:"AddTime"`
			AvailableStatus  string `json:"AvailableStatus"`
			BindStatus       string `json:"BindStatus"`
			CustomArgument   string `json:"CustomArgument"`
			Enable           int    `json:"Enable"`
			ModTime          string `json:"ModTime"`
			ResourceId       string `json:"ResourceId"`
			TriggerAttribute string `json:"TriggerAttribute"`
			TriggerDesc      string `json:"TriggerDesc"`
			TriggerName      string `json:"TriggerName"`
			Type             string `json:"Type"`
		} `json:"TriggerInfo"`
		Error struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`
	} `json:"Response"`
}

func (c *Client) CreateHTTPTrigger(opts CreateHTTPTriggerOptions) (*CreateTriggerResponse, error) {
	requestBody := CreateHTTPTriggerRequest{
		FunctionName: opts.FunctionName,
		TriggerName:  opts.TriggerName,
		Type:         "apigw",
		TriggerDesc: `{
    "api": {
        "authRequired": "FALSE",
        "requestConfig": {
            "method": "ANY"
        },
        "isIntegratedResponse": "FALSE"
    },
    "service": {
        "serviceName": "SCF_API_SERVICE"
    },
    "release": {
        "environmentName": "release"
    }
}`,
	}
	resp, err := c.request(http.MethodPost, "CreateTrigger", requestBody)
	if err != nil {
		return nil, err
	}

	var respJSON CreateTriggerResponse
	if err := resp.ToJSON(&respJSON); err != nil {
		return nil, errors.Wrap(err, "json decode")
	}
	if respJSON.Response.Error.Code != "" {
		return nil, errors.Errorf("%s: %s", respJSON.Response.Error.Code, respJSON.Response.Error.Message)
	}
	return &respJSON, nil
}
