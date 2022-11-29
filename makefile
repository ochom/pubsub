SHELL=/bin/bash

build:
	@echo "building dev ..."
	@go build -o dist/pub examples/publisher/main.go
	@go build -o dist/sub examples/consumer/main.go

pub:
	@echo "Running publisher ..."
	@make dist
	@source ./env.sh && ./dist/pub

sub:
	@echo "Running consumer ..."
	@make dist
	@source ./env.sh && ./dist/sub

install:
	go mod tidy
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/kisielk/errcheck@latest
	go install golang.org/x/lint/golint@latest
	go install github.com/axw/gocov/gocov@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/client9/misspell/cmd/misspell@latest

lint:
	staticcheck ./...
	go fmt $(go list ./... | grep -v /vendor/)
	go vet $(go list ./... | grep -v /vendor/)
	golint -set_exit_status $(go list ./... | grep -v /vendor/)
	errcheck -ignore 'os:.*,' $(go list ./... | grep -v /vendor/)
	misspell -error .
	gosec ./...
