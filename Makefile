VERSION     ?= 1.0.0
BINARY      := convert
MAIN        := ./cmd/convert
LDFLAGS     := -s -w -X main.version=$(VERSION)
BUILD_FLAGS := -ldflags "$(LDFLAGS)"

.PHONY: build test clean lint cross

build:
	go build $(BUILD_FLAGS) -o $(BINARY) $(MAIN)

test:
	go test ./... -race -count=1

lint:
	go vet ./...

clean:
	rm -f $(BINARY) convert-linux-amd64 convert-darwin-amd64 convert-darwin-arm64 convert-windows-amd64.exe

# Cross-platform builds
cross: cross-linux cross-darwin-amd64 cross-darwin-arm64 cross-windows

cross-linux:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o convert-linux-amd64 $(MAIN)

cross-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o convert-darwin-amd64 $(MAIN)

cross-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o convert-darwin-arm64 $(MAIN)

cross-windows:
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o convert-windows-amd64.exe $(MAIN)
