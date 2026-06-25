BINARY   := taskman
BUILD_DIR := ./bin
MAIN     := ./main.go
GOBIN    := $(shell go env GOPATH)/bin

.PHONY: all build run test lint clean install

all: build

## build: compile the binary into ./bin/
build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY) $(MAIN)
	@echo "Built $(BUILD_DIR)/$(BINARY)"

## install: install the binary to GOPATH/bin
install:
	@#region agent log
	@bash -c 'ts=$$(date +%s%3N 2>/dev/null || date +%s); echo "{\"sessionId\":\"a57069\",\"runId\":\"verify\",\"hypothesisId\":\"A\",\"location\":\"Makefile:install\",\"message\":\"install target invoked\",\"data\":{\"gopath\":\"$$(go env GOPATH)\"},\"timestamp\":$$ts}" >> .cursor/debug-a57069.log'
	@#endregion
	go build -ldflags="-s -w" -o $(GOBIN)/$(BINARY) $(MAIN)
	@#region agent log
	@bash -c 'ts=$$(date +%s%3N 2>/dev/null || date +%s); bin="$$(go env GOPATH)/bin/taskman"; ex=false; [ -f "$$bin" ] && ex=true; echo "{\"sessionId\":\"a57069\",\"runId\":\"verify\",\"hypothesisId\":\"E\",\"location\":\"Makefile:install\",\"message\":\"install finished\",\"data\":{\"binary\":\"$$bin\",\"exists\":$$ex},\"timestamp\":$$ts}" >> .cursor/debug-a57069.log'
	@#endregion

## run: build and run with optional ARGS (e.g. make run ARGS="list --all")
run: build
	$(BUILD_DIR)/$(BINARY) $(ARGS)

## test: run all tests with race detector
test:
	go test -race ./...

## test-cover: run tests and open HTML coverage report
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## lint: run golangci-lint (must be installed)
lint:
	golangci-lint run ./...

## clean: remove build artifacts
clean:
	rm -rf $(BUILD_DIR) coverage.out