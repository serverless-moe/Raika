// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package aliyun

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"io"
	"net/http"
	"sort"
	"strings"
)

// GetAuthorizationHeader returns the authorization header.
func (c *Client) GetAuthorizationHeader(req *http.Request) string {
	return "FC " + c.accessKeyID + ":" + c.GetSignature(req)
}

// GetSignature returns the signature string.
func (c *Client) GetSignature(req *http.Request) string {
	// Sort fcHeaders.
	headers := &fcHeaders{}
	for k := range req.Header {
		if strings.HasPrefix(strings.ToLower(k), "x-fc-") {
			headers.Keys = append(headers.Keys, strings.ToLower(k))
			headers.Values = append(headers.Values, req.Header.Get(k))
		}
	}
	sort.Sort(headers)
	fcHeaders := ""
	for i := range headers.Keys {
		fcHeaders += headers.Keys[i] + ":" + headers.Values[i] + "\n"
	}

	httpMethod := req.Method
	contentMd5 := req.Header.Get("Content-MD5")
	contentType := "application/json"
	date := req.Header.Get("Date")
	fcResource := ApiVersion + req.URL.Path

	signStr := httpMethod + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + fcHeaders + fcResource

	h := hmac.New(func() hash.Hash { return sha256.New() }, []byte(c.accessKeySecret))
	_, _ = io.WriteString(h, signStr)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type fcHeaders struct {
	Keys   []string
	Values []string
}

func (f *fcHeaders) Len() int { return len(f.Values) }
func (f *fcHeaders) Less(i, j int) bool {
	return bytes.Compare([]byte(f.Keys[i]), []byte(f.Keys[j])) < 0
}
func (f *fcHeaders) Swap(i, j int) {
	f.Values[i], f.Values[j] = f.Values[j], f.Values[i]
	f.Keys[i], f.Keys[j] = f.Keys[j], f.Keys[i]
}
