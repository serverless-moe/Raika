// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tencentcloud

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GetAuthorizationHeader returns the authorization header.
func (c *Client) GetAuthorizationHeader(req *http.Request, body []byte) string {
	t := time.Now().UTC()
	date := t.Format("2006-01-02")
	req.Header.Set("x-tc-timestamp", strconv.Itoa(int(t.Unix())))
	credentialScope := date + "/scf/tc3_request"

	var headerKeys []string
	for k := range req.Header {
		k = strings.ToLower(k)
		if strings.HasPrefix(k, "x-tc-") {
			continue
		}
		headerKeys = append(headerKeys, k)
	}
	sort.Strings(headerKeys)
	signedHeaders := strings.Join(headerKeys, ";")

	var canonicalHeaders string
	for _, k := range headerKeys {
		v := req.Header.Get(k)
		canonicalHeaders += strings.ToLower(k) + ":" + v + "\n"
	}

	hashedCanonicalRequest := sha256Hex(
		[]byte(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
			req.Method,
			req.URL.Path,
			req.URL.RawQuery,
			canonicalHeaders,
			signedHeaders,
			sha256Hex(body),
		)),
	)

	// Signature
	secretDate := hmacSha256(date, "TC3"+c.secretKey)
	secretService := hmacSha256("scf", secretDate)
	secretSigning := hmacSha256("tc3_request", secretService)
	signature := hex.EncodeToString([]byte(hmacSha256(
		fmt.Sprintf("%s\n%d\n%s\n%s",
			"TC3-HMAC-SHA256",
			t.Unix(),
			credentialScope,
			hashedCanonicalRequest,
		),
		secretSigning),
	))

	return fmt.Sprintf("TC3-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s", c.secretID, credentialScope, signedHeaders, signature)
}

func hmacSha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	_, _ = hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

func sha256Hex(s []byte) string {
	b := sha256.Sum256(s)
	return hex.EncodeToString(b[:])
}
