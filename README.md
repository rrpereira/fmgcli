# fmgcli

Go client for FortiManager JSON-RPC APIs.

## Features

- User/password login flow and API key flow
- Session-based request support
- FortiManager object management helpers for:
  - Addresses
  - Services
  - Policies
  - Workspace lock/unlock and commit
- Test suite based on `httptest`

## Installation

```bash
go get github.com/rrpereira/fmgcli
```

## Quick Start

### User Login

```go
package main

import (
    "log"

    "github.com/rrpereira/fmgcli"
)

func main() {
    client := fmgcli.NewUserClient("https://fortimanager.example.com", "my-user", "my-password")

    if err := client.Login(); err != nil {
        log.Fatal(err)
    }
    defer client.Logout()

    if err := client.CreateSubnetAddress("root", "host-10-0-0-1", "10.0.0.1", "255.255.255.255", "managed by fmgcli"); err != nil {
        log.Fatal(err)
    }
}
```

### API Key

```go
package main

import (
    "log"

    "github.com/rrpereira/fmgcli"
)

func main() {
    client := fmgcli.NewAPIClient("https://fortimanager.example.com", "my-api-key")

    if _, err := client.GetServiceByName("root", "HTTPS"); err != nil {
        log.Fatal(err)
    }
}
```

## Compatibility

- Requires Go 1.23+
- Designed for FortiManager JSON-RPC endpoints (`/jsonrpc`)

## Development

```bash
go test ./...
```

### End-to-End Tests

The repository includes a full template at `.env.example` with every variable referenced by e2e tests.

```bash
cp .env.example .env
# Edit .env with values from your FortiManager environment
go test -tags=e2e ./...
```

Environment variables used by e2e tests:

- `FMG_E2E_HOST`
- `FMG_E2E_USER`
- `FMG_E2E_PASSWORD`
- `FMG_E2E_TOKEN`
- `FMG_E2E_ADOM`
- `FMG_E2E_ADDRESS_NAME`
- `FMG_E2E_ADDRESS_NAMES`
- `FMG_E2E_ADDRESS_METAFIELD_KEY`
- `FMG_E2E_ADDRESS_METAFIELD_VALUE`
- `FMG_E2E_ADDRESS_METAFIELD_VALUES`
- `FMG_E2E_ADDRESS_NAME_IP_NETMASK_NAME`
- `FMG_E2E_ADDRESS_NAME_IP_NETMASK_IP`
- `FMG_E2E_ADDRESS_NAME_IP_NETMASK_NETMASK`
- `FMG_E2E_PKG`
- `FMG_E2E_POLICY_ID`
- `FMG_E2E_POLICY_METAFIELD_KEY`
- `FMG_E2E_POLICY_METAFIELD_VALUE`
- `FMG_E2E_POLICY_METAFIELD_VALUES`
- `FMG_E2E_SERVICE_NAME`
- `FMG_E2E_SERVICE_PROTOCOL`
- `FMG_E2E_SERVICE_MIN_PORT`
- `FMG_E2E_SERVICE_MAX_PORT`
- `FMG_E2E_SERVICE_NAMES`
- `FMG_E2E_SERVICE_METAFIELD_KEY`
- `FMG_E2E_SERVICE_METAFIELD_VALUE`
- `FMG_E2E_SERVICE_METAFIELD_VALUES`
- `FMG_E2E_GROUP_METAFIELD_KEY`

## Versioning

This project follows semantic versioning.

## License

This project is licensed under the MIT License. See `LICENSE`.

## Disclaimer

Fortinet and FortiManager are trademarks of Fortinet, Inc. This project is an independent client and is not officially affiliated with Fortinet.
