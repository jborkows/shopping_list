.PHONY: build run dev test fmt migrate-up migrate-down clean

APP_NAME := shopping
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)

GOCACHE ?= $(CURDIR)/tmp/gocache

export GOCACHE

build:
	mkdir -p $(BIN_DIR) $(GOCACHE)
	go build -o $(BIN) ./cmd/shopping

run:
	mkdir -p $(GOCACHE)
	go run ./cmd/shopping

dev:
	air

test:
	mkdir -p $(GOCACHE)
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

clean:
	rm -rf ./tmp ./bin

migrate-up:
	migrate -path ./migrations -database "sqlite3://data/shopping.db" up

migrate-down:
	migrate -path ./migrations -database "sqlite3://data/shopping.db" down 1
