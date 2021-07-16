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

	var functionStatus string
	for functionStatus != "Active" {
		time.Sleep(2 * time.Second)
		functionInfo, err := c.GetFunction(opts.Name)
		if err != nil {
			return "", errors.Wrap(err, "get function")
		}
		functionStatus = functionInfo.Response.Status
	}

	// Create HTTP trigger for function.
	resp, err := c.CreateHTTPTrigger(CreateHTTPTriggerOptions{
		TriggerName:  platform.TriggerName,
		FunctionName: opts.Name,
	})
	if err != nil {
		return "", errors.Wrap(err, "create HTTP trigger")
	}
	return resp.Service.SubDomain, nil
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

type GetFunctionRequest struct {
	FunctionName string `json:"FunctionName"`
}

type GetFunctionResponse struct {
	Response struct {
		Qualifier   string `json:"Qualifier"`
		Description string `json:"Description"`
		Timeout     int    `json:"Timeout"`
		InitTimeout int    `json:"InitTimeout"`
		MemorySize  int    `json:"MemorySize"`
		Runtime     string `json:"Runtime"`
		VpcConfig   struct {
			VpcId    string `json:"VpcId"`
			SubnetId string `json:"SubnetId"`
		} `json:"VpcConfig"`
		Environment struct {
			Variables []interface{} `json:"Variables"`
		} `json:"Environment"`
		Handler           string `json:"Handler"`
		UseGpu            string `json:"UseGpu"`
		Role              string `json:"Role"`
		CodeSize          int    `json:"CodeSize"`
		FunctionVersion   string `json:"FunctionVersion"`
		FunctionName      string `json:"FunctionName"`
		Namespace         string `json:"Namespace"`
		InstallDependency string `json:"InstallDependency"`
		Status            string `json:"Status"`
		AvailableStatus   string `json:"AvailableStatus"`
		StatusDesc        string `json:"StatusDesc"`
		FunctionId        string `json:"FunctionId"`
		L5Enable          string `json:"L5Enable"`
		EipConfig         struct {
			EipFixed string        `json:"EipFixed"`
			Eips     []interface{} `json:"Eips"`
		} `json:"EipConfig"`
		ModTime          string        `json:"ModTime"`
		AddTime          string        `json:"AddTime"`
		Layers           []interface{} `json:"Layers"`
		DeadLetterConfig struct {
			Type       string `json:"Type"`
			Name       string `json:"Name"`
			FilterType string `json:"FilterType"`
		} `json:"DeadLetterConfig"`
		OnsEnable       string `json:"OnsEnable"`
		PublicNetConfig struct {
			PublicNetStatus string `json:"PublicNetStatus"`
			EipConfig       struct {
				EipStatus  string        `json:"EipStatus"`
				EipAddress []interface{} `json:"EipAddress"`
			} `json:"EipConfig"`
		} `json:"PublicNetConfig"`
		DeployMode  string        `json:"DeployMode"`
		Triggers    []interface{} `json:"Triggers"`
		ClsLogsetId string        `json:"ClsLogsetId"`
		ClsTopicId  string        `json:"ClsTopicId"`
		CodeInfo    string        `json:"CodeInfo"`
		CodeResult  string        `json:"CodeResult"`
		CodeError   string        `json:"CodeError"`
		ErrNo       int           `json:"ErrNo"`
		Tags        []interface{} `json:"Tags"`
		AccessInfo  struct {
			Host string `json:"Host"`
			Vip  string `json:"Vip"`
		} `json:"AccessInfo"`
		Type      string `json:"Type"`
		CfsConfig struct {
			CfsInsList []interface{} `json:"CfsInsList"`
		} `json:"CfsConfig"`
		StatusReasons  []interface{} `json:"StatusReasons"`
		AsyncRunEnable string        `json:"AsyncRunEnable"`
		TraceEnable    string        `json:"TraceEnable"`
		LogType        string        `json:"LogType"`
		RequestId      string        `json:"RequestId"`
	} `json:"Response"`
}

func (c *Client) GetFunction(functionName string) (*GetFunctionResponse, error) {
	resp, err := c.request(http.MethodPost, "GetFunction", GetFunctionRequest{
		FunctionName: functionName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "get function")
	}

	var respJSON GetFunctionResponse
	return &respJSON, resp.ToJSON(&respJSON)
}
