GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.17","$(shell printf "$(GO_VERSION_SHORT)\n1.17" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.17. Found: $(GO_VERSION_SHORT))
endif

.PHONY: all test lint tidy style cover reqs

all: build

test:
	go test -v -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print ($$3)}'

lint:
	@command -v golangci-lint 2>&1 > /dev/null || (echo "Install golangci-lint" && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(shell go env GOPATH)/bin" v1.42.1)
	~/go/bin/golangci-lint run ./...

tidy:
	go mod tidy

style:
	find . -iname *.go | xargs gofmt -w

cover:
	go tool cover -html cover.out

reqs:
	go get -d github.com/golang/mock/mockgen@latest

build:
	go mod download && \
	  CGO_ENABLED=0 go build -o bin/act-transition-api cmd/act-transition-api/main.go
