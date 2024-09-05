// Copyright 2023 Specter Ops, Inc.
//
// Licensed under the Apache License, Version 2.0
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package oapiclient

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"time"
)

type HMACCredentials struct {
	TokenKey string
	TokenID  string
}

func NewSecurityProviderHMACCredentials(token string, token_id string) (*HMACCredentials, error) {
	return &HMACCredentials{
		TokenKey: token,
		TokenID:  token_id,
	}, nil
}

// Based on python example
func (c *HMACCredentials) Intercept(ctx context.Context, req *http.Request) error {
	// Digester is initialized with HMAC-SHA-256 using the token key as the HMAC digest key.
	digester := hmac.New(sha256.New, []byte(c.TokenKey))

	// OperationKey is the first HMAC digest link in the signature chain. This prevents replay attacks that seek to
	// modify the request method or URI. It is composed of concatenating the request method and the request URI with
	// no delimiter and computing the HMAC digest using the token key as the digest secret.
	//
	// Example: GET /api/v2/test/resource HTTP/1.1
	// Signature Component: GET/api/v2/test/resource
	digester.Write([]byte(req.Method + req.URL.RequestURI()))

	// Update the digester for further chaining
	digester = hmac.New(sha256.New, digester.Sum(nil))

	// DateKey is the next HMAC digest link in the signature chain. This encodes the RFC3339 formatted datetime
	// value as part of the signature to the hour to prevent replay attacks that are older than max two hours. This
	// value is added to the signature chain by cutting off all values from the RFC3339 formatted datetime from the
	// hours value forward:
	//
	// Example: 2020-12-01T23:59:60Z
	// Signature Component: 2020-12-01T23
	// Format the current time as RFC3339
	datetimeFormatted := time.Now().UTC().Format(time.RFC3339)
	digester.Write([]byte(datetimeFormatted[:13]))

	// Update the digester for further chaining
	digester = hmac.New(sha256.New, digester.Sum(nil))

	// Body signing is the last HMAC digest link in the signature chain. This encodes the request body as part of
	// the signature to prevent replay attacks that seek to modify the payload of a signed request. In the case
	// where there is no body content the HMAC digest is computed anyway, simply with no values written to the
	// digester.
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}

		// Interceptors modify request in place (sigh)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		digester.Write(bodyBytes)
	}

	// Perform the request with the signed and expected headers
	req.Header.Set("User-Agent", "bhe-go-sdk 0001")
	req.Header.Set("Authorization", "bhesignature "+c.TokenID)
	req.Header.Set("RequestDate", datetimeFormatted)
	req.Header.Set("Signature", base64.StdEncoding.EncodeToString(digester.Sum(nil)))

	return nil
}
