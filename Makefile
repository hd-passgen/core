BINARY_NAME=hd-passgen

GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif

build:
	CGO_ENABLED=0 go build -o $(BINARY_NAME) .

install: build
	install $(BINARY_NAME) ${GOBIN}/$(BINARY_NAME)
	rm -f $(BINARY_NAME)
	
lint:
	golangci-lint run ./...
