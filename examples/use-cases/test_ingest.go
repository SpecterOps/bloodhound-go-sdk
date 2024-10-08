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
	"flag"
	"fmt"
	. "github.com/SpecterOps/bloodhound-go-sdk/sdk"
	"log"
	"os"
)

func parseArgs() *string {
	flag.Args()
	filePath := flag.String("file", "", "Path to the sample zip file")
	flag.Parse()
	if *filePath == "" {
		log.Fatal("Specify a file path to a sample zip file")
	}
	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s\n", filePath)
	}
	return filePath
}

func ingestFile(path string, client *ClientWithResponses) (err error) {
	var params = &CreateFileUploadJobParams{}

	createFileUploadResponse, err := client.CreateFileUploadJobWithResponse(context.Background(), params)
	if err != nil {
		log.Print("Error while creating a file upload job id", err)
		return err
	}
	if createFileUploadResponse.StatusCode() != 201 {
		return fmt.Errorf("error getting api version %v", createFileUploadResponse.StatusCode())
	}

	var job_id = *createFileUploadResponse.JSON201.Data.Id
	var content_type = "application/zip"

	test_file, err := os.Open(path)
	if err != nil {
		log.Print("Error opening file", err)
		return err
	}

	var uploadParams = &UploadFileToJobParams{
		ContentType: "application/zip",
	}

	response, err := client.UploadFileToJobWithBodyWithResponse(context.Background(), job_id, uploadParams, content_type, test_file)
	if err != nil {
		log.Print("Error uploading file", err)
		return err
	}

	if response.StatusCode() != 202 { // Accepted
		return fmt.Errorf("error getting api version %v", createFileUploadResponse.StatusCode())
	}

	job, err := client.EndFileUploadJob(context.Background(), job_id, nil)
	if err != nil || job.StatusCode != 200 {
		return err
	}

	return nil
}

func main() {
	sample_zip_file := parseArgs()

	customHttpClient, rerr := GetLocalhostWithSubdomainHttpClient()
	if rerr != nil {
		log.Fatal("Ooof cant make bloodhound.localhost resolving http.Client", rerr)
	}

	// API token and token id obtained from environment variables
	key := os.Getenv("API_TOKEN_KEY")
	id := os.Getenv("API_TOKEN_ID")
	if key == "" || id == "" {
		log.Fatal("You must set API_TOKEN_KEY and API_TOKEN_ID environment variables to your API key and id values.")
	}

	// server URL obtained from environment variables
	server := os.Getenv("BLOODHOUND_SERVER")
	if server == "" {
		log.Fatal("You must set BLOODHOUND_SERVER environment variable to the URL of the bloodhound server")
	}

	// HMAC Security Provider
	var hmacTokenProvider, serr = NewSecurityProviderHMACCredentials(key, id)

	if serr != nil {
		log.Fatal("Error creating bearer token middleware", serr)
	}
	client, crerr := NewClientWithResponses(
		server,
		WithRequestEditorFn(hmacTokenProvider.Intercept),
		WithHTTPClient(customHttpClient))
	if crerr != nil {
		log.Fatal("Error creating client", crerr)
	}

	err := ingestFile(*sample_zip_file, client)
	if err != nil {
		log.Fatal("Error uploading file", err)
	} else {
		log.Print("File uploaded successfully")
	}
}
