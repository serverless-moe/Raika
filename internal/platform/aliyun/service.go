// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package aliyun

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

var ErrRaikaServiceNotFound = errors.New("Raika service not found")

type Services []*Service
type Service struct {
	ID        string `json:"serviceId"`
	Name      string `json:"serviceName"`
	Desc      string `json:"description"`
	CreatedAt string `json:"createdTime"`
}

type ListServicesResponse struct {
	Services  Services `json:"services"`
	NextToken string   `json:"nextToken"`
}

func (c *Client) ListServices() (Services, error) {
	var response ListServicesResponse
	resp, err := c.request(http.MethodGet, "/services")
	if err != nil {
		return nil, err
	}

	if err := resp.ToJSON(&response); err != nil {
		return nil, errors.Wrap(err, "JSON decode")
	}
	return response.Services, nil
}

func (c *Client) CreateService(name, desc string) (*Service, error) {
	var response Service
	resp, err := c.request(http.MethodPost, "/services", map[string]interface{}{
		"serviceName": name,
		"description": desc,
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.ToString())
	}
	return &response, resp.ToJSON(&response)
}

func (c *Client) GetRaikaService() (*Service, error) {
	services, err := c.ListServices()
	if err != nil {
		return nil, errors.Wrap(err, "list services")
	}

	for _, s := range services {
		if strings.HasPrefix(s.Name, "Raika-") {
			return s, nil
		}
	}
	return nil, ErrRaikaServiceNotFound
}
