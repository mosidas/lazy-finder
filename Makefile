BINARY := lazy-finder
BIN_DIR := bin

.PHONY: all build run test fmt vet tidy clean install

all: build

build:
	go build -o $(BIN_DIR)/$(BINARY) .

run:
	go run . $(ARGS)

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

tidy:
	go mod tidy

install:
	go install .

clean:
	rm -rf $(BIN_DIR)
