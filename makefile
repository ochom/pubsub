SHELL=/bin/bash
producer:
	go build -o pub publisher/main.go

subscriber:
	go build -o sub consumer/main.go