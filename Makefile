BINARY_NAME=fconv
BUILD_DIR=./bin
MAIN_PATH=./cmd/file-converter/main.go
GOFLAGS=-ldflags="-s -w"

.PHONY: build clean install test deps run

build:
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

clean:
	rm -rf cmd/$(BUILD_DIR)
	go clean -cache

test:
	go test -v ./...

deps:
	go mod download
	go mod tidy

run: build
	$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

.PHONY: fmt lint

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

