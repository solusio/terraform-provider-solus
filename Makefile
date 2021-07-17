NAME=solusio
BINARY=terraform-provider-${NAME}
HOOK=hooks/pre-commit/main.go

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

.PHONY: init
init: init/hook

init/hook: ${HOOK}
	go build -o .git/hooks/pre-commit ${HOOK}
