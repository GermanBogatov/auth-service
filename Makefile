lint:
	@go clean -cache
	golangci-lint run ./...
.PHONY: lint

install-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
.PHONY: install-lint

compose-up:
	docker-compose up
.PHONY: compose-up