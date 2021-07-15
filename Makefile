NAME=solus
BINARY=terraform-provider-${NAME}

.PHONY: all
all: fmt lint test build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -mod=vendor -race -cover $$(go list ./...)

.PHONY: build
build:
	go build -o ${BINARY}
