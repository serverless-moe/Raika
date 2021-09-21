// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package aws

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/platform"
)

func (c *Client) CreateFunction(opts platform.CreateFunctionOptions) (string, error) {
	fmt.Println(c.accessKey, c.secretKey, c.regionID)
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(c.accessKey, c.secretKey, ""),
		Region:      &c.regionID,
	})
	if err != nil {
		return "", errors.Wrap(err, "new session")
	}

	zipFileBase64, err := packFile(opts.File)
	if err != nil {
		return "", errors.Wrap(err, "pack file")
	}

	environmentVariables := make(map[string]*string)
	for k, v := range opts.EnvironmentVariables {
		value := v
		environmentVariables[k] = &value
	}

	lamb := lambda.New(sess)
	resp, err := lamb.CreateFunction(&lambda.CreateFunctionInput{
		Code: &lambda.FunctionCode{
			ZipFile: zipFileBase64,
		},
		Description:  &opts.Description,
		Environment:  &lambda.Environment{Variables: environmentVariables},
		FunctionName: &opts.Name,
		Handler:      aws.String("bootstrap"),
		MemorySize:   &opts.MemorySize,
		Role:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", c.accountID, c.roleName)),
		Runtime:      aws.String("provided"),
		Timeout:      aws.Int64(int64(opts.RuntimeTimeout / time.Second)),
	})
	if err != nil {
		return "", err
	}
	log.Trace("%+v", resp)
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
