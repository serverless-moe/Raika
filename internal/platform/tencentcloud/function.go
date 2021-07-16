// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tencentcloud

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/platform"
)

type CreateFunctionRequest struct {
	Name        string `json:"FunctionName"`
	Description string `json:"Description"`
	Code        struct {
		ZipFile []byte `json:"ZipFile"`
	} `json:"Code"`
	Runtime         string                 `json:"Runtime"`
	MemorySize      int64                  `json:"MemorySize"`
	Environment     map[string]string      `json:"Environment,omitempty"`
	InitTimeout     int                    `json:"InitTimeout"`
	Timeout         int                    `json:"Timeout"`
	Type            string                 `json:"Type"`
	PublicNetConfig map[string]interface{} `json:"PublicNetConfig"`
}

func (c *Client) CreateFunction(opts platform.CreateFunctionOptions) (string, error) {
	zipFile, err := packFile(opts.File)
	if err != nil {
		return "", errors.Wrap(err, "pack file")
	}

	_, err = c.request(http.MethodPost, "CreateFunction", CreateFunctionRequest{
		Name:        opts.Name,
		Description: opts.Description,
		Code: struct {
			ZipFile []byte `json:"ZipFile"`
		}{ZipFile: zipFile},
		Runtime:     "Go1",
		MemorySize:  opts.MemorySize,
		Environment: opts.Environment,
		InitTimeout: int(opts.InitializationTimeout / time.Second),
		Timeout:     int(opts.RuntimeTimeout / time.Second),
		Type:        "HTTP",
		PublicNetConfig: map[string]interface{}{
			"PublicNetStatus": "ENABLE",
			"EipConfig": map[string]interface{}{
				"EipStatus": "DISABLE",
			},
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "create function")
	}

	// Create HTTP trigger for function.
	resp, err := c.CreateHTTPTrigger(CreateHTTPTriggerOptions{
		TriggerName:  platform.TriggerName,
		FunctionName: opts.Name,
	})
	if err != nil {
		return "", errors.Wrap(err, "create HTTP trigger")
	}

	log.Trace("%+v", resp.Response)
	return resp.Response.TriggerInfo.ResourceId, nil
}

func packFile(path string) ([]byte, error) {
	output := new(bytes.Buffer)
	zipWriter := zip.NewWriter(output)

	bootstrapHeader := &zip.FileHeader{
		Name:     "scf_bootstrap",
		Modified: time.Now(),
	}
	bootstrapHeader.SetMode(0777)
	zipEntry, err := zipWriter.CreateHeader(bootstrapHeader)
	if err != nil {
		return nil, errors.Wrap(err, "create bootstrap file header")
	}

	_, err = zipEntry.Write([]byte("#!/bin/bash\n./bootstrap"))
	if err != nil {
		return nil, errors.Wrap(err, "write bootstrap file")
	}

	binaryHeader := &zip.FileHeader{
		Name:     "bootstrap",
		Modified: time.Now(),
	}
	binaryHeader.SetMode(0777)
	zipEntry, err = zipWriter.CreateHeader(binaryHeader)
	if err != nil {
		return nil, errors.Wrap(err, "create binary file header")
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
