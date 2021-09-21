// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package aliyun

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/platform"
)

type CreateFunctionRequest struct {
	Name        string `json:"functionName"`
	Description string `json:"description"`
	Code        struct {
		ZipBase64 []byte `json:"zipFile"`
	} `json:"code"`
	Handler               string            `json:"handler"`
	Runtime               string            `json:"runtime"`
	MemorySize            int64             `json:"memorySize"`
	InitializationTimeout int               `json:"initializationTimeout"`
	Timeout               int               `json:"timeout"`
	CAPort                int               `json:"caPort"`
	EnvironmentVariables  map[string]string `json:"environmentVariables,omitempty"`
}

func (c *Client) CreateFunction(opts platform.CreateFunctionOptions) (string, error) {
	// Create service if not exists.
	_, err := c.GetRaikaService()
	if err == ErrRaikaServiceNotFound {
		log.Trace("Raika service not found on aliyun, create...")
		_, err = c.CreateService(ServiceName, "Service for Raika.")
		if err != nil {
			return "", errors.Wrap(err, "create service")
		}
	} else if err != nil {
		return "", err
	}

	if opts.MemorySize < 128 || opts.MemorySize > 3072 || opts.MemorySize%64 != 0 {
		return "", errors.Errorf("wrong memory size: %d", opts.MemorySize)
	}

	zipFile, err := packFile(opts.File)
	if err != nil {
		return "", errors.Wrap(err, "pack file")
	}

	// Check current function name exists.
	_, err = c.GetFunction(ServiceName, opts.Name)
	if err != nil && err != ErrFunctionNotExists {
		return "", errors.Wrap(err, "get function")
	} else if err == nil {
		// Function exists, delete it.
		log.Trace("Function %q exists on aliyun, replace...", opts.Name)

		// Delete the triggers under the function.
		triggers, err := c.ListTriggers(ServiceName, opts.Name)
		if err != nil {
			return "", errors.Wrap(err, "list triggers")
		}
		for _, trigger := range triggers.Triggers {
			log.Trace("Delete trigger: %q...", trigger.TriggerName)
			if err := c.DeleteTrigger(ServiceName, opts.Name, trigger.TriggerName); err != nil {
				return "", errors.Wrapf(err, "delete trigger: %q", opts.Name)
			}
		}

		log.Trace("Delete function: %q...", opts.Name)
		if err := c.DeleteFunction(ServiceName, opts.Name); err != nil && err != ErrFunctionNotExists {
			return "", errors.Wrap(err, "delete function")
		}
	}

	requestBody := CreateFunctionRequest{
		Name:        opts.Name,
		Description: opts.Description,
		Code: struct {
			ZipBase64 []byte `json:"zipFile"`
		}{
			ZipBase64: zipFile,
		},
		Handler:               "index.handler",
		Runtime:               "custom",
		MemorySize:            opts.MemorySize,
		InitializationTimeout: int(opts.InitializationTimeout / time.Second),
		Timeout:               int(opts.RuntimeTimeout / time.Second),
		CAPort:                opts.HTTPPort,
		EnvironmentVariables:  opts.EnvironmentVariables,
	}

	log.Trace("Deploy function: %q...", opts.Name)
	_, err = c.request(http.MethodPost, fmt.Sprintf("/services/%s/functions", ServiceName), requestBody)
	if err != nil {
		return "", errors.Wrap(err, "create function")
	}

	if opts.TriggerType == "http" {
		// Create HTTP trigger for function.
		err = c.CreateHTTPTrigger(CreateHTTPTriggerOptions{
			TriggerName:  platform.HTTPTriggerName,
			ServiceName:  ServiceName,
			FunctionName: opts.Name,
		})
		if err != nil {
			return "", errors.Wrap(err, "create HTTP trigger")
		}

		return fmt.Sprintf("https://%s.%s.fc.aliyuncs.com/2016-08-15/proxy/%s/%s/", c.accountID, c.regionID, ServiceName, opts.Name), nil
	} else if opts.TriggerType == "cron" {
		err = c.CreateCronTrigger(CreateCronTriggerOptions{
			TriggerName:  platform.CronTriggerName,
			ServiceName:  ServiceName,
			FunctionName: opts.Name,
			CronString:   opts.CronString,
		})
		if err != nil {
			return "", errors.Wrap(err, "create timer trigger")
		}
	} else {
		return "", errors.Errorf("unexpected trigger type %q", opts.TriggerType)
	}
	return "", nil
}

func packFile(path string) ([]byte, error) {
	output := new(bytes.Buffer)
	zipWriter := zip.NewWriter(output)

	header := &zip.FileHeader{
		Name:     "bootstrap",
		Modified: time.Now(),
	}
	header.SetMode(0777)
	zipEntry, err := zipWriter.CreateHeader(header)
	if err != nil {
		return nil, errors.Wrap(err, "create header")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	if _, err := io.Copy(zipEntry, file); err != nil {
		return nil, errors.Wrap(err, "copy")
	}
	_ = zipWriter.Close()

	return output.Bytes(), nil
}

type GetFunctionResponse struct {
	CodeChecksum          string    `json:"codeChecksum"`
	CodeSize              int       `json:"codeSize"`
	CreatedTime           time.Time `json:"createdTime"`
	Description           string    `json:"description"`
	FunctionId            string    `json:"functionId"`
	FunctionName          string    `json:"functionName"`
	Handler               string    `json:"handler"`
	MemorySize            int       `json:"memorySize"`
	Runtime               string    `json:"runtime"`
	Timeout               int       `json:"timeout"`
	InitializationTimeout int       `json:"initializationTimeout"`
	Initializer           string    `json:"initializer"`
	CaPort                int       `json:"caPort"`
	CustomContainerConfig struct {
		Args             string `json:"args"`
		Command          string `json:"command"`
		Image            string `json:"image"`
		AccelerationType string `json:"accelerationType"`
		AccelerationInfo struct {
			Status string `json:"status"`
		} `json:"accelerationInfo"`
	} `json:"customContainerConfig"`
	Layers []string `json:"layers"`
}

var ErrFunctionNotExists = errors.New("function not found")

func (c *Client) GetFunction(serviceName, functionName string) (*GetFunctionResponse, error) {
	resp, err := c.request(http.MethodGet, fmt.Sprintf("/services/%s/functions/%s", serviceName, functionName))
	if err != nil {
		return nil, errors.Wrap(err, "create function")
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrFunctionNotExists
		}
		return nil, errors.Errorf("unexpected status code %d: %v", resp.StatusCode, resp.ToString())
	}

	var respJSON GetFunctionResponse
	return &respJSON, resp.ToJSON(&respJSON)
}

func (c *Client) DeleteFunction(serviceName, functionName string) error {
	resp, err := c.request(http.MethodDelete, fmt.Sprintf("/services/%s/functions/%s", serviceName, functionName))
	if err != nil {
		return errors.Wrap(err, "delete function")
	}

	if resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusNotFound {
			return ErrFunctionNotExists
		}
		return errors.Errorf("unexpected status code %d: %v", resp.StatusCode, resp.ToString())
	}
	return nil
}
