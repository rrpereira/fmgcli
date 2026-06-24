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

Set these environment variables depending on the test:

- `FMG_E2E_HOST` — required for all e2e tests
- `FMG_E2E_USER` and `FMG_E2E_PASSWORD` — required for `TestE2E_LoginLogout`
- `FMG_E2E_TOKEN` — required for `TestE2E_LockUnlock`
- `FMG_E2E_ADOM` — optional for `TestE2E_LockUnlock` (defaults to `root`)

You can set these in a `.env` file in the project root, or export them in your shell before running tests:

```bash
export FMG_E2E_HOST=https://your-fortimanager.example.com
export FMG_E2E_TOKEN=your-api-token
export FMG_E2E_USER=your-username
export FMG_E2E_PASSWORD=your-password
export FMG_E2E_ADOM=root
go test -tags=e2e ./...
```

Or use inline environment variables:

```bash
FMG_E2E_HOST=https://your-fortimanager.example.com FMG_E2E_TOKEN=your-api-token FMG_E2E_USER=your-username FMG_E2E_PASSWORD=your-password FMG_E2E_ADOM=root go test -tags=e2e ./...
```

## Versioning

This project follows semantic versioning.

## License

This project is licensed under the MIT License. See `LICENSE`.

## Disclaimer

Fortinet and FortiManager are trademarks of Fortinet, Inc. This project is an independent client and is not officially affiliated with Fortinet.
