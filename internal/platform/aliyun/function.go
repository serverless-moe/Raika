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
	Handler               string `json:"handler"`
	Runtime               string `json:"runtime"`
	MemorySize            int64  `json:"memorySize"`
	InitializationTimeout int    `json:"initializationTimeout"`
	Timeout               int    `json:"timeout"`
	CAPort                int    `json:"caPort"`
}

func (c *Client) CreateFunction(opts platform.CreateFunctionOptions) (string, error) {
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
	}

	_, err = c.request(http.MethodPost, fmt.Sprintf("/services/%s/functions", ServiceName), requestBody)
	if err != nil {
		return "", errors.Wrap(err, "create function")
	}

	// Create HTTP trigger for function.
	err = c.CreateHTTPTrigger(CreateHTTPTriggerOptions{
		TriggerName:  platform.TriggerName,
		ServiceName:  ServiceName,
		FunctionName: opts.Name,
	})
	if err != nil {
		return "", errors.Wrap(err, "create HTTP trigger")
	}

	return fmt.Sprintf("https://%s.%s.fc.aliyuncs.com/2016-08-15/proxy/%s/%s/", c.accountID, c.regionID, ServiceName, opts.Name), nil
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
