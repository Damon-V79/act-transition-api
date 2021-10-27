.PHONY: all test lint tidy style cover
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

build:
	go build -o bin/act-transition-retranslator cmd/act-transition-api/main.go
