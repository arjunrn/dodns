.PHONY: lint build test
all: build

lint:
	golangci-lint run -c .golangci.yaml 

build:
	go build -o bin/dodns main.go

test:
	go test ./...
