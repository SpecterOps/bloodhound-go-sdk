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

package main

import (
	context "context"
	"fmt"
	. "github.com/SpecterOps/bloodhound-go-sdk/sdk"
	. "github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"log"
)

func main() {
	// httpClient that handles localhost with subdomains (bloodhound.localhost)
	var customHttpClient, rerr = GetLocalhostWithSubdomainHttpClient()
	if rerr != nil {
		log.Fatal("Ooof cant make bloodhound.localhost resolving http.Client", rerr)
	}

	// The bearer_token
	var bearer_token = "<YOUR BEARER TOKEN>"

	// Bearer Token Security Provider
	var bearerTokenProvider, serr = NewSecurityProviderBearerToken(bearer_token)
	if serr != nil {
		log.Fatal("Error creating bearer token middleware", serr)
	}

	client, crerr := NewClientWithResponses(
		"http://bloodhound.localhost/",
		WithRequestEditorFn(bearerTokenProvider.Intercept),
		WithBaseURL("http://bloodhound.localhost/"),
		WithHTTPClient(customHttpClient))
	if crerr != nil {
		log.Fatal("Error creating client", crerr)
	}

	// Get the API Version from the server
	var params = &GetApiVersionParams{}
	version, err := client.GetApiVersionWithResponse(context.Background(), params)
	if err != nil {
		log.Print("Error while getting api version", err)
		return
	}
	if version.StatusCode() == 200 {
		fmt.Printf("Version: %s\n", *version.JSON200.Data.ServerVersion)
	} else {
		log.Fatal("Error getting api version", version.StatusCode())
	}
}
