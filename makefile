SHELL=/bin/bash

dev:
	@echo "Building dev ..."
	@go build -o build/pub examples/publisher/main.go
	@go build -o build/sub examples/consumer/main.go