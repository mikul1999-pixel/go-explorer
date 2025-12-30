APP_NAME := go-explorer
CMD_DIR := ./cmd/explorer
PORT := 3030
ROOT_DIR := /tmp/go-explorer-data

.PHONY: run test build clean

run:
	@mkdir -p $(ROOT_DIR)
	@echo "Running $(APP_NAME) on :$(PORT)"
	ROOT_DIR=$(ROOT_DIR) LISTEN_ADDR=:$(PORT) go run $(CMD_DIR)

test:
	go test ./...

build:
	go build -o bin/$(APP_NAME) $(CMD_DIR)

clean:
	rm -rf bin
