# OAPI Generated Bloodhound SDK

## Overview

This is a Go SDK for BloodHound version `<version>`.

## OpenAPI spec

This client SDK is generated from `openapi/openapi.yaml`.

TODO: 
- Include versioning scheme*
- Add generation information

## Building and Running

### Dependencies TODO:

## Examples

For each of the examples, you must set the following environment variables:

| Name | Value | Example                               |
|------|-------|---------------------------------------|
 | API_TOKEN | Generated API token | hk...jgfZCQ==                         |
 | API_TOKEN_ID | Id of generated API token | 467e-bb1f-dc29...5bfc                 |
 | BLOODHOUND_SERVER | Server URL | https://demo.bloodhoundenterprise.io/ |
### Authentication examples.  

`test_bearer_token_client.go` demonstrates how to use the SDK with bearer token authentication.

`test_hmac_token_client.go` demonstrates how to use the SDK with API token authentication.

### More complex cases

`test_ingest.go` demonstrates slightly more complex use of the SDK.

## Build And Run Examples

### Bearer Token Authentication

```bash
cd examples/bearer-authentication
go run ./test_bearer_token_client.go
```

### HMAC Token Authentication

```bash
cd examples/hmac-authentication
go run ./test_hmac_token_client.go
```

### Ingest Example

```bash
cd examples/use-cases
go run ./test_ingest.go --file <path to zip file>
```

## Contact

Please check out the [Contact page](https://github.com/SpecterOps/BloodHound/wiki/Contact) in our wiki for details on how to reach out with questions and suggestions.

## Licensing

```
Copyright 2024 Specter Ops, Inc.

Licensed under the Apache License, Version 2.0
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

Unless otherwise annotated by a lower-level LICENSE file or license header, all files in this repository are released
under the `Apache-2.0` license. A full copy of the license may be found in the top-level [LICENSE](LICENSE) file.
rwise annotated by a lower-level LICENSE file or license header, all files in this repository are released under the Apache-2.0 license. A full copy of the license may be found in the top-level LICENSE file.