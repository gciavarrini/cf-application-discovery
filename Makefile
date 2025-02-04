.DEFAULT_GOAL := build

.PHONY:fmt vet build lint

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build

lint:
	staticcheck ./...
	go vet ./...