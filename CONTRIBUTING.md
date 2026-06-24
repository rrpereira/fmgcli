# Contributing

Thanks for contributing to `fmgcli`.

## Development setup

1. Install Go 1.23 or newer.
2. Clone the repository.
3. Run tests:

```bash
go test ./...
```

## Pull request guidelines

- Keep changes focused and small.
- Add or update tests for behavior changes.
- Do not commit credentials or FortiManager session tokens.
- Ensure `go test ./...` passes before opening a PR.

## Code style

- Follow standard Go formatting (`gofmt`).
- Prefer clear, explicit error messages.
