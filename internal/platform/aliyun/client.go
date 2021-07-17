// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package aliyun

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/thanhpk/randstr"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/platform"
)

var _ platform.Cloud = (*Client)(nil)

type Client struct {
	regionID                                string
	accountID, accessKeyID, accessKeySecret string
}

func New(opts platform.AuthenticateOptions) *Client {
	return &Client{
		regionID:        opts[RegionIDField],
		accountID:       opts[AccountIDField],
		accessKeyID:     opts[AccessKeyIDField],
		accessKeySecret: opts[AccessKeySecretField],
	}
}

func (c *Client) String() string {
	return "aliyun"
}

func (c *Client) Authenticate() error {
	u, err := url.Parse(fmt.Sprintf("https://ecs-%s.aliyuncs.com/", c.regionID))
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	query := url.Values{}
	query.Set("AccessKeyId", c.accessKeyID)
	query.Set("Action", "DescribeRegions")
	query.Set("Format", "JSON")
	query.Set("RegionId", c.regionID)
	query.Set("SignatureMethod", "HMAC-SHA1")
	query.Set("SignatureNonce", strings.ToUpper(randstr.String(24)))
	query.Set("SignatureType", "")
	query.Set("SignatureVersion", "1.0")
	query.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	query.Set("Version", "2014-05-26")
	u.RawQuery = query.Encode()

	// Generate signature
	hashSign := hmac.New(sha1.New, []byte(c.accessKeySecret+"&"))
	hashSign.Write([]byte(http.MethodGet + "&%2F&" + url.QueryEscape(u.RawQuery)))
	signature := base64.StdEncoding.EncodeToString(hashSign.Sum(nil))
	u.RawQuery += "&Signature=" + url.QueryEscape(signature)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "new request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var respJSON struct {
			Code      string
			Message   string
			Recommend string
		}
		if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
			return errors.Wrap(err, "JSON decode")
		}

		log.Error("Failed to authenticate to aliyun.")
		log.Warn("[ %s ] %s", respJSON.Code, respJSON.Message)
		log.Warn("Recommend: %s", respJSON.Recommend)
		return errors.New(respJSON.Code)
	}
	return nil
}

func (c *Client) request(method string, baseURL string, requestBody ...interface{}) (*response, error) {
	u := fmt.Sprintf("https://%s.%s.fc.aliyuncs.com/%s/%s", c.accountID, c.regionID, ApiVersion, strings.TrimLeft(baseURL, "/"))
	var body io.Reader
	if len(requestBody) == 1 {
		repBody, err := json.Marshal(requestBody[0])
		if err != nil {
			return nil, errors.Wrap(err, "JSON encode")
		}
		body = bytes.NewReader(repBody)
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	req.Header.Set("Authorization", c.GetAuthorizationHeader(req))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}

	return &response{
		Response: resp,
	}, nil
}

type response struct {
	*http.Response
}

func (r *response) ToJSON(v interface{}) error {
	defer func() { _ = r.Body.Close() }()
	return json.NewDecoder(r.Body).Decode(v)
}

func (r *response) ToString() string {
	defer func() { _ = r.Body.Close() }()
	resp, _ := io.ReadAll(r.Body)
	return string(resp)
}
