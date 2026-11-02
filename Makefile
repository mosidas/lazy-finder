BINARY := lazy-finder
BIN_DIR := bin
DIST_DIR := dist
PLATFORMS := darwin/arm64 darwin/amd64 linux/arm64 linux/amd64

.PHONY: all build run test fmt vet tidy clean install dist

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

dist:
	@if [ -z "$(TAG)" ]; then echo "make dist TAG=vX.Y.Z の形で指定する" >&2; exit 1; fi
	rm -rf $(DIST_DIR)
	@set -e; for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; out=$(DIST_DIR)/$$os-$$arch; \
		echo "build $$os/$$arch"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -o $$out/$(BINARY) .; \
		cp LICENSE $$out/LICENSE; \
		tar -czf $(DIST_DIR)/$(BINARY)-$(TAG)-$$os-$$arch.tar.gz -C $$out $(BINARY) LICENSE; \
	done

clean:
	rm -rf $(BIN_DIR) $(DIST_DIR)
