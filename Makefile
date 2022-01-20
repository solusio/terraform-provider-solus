NAME=solus
BINARY=terraform-provider-${NAME}
HOOK=hooks/pre-commit/main.go
LIST=`go list ./... | grep -v /hooks/pre-commit`

ifneq (,$(wildcard ./.testacc.env))
	include .testacc.env
	export
endif

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
	go test -mod=vendor -race -cover $(LIST) $(TESTARGS)

.PHONY: testacc
testacc:
	TF_ACC=1 go test -mod=vendor -race -cover $(LIST) $(TESTARGS) -timeout 120m

.PHONY: build
build:
	go build -o ${BINARY}

.PHONY: init
init: init/hook

init/hook: ${HOOK}
	go build -o .git/hooks/pre-commit ${HOOK}
