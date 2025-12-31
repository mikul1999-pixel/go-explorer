BIN_DIR  := ./bin
APP_NAME := go-explorer
CMD_DIR := ./cmd/explorer
PORT := 3030
ROOT_DIR := /home/mikul/data

ENV_FILE := .env

ifneq (,$(wildcard $(ENV_FILE)))
	include $(ENV_FILE)
	export
endif

.PHONY: help run test build clean

help:
	@echo "make:"
	@echo "  run         Run locally (go run)"
	@echo "  test        Run tests"
	@echo "  build       Build binary to ./bin"
	@echo "  clean       Remove build artifacts"

run:
	@echo "Running $(APP_NAME) on :$(PORT)"
	ROOT_DIR=$(ROOT_DIR) LISTEN_ADDR=:$(PORT) go run $(CMD_DIR)

test:
	go test ./... -v

build:
	go build -o bin/$(APP_NAME) $(CMD_DIR)

clean:
	rm -rf $(BIN_DIR)

