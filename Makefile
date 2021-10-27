.PHONY: build
build:
	go build cmd/act-transition-api/main.go

.PHONY: test
test:
	go test -v ./...