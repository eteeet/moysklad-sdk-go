GOLANGCI_LINT_VERSION := v2.15.0

.PHONY: test lint gen

test:
	go test ./...

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

gen:
	go generate ./...
