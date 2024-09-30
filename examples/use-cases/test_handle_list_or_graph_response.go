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
	"context"
	"encoding/json"
	"io"
	"log"
	. "oapi-client/sdk"
	"strings"
)

func main() {

	customHttpClient, rerr := GetLocalhostWithSubdomainHttpClient()
	if rerr != nil {
		log.Fatal("Ooof cant make bloodhound.localhost resolving http.Client", rerr)
	}

	// API token
	var token = "/2lJ0Gf3qjb6zvZw01EoFmK34qtwkiOOeUgcofMXh5dyYEF61LwRvQ=="
	var token_id = "8efd2add-c453-4a01-a3e0-2775c9b9cbc4"

	// HMAC Security Provider
	var hmacTokenProvider, serr = NewSecurityProviderHMACCredentials(token, token_id)

	if serr != nil {
		log.Fatal("Error creating bearer token middleware", serr)
	}
	client, crerr := NewClientWithResponses(
		"http://bloodhound.localhost/",
		WithRequestEditorFn(hmacTokenProvider.Intercept),
		WithBaseURL("http://bloodhound.localhost/"),
		WithHTTPClient(customHttpClient))
	if crerr != nil {
		log.Fatal("Error creating client", crerr)
	}

	// Get the Domain objects
	response, err := client.GetAvailableDomainsWithResponse(context.Background(), nil)
	if err != nil {
		log.Fatal("Error getting available domains", err)
	}
	if response.StatusCode() != 200 {
		log.Fatal("Error getting available domains", response.StatusCode())
	}

	// if success we get domainentity computers and we get this as a list
	for _, y := range *response.JSON200.Data {
		log.Printf("Domain name: %s id: %s type: %s", *y.Name, *y.Id, *y.Type)
		listReturnType := "list"
		computersResponse, err := client.GetDomainEntityComputersWithResponse(context.Background(),
			*y.Id,
			&GetDomainEntityComputersParams{
				Type: (*GetDomainEntityComputersParamsType)(&listReturnType),
			},
		)
		if computersResponse.StatusCode() != 200 {
			log.Println("Error getting domain computersResponse", computersResponse.StatusCode())
			continue
		}

		// For reach computer we get the computer entity controllables as a graph
		graphReturnType := "graph"
		for _, value := range *computersResponse.JSON200.Data {
			log.Printf("\tComputer name: %s label: %s id: %s", *y.Name, *value.Label, *value.ObjectID)
			resp, err := client.GetComputerEntityControllables(context.Background(),
				*value.ObjectID,
				&GetComputerEntityControllablesParams{
					Type: (*GetComputerEntityControllablesParamsType)(&graphReturnType),
				})
			if err != nil {
				log.Fatal("Error getting computer entity controllables", err)
			}
			if resp.StatusCode != 200 {
				log.Println("Error getting resp", resp.StatusCode)
				continue
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error getting resp", resp.StatusCode)
				continue
			}
			defer func() { _ = resp.Body.Close() }()
			log.Println(string(body))
			var nodes map[string]json.RawMessage
			nodes_err := json.Unmarshal(body, &nodes)
			if nodes_err != nil {
				log.Println("Error getting resp", resp.StatusCode)
				continue
			}
			// Make
			for key, v := range nodes {
				if strings.HasPrefix(key, "rel_") {
					var graphEdge *ModelBhGraphEdge
					err := json.Unmarshal(v, &graphEdge)
					if err != nil {
						log.Println("Error unmarshalling json", string(v))
						continue
					}
					log.Println("Graph Edge", graphEdge)

				} else {
					var graphNode *ModelBhGraphNode
					err := json.Unmarshal(v, &graphNode)
					if err != nil {
						log.Println("Error unmarshalling json", string(v))
						continue
					}
					log.Println("Graph Node", graphNode)
				}
				log.Println("graph node ", key, value.Label)
				log.Println("graph node ", key, v)
			}
			//var graphItem *ModelBhGraphItem
			//json.Unmarshal(body, &graphItem)
			//log.Println(graphItem)
			//var graphData map[string]json.RawMessage
			//err = json.Unmarshal(body, &graphData)
			//if err != nil {
			//	log.Println("Error getting resp", resp.StatusCode)
			//	continue
			//}
			//for key, value := range graphData {
			//	var foo ModelBhGraphGraph
			//	jerr := json.Unmarshal(value, &foo)
			//	if jerr != nil {
			//		var prettyJSON bytes.Buffer
			//		json.Indent(&prettyJSON, value, "", "  ")
			//		log.Printf("Error unmarshalling key: *%s json: %s\n", key, prettyJSON.String())
			//	} else {
			//		log.Println("we unmarshalled a ModelBhGraphGraph")
			//	}
			//
			//}
			//log.Println("Computers response", body)
		}
		if err != nil {
			return
		}
	}
}
