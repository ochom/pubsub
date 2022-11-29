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