BIN:= discord-version-api

all: build

build:
	@go mod tidy
	@go build -ldflags "-s -w"

check:
	@gofmt -w -s *.go
	@go test

clean:
	@go clean

.PHONY: build install check uninstall clean
