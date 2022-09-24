SHELL=/bin/bash
producer:
	go run "publisher/main.go"

subscriber:
	go run "consumer/main.go"