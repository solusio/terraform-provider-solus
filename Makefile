NAME=solusio
BINARY=terraform-provider-${NAME}
HOOK=hooks/pre-commit/main.go
LIST=`go list ./... | grep -v /hooks/pre-commit`

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
	go test $(TESTARGS) -mod=vendor -race -cover $(LIST)

.PHONY: testacc
testacc:
	TF_ACC=1 go test $(TESTARGS) -mod=vendor -race -cover $(LIST)

.PHONY: build
build:
	go build -o ${BINARY}

.PHONY: init
init: init/hook

init/hook: ${HOOK}
	go build -o .git/hooks/pre-commit ${HOOK}
