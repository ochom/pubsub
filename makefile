SHELL=/bin/bash

build:
	@echo "Building dev ..."
	@go build -o build/pub examples/publisher/main.go
	@go build -o build/sub examples/consumer/main.go

pub:
	@echo "Running publisher ..."
	@make build
	@source ./env.sh && ./build/pub

sub:
	@echo "Running consumer ..."
	@make build
	@source ./env.sh && ./build/sub