.DEFAULT_GOAL := build

.PHONY:fmt vet build

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build

# Variables
SCHEMA_FILES := ./resources/cf-schema.json
OUTPUT_DIR := ./generated
PACKAGE_NAME := generated
GO_JSONSCHEMA := github.com/atombender/go-jsonschema

# Target to ensure go-jsonschema is installed
check-go-jsonschema:
	@go list -m $(GO_JSONSCHEMA) >/dev/null 2>&1 || ( \
		echo "Installing $(GO_JSONSCHEMA)..." && \
		go get $(GO_JSONSCHEMA)/... \
	)

# Target to generate Go structs from JSON schemas
generate: check-go-jsonschema
	go-jsonschema -p $(PACKAGE_NAME) $(SCHEMA_FILES) -o $(OUTPUT_DIR)
.PHONY: generate