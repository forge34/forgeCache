BINARY_NAME=redis-clone
MAIN_PATH=cmd/server/main.go

.PHONY: all build run clean

all: build

build:
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

run: build
	@echo "Starting server..."
	@./bin/$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean
